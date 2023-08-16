package config

import (
	"log/slog"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Health   Health     `yaml:"health"`
	Backends []*Backend `yaml:"backends"`
}

type Backend struct {
	Addr string `yaml:"addr"`
	Name string `yaml:"name"`
}

type Health struct {
	Interval int `yaml:"interval"`
	Timeout  int `yaml:"timeout"`
}

type Server struct {
	Port           int `yaml:"port"`
	MaxConnections int `yaml:"max_connections"`
}

func LoadConfig(filePath string) error {
	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error reading config file")
	}

	viper.WatchConfig()

	return nil
}

func GetHealth() *Health {
	health := &Health{
		Interval: viper.GetInt("health.interval"),
		Timeout:  viper.GetInt("health.timeout"),
	}

	return health
}

func GetBackends() []*Backend {
	var backends []*Backend

	err := viper.UnmarshalKey("backends", &backends)
	if err != nil {
		slog.Error(err.Error(), "msg", "error unmarshaling backends")
	}

	return backends
}

func GetServerConfig() *Server {
	server := &Server{
		Port:           viper.GetInt("server.port"),
		MaxConnections: viper.GetInt("server.max_connections"),
	}

	return server
}
