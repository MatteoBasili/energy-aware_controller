package optimizer

import (
	"log"

	"ea-controller/internal/model"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
)

// fallbackMinEnergy routes all traffic to the most energy-efficient variant
func fallbackMinEnergy(s *model.Service) []float64 {
	n := len(s.Variants)
	x := make([]float64, n)

	if n == 0 {
		return x
	}

	minIdx := 0
	minC := s.Variants[0].Cij

	for i := 1; i < n; i++ {
		if s.Variants[i].Cij < minC {
			minC = s.Variants[i].Cij
			minIdx = i
		}
	}

	log.Printf("[SERVICE %s] FALLBACK → routing 100%% to most energy-efficient variant (%s, Cij=%.6f)",
		s.Name,
		s.Variants[minIdx].Name,
		minC,
	)

	x[minIdx] = 1.0
	return x
}

func buildConstraints(n int, variants []model.Variant, R float64, availableBudget float64) (*mat.Dense, []float64) {
	G := mat.NewDense(n+1, n, nil)
	h := make([]float64, n+1)

	for i := 0; i < n; i++ {
		G.Set(0, i, variants[i].Cij*R)
	}
	h[0] = availableBudget

	for i := 0; i < n; i++ {
		G.Set(i+1, i, 1)
		h[i+1] = 1
	}

	return G, h
}

// Optimize solves the LP problem for a service
func Optimize(s *model.Service) ([]float64, bool) {
	n := len(s.Variants)

	totalIdle := 0.0
	for _, v := range s.Variants {
		totalIdle += v.PowerIdle
	}

	availableBudget := s.Budget - totalIdle
	if availableBudget <= 0 {
		log.Printf("[SERVICE %s] ENERGY VIOLATION: idle consumption (%.3f) exceeds budget (%.3f)",
			s.Name, totalIdle, s.Budget)

		return fallbackMinEnergy(s), true
	}

	c := make([]float64, n)
	for i := range c {
		c[i] = -s.Variants[i].Utility
	}

	G, h := buildConstraints(n, s.Variants, s.R, availableBudget)

	A := mat.NewDense(1, n, nil)
	for i := 0; i < n; i++ {
		A.Set(0, i, 1)
	}
	b := []float64{1}

	cNew, ANew, bNew := lp.Convert(c, G, h, A, b)
	_, x, err := lp.Simplex(cNew, ANew, bNew, 1e-6, nil)

	if err != nil {
		log.Printf("[SERVICE %s] ENERGY CONSTRAINT NOT SATISFIABLE (LP infeasible): %v", s.Name, err)

		for _, v := range s.Variants {
			log.Printf("  - Variant %s -> Cij=%.6f, R=%.3f, estimated consumption=%.3f",
				v.Name, v.Cij, s.R, v.Cij*s.R)
		}

		log.Printf("  Available budget: %.3f", availableBudget)

		return fallbackMinEnergy(s), true
	}

	return x[:n], false
}
