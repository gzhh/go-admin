package app

import (
	"fmt"
	"go-admin/app/admin/server"
	"go-admin/internal/lib/config"
	"os"

	"github.com/json-iterator/go/extra"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  ``,
	Run:   server.Run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	extra.RegisterFuzzyDecoders()

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	config.Init()
}
