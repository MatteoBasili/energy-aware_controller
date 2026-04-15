package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"ea-controller/internal/constants"
)

type Config struct {
	PrometheusURL string
	Namespace     string
	Interval      time.Duration
	Epsilon       float64
	RateInterval  string
        MaxWorkers    int
}

func Load() Config {
	cfg := Config{}

	cfg.PrometheusURL = getEnv("PROMETHEUS_URL", constants.DefaultPrometheusURL)
	cfg.Namespace = getEnv("NAMESPACE", constants.DefaultNamespace)

	interval := getEnvInt("CONTROL_LOOP_INTERVAL", constants.DefaultIntervalSec)
	cfg.Interval = time.Duration(interval) * time.Second

	cfg.Epsilon = getEnvFloat("EPSILON", constants.DefaultEpsilon)

	cfg.RateInterval = getEnv("PROM_RATE_INTERVAL", constants.DefaultRateInterval)

        cfg.MaxWorkers = getEnvInt("MAX_WORKERS", constants.DefaultMaxWorkers)

	log.Printf(
		"CONFIG -> Prometheus: %s | Namespace: %s | Interval: %v | Epsilon: %f | RateInterval: %s | MaxWorkers: %d\n",
		cfg.PrometheusURL,
		cfg.Namespace,
		cfg.Interval,
		cfg.Epsilon,
		cfg.RateInterval,
                cfg.MaxWorkers,
	)

	return cfg
}

func getEnv(key, def string) string {
        if val := os.Getenv(key); val != "" {
                return val
        }
        return def
}

func getEnvInt(key string, def int) int {
        val := os.Getenv(key)
        if val == "" {
                return def
        }
        v, err := strconv.Atoi(val)
        if err != nil {
                return def
        }
        return v
}

func getEnvFloat(key string, def float64) float64 {
        val := os.Getenv(key)
        if val == "" {
                return def
        }
        v, err := strconv.ParseFloat(val, 64)
        if err != nil {
                return def
        }
        return v
}
