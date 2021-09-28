package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	SecretKey     string `mapstructure:"SECRET_KEY"`
	ServerAddress string `mapstructure:"ADDRESS_SERVER"`
	DataBaseURL string `mapstructure:"DATABASE_URL"`
}

var ConfigApp Config

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&ConfigApp)
	if err != nil {
		return err
	}
	return nil
}
