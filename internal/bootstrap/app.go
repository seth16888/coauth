package bootstrap

import (
	"coauth/internal/di"
	"coauth/internal/server"
)

func StartApp(configFile string) error {
	appDeps := di.NewContainer(configFile)

	// 启动应用
	return server.Start(appDeps)
}
