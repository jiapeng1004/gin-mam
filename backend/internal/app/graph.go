package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-mam/backend/internal/domain/asset"
	assetmodel "github.com/gin-mam/backend/internal/domain/asset/model"
	"github.com/gin-mam/backend/internal/domain/sys"
	"github.com/gin-mam/backend/internal/domain/sys/model"
	infraconfig "github.com/gin-mam/backend/internal/infra/config"
	"github.com/gin-mam/backend/internal/infra/mysql"
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
	assetSvc := asset.NewService(assetRepo, configSvc, txMgr)
	assetHandler := asset.NewHandler(assetSvc)
	router := NewRouter(sysHandler, assetHandler, cfg.JWT.Secret)

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
		Router:       router,
	}, nil
}