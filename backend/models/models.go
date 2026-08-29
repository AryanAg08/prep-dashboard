package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Activity is a single daily study/activity log entry.
type Activity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date        string             `bson:"date" json:"date"`               // YYYY-MM-DD
	Topic       string             `bson:"topic" json:"topic"`             // DSA, LLD, etc.
	Title       string             `bson:"title" json:"title"`             // what you did
	Description string             `bson:"description" json:"description"` // details
	Hours       float64            `bson:"hours" json:"hours"`             // time spent
	Done        bool               `bson:"done" json:"done"`               // completed?
	CreatedAt   int64              `bson:"createdAt" json:"createdAt"`
}

// Note is a free-form note attached to a topic.
type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Topic     string             `bson:"topic" json:"topic"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Tags      []string           `bson:"tags" json:"tags"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt" json:"updatedAt"`
}

// Resource is a curated study material for a topic.
type Resource struct {
	Type  string `json:"type"` // book, playlist, sheet, course, article
	Title string `json:"title"`
	URL   string `json:"url"`
	Note  string `json:"note"`
}

// TopicGuide bundles curated resources + roadmap for one interview topic.
type TopicGuide struct {
	Key       string     `json:"key"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	Summary   string     `json:"summary"`
	Roadmap   []string   `json:"roadmap"`
	Resources []Resource `json:"resources"`
}

// LinkPreview is the metadata scraped from a DSA question URL.
type LinkPreview struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	SiteName    string `json:"siteName"`
	Favicon     string `json:"favicon"`
}
