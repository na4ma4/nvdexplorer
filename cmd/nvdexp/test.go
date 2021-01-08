package main

import (
	"github.com/na4ma4/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nolint: gochecknoglobals // cobra uses globals in main
var cmdTest = &cobra.Command{
	Use:   "test",
	Short: "Test Command",
	Run:   testCommand,
	Args:  cobra.NoArgs,
}

// nolint:gochecknoinits // init is used in main for cobra
func init() {
	rootCmd.AddCommand(cmdTest)
}

func testCommand(cmd *cobra.Command, args []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "nvdexplorer")

	logger, _ := cfg.ZapConfig().Build()
	defer logger.Sync() //nolint: errcheck

}
