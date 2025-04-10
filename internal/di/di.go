package di

import (
	"coauth/internal/biz"
	"coauth/internal/config"
	"coauth/internal/data"
	"coauth/internal/database"
	"coauth/internal/service"
	"coauth/pkg/captcha"
	"coauth/pkg/jwt"
	"coauth/pkg/logger"
	"coauth/pkg/validator"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DI *Container
)

type Container struct {
	Config       *config.Conf      // 配置文件
	DB           *gorm.DB          // 数据库连接
	DbRepo       biz.AuthorizeRepo // 数据库仓库，用于操作数据库
	Log          *zap.Logger
	AuthUsecase  *biz.AuthorizeUseCase
	Svc          *service.CoauthService
	LoginUsecase *biz.LoginUsecase
}

func NewContainer(configFile string) *Container {
	conf := config.ReadConfigFromFile(configFile)
	log := logger.InitLogger(conf.Log)

	db, err := database.NewDB(conf.DB)
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}
	dbRepo := data.NewAuthorizeData(db, log)
	authUsecase := biz.NewAuthorizeUseCase(dbRepo, log)
	captchaStore := captcha.NewMemoryStore()
	captchaUsecase := biz.NewCaptchaUsecase(captchaStore, log)

	jwtExp := time.Duration(conf.Jwt.ExpireTime) * time.Minute
	jwtMaxRefresh := time.Duration(conf.Jwt.MaxRefreshTime) * time.Minute
	jwtSvc := jwt.NewJWTService(
		conf.Jwt.SignKey,
		conf.Jwt.Issuer,
		jwtExp,
		jwtMaxRefresh,
		time.Local,
		log,
	)

	akBlacklistRepo := data.NewTokenBlacklistData()
	akBlacklist := biz.NewTokenBlacklistUsecase(
		akBlacklistRepo,
		log,
	)

	loginRepo := data.NewLoginData(db, log)
	loginUsecase := biz.NewLoginUsecase(
		loginRepo,
		log,
		jwtSvc,
		akBlacklist,
	)

	svc := service.NewCoauthService(
		authUsecase,
		log,
		captchaUsecase,
		loginUsecase,
		validator.NewValidator(),
	)

	DI = &Container{
		Config:       conf,
		DB:           db,
		DbRepo:       dbRepo,
		Log:          log,
		AuthUsecase:  authUsecase,
		Svc:          svc,
		LoginUsecase: loginUsecase,
	}
	return DI
}
