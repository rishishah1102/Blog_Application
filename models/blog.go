package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogPost is the blog model
type BlogPost struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description" binding:"required"`
	Image       string             `json:"image" bson:"image" binding:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
