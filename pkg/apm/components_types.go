package apm

import "time"

// Component represents information about a New Relic Plugins component.
type Component struct {
	ID             int             `json:"id"`
	Name           string          `json:"name,omitempty"`
	HealthStatus   string          `json:"health_status,omitempty"`
	SummaryMetrics []SummaryMetric `json:"summary_metrics"`
}

// ComponentMetric represents metric information for a specific component.
type ComponentMetric struct {
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values"`
}

// SummaryMetric represents summary information for a specific metric.
type SummaryMetric struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Metric        string          `json:"metric"`
	ValueFunction string          `json:"value_function"`
	Thresholds    MetricThreshold `json:"thresholds"`
	Values        MetricValue     `json:"values"`
}

// MetricThreshold represents a threshold value for a metric.
type MetricThreshold struct {
	Caution  float64 `json:"caution"`
	Critical float64 `json:"critical"`
}

// MetricValue represents the observed value of a metric.
type MetricValue struct {
	Raw       float64 `json:"raw"`
	Formatted string  `json:"formatted"`
}

// Metric represents data for a specific metric.
type Metric struct {
	Name       string            `json:"name"`
	Timeslices []MetricTimeslice `json:"timeslices"`
}

// MetricTimeslice represents the values of a metric over a given time.
type MetricTimeslice struct {
	From   *time.Time         `json:"from,omitempty"`
	To     *time.Time         `json:"to,omitempty"`
	Values map[string]float64 `json:"values,omitempty"`
}
