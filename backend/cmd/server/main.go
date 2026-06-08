package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/app"
)

func main() {
	configPath := flag.String("config", "./config.yaml", "path to config file")
	flag.Parse()

	cfg, err := app.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	graph, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("build graph: %v", err)
	}

	ctx := context.Background()
	if err := graph.SysSvc.EnsureSeedAdmin(ctx); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("listening on %s", addr)
	if err := graph.Router.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}