package cmd

import (
	"coauth/internal/bootstrap"

	"github.com/spf13/cobra"
)

var (
	configFile string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "conf", "c",
		"conf/conf.yaml", "--conf config file (default is conf/conf.yaml)")
}

var rootCmd = &cobra.Command{
	Use:   "coauth [command] [flags] [args]",
	Short: "coauth is a OAuth2 server",
	Long:  `coauth is a OAuth2 server. It provides api for client to authorize and get access token.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bootstrap.StartApp(configFile)
	},
}

func Execute() error {
	return rootCmd.Execute()
}
