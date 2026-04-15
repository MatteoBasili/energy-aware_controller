package controller

import (
	"log"
	"sync"
	"time"

	"ea-controller/internal/config"
	"ea-controller/internal/istio"
	"ea-controller/internal/monitoring"
	"ea-controller/internal/optimizer"
	"ea-controller/internal/store"
	"ea-controller/internal/model"
	"ea-controller/pkg/utils"

	istioclient "istio.io/client-go/pkg/clientset/versioned"
)

func Run(cfg config.Config, st *store.Store, client *istioclient.Clientset) {
	sem := make(chan struct{}, cfg.MaxWorkers)

	for {
		var wg sync.WaitGroup

		services := st.List()

		for _, s := range services {
			wg.Add(1)
			sem <- struct{}{}

			go func(s *model.Service) {
				defer wg.Done()
				defer func() { <-sem }()

				monitoring.Run(s, cfg)

				x, fallback := optimizer.Optimize(s)

				var skip bool

				if !fallback {
					x, skip = utils.Stabilize(s.PrevX, x, cfg.Epsilon)
				}

				if skip {
        				log.Printf(
						"[SERVICE %s] Stabilized (Δ < %.4f) → skipping update",
						s.Name,
						cfg.Epsilon,
					)
        				return
				}

				w := utils.Quantize(s, x)

				istio.Apply(s, w, client)

				s.PrevX = x
			}(s)
		}

		wg.Wait()
		time.Sleep(cfg.Interval)
	}
}
