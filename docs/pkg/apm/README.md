# apm
--
    import "."


## Usage

#### type APM

```go
type APM struct {
}
```

APM is used to communicate with the New Relic APM product.

#### func  New

```go
func New(config config.Config) APM
```
New is used to create a new APM client instance.

#### func (*APM) DeleteApplication

```go
func (apm *APM) DeleteApplication(applicationID int) (*Application, error)
```
DeleteApplication is used to delete a New Relic application. This process will
only succeed if the application is no longer reporting data.

#### func (*APM) GetApplication

```go
func (apm *APM) GetApplication(applicationID int) (*Application, error)
```
GetApplication is used to retrieve a single New Relic application.

#### func (*APM) ListApplications

```go
func (apm *APM) ListApplications(params *ListApplicationsParams) ([]Application, error)
```
ListApplications is used to retrieve New Relic applications.

#### func (*APM) UpdateApplication

```go
func (apm *APM) UpdateApplication(applicationID int, params UpdateApplicationParams) (*Application, error)
```
UpdateApplication is used to update a New Relic application's name and/or
settings.

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

#### type UpdateApplicationParams

```go
type UpdateApplicationParams struct {
	Name     string
	Settings ApplicationSettings
}
```

UpdateApplicationParams represents a set of parameters to be used when updating
New Relic applications.
