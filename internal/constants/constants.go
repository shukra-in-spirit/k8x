package constants

const (
	DefaultPrometheusTimeRange = 3 * 60 * 60
	// PrometheusMode determines whether a 'local' csv file is used or data is read from 'prometheus'
	PrometheusMode = "prometheus"

	// TrainingDataDuration is the number of days the data is collected from prometheus
	TrainingDataDuration = -14

	// KubeClientMode is the flag to decide if the service is running in 'local' or in 'cluster'
	KubeClientMode = "local"

	Local = "local"

	// PrometheusRequestTimeOut is the timeout period for the prometheus call
	PrometheusRequestTimeOut = 5

	// StepsMinutesInterval is the steps for the prometheus query in minutes
	StepsMinutesInterval = 2
)
