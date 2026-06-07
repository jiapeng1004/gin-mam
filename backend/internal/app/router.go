package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/sys"
	"github.com/gin-mam/backend/internal/pkg/httpx"
)

func NewRouter(sysHandler *sys.Handler, assetHandler *asset.Handler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(httpx.ErrorMiddleware(), gin.Logger())

	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", sysHandler.Login)

	authed := v1.Group("")
	authed.Use(httpx.JWTMiddleware(jwtSecret))
	authed.GET("/sys/user/page", sysHandler.UserPage)

	authed.POST("/asset/page", assetHandler.Page)
	authed.GET("/asset/:id", assetHandler.Get)
	authed.POST("/asset", assetHandler.Create)
	authed.PUT("/asset/:id", assetHandler.Update)
	authed.DELETE("/asset/:id", assetHandler.Delete)

	return r
}