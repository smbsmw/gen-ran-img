package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	GRPC_PORT = "GRPC_PORT"
	HTTP_PORT = "HTTP_PORT"
)

type Config struct {
	AppVersion string
	Server     Server
	Logger     Logger
	Kafka      Kafka
	Minio      Minio
	Http       Http
}

type Server struct {
	Port              string
	Development       bool
	Timeout           time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxConnectionIdle time.Duration
	MaxConnectionAge  time.Duration
	Kafka             Kafka
}

type Http struct {
	Port              string
	PprofPort         string
	Timeout           time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CookieLifeTime    int
	SessionCookieName string
}

// Logger config
type Logger struct {
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type Kafka struct {
	Brokers []string
}

type Minio struct {
	Addr string
}

func exportConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if os.Getenv("MODE") == "DOCKER" {
		viper.SetConfigName("config-docker.yaml")
	} else {
		viper.SetConfigName("config.yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// ParseConfig parses config file
func ParseConfig() (*Config, error) {
	if err := exportConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	gRPCPort := os.Getenv(GRPC_PORT)
	if gRPCPort != "" {
		cfg.Server.Port = gRPCPort
	}

	httpPort := os.Getenv(HTTP_PORT)
	if httpPort != "" {
		cfg.Http.Port = httpPort
	}

	return &cfg, nil
}
