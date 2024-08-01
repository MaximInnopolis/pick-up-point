package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var defaultGrpcPort = "50051"
var defaultPrometheusPort = "9090"

type KafkaConfig struct {
	BrokerList []string
	Topic      string
}

type ServerConfig struct {
	GrpcPort string
}

type PrometheusConfig struct {
	PrometheusPort string
}

type Config struct {
	DbUrl            string
	OutputMode       string
	CacheTTL         time.Duration
	KafkaConfig      KafkaConfig
	ServerConfig     ServerConfig
	PrometheusConfig PrometheusConfig
}

func New() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL не задан")
	}

	strCacheTTL := os.Getenv("CACHE_TTL")
	if strCacheTTL == "" {
		return nil, fmt.Errorf("CACHE_TTL не задан")
	}
	cacheTTL, err := time.ParseDuration(strCacheTTL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге CACHE_TTL: %w", err)
	}

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS не задано")
	}

	brokerList := strings.Split(brokers, ",")
	topic := os.Getenv("KAFKA_TOPIC")

	outputMode := os.Getenv("OUTPUT_MODE")
	if outputMode == "" {
		outputMode = "stdout"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = defaultGrpcPort
	}

	promPort := os.Getenv("PROMETHEUS_PORT")
	if promPort == "" {
		promPort = defaultPrometheusPort
	}

	return &Config{
		DbUrl: dbURL,
		KafkaConfig: KafkaConfig{
			BrokerList: brokerList,
			Topic:      topic,
		},
		OutputMode: outputMode,
		CacheTTL:   cacheTTL,
		ServerConfig: ServerConfig{
			GrpcPort: grpcPort,
		},
		PrometheusConfig: PrometheusConfig{
			PrometheusPort: promPort,
		},
	}, nil

}
