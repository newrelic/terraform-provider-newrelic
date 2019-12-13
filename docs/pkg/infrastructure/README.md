# infrastructure
--
    import "."


## Usage

<<<<<<< HEAD
#### type ErrorDetail

```go
=======
#### type AlertCondition

```go
type AlertCondition struct {
	Comparison          string              `json:"comparison,omitempty"`
	CreatedAt           serialization.Epoch `json:"created_at_epoch_millis,omitempty"`
	Critical            *Threshold          `json:"critical_threshold,omitempty"`
	Enabled             bool                `json:"enabled"`
	Event               string              `json:"event_type,omitempty"`
	ID                  int                 `json:"id,omitempty"`
	IntegrationProvider string              `json:"integration_provider,omitempty"`
	Name                string              `json:"name,omitempty"`
	PolicyID            int                 `json:"policy_id,omitempty"`
	ProcessWhere        string              `json:"process_where_clause,omitempty"`
	RunbookURL          string              `json:"runbook_url,omitempty"`
	Select              string              `json:"select_value,omitempty"`
	Type                string              `json:"type,omitempty"`
	UpdatedAt           serialization.Epoch `json:"updated_at_epoch_millis,omitempty"`
	ViolationCloseTimer *int                `json:"violation_close_timer,omitempty"`
	Warning             *Threshold          `json:"warning_threshold,omitempty"`
	Where               string              `json:"where_clause,omitempty"`
}
```

AlertCondition represents a New Relic Infrastructure alert condition.

#### type ErrorDetail

```go
>>>>>>> feat: add ListAlertConditions for infrastructure
type ErrorDetail struct {
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}
```

ErrorDetail represents the details of an error response from New Relic
Infrastructure.

#### type ErrorResponse

```go
type ErrorResponse struct {
	Errors []*ErrorDetail `json:"errors,omitempty"`
}
```

ErrorResponse represents the an error response from New Relic Infrastructure.

#### func (*ErrorResponse) Error

```go
func (e *ErrorResponse) Error() string
```
Error surfaces an error message from the Infrastructure error response

#### type Infrastructure

```go
type Infrastructure struct {
}
```


#### func  New

```go
func New(config config.Config) Infrastructure
```
New is used to create a new Synthetics client instance.

#### func (*Infrastructure) ListAlertConditions

```go
func (i *Infrastructure) ListAlertConditions(policyID int) ([]AlertCondition, error)
```
ListAlertConditions is used to retrieve New Relic Infrastructure alert
conditions.

#### type Threshold

```go
type Threshold struct {
	Duration int    `json:"duration_minutes,omitempty"`
	Function string `json:"time_function,omitempty"`
	Value    int    `json:"value,omitempty"`
}
```

Threshold represents an New Relic Infrastructure alert condition threshold.
