package main

import (
	"log"

	"github.com/Anjsvf/api-upload/config"
	"github.com/Anjsvf/api-upload/database"
	"github.com/Anjsvf/api-upload/handlers"
	"github.com/Anjsvf/api-upload/routes"
	"github.com/Anjsvf/api-upload/seed"
)

func main() {

	cfg := config.Load()

	
	database.Connect(cfg.MongoURI, cfg.MongoDBName)

	
	handlers.SetYoutubeAPIKey(cfg.YoutubeAPIKey)

	
	seed.RunTags()

	
	r := routes.Setup()
	log.Printf("Servidor rodando na alturas na  porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}