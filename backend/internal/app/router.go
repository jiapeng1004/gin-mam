package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/catalog"
	"github.com/gin-mam/backend/internal/domain/search"
	"github.com/gin-mam/backend/internal/domain/sys"
	"github.com/gin-mam/backend/internal/domain/transcode"
	"github.com/gin-mam/backend/internal/domain/workflow"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter 注册全部 HTTP 路由与 Swagger UI。
func NewRouter(sysHandler *sys.Handler, assetHandler *asset.Handler, catalogHandler *catalog.Handler, searchHandler *search.Handler, workflowHandler *workflow.Handler, transcodeHandler *transcode.Handler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(httpx.ErrorMiddleware(), gin.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", sysHandler.Login)
	v1.POST("/transcode/callback", transcodeHandler.Callback)

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

	authed.POST("/asset-workflow/submit", workflowHandler.Submit)
	authed.POST("/asset-workflow/audit", workflowHandler.Audit)
	authed.POST("/asset-workflow/multi-audit", workflowHandler.MultiAudit)
	authed.POST("/asset-workflow/revoke", workflowHandler.Revoke)
	authed.POST("/asset-workflow/assigned-to-me", workflowHandler.AssignedToMe)
	authed.POST("/asset-workflow/created-by-me", workflowHandler.CreatedByMe)
	authed.POST("/asset-workflow/audited-by-me", workflowHandler.AuditedByMe)
	authed.POST("/asset-workflow/all", workflowHandler.All)

	authed.POST("/workflow/def/page", workflowHandler.DefPage)
	authed.POST("/workflow/def", workflowHandler.CreateDef)

	authed.POST("/transcode/group/page", transcodeHandler.GroupPage)
	authed.POST("/transcode/group", transcodeHandler.CreateGroup)
	authed.PUT("/transcode/group/:id", transcodeHandler.UpdateGroup)
	authed.DELETE("/transcode/group/:id", transcodeHandler.DeleteGroup)
	authed.POST("/transcode/task", transcodeHandler.CreateTask)
	authed.GET("/transcode/task/:id", transcodeHandler.GetTask)
	authed.POST("/transcode/task/page", transcodeHandler.TaskPage)

	return r
}
