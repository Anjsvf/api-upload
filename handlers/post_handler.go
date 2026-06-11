package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Anjsvf/api-upload/database"
	"github.com/Anjsvf/api-upload/filters"
	"github.com/Anjsvf/api-upload/models"
	"github.com/Anjsvf/api-upload/validators"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var youtubeAPIKey string
 
func SetYoutubeAPIKey(key string) {
	youtubeAPIKey = key
}
 

func CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}
 
	
	if err := filters.CheckFields(req.Title, req.Caption); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
 
	
	if req.Type != models.PostTypeYouTube && req.Type != models.PostTypeWhatsApp {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tipo inválido. Use 'youtube' ou 'whatsapp'",
		})
		return
	}
 
	
	var normalizedLink string
	var linkErr error
 
	switch req.Type {
	case models.PostTypeYouTube:
		normalizedLink, linkErr = validators.ValidateYouTubeLink(req.Link, youtubeAPIKey)
	case models.PostTypeWhatsApp:
		normalizedLink, linkErr = validators.ValidateWhatsAppLink(req.Link)
	}
 
	if linkErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": linkErr.Error()})
		return
	}
 
	
	tagCollection := database.GetCollection("tags")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
 
	var tag models.Tag
 
	
	err := tagCollection.FindOne(ctx, bson.M{"slug": req.TagID}).Decode(&tag)
	if err != nil {
	
		tagID, idErr := primitive.ObjectIDFromHex(req.TagID)
		if idErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tag de vício não encontrada. Use o slug (ex: 'jogos') ou o ID da tag"})
			return
		}
		if err := tagCollection.FindOne(ctx, bson.M{"_id": tagID}).Decode(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tag de vício não encontrada"})
			return
		}
	}
 
	
	clientIP, _ := c.Get("client_ip")
 
	post := models.Post{
		ID:        primitive.NewObjectID(),
		Title:     req.Title,
		Caption:   req.Caption,
		Link:      normalizedLink,
		Type:      req.Type,
		TagID:     tag.ID,
		TagName:   tag.Name,
		TagSlug:   tag.Slug,
		TagEmoji:  tag.Emoji,
		IP:        clientIP.(string),
		CreatedAt: time.Now(),
	}
 
	postCollection := database.GetCollection("posts")
	if _, err := postCollection.InsertOne(ctx, post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar postagem"})
		return
	}
 
	c.JSON(http.StatusCreated, gin.H{
		"message": "Postagem criada com sucesso!",
		"post":    post,
	})
}
 

func GetFeed(c *gin.Context) {
	var filter models.FeedFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros inválidos"})
		return
	}
 

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 50 {
		filter.Limit = 20
	}
 

	query := bson.M{}
 
	if filter.TagSlug != "" {
		query["tag_slug"] = filter.TagSlug
	}
 
	if filter.Type != "" {
		if filter.Type != models.PostTypeYouTube && filter.Type != models.PostTypeWhatsApp {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tipo inválido. Use 'youtube' ou 'whatsapp'",
			})
			return
		}
		query["type"] = filter.Type
	}
 
	collection := database.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
 
	skip := (filter.Page - 1) * filter.Limit
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(filter.Limit)
 
	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar postagens"})
		return
	}
	defer cursor.Close(ctx)
 
	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar postagens"})
		return
	}
 
	
	total, _ := collection.CountDocuments(ctx, query)
 
	if posts == nil {
		posts = []models.Post{}
	}
 
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":  filter.Page,
			"limit": filter.Limit,
			"total": total,
		},
	})
}