package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"prep-dashboard/backend/db"
	"prep-dashboard/backend/models"
)

// ListActivities returns activities, newest first, optionally filtered by ?topic= or ?date=.
func ListActivities(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if t := c.Query("topic"); t != "" {
		filter["topic"] = t
	}
	if d := c.Query("date"); d != "" {
		filter["date"] = d
	}

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := db.Collection("activities").Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	activities := []models.Activity{}
	if err := cur.All(ctx, &activities); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

// CreateActivity inserts a new activity log entry.
func CreateActivity(c *gin.Context) {
	var a models.Activity
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if a.Date == "" {
		a.Date = time.Now().Format("2006-01-02")
	}
	a.ID = primitive.NewObjectID()
	a.CreatedAt = time.Now().UnixMilli()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := db.Collection("activities").InsertOne(ctx, a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

// UpdateActivity edits fields of an existing activity (used to toggle Done, edit text).
func UpdateActivity(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	delete(body, "id")
	delete(body, "_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.Collection("activities").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DeleteActivity removes an activity.
func DeleteActivity(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := db.Collection("activities").DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ActivityStats returns aggregate counts for the dashboard header.
func ActivityStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection("activities")
	total, _ := coll.CountDocuments(ctx, bson.M{})
	done, _ := coll.CountDocuments(ctx, bson.M{"done": true})
	today, _ := coll.CountDocuments(ctx, bson.M{"date": time.Now().Format("2006-01-02")})

	// Sum hours + per-topic breakdown via aggregation.
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$topic",
			"count": bson.M{"$sum": 1},
			"hours": bson.M{"$sum": "$hours"},
		}},
	}
	cur, err := coll.Aggregate(ctx, pipeline)
	perTopic := []bson.M{}
	var totalHours float64
	if err == nil {
		_ = cur.All(ctx, &perTopic)
		cur.Close(ctx)
		for _, row := range perTopic {
			if h, ok := row["hours"].(float64); ok {
				totalHours += h
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total":      total,
		"done":       done,
		"today":      today,
		"totalHours": totalHours,
		"perTopic":   perTopic,
	})
}
