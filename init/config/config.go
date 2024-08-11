package config

import (
	"github.com/spf13/viper"

	"news-service/init/logger"
	"news-service/pkg/constants"
)

var ServerConfig Config

type Config struct {
	APIPort  int    `mapstructure:"API_PORT"`
	APIDebug bool   `mapstructure:"API_DEBUG"`
	APIEntry string `mapstructure:"API_ENTRY"`

	PostgresDSN string `mapstructure:"POSTGRESQL_DSN"`
}

func InitConfig() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(err.Error(), constants.LoggerConfig)

		return err
	}

	if err := viper.Unmarshal(&ServerConfig); err != nil {
		logger.Error(err.Error(), constants.LoggerConfig)

		return err
	}

	if ServerConfig.APIPort == 0 {
		logger.Error(constants.EmptyRequiredVar.Error(), constants.LoggerConfig)

		return constants.EmptyRequiredVar
	}

	return nil
}
