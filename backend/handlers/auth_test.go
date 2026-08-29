package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func postLogin(t *testing.T, r *gin.Engine, body string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

func loginRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/login", Login("hunter2"))
	return r
}

func TestLogin(t *testing.T) {
	cases := []struct {
		name string
		body string
		want int
	}{
		{"correct password", `{"password":"hunter2"}`, http.StatusOK},
		{"wrong password", `{"password":"nope"}`, http.StatusUnauthorized},
		{"empty password", `{"password":""}`, http.StatusUnauthorized},
		{"missing field", `{}`, http.StatusUnauthorized},
		{"malformed json", `not json`, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := postLogin(t, loginRouter(), tc.body); got != tc.want {
				t.Errorf("status = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestLoginRateLimited(t *testing.T) {
	r := loginRouter()

	// The limiter allows 10 attempts per minute per IP; httptest reuses one.
	for i := 0; i < 10; i++ {
		if got := postLogin(t, r, `{"password":"nope"}`); got != http.StatusUnauthorized {
			t.Fatalf("attempt %d: status = %d, want %d", i+1, got, http.StatusUnauthorized)
		}
	}

	if got := postLogin(t, r, `{"password":"nope"}`); got != http.StatusTooManyRequests {
		t.Errorf("11th attempt: status = %d, want %d", got, http.StatusTooManyRequests)
	}

	// The correct password is refused too while throttled -- otherwise the
	// limiter would be trivially bypassed by guessing alongside a valid login.
	if got := postLogin(t, r, `{"password":"hunter2"}`); got != http.StatusTooManyRequests {
		t.Errorf("throttled valid login: status = %d, want %d", got, http.StatusTooManyRequests)
	}
}
