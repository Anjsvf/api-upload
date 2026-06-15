package seed

import (
	"context"
	"log"
	"time"

	"github.com/Anjsvf/api-upload/database"
	"github.com/Anjsvf/api-upload/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var defaultTags = []models.Tag{
	
	{Name: "Jogos", Slug: "jogos", Emoji: "#"},
	{Name: "nofap/pornografia", Slug: "nofap-pornografia", Emoji: "#"},
	{Name: "Redes Sociais", Slug: "redes-sociais", Emoji: "#"},
	{Name: "Procastinação", Slug: "procastinacao", Emoji: "#"},

	
	// Vícios sexuais
	{Name: "Masturbação", Slug: "masturbacao", Emoji: "#"},
	{Name: "Sexo Compulsivo", Slug: "sexo-compulsivo", Emoji: "#"},
	{Name: "Motivação", Slug: "motivacao", Emoji: "#"},
 
	// Vícios em substâncias
	{Name: "Álcool", Slug: "alcool", Emoji: "#"},
	{Name: "Cigarro / Nicotina", Slug: "cigarro", Emoji: "#"},
	{Name: "Drogas", Slug: "drogas", Emoji: "#"},
	{Name: "Maconha", Slug: "maconha", Emoji: "#"},
	{Name: "Cafeína", Slug: "cafeina", Emoji: "#"},
	{Name: "Energéticos", Slug: "energeticos", Emoji: "#"},
 
	// Vícios comportamentais
	{Name: "Apostas", Slug: "apostas", Emoji: "#"},
	{Name: "Compras Compulsivas", Slug: "compras", Emoji: "#"},
	{Name: "Compulsão Alimentar", Slug: "alimentar", Emoji: "#"},
	{Name: "roer unhas", Slug: "roer-unhas", Emoji: "#"},
	// Outros
	{Name: "Outros", Slug: "outros", Emoji: "#"},
}


func RunTags() {
	collection := database.GetCollection("tags")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Erro ao verificar tags: %v", err)
		return
	}

	if count > 0 {
		log.Printf("✅ Tags já existem no banco (%d tags). Seed ignorado.", count)
		return
	}

	var docs []interface{}
	for _, tag := range defaultTags {
		tag.ID = primitive.NewObjectID()
		tag.CreatedAt = time.Now()
		docs = append(docs, tag)
	}

	if _, err := collection.InsertMany(ctx, docs); err != nil {
		log.Printf("Erro ao inserir tags: %v", err)
		return
	}

	log.Printf("✅ %d tags inseridas com sucesso!", len(defaultTags))
}