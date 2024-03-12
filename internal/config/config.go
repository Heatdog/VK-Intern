package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ConfigStorage struct {
	Server      ServerListen   `mapstructure:"server_listen"`
	Postgre     PostgreStorage `mapstructure:"postgre_settings"`
	PasswordKey string         `mapstructure:"password_key"`
	Redis       RedisStorage   `mapstructure:"redis_settings"`
}

type ServerListen struct {
	IP   string `mapstructure:"ip"`
	Port string `mapstructure:"port"`
}

type PostgreStorage struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RedisStorage struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Password        string `mapstructure:"password"`
	TokenExparation int    `mapstructure:"token_expiration_days"`
}

func NewConfigStorage(logger *logrus.Logger) *ConfigStorage {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatalf("config file reading error: %s", err.Error())
	}
	res := &ConfigStorage{}
	if err := viper.Unmarshal(res); err != nil {
		logger.Fatalf("unmarshaling config file error: %s", err.Error())
	}
	return res
}
