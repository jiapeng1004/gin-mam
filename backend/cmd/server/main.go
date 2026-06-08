// GIN-MAM 媒资管理系统 HTTP 服务入口。
//
//	@title						GIN-MAM API
//	@version					1.0
//	@description				媒资管理系统 REST API。成功响应直接返回业务 JSON（camelCase）；失败返回 { err_code, err_msg }。
//	@host						localhost:8080
//	@BasePath					/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Bearer Token，格式：Bearer {token}
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/app"

	_ "github.com/gin-mam/backend/docs" // swagger 生成文档
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
