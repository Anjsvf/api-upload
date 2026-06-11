package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostType define se a postagem é um vídeo do YouTube ou grupo do WhatsApp
type PostType string

const (
	PostTypeYouTube  PostType = "youtube"
	PostTypeWhatsApp PostType = "whatsapp"
)

type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Caption   string             `json:"caption" bson:"caption"`
	Link      string             `json:"link" bson:"link"`
	Type      PostType           `json:"type" bson:"type"`
	TagID     primitive.ObjectID `json:"tag_id" bson:"tag_id"`
	TagName   string             `json:"tag_name" bson:"tag_name"`
	TagSlug   string             `json:"tag_slug" bson:"tag_slug"`
	TagEmoji  string             `json:"tag_emoji" bson:"tag_emoji"`
	IP        string             `json:"-" bson:"ip"` 
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}


type CreatePostRequest struct {
	Title   string   `json:"title" binding:"required,min=3,max=100"`
	Caption string   `json:"caption" binding:"required,min=3,max=500"`
	Link    string   `json:"link" binding:"required"`
	Type    PostType `json:"type" binding:"required"`
	TagID   string   `json:"tag_id" binding:"required"`
}

// FeedFilter é usado para filtrar o feed
type FeedFilter struct {
	TagSlug string   `form:"tag"`
	Type    PostType `form:"type"`
	Page    int64    `form:"page"`
	Limit   int64    `form:"limit"`
}