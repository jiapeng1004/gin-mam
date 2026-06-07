package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/catalog"
	"github.com/gin-mam/backend/internal/domain/search"
	"github.com/gin-mam/backend/internal/domain/sys"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

func NewRouter(sysHandler *sys.Handler, assetHandler *asset.Handler, catalogHandler *catalog.Handler, searchHandler *search.Handler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(httpx.ErrorMiddleware(), gin.Logger())

	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", sysHandler.Login)

	authed := v1.Group("")
	authed.Use(httpx.JWTMiddleware(jwtSecret))
	authed.GET("/sys/user/page", sysHandler.UserPage)

	authed.GET("/catalog/tree", catalogHandler.Tree)
	authed.POST("/catalog", catalogHandler.Create)
	authed.PUT("/catalog/:id", catalogHandler.Update)
	authed.DELETE("/catalog/:id", catalogHandler.Delete)
	authed.GET("/catalog/:id/config", catalogHandler.ListConfig)
	authed.PUT("/catalog/:id/config", catalogHandler.BatchUpdateConfig)

	authed.POST("/asset/page", assetHandler.Page)
	authed.GET("/asset/:id", assetHandler.Get)
	authed.POST("/asset", assetHandler.Create)
	authed.PUT("/asset/:id", assetHandler.Update)
	authed.DELETE("/asset/:id", assetHandler.Delete)

	authed.POST("/asset/upload/init", assetHandler.UploadInit)
	authed.POST("/asset/upload/chunk", assetHandler.UploadChunk)
	authed.POST("/asset/upload/complete", assetHandler.UploadComplete)

	authed.POST("/search/asset", searchHandler.SearchAsset)

	return r
}
