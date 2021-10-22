package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	SecretKey     string `mapstructure:"SECRET_KEY"`
	ServerAddress string `mapstructure:"ADDRESS_SERVER"`
	DataBaseURL   string `mapstructure:"DATABASE_URL"`
	AwsAccessKey  string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey  string `mapstructure:"AWS_SECRET_KEY"`
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
