package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"prep-dashboard/backend/db"
	"prep-dashboard/backend/handlers"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	// Load .env if present (ignored in prod where real env vars are set).
	_ = godotenv.Load()

	mongoURI := env("MONGO_URI", "mongodb://localhost:27017")
	dbName := env("DB_NAME", "prep_dashboard")
	port := env("PORT", "9002")

	db.Connect(mongoURI, dbName)

	r := gin.Default()

	// Allow the React dev server (and configurable origins) to call the API.
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{env("CORS_ORIGIN", "http://localhost:9001")}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsCfg))

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
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
