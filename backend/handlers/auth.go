package handlers

import (
	"crypto/subtle"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Password string `json:"password"`
}

// loginLimiter throttles password guesses per client IP. The login route is
// necessarily unauthenticated, so without this a human-chosen password is
// brute-forceable at network speed.
type loginLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	window   time.Duration
	max      int
}

func newLoginLimiter(max int, window time.Duration) *loginLimiter {
	return &loginLimiter{
		attempts: make(map[string][]time.Time),
		window:   window,
		max:      max,
	}
}

// allow records an attempt from ip and reports whether it stays under the cap.
func (l *loginLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := time.Now().Add(-l.window)
	kept := l.attempts[ip][:0]
	for _, t := range l.attempts[ip] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}

	if len(kept) >= l.max {
		l.attempts[ip] = kept
		return false
	}

	l.attempts[ip] = append(kept, time.Now())
	return true
}

// Login verifies a password against the configured secret. On success the
// client stores that same password and sends it as X-API-Key on every later
// request -- there is no session, so nothing to expire or revoke server-side.
func Login(password string) gin.HandlerFunc {
	want := []byte(password)
	limiter := newLoginLimiter(10, time.Minute)

	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many attempts, try again shortly",
			})
			return
		}

		var body loginRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if subtle.ConstantTimeCompare([]byte(body.Password), want) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
