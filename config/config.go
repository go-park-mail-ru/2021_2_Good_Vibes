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

	MainConfig struct {
		SecretKey     string `mapstructure:"SECRET_KEY"`
		ServerAddress string `mapstructure:"ADDRESS_SERVER"`
	} `mapstructure:"MAIN_CONFIG"`

	DataBase struct {
		User     string `mapstructure:"USER"`
		DBName   string `mapstructure:"DBNAME"`
		Password string `mapstructure:"PASSWORD"`
		Host     string `mapstructure:"HOST"`
		Port     string `mapstructure:"PORT"`
	} `mapstructure:"DATA_BASE_CONFIG"`

	Email struct {
		Server     string `mapstructure:"EMAIL_SERVER"`
		ServerPort int    `mapstructure:"EMAIL_SERVER_PORT"`
		Address    string `mapstructure:"EMAIL_ADDRESS"`
		Password   string `mapstructure:"EMAIL_ADDRESS_PASSWORD"`
	} `mapstructure:"EMAIL_CONFIG"`
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
