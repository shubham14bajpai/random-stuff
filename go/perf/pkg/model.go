package perf

type Config struct {
	Sample   []string
	Interval []string
	Instance []string
	Metrics
}

type Metrics struct {
	CPU    []Metric
	Memory []Metric
}

type Metric struct {
	Name  string
	Query string
}
