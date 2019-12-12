# apm
--
    import "."


## Usage

#### type APM

```go
type APM struct {
}
```


#### func  New

```go
func New(config newrelic.Config) APM
```

#### func (*APM) ListApplications

```go
func (apm *APM) ListApplications(params *ListApplicationsParams) ([]Application, error)
```
ListApplications is used to retrieve New Relic applications.

#### type Application

```go
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
```

Application represents information about a New Relic application.

#### type ApplicationEndUserSummary

```go
type ApplicationEndUserSummary struct {
	ResponseTime float64 `json:"response_time"`
	Throughput   float64 `json:"throughput"`
	ApdexTarget  float64 `json:"apdex_target"`
	ApdexScore   float64 `json:"apdex_score"`
}
```

ApplicationEndUserSummary represents performance information about a New Relic
application.

#### type ApplicationLinks

```go
type ApplicationLinks struct {
	ServerIDs     []int `json:"servers,omitempty"`
	HostIDs       []int `json:"application_hosts,omitempty"`
	InstanceIDs   []int `json:"application_instances,omitempty"`
	AlertPolicyID int   `json:"alert_policy"`
}
```

ApplicationLinks represents all the links for a New Relic application.

#### type ApplicationSettings

```go
type ApplicationSettings struct {
	AppApdexThreshold        float64 `json:"app_apdex_threshold,omitempty"`
	EndUserApdexThreshold    float64 `json:"end_user_apdex_threshold,omitempty"`
	EnableRealUserMonitoring bool    `json:"enable_real_user_monitoring"`
	UseServerSideConfig      bool    `json:"use_server_side_config"`
}
```

ApplicationSettings represents some of the settings of a New Relic application.

#### type ApplicationSummary

```go
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
```

ApplicationSummary represents performance information about a New Relic
application.

#### type ListApplicationsParams

```go
type ListApplicationsParams struct {
	Name     *string
	Host     *string
	IDs      []int
	Language *string
}
```

ListApplicationsParams represents a set of filters to be used when querying New
Relic applications.
