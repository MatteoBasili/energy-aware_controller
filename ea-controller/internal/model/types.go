package model

type Variant struct {
	Name       string
	Workload   string
	Utility    float64
	Throughput float64

	PowerDynamic float64
	PowerIdle    float64

	Cij float64

	PrevCij   float64
	DefaultCj float64
}

type Service struct {
	Name      string
	Namespace string
	Variants  []Variant
	Budget    float64
	R         float64
	PrevX     []float64
}
