package di

import (
	"coauth/internal/biz"
	"coauth/internal/config"
	"coauth/internal/data"
	"coauth/internal/database"
	"coauth/internal/service"
	"coauth/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DI *Container
)

type Container struct {
	Config      *config.Conf      // 配置文件
	DB          *gorm.DB          // 数据库连接
	DbRepo      biz.AuthorizeRepo // 数据库仓库，用于操作数据库
	Log         *zap.Logger
	AuthUsecase *biz.AuthorizeUseCase
	Svc         *service.CoauthService
}

func NewContainer(configFile string) *Container {
	conf := config.ReadConfigFromFile(configFile)
	log := logger.InitLogger(conf.Log)
	db, err := database.NewDB(conf.DB)
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}
	dbRepo := data.NewAuthorizeData(db)
	authUsecase := biz.NewAuthorizeUseCase(dbRepo, log)
	svc := service.NewCoauthService(authUsecase, log)

	DI = &Container{
		Config:      conf,
		DB:          db,
		DbRepo:      dbRepo,
		Log:         log,
		AuthUsecase: authUsecase,
		Svc:         svc,
	}
	return DI
}
