package controllers

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

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

func (csvFile *CSVFile) GetCSVDataCPUandMem(ctx context.Context) (*models.PrometheusDataSetResponse, *models.PrometheusDataSetResponse, error) {
	file, err := os.Open(csvFile.filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("Error reading CSV: %v", err)
	}

	var cpuData []models.PrometheusDataSetResponseItem
	var memData []models.PrometheusDataSetResponseItem
	for _, record := range records {

		cpu, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			return nil, nil, fmt.Errorf("Error parsing CPU value: %v", err)
		}

		memory, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			return nil, nil, fmt.Errorf("Error parsing memory value: %v", err)
		}

		// Create a Data struct and append it to the slice
		cpuDataPoint := models.PrometheusDataSetResponseItem{
			Metric: float32(cpu),
		}
		cpuData = append(cpuData, cpuDataPoint)

		memDataPoint := models.PrometheusDataSetResponseItem{
			Metric: float32(memory),
		}
		memData = append(memData, memDataPoint)
	}

	cpuDataFrame := &models.PrometheusDataSetResponse{
		PromItemList: cpuData,
	}
	memDataFrame := &models.PrometheusDataSetResponse{
		PromItemList: memData,
	}
	return cpuDataFrame, memDataFrame, nil
}

func GetCSVData(ctx context.Context, filePath string) (*models.PrometheusDataSetResponse, error) {
	file, err := os.Open(filePath)
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

		metric, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing memory value: %v", err)
		}

		// Create a Data struct and append it to the slice
		dataPoint := models.PrometheusDataSetResponseItem{
			Metric: float32(metric),
		}
		data = append(data, dataPoint)

	}

	dataFrame := &models.PrometheusDataSetResponse{
		PromItemList: data,
	}
	return dataFrame, nil
}
