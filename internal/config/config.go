package config

import (
	"fmt"

	"github.com/seth16888/coauth/internal/database"
	"github.com/seth16888/coauth/pkg/logger"

	"sync"

	"github.com/spf13/viper"
)

var (
	confVar   *Conf
	confMutex sync.Mutex
)

type Conf struct {
	Server Server
	DB     *database.DatabaseConfig
	Log    *logger.LogConfig
	Jwt    *Jwt
}

func GetConf() *Conf {
	confMutex.Lock()
	defer confMutex.Unlock()

	if confVar != nil {
		return confVar
	}
	return ReadConfigFromFile("")
}

type Jwt struct {
	Issuer         string
	ExpireTime     int
	MaxRefreshTime int
	SignKey        string
}

type Server struct {
	Grpc Grpc `json:"grpc"`
}

type Grpc struct {
	Addr    string `json:"addr"`
	Timeout int    `json:"timeout"`
}

func ReadConfigFromFile(file string) *Conf {
	if file == "" {
		file = "conf.yaml"
	}
	fmt.Println("read config from file: ", file)

	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath("~")

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	confVar = &Conf{}
	if err := viper.Unmarshal(confVar); err != nil {
		panic(err)
	}

	// watch
	viper.WatchConfig()

	return confVar
}
