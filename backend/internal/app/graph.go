package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/catalog"
	catalogmodel "github.com/gin-mam/backend/internal/domain/catalog/model"
	assetmodel "github.com/gin-mam/backend/internal/domain/asset/model"
	"github.com/gin-mam/backend/internal/domain/sys"
	"github.com/gin-mam/backend/internal/domain/sys/model"
	infraconfig "github.com/gin-mam/backend/internal/infra/config"
	"github.com/gin-mam/backend/internal/infra/mysql"
	"github.com/gin-mam/backend/internal/infra/storage"
	redispkg "github.com/gin-mam/backend/internal/infra/redis"
	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Graph struct {
	DB            *gorm.DB
	Redis         *redis.Client
	TxMgr         tx.Manager
	ConfigSvc     infraconfig.Service
	SysRepo       sys.Repository
	SysSvc        sys.Service
	SysHandler    *sys.Handler
	AssetRepo     asset.Repository
	AssetSvc      asset.Service
	AssetHandler  *asset.Handler
	CatalogRepo    catalog.Repository
	CatalogSvc     catalog.Service
	CatalogHandler *catalog.Handler
	AssetUpload   asset.UploadService
	ObjectStorage storage.Storage
	Router        *gin.Engine
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
	objectStorage, err := storage.NewS3Storage(context.Background(), configSvc)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	assetSvc := asset.NewService(assetRepo, configSvc, txMgr)
	assetUpload := asset.NewUploadService(asset.NewRedisUploadStore(rdb), configSvc, objectStorage, assetSvc)
	assetHandler := asset.NewHandler(assetSvc, assetUpload)
	catalogRepo := catalog.NewMySQLRepository(db)
	catalogSvc := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogSvc)
	router := NewRouter(sysHandler, assetHandler, catalogHandler, cfg.JWT.Secret)

	return &Graph{
		DB:           db,
		Redis:        rdb,
		TxMgr:        txMgr,
		ConfigSvc:    configSvc,
		SysRepo:      sysRepo,
		SysSvc:       sysSvc,
		SysHandler:   sysHandler,
		AssetRepo:    assetRepo,
		AssetSvc:     assetSvc,
		AssetHandler: assetHandler,
		AssetUpload:  assetUpload,
		CatalogRepo:    catalogRepo,
		CatalogSvc:     catalogSvc,
		CatalogHandler: catalogHandler,
		ObjectStorage: objectStorage,
		Router:       router,
	}, nil
}