package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func main() {
	// Prometheus server URL
	prometheusURL := "http://localhost:9090" // Change this to your Prometheus server URL

	// Create a new Prometheus API client
	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating client:", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// PromQL query
	query := `rate(eagle_pod_container_resource_usage_memory_bytes{exported_container="prometheus",exported_namespace="default"}[2m])`

	// Query Prometheus
	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error querying Prometheus:", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Fprintln(os.Stderr, "Warnings:", warnings)
	}

	// Print the result
	fmt.Printf("Result:\n%s\n", result.String())

	// To handle the result as a vector (list of samples), you can cast it like this:
	vectorVal, ok := result.(model.Vector)
	if !ok {
		fmt.Fprintln(os.Stderr, "Query result is not a vector")
		os.Exit(1)
	}

	for _, sample := range vectorVal {
		fmt.Printf("Metric: %s, Value: %f, Time: %v\n", sample.Metric, sample.Value, sample.Timestamp.Time())
	}
}
