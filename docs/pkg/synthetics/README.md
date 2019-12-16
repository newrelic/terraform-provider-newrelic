# synthetics
--
    import "github.com/newrelic/newrelic-client-go/pkg/synthetics"


## Usage

#### type Monitor

```go
type Monitor struct {
	ID           string         `json:"id,omitempty"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Frequency    uint           `json:"frequency"`
	URI          string         `json:"uri"`
	Locations    []string       `json:"locations"`
	Status       string         `json:"status"`
	SLAThreshold float64        `json:"slaThreshold"`
	UserID       uint           `json:"userId,omitempty"`
	APIVersion   string         `json:"apiVersion,omitempty"`
	ModifiedAt   time.Time      `json:"modified_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	Options      MonitorOptions `json:"options,omitempty"`
}
```

Monitor represents a New Relic Synthetics monitor.

#### type MonitorOptions

```go
type MonitorOptions struct {
	ValidationString       string `json:"validationString,omitempty"`
	VerifySSL              bool   `json:"verifySSL,omitempty"`
	BypassHEADRequest      bool   `json:"bypassHEADRequest,omitempty"`
	TreatRedirectAsFailure bool   `json:"treatRedirectAsFailure,omitempty"`
}
```

MonitorOptions represents the options for a New Relic Synthetics monitor.

#### type Synthetics

```go
type Synthetics struct {
}
```

Synthetics is used to communicate with the New Relic Synthetics product.

#### func  New

```go
func New(config config.Config) Synthetics
```
New is used to create a new Synthetics client instance.

#### func (*Synthetics) ListMonitors

```go
func (s *Synthetics) ListMonitors() ([]Monitor, error)
```
ListMonitors is used to retrieve New Relic Synthetics monitors.
