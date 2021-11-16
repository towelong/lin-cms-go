package config

import (
	"fmt"
	"github.com/spf13/viper"
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
	}
	Lin struct {
		CMS struct {
			LoggerEnabled      bool   `mapstructure:"loggerEnabled"`
			TokenAccessExpire  int64  `mapstructure:"tokenAccessExpire"`
			TokenRefreshExpire int64  `mapstructure:"tokenRefreshExpire"`
			TokenSecret        string `mapstructure:"tokenSecret"`
		}
	}
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("parse config file: %s ", err))
	}
	viper.UnmarshalKey("server", &config.Server)
	viper.UnmarshalKey("mysql", &config.MySQL)
	viper.UnmarshalKey("lin", &config.Lin)
	viper.UnmarshalKey("cms", &config.Lin.CMS)
	return
}
