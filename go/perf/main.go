package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	perf "github.com/shubham14bajpai/perf/pkg"
)

const (
	cpuPercentile99Query = `quantile_over_time(0.99, ((sum(rate(node_cpu_seconds_total{instance="{{.Instance}}",job="node"}[{{.Interval}}])) - sum(rate(node_cpu_seconds_total{instance="{{.Instance}}",job="node",mode="idle"}[{{.Interval}}]))) * 100 )[{{.Range}}:])`
	cpuPercentile95Query = `quantile_over_time(0.95, ((sum(rate(node_cpu_seconds_total{instance="{{.Instance}}",job="node"}[{{.Interval}}])) - sum(rate(node_cpu_seconds_total{instance="{{.Instance}}",job="node",mode="idle"}[{{.Interval}}]))) * 100 )[{{.Range}}:])`
	cpuAvergeQuery       = `avg_over_time(((sum(rate(node_cpu_seconds_total{instance="{{.Instance}}",job="node"}[{{.Interval}}])) - sum(rate(node_cpu_seconds_total{instance="{{.Instance}}",job="node",mode="idle"}[{{.Interval}}]))) * 100 )[{{.Range}}:])`
	memPercentile99Query = `quantile_over_time(0.99, (node_memory_MemTotal_bytes{instance="{{.Instance}}",job="node"} - node_memory_MemFree_bytes{instance="{{.Instance}}",job="node"} - (node_memory_Cached_bytes{instance="{{.Instance}}",job="node"} + node_memory_Buffers_bytes{instance="{{.Instance}}",job="node"} + node_memory_SReclaimable_bytes{instance="{{.Instance}}",job="node"})) [{{.Range}}:])`
	memPercentile95Query = `quantile_over_time(0.95, (node_memory_MemTotal_bytes{instance="{{.Instance}}",job="node"} - node_memory_MemFree_bytes{instance="{{.Instance}}",job="node"} - (node_memory_Cached_bytes{instance="{{.Instance}}",job="node"} + node_memory_Buffers_bytes{instance="{{.Instance}}",job="node"} + node_memory_SReclaimable_bytes{instance="{{.Instance}}",job="node"})) [{{.Range}}:])`
	memAverageQuery      = `avg_over_time((node_memory_MemTotal_bytes{instance="{{.Instance}}",job="node"} - node_memory_MemFree_bytes{instance="{{.Instance}}",job="node"} - (node_memory_Cached_bytes{instance="{{.Instance}}",job="node"} + node_memory_Buffers_bytes{instance="{{.Instance}}",job="node"} + node_memory_SReclaimable_bytes{instance="{{.Instance}}",job="node"})) [{{.Range}}:])`
)

func main() {
	client, err := api.NewClient(api.Config{
		Address: "http://10.0.2.15:9090",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := perf.Config{
		Interval: []string{"1h"},
		Sample:   []string{"60s"},
		Instance: []string{"bench"},
		Metrics: perf.Metrics{
			CPU: []perf.Metric{
				{Name: "99th Percen", Query: cpuPercentile99Query},
				{Name: "95th Percen", Query: cpuPercentile95Query},
				{Name: "Average", Query: cpuAvergeQuery},
			},
			Memory: []perf.Metric{
				{Name: "99th Percen", Query: memPercentile99Query},
				{Name: "95th Percen", Query: memPercentile95Query},
				{Name: "Average", Query: memAverageQuery},
			},
		},
	}

	perf.Stats(ctx, v1api, cfg)
}
