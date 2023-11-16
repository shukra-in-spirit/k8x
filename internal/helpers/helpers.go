package helpers

import "strings"

func DecomposeServiceID(service_id string) (string, string) {
	// Split the input string based on the '-' delimiter
	parts := strings.Split(service_id, "-")

	// Check if there are at least two parts (serviceName and namespace)
	if len(parts) < 2 {
		return "", ""
	}

	// Extract serviceName and namespace
	serviceName := parts[0]
	namespace := parts[1]

	return serviceName, namespace
}
