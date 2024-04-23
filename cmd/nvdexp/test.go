package main

import (
	"github.com/na4ma4/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdTest = &cobra.Command{
	Use:   "test",
	Short: "Test Command",
	Run:   testCommand,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(cmdTest)
}

func testCommand(_ *cobra.Command, _ []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "nvdexplorer")

	logger, _ := cfg.ZapConfig().Build()
	defer logger.Sync()
}
