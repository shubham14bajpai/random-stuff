package perf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type Input struct {
	Instance string
	Range    string
	Interval string
}

func Stats(ctx context.Context, v1api v1.API, cfg Config) {
	for _, interval := range cfg.Interval {
		for _, sample := range cfg.Sample {
			for _, instance := range cfg.Instance {
				input := Input{
					Range:    interval,
					Interval: sample,
					Instance: instance,
				}
				fmt.Printf("\n%+v\n\n", input)
				header := ""
				values := ""
				for _, query := range cfg.CPU {
					var buf bytes.Buffer
					t := template.Must(template.New("query").Parse(query.Query))
					err := t.Execute(&buf, input)
					if err != nil {
						panic(err)
					}
					q := buf.String()
					buf.Reset()
					// fmt.Println(q)
					result, warnings, err := v1api.Query(ctx, q, time.Now())
					if err != nil {
						fmt.Printf("Error querying Prometheus: %v\n", err)
						os.Exit(1)
					}
					if len(warnings) > 0 {
						fmt.Printf("Warnings: %v\n", warnings)
					}
					header += query.Name + "\t"
					values += fmt.Sprintf("%.2f\t\t", format(result.String()))
				}
				fmt.Println("CPU")
				fmt.Println("\t", header, "\n\t", values)

				header = ""
				values = ""
				for _, query := range cfg.Memory {
					var buf bytes.Buffer
					t := template.Must(template.New("query").Parse(query.Query))
					err := t.Execute(&buf, input)
					if err != nil {
						panic(err)
					}
					q := buf.String()
					buf.Reset()
					// fmt.Println(q)
					result, warnings, err := v1api.Query(ctx, q, time.Now())
					if err != nil {
						fmt.Printf("Error querying Prometheus: %v\n", err)
						os.Exit(1)
					}
					if len(warnings) > 0 {
						fmt.Printf("Warnings: %v\n", warnings)
					}
					header += query.Name + "\t"
					values += fmt.Sprintf("%.0f\t\t", (format(result.String()) / 1000000))
				}
				fmt.Println("Memory")
				fmt.Println("\t", header, "\n\t", values)
			}
			fmt.Print("-------------------------------------------\n\n")
		}
	}
}

func format(s string) float64 {
	// fmt.Println("print:", s)
	re := regexp.MustCompile(`> (.*) @`)
	d, err := strconv.ParseFloat(strings.TrimLeft(strings.TrimRight(re.FindString(s), " @"), "> "), 64)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)

	}
	return math.Round(d*100) / 100
}
