// Package middleware holds gin middleware shared across routes.
package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HeaderName is the request header carrying the shared secret.
const HeaderName = "X-API-Key"

// APIKeyAuth rejects any request whose X-API-Key header does not match the
// configured secret. The comparison is constant-time so that response timing
// does not leak how much of the key a guess got right.
//
// Preflight requests are let through untouched: browsers never attach custom
// headers to an OPTIONS preflight, so checking one here would fail every
// cross-origin write before the real request is ever sent.
func APIKeyAuth(key string) gin.HandlerFunc {
	want := []byte(key)

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		got := []byte(c.GetHeader(HeaderName))
		if subtle.ConstantTimeCompare(got, want) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or missing " + HeaderName,
			})
			return
		}

		c.Next()
	}
}
