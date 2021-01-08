package mainconfig

import (
	"os"
	"runtime"
	"strings"

	"github.com/na4ma4/config"
	"github.com/spf13/viper"
)

// ConfigInit is the common config initialisation for the commands.
func ConfigInit() {
	viper.SetConfigName("nvmexplorer")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./artifacts")
	viper.AddConfigPath("./test")
	viper.AddConfigPath("$HOME/.nvmexplorer")
	viper.AddConfigPath("/etc/nvmexplorer")
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath("/usr/local/nvmexplorer/etc")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("/run/secrets")
	viper.AddConfigPath("/etc/nsca")
	viper.AddConfigPath("/etc/nagios")
	viper.AddConfigPath(".")

	viper.SetDefault("nvd.cache", "$HOME/.nvdexplorer/cache")

	_ = viper.ReadInConfig()
}

// UserHomeDir returns the users home directory.
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// ResolveEnvironment resolves environment variables in path options.
func ResolveEnvironment(cfg config.Conf) {
	inPath := cfg.GetString("nvd.cache")
	if inPath == "$HOME" || strings.HasPrefix(inPath, "$HOME"+string(os.PathSeparator)) {
		cfg.SetString("nvd.cache", UserHomeDir()+inPath[5:])
	}
}
