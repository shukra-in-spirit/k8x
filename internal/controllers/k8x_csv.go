package controllers

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/shukra-in-spirit/k8x/internal/models"
)

// based on https://prometheus.io/docs/prometheus/latest/querying/api/

type CSVFunctions interface {
	GetCSVData(ctx context.Context) (*models.PrometheusDataSetResponse, error)
}

type CSVFile struct {
	filePath string
}

func LoadCSVFile(path string) *CSVFile {
	return &CSVFile{
		filePath: path,
	}
}

func (csvFile *CSVFile) GetCSVData(ctx context.Context) (*models.PrometheusDataSetResponse, error) {
	file, err := os.Open(csvFile.filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading CSV: %v", err)
	}

	var data []models.PrometheusDataSetResponseItem
	for _, record := range records {
		// Parse timestamp, CPU, and memory values from the CSV record
		timestamp, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			return nil, fmt.Errorf("Error parsing timestamp: %v", err)
		}

		cpu, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing CPU value: %v", err)
		}

		memory, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing memory value: %v", err)
		}

		// Create a Data struct and append it to the slice
		dataPoint := models.PrometheusDataSetResponseItem{
			Timestamp: timestamp,
			CPU:       float32(cpu),
			Memory:    float32(memory),
		}
		data = append(data, dataPoint)
	}

	dataFrame := &models.PrometheusDataSetResponse{
		PromItemList: data,
	}
	return dataFrame, nil
}