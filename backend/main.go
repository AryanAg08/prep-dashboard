package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"prep-dashboard/backend/db"
	"prep-dashboard/backend/handlers"
	"prep-dashboard/backend/middleware"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// splitList parses a comma-separated env value into trimmed, non-empty items.
func splitList(v string) []string {
	var out []string
	for _, part := range strings.Split(v, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func main() {
	// Load .env if present (ignored in prod where real env vars are set).
	_ = godotenv.Load()

	mongoURI := env("MONGO_URI", "mongodb://localhost:27017")
	dbName := env("DB_NAME", "prep_dashboard")
	port := env("PORT", "9002")

	// Fail fast rather than booting an unprotected API: an empty password
	// would make the constant-time compare match a request that sends none.
	password := os.Getenv("APP_PASSWORD")
	if password == "" {
		log.Fatal("APP_PASSWORD is required (generate one with: openssl rand -hex 32)")
	}

	db.Connect(mongoURI, dbName)

	r := gin.Default()

	// Allow the configured origins to call the API. Note that the browser's
	// Origin survives the vite/nginx proxy hop (changeOrigin rewrites Host,
	// not Origin), so the deployed domain must be listed here or gin-contrib
	// /cors answers every request with a bare 403.
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = splitList(env("CORS_ORIGIN", "http://localhost:9001"))
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", middleware.HeaderName}
	r.Use(cors.New(corsCfg))

	// Health check stays open so uptime probes don't need the key.
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Login is necessarily open -- it is how a client obtains the secret to
	// use everywhere else. It is rate limited per IP.
	r.POST("/api/login", handlers.Login(password))

	// Everything else requires the shared secret, reads included.
	api := r.Group("/api", middleware.APIKeyAuth(password))
	{
		api.GET("/guides", handlers.GetGuides)

		api.GET("/activities", handlers.ListActivities)
		api.POST("/activities", handlers.CreateActivity)
		api.PATCH("/activities/:id", handlers.UpdateActivity)
		api.DELETE("/activities/:id", handlers.DeleteActivity)
		api.GET("/stats", handlers.ActivityStats)

		api.GET("/notes", handlers.ListNotes)
		api.POST("/notes", handlers.CreateNote)
		api.PATCH("/notes/:id", handlers.UpdateNote)
		api.DELETE("/notes/:id", handlers.DeleteNote)

		api.POST("/link-preview", handlers.LinkPreview)
	}

	log.Printf("server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
