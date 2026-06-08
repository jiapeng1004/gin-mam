package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/asset"
	assetmodel "github.com/gin-mam/backend/internal/domain/asset/model"
	"github.com/gin-mam/backend/internal/domain/catalog"
	catalogmodel "github.com/gin-mam/backend/internal/domain/catalog/model"
	"github.com/gin-mam/backend/internal/domain/search"
	"github.com/gin-mam/backend/internal/domain/sys"
	"github.com/gin-mam/backend/internal/domain/sys/model"
	"github.com/gin-mam/backend/internal/domain/transcode"
	transcodeclient "github.com/gin-mam/backend/internal/domain/transcode/client"
	transcodemodel "github.com/gin-mam/backend/internal/domain/transcode/model"
	"github.com/gin-mam/backend/internal/domain/workflow"
	workflowmodel "github.com/gin-mam/backend/internal/domain/workflow/model"
	infraconfig "github.com/gin-mam/backend/internal/infra/config"
	"github.com/gin-mam/backend/internal/infra/elasticsearch"
	"github.com/gin-mam/backend/internal/infra/lock"
	"github.com/gin-mam/backend/internal/infra/mysql"
	redispkg "github.com/gin-mam/backend/internal/infra/redis"
	"github.com/gin-mam/backend/internal/infra/storage"
	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Graph struct {
	DB               *gorm.DB
	Redis            *redis.Client
	TxMgr            tx.Manager
	ConfigSvc        infraconfig.Service
	SysRepo          sys.Repository
	SysSvc           sys.Service
	SysHandler       *sys.Handler
	AssetRepo        asset.Repository
	AssetSvc         asset.Service
	AssetHandler     *asset.Handler
	CatalogRepo      catalog.Repository
	CatalogSvc       catalog.Service
	CatalogHandler   *catalog.Handler
	SearchSvc        search.Service
	SearchHandler    *search.Handler
	WorkflowRepo     workflow.Repository
	WorkflowSvc      workflow.Service
	WorkflowHandler  *workflow.Handler
	TranscodeRepo     transcode.Repository
	TranscodeSvc      transcode.Service
	TranscodeHandler  *transcode.Handler
	AssetUpload      asset.UploadService
	ObjectStorage    storage.Storage
	Router           *gin.Engine
}

func Build(cfg *Config) (*Graph, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	db, err := mysql.NewMySQL(cfg.MySQL.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Org{},
		&model.Menu{},
		&infraconfig.SysConfig{},
		&assetmodel.Asset{},
		&assetmodel.AssetFile{},
		&assetmodel.AssetMetadata{},
		&catalogmodel.Catalog{},
		&catalogmodel.CatalogConfig{},
		&workflowmodel.WorkflowDef{},
		&workflowmodel.WorkflowLevelUser{},
		&workflowmodel.WorkflowInstance{},
		&workflowmodel.WorkflowInstanceLevelUser{},
		&workflowmodel.WorkflowOperate{},
		&transcodemodel.TranscodeGroup{},
		&transcodemodel.TranscodeProfile{},
		&transcodemodel.TranscodeTask{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	rdb := redispkg.NewRedis(cfg.Redis.Addr, cfg.Redis.DB)
	txMgr := tx.NewManager(db)
	configSvc := infraconfig.NewService(db, rdb)
	sysRepo := sys.NewMySQLRepository(db)
	sysSvc := sys.NewService(sysRepo, sys.JWTConfig{
		Secret:      cfg.JWT.Secret,
		ExpireHours: cfg.JWT.ExpireHours,
	})
	sysHandler := sys.NewHandler(sysSvc)
	assetRepo := asset.NewMySQLRepository(db)
	objectStorage := storage.NewS3Storage(configSvc)
	assetSvc := asset.NewService(assetRepo, configSvc, txMgr, nil)
	esClient, err := elasticsearch.NewClient(cfg.Elasticsearch.Addresses)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: %w", err)
	}
	searchSvc := search.NewService(esClient, assetSvc)
	asset.WireIndexer(assetSvc, searchSvc)
	searchHandler := search.NewHandler(searchSvc)
	assetUpload := asset.NewUploadService(asset.NewRedisUploadStore(rdb), configSvc, objectStorage, assetSvc)
	assetHandler := asset.NewHandler(assetSvc, assetUpload)
	catalogRepo := catalog.NewMySQLRepository(db)
	catalogSvc := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogSvc)
	workflowRepo := workflow.NewMySQLRepository(db)
	workflowLocker := lock.NewRedisLocker(rdb)
	workflowSvc := workflow.NewService(workflowRepo, assetRepo, assetSvc, txMgr, workflowLocker)
	workflowHandler := workflow.NewHandler(workflowSvc)
	transcodeRepo := transcode.NewMySQLRepository(db)
	transcodeClient := transcodeclient.NewHTTPClient()
	transcodeSvc := transcode.NewService(transcodeRepo, configSvc, txMgr, transcodeClient)
	transcodeHandler := transcode.NewHandler(transcodeSvc)
	router := NewRouter(sysHandler, assetHandler, catalogHandler, searchHandler, workflowHandler, transcodeHandler, cfg.JWT.Secret)

	return &Graph{
		DB:              db,
		Redis:           rdb,
		TxMgr:           txMgr,
		ConfigSvc:       configSvc,
		SysRepo:         sysRepo,
		SysSvc:          sysSvc,
		SysHandler:      sysHandler,
		AssetRepo:       assetRepo,
		AssetSvc:        assetSvc,
		AssetHandler:    assetHandler,
		AssetUpload:     assetUpload,
		CatalogRepo:     catalogRepo,
		CatalogSvc:      catalogSvc,
		CatalogHandler:  catalogHandler,
		SearchSvc:       searchSvc,
		SearchHandler:   searchHandler,
		WorkflowRepo:    workflowRepo,
		WorkflowSvc:     workflowSvc,
		WorkflowHandler: workflowHandler,
		TranscodeRepo:    transcodeRepo,
		TranscodeSvc:     transcodeSvc,
		TranscodeHandler: transcodeHandler,
		ObjectStorage:   objectStorage,
		Router:          router,
	}, nil
}
