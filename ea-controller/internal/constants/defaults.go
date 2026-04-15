package constants

const (
        // --- CONFIG DEFAULTS ---

        DefaultPrometheusURL = "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
        DefaultNamespace     = "approx"

        DefaultIntervalSec  = 15
        DefaultEpsilon      = 0.05
        DefaultRateInterval = "1m"
        DefaultMaxWorkers   = 10
)
