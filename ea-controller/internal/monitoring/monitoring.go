package monitoring

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"sync"

	"ea-controller/internal/config"
	"ea-controller/internal/constants"
	"ea-controller/internal/model"
	"ea-controller/internal/prometheus"
)

func Run(s *model.Service, cfg config.Config) {
	client := prometheus.New(cfg.PrometheusURL)

	var wg sync.WaitGroup

	var throughput float64
	var tMap map[string]float64
	var dynMap map[string]float64
	var idleMap map[string]float64

	runAsync(&wg, func() {
		throughput = getServiceThroughput(client, s, cfg)
	})

	runAsync(&wg, func() {
		tMap = getVariantThroughput(client, s, cfg)
	})

	runAsync(&wg, func() {
		dynMap = getVariantPower(client, cfg, constants.ModeDynamic)
	})

	runAsync(&wg, func() {
		idleMap = getVariantPower(client, cfg, constants.ModeIdle)
	})

	wg.Wait()

	s.R = throughput
	log.Printf("[SERVICE %s] R=%.4f\n", s.Name, s.R)

	for i := range s.Variants {
		v := &s.Variants[i]

		v.PowerDynamic = sumPowerByWorkload(dynMap, v.Workload)
		v.PowerIdle = sumPowerByWorkload(idleMap, v.Workload)

		v.Throughput = tMap[v.Workload]

		if v.Throughput > 0 {
			v.Cij = v.PowerDynamic / v.Throughput
			v.PrevCij = v.Cij
		} else if v.PrevCij > 0 {
			v.Cij = v.PrevCij
		} else {
			v.Cij = v.DefaultCj
		}

		log.Printf("[SERVICE %s] VARIANT %s -> R=%.3f Dyn=%.3f Idle=%.3f C=%.6f\n",
			s.Name, v.Name, v.Throughput, v.PowerDynamic, v.PowerIdle, v.Cij)
	}
}

func runAsync(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fn()
	}()
}

func getServiceThroughput(c *prometheus.Client, s *model.Service, cfg config.Config) float64 {
	query := fmt.Sprintf(`
	sum(rate(%s{
		%s="%s",
		%s="%s",
		%s="%s"
	}[%s]))
	`, 	constants.MetricIstioRequests,
        	constants.LabelServiceName, s.Name,
        	constants.LabelServiceNamespace, cfg.Namespace,
        	constants.LabelReporter, constants.ReporterDestination,
        	cfg.RateInterval,
	)

	return c.Scalar(query)
}

func getVariantThroughput(c *prometheus.Client, s *model.Service, cfg config.Config) map[string]float64 {
	query := fmt.Sprintf(`
	sum by (%s)(
		rate(%s{
			%s="%s",
			%s="%s",
			%s="%s"
		}[%s])
	)
	`, 	constants.LabelWorkload,
        	constants.MetricIstioRequests,
        	constants.LabelServiceName, s.Name,
        	constants.LabelServiceNamespace, cfg.Namespace,
        	constants.LabelReporter, constants.ReporterDestination,
        	cfg.RateInterval,
	)

	res, _ := c.Query(query)

	out := map[string]float64{}
	for _, r := range res {
		workload := r.Metric[constants.LabelWorkload]
		val := toFloat(r.Value[1])
		out[workload] = val
	}
	return out
}

func getVariantPower(c *prometheus.Client, cfg config.Config, mode string) map[string]float64 {
	query := fmt.Sprintf(`
	sum by (%s)(
		rate(%s{
			%s="%s",
			%s="%s",
			%s!~"%s|%s"
		}[%s])
	)
	`,	constants.LabelPod,
        	constants.MetricKeplerEnergy,
        	constants.LabelContainerNS, cfg.Namespace,
        	constants.LabelMode, mode,
                constants.LabelContainerName,
                constants.ContainerIstioProxy, constants.ContainerIstioInit,
        	cfg.RateInterval,
	)

	res, _ := c.Query(query)

	out := map[string]float64{}
	for _, r := range res {
		pod := r.Metric[constants.LabelPod]
		val := toFloat(r.Value[1])
		out[pod] = val
	}
	return out
}

func sumPowerByWorkload(powerMap map[string]float64, workload string) float64 {
	sum := 0.0
	for pod, val := range powerMap {
		if workloadBelongsToPod(workload, pod) {
			sum += val
		}
	}
	return sum
}

func workloadBelongsToPod(workload, pod string) bool {
        return strings.HasPrefix(pod, workload)
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case int:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}
