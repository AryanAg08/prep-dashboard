package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// serve runs one request through APIKeyAuth and reports the status code plus
// whether the guarded handler actually ran.
func serve(t *testing.T, method, header string) (int, bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	reached := false
	r := gin.New()
	r.Use(APIKeyAuth("correct-secret"))
	r.Handle(method, "/api/thing", func(c *gin.Context) {
		reached = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(method, "/api/thing", nil)
	if header != "" {
		req.Header.Set(HeaderName, header)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w.Code, reached
}

func TestAPIKeyAuth(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		header     string
		wantStatus int
		wantReach  bool
	}{
		{"correct key passes", http.MethodGet, "correct-secret", http.StatusOK, true},
		{"missing header rejected", http.MethodGet, "", http.StatusUnauthorized, false},
		{"wrong key rejected", http.MethodGet, "nope", http.StatusUnauthorized, false},
		{"prefix of key rejected", http.MethodGet, "correct", http.StatusUnauthorized, false},
		{"key with suffix rejected", http.MethodGet, "correct-secretX", http.StatusUnauthorized, false},
		{"writes are guarded too", http.MethodPost, "", http.StatusUnauthorized, false},
		{"preflight passes without key", http.MethodOptions, "", http.StatusOK, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, reached := serve(t, tc.method, tc.header)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d", status, tc.wantStatus)
			}
			if reached != tc.wantReach {
				t.Errorf("handler reached = %v, want %v", reached, tc.wantReach)
			}
		})
	}
}
