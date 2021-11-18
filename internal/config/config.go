package config

import (
	"fmt"
	"path"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/towelong/lin-cms-go/internal/pkg/log"
	"github.com/towelong/lin-cms-go/pkg"
)

type Config struct {
	Server struct {
		Port      int `mapstructure:"port"`
		PprofPort int `mapstructure:"pprof_port"`
	}
	MySQL struct {
		Port     int    `mapstructure:"port"`
		Account  string `mapstructure:"account"`
		Password int    `mapstructure:"password"`
		Addr     string `mapstructure:"addr"`
		Database string `mapstructure:"database"`
	}
	Lin struct {
		CMS struct {
			LoggerEnabled      bool   `mapstructure:"loggerEnabled"`
			TokenAccessExpire  int64  `mapstructure:"tokenAccessExpire"`
			TokenRefreshExpire int64  `mapstructure:"tokenRefreshExpire"`
			TokenSecret        string `mapstructure:"tokenSecret"`
		}
		File struct {
			Domain      string   `mapstructure:"domain"`
			Exclude     []string `mapstructure:"exclude"`
			Include     []string `mapstructure:"include"`
			Nums        int      `mapstructure:"nums"`
			SingleLimit int      `mapstructure:"singleLimit"`
			StoreDir    string   `mapstructure:"storeDir"`
		}
	}
}

func LoadConfig() Config {
	pflag.String("env", "dev", "application environment")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AddConfigPath(path.Join(pkg.GetCurrentAbPath(), "/config"))
	if viper.GetString("env") == "dev" {
		log.Logger.Info("当前环境为开发环境")
		viper.SetConfigName("config.dev")
	}
	if viper.GetString("env") == "prod" {
		log.Logger.Info("当前环境为生产环境")
		viper.SetConfigName("config.prod")
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("parse config file: %s ", err))
	}
	var config Config
	viper.UnmarshalKey("server", &config.Server)
	viper.UnmarshalKey("mysql", &config.MySQL)
	viper.UnmarshalKey("lin", &config.Lin)
	viper.UnmarshalKey("cms", &config.Lin.CMS)
	viper.UnmarshalKey("file", &config.Lin.File)
	return config
}
