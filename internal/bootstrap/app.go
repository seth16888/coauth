package bootstrap

import (
	"github.com/seth16888/coauth/internal/di"
	"github.com/seth16888/coauth/internal/server"
)

func StartApp(configFile string) error {
	appDeps := di.NewContainer(configFile)

	// 启动应用
	return server.Start(appDeps)
}
