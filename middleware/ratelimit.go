package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Anjsvf/api-upload/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const maxPostsPerDay = 4


func GetClientIP(c *gin.Context) string {
	
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
	
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return strings.TrimSpace(ip)
	}

	return c.ClientIP()
}


func RateLimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetClientIP(c)

		collection := database.GetCollection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Conta postagens do IP nas últimas 24 horas
		since := time.Now().Add(-24 * time.Hour)
		count, err := collection.CountDocuments(ctx, bson.M{
			"ip":         ip,
			"created_at": bson.M{"$gte": since},
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro interno ao verificar limite de postagens",
			})
			c.Abort()
			return
		}

		if count >= maxPostsPerDay {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Limite de postagens atingido",
				"message": "Você já fez 4 postagens hoje. Tente novamente amanhã.",
				"limit":   maxPostsPerDay,
				"used":    count,
			})
			c.Abort()
			return
		}

		// Passa o IP adiante no contexto para o handler usar
		c.Set("client_ip", ip)
		c.Next()
	}
}