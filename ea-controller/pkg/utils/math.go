package utils

import (
	"math"
	"sort"

	"ea-controller/internal/model"
)

// Stabilize avoids oscillations in LP output
func Stabilize(prev, curr []float64, eps float64) ([]float64, bool) {
	if prev == nil {
		return curr, false
	}

	diff := 0.0
	for i := range curr {
		diff += math.Abs(curr[i] - prev[i])
	}

	if diff < eps {
		return prev, true // STABLE → skip
	}
	return curr, false
}

type item struct {
	idx     int
	cij     float64
	utility float64
}

func floorWeights(x []float64) ([]int, int) {
	n := len(x)
	w := make([]int, n)
	sum := 0

	for i := range x {
		val := 100 * x[i]
		w[i] = int(math.Floor(val))
		sum += w[i]
	}
	return w, 100 - sum
}

func distribute(w []int, variants []model.Variant, delta int) []int {
	n := len(w)

	items := make([]item, n)
	for i := range items {
		items[i] = item{
			idx:     i,
			cij:     variants[i].Cij,
			utility: variants[i].Utility,
		}
	}

	sort.Slice(items, func(i, j int) bool {
		if math.Abs(items[i].cij-items[j].cij) < 1e-9 {
			return items[i].utility > items[j].utility
		}
		return items[i].cij < items[j].cij
	})

	for k := 0; k < delta; k++ {
		w[items[k%n].idx]++
	}

	return w
}

// Quantize converts probabilities into integer weights (sum=100)
func Quantize(s *model.Service, x []float64) []int {
	w, delta := floorWeights(x)

	if delta <= 0 {
		return w
	}

	return distribute(w, s.Variants, delta)
}
