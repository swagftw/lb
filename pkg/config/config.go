package config

import (
	"log/slog"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// reading config is not thread safe
// using mutexes to use locks
var lock = sync.Mutex{}

type Backend struct {
	Addr string `yaml:"addr"`
	ID   int    `yaml:"id"`
}

type Health struct {
	Interval int `yaml:"interval"`
	Timeout  int `yaml:"timeout"`
}

type Server struct {
	Port  int  `yaml:"port"`
	Debug bool `yaml:"debug"`
}

// LoadConfig loads the config from config.yaml file located at the root
func LoadConfig(filePath string) error {
	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("failed to read config", "err", err)

		return errors.Wrap(err, "error reading config file")
	}

	// keep watching for config changes
	//
	// we are doing this to observe changes mainly from updating backends config
	viper.WatchConfig()

	return nil
}

// GetHealth gets health check config
func GetHealth() *Health {
	lock.Lock()
	defer lock.Unlock()

	health := &Health{
		Interval: viper.GetInt("health.interval"),
		Timeout:  viper.GetInt("health.timeout"),
	}

	return health
}

// GetBackends gets the list of backends
func GetBackends() []*Backend {
	lock.Lock()
	defer lock.Unlock()

	var backends []*Backend

	err := viper.UnmarshalKey("backends", &backends)
	if err != nil {
		slog.Error("failed to unmarshall backends config", "msg", err)

		return nil
	}

	return backends
}

// GetBalancerConfig gets the config for running load balancer
func GetBalancerConfig() *Server {
	lock.Lock()
	defer lock.Unlock()

	server := &Server{
		Port:  viper.GetInt("balancer.port"),
		Debug: viper.GetBool("balancer.debug"),
	}

	return server
}
