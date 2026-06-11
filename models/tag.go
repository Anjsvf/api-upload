package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Slug      string             `json:"slug" bson:"slug"`
	Emoji     string             `json:"emoji" bson:"emoji"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}