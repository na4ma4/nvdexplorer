package main

import (
	"os"

	"github.com/na4ma4/config"
	"github.com/na4ma4/nvdexplorer/internal/mainconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cmdSync = &cobra.Command{
	Use:   "sync",
	Short: "Synchronise NVD Databases",
	Run:   syncCommand,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(cmdSync)
}

func syncCommand(_ *cobra.Command, _ []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "nvdexplorer")
	mainconfig.ResolveEnvironment(cfg)

	logger, _ := cfg.ZapConfig().Build()
	defer logger.Sync()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	logger.Info("cache dir", zap.String("nvd.cache", cfg.GetString("nvd.cache")))
	if err := os.MkdirAll(cfg.GetString("nvd.cache"), 0755); err != nil {
		logger.Panic("unable to create cache directory", zap.Error(err))
	}

	// TODO fix this
	// t := time.Now()
	// for i := 2002; i <= t.Year(); i++ {
	// 	logger.Info("processing year", zap.Int("year", i))
	// 	n := nvdcache.NewNVDByYear(cfg, i)
	// 	if err := n.Download(ctx); err != nil {
	// 		logger.Panic("unable to update cached data", zap.Error(err))
	// 	}
	// 	os.Exit(0) //TODO testing
	// }
}
