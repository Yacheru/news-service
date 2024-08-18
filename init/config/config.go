package config

import (
	"github.com/disgoorg/snowflake/v2"
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

	WebhookID    snowflake.ID `mapstructure:"WEBHOOK_ID"`
	WebhookToken string       `mapstructure:"WEBHOOK_TOKEN"`
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

	if ServerConfig.APIPort == 0 || ServerConfig.APIEntry == "" ||
		ServerConfig.PostgresDSN == "" || ServerConfig.WebhookID == 0 || ServerConfig.WebhookToken == "" {
		logger.Error(constants.EmptyRequiredVar.Error(), constants.LoggerConfig)

		return constants.EmptyRequiredVar
	}

	return nil
}
