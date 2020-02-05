package apm

import "time"

// Application represents information about a New Relic application.
type Application struct {
	ID             int                       `json:"id,omitempty"`
	Name           string                    `json:"name,omitempty"`
	Language       string                    `json:"language,omitempty"`
	HealthStatus   string                    `json:"health_status,omitempty"`
	Reporting      bool                      `json:"reporting"`
	LastReportedAt string                    `json:"last_reported_at,omitempty"`
	Summary        ApplicationSummary        `json:"application_summary,omitempty"`
	EndUserSummary ApplicationEndUserSummary `json:"end_user_summary,omitempty"`
	Settings       ApplicationSettings       `json:"settings,omitempty"`
	Links          ApplicationLinks          `json:"links,omitempty"`
}

// ApplicationSummary represents performance information about a New Relic application.
type ApplicationSummary struct {
	ResponseTime            float64 `json:"response_time"`
	Throughput              float64 `json:"throughput"`
	ErrorRate               float64 `json:"error_rate"`
	ApdexTarget             float64 `json:"apdex_target"`
	ApdexScore              float64 `json:"apdex_score"`
	HostCount               int     `json:"host_count"`
	InstanceCount           int     `json:"instance_count"`
	ConcurrentInstanceCount int     `json:"concurrent_instance_count"`
}

// ApplicationEndUserSummary represents performance information about a New Relic application.
type ApplicationEndUserSummary struct {
	ResponseTime float64 `json:"response_time"`
	Throughput   float64 `json:"throughput"`
	ApdexTarget  float64 `json:"apdex_target"`
	ApdexScore   float64 `json:"apdex_score"`
}

// ApplicationSettings represents some of the settings of a New Relic application.
type ApplicationSettings struct {
	AppApdexThreshold        float64 `json:"app_apdex_threshold,omitempty"`
	EndUserApdexThreshold    float64 `json:"end_user_apdex_threshold,omitempty"`
	EnableRealUserMonitoring bool    `json:"enable_real_user_monitoring"`
	UseServerSideConfig      bool    `json:"use_server_side_config"`
}

// ApplicationLinks represents all the links for a New Relic application.
type ApplicationLinks struct {
	ServerIDs     []int `json:"servers,omitempty"`
	HostIDs       []int `json:"application_hosts,omitempty"`
	InstanceIDs   []int `json:"application_instances,omitempty"`
	AlertPolicyID int   `json:"alert_policy"`
}

// Deployment represents information about a New Relic application deployment.
type Deployment struct {
	Links       *DeploymentLinks `json:"links,omitempty"`
	ID          int              `json:"id,omitempty"`
	Revision    string           `json:"revision"`
	Changelog   string           `json:"changelog,omitempty"`
	Description string           `json:"description,omitempty"`
	User        string           `json:"user,omitempty"`
	Timestamp   string           `json:"timestamp,omitempty"`
}

// DeploymentLinks contain the application ID for the deployment.
type DeploymentLinks struct {
	ApplicationID int `json:"application,omitempty"`
}

// KeyTransaction represents information about a New Relic key transaction.
type KeyTransaction struct {
	ID              int                       `json:"id,omitempty"`
	Name            string                    `json:"name,omitempty"`
	TransactionName string                    `json:"transaction_name,omitempty"`
	HealthStatus    string                    `json:"health_status,omitempty"`
	LastReportedAt  string                    `json:"last_reported_at,omitempty"`
	Reporting       bool                      `json:"reporting"`
	Summary         ApplicationSummary        `json:"application_summary,omitempty"`
	EndUserSummary  ApplicationEndUserSummary `json:"end_user_summary,omitempty"`
	Links           KeyTransactionLinks       `json:"links,omitempty"`
}

// KeyTransactionLinks represents associations for a key transaction.
type KeyTransactionLinks struct {
	Application int `json:"application,omitempty"`
}

// Label represents a New Relic label.
type Label struct {
	Key      string     `json:"key,omitempty"`
	Category string     `json:"category,omitempty"`
	Name     string     `json:"name,omitempty"`
	Links    LabelLinks `json:"links,omitempty"`
}

// LabelLinks represents external references on the Label.
type LabelLinks struct {
	Applications []int `json:"applications"`
	Servers      []int `json:"servers"`
}

// ListApplicationsParams represents a set of filters to be
// used when querying New Relic applications.
type ListApplicationsParams struct {
	Name     string `url:"filter[name],omitempty"`
	Host     string `url:"filter[host],omitempty"`
	IDs      []int  `url:"filter[ids],omitempty,comma"`
	Language string `url:"filter[language],omitempty"`
}

// UpdateApplicationParams represents a set of parameters to be
// used when updating New Relic applications.
type UpdateApplicationParams struct {
	Name     string
	Settings ApplicationSettings
}

// MetricNamesParams are the request parameters for the /metrics.json endpoint.
type MetricNamesParams struct {
	Name string `url:"name,omitempty"`
}

// MetricDataParams are the request parameters for the /metrics/data.json endpoint.
type MetricDataParams struct {
	Names     []string   `url:"names,omitempty"`
	Values    []string   `url:"values,omitempty"`
	From      *time.Time `url:"from,omitempty"`
	To        *time.Time `url:"to,omitempty"`
	Period    int        `url:"period,omitempty"`
	Summarize bool       `url:"summarize,omitempty"`
	Raw       bool       `url:"raw,omitempty"`
}

// MetricName is the name of a metric, and the names of the values that can be retrieved.
type MetricName struct {
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}

// MetricData is the series of time windows and the data therein, for a given metric name.
type MetricData struct {
	Name       string            `json:"name,omitempty"`
	Timeslices []MetricTimeslice `json:"timeslices,omitempty"`
}

// MetricTimeslice is a single window of time for a given metric, with the associated metric data.
type MetricTimeslice struct {
	From   *time.Time            `json:"from"`
	To     *time.Time            `json:"to"`
	Values MetricTimesliceValues `json:"values"`
}

//MetricTimesliceValues is the collection of metric values for a single time slice.
type MetricTimesliceValues struct {
	AsPercentage           float64 `json:"as_percentage"`
	AverageTime            float64 `json:"average_time"`
	CallsPerMinute         float64 `json:"calls_per_minute"`
	MaxValue               float64 `json:"max_value"`
	TotalCallTimePerMinute float64 `json:"total_call_time_per_minute"`
	Utilization            float64 `json:"utilization"`
}
