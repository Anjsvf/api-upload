package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Anjsvf/api-upload/database"
	"github.com/Anjsvf/api-upload/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)


func GetTags(c *gin.Context) {
	collection := database.GetCollection("tags")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar tags"})
		return
	}
	defer cursor.Close(ctx)

	var tags []models.Tag
	if err := cursor.All(ctx, &tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar tags"})
		return
	}

	if tags == nil {
		tags = []models.Tag{}
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}