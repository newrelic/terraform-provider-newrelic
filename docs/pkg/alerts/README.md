# alerts
--
    import "github.com/newrelic/newrelic-client-go/pkg/alerts"


## Usage

#### func  TestIntegrationListAlertPolicies

```go
func TestIntegrationListAlertPolicies(t *testing.T)
```

#### type AlertPolicy

```go
type AlertPolicy struct {
	ID                 int    `json:"id,omitempty"`
	IncidentPreference string `json:"incident_preference,omitempty"`
	Name               string `json:"name,omitempty"`
	CreatedAt          int64  `json:"created_at,omitempty"`
	UpdatedAt          int64  `json:"updated_at,omitempty"`
}
```

AlertPolicy represents a New Relic alert policy.

#### type Alerts

```go
type Alerts struct {
}
```

Alerts is used to communicate with New Relic Alerts.

#### func  New

```go
func New(config config.Config) Alerts
```
New is used to create a new Alerts client instance.

#### func (*Alerts) CreateAlertPolicy

```go
func (alerts *Alerts) CreateAlertPolicy(policy AlertPolicy) (*AlertPolicy, error)
```
CreateAlertPolicy creates a new alert policy for a given account.

#### func (*Alerts) GetAlertPolicy

```go
func (alerts *Alerts) GetAlertPolicy(id int) (*AlertPolicy, error)
```
GetAlertPolicy returns a specific alert policy by ID for a given account.

#### func (*Alerts) ListAlertPolicies

```go
func (alerts *Alerts) ListAlertPolicies(params *ListAlertPoliciesParams) ([]AlertPolicy, error)
```
ListAlertPolicies returns a list of Alert Policies for a given account.

#### func (*Alerts) UpdateAlertPolicy

```go
func (alerts *Alerts) UpdateAlertPolicy(policy AlertPolicy) (*AlertPolicy, error)
```
UpdateAlertPolicy update an alert policy for a given account.

#### type ListAlertPoliciesParams

```go
type ListAlertPoliciesParams struct {
	Name *string
}
```

ListAlertPoliciesParams represents a set of filters to be used when querying New
Relic alert policies.
