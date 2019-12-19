# dashboards
--
    import "."


## Usage

```go
var (
	// Visibility specifies the possible options for a dashboard's visibility.
	Visibility = struct {
		Owner VisibilityType
		All   VisibilityType
	}{
		Owner: "owner",
		All:   "all",
	}

	// Editable specifies the possible options for who can edit a dashboard.
	Editable = struct {
		Owner    EditableType
		All      EditableType
		ReadOnly EditableType
	}{
		Owner:    "editable_by_owner",
		All:      "editable_by_all",
		ReadOnly: "read_only",
	}
)
```

#### type Dashboard

```go
type Dashboard struct {
	ID         int               `json:"id"`
	Title      string            `json:"title,omitempty"`
	Icon       string            `json:"icon,omitempty"`
	CreatedAt  time.Time         `json:"created_at,omitempty"`
	UpdatedAt  time.Time         `json:"updated_at,omitempty"`
	Visibility VisibilityType    `json:"visibility,omitempty"`
	Editable   EditableType      `json:"editable,omitempty"`
	UIURL      string            `json:"ui_url,omitempty"`
	APIURL     string            `json:"api_url,omitempty"`
	OwnerEmail string            `json:"owner_email,omitempty"`
	Metadata   DashboardMetadata `json:"metadata"`
	Filter     DashboardFilter   `json:"filter,omitempty"`
	Widgets    []DashboardWidget `json:"widgets,omitempty"`
}
```

Dashboard represents information about a New Relic dashboard.

#### type DashboardFilter

```go
type DashboardFilter struct {
	EventTypes []string `json:"event_types,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
}
```

DashboardFilter represents the filter in a dashboard.

#### type DashboardMetadata

```go
type DashboardMetadata struct {
	Version int `json:"version"`
}
```

DashboardMetadata represents metadata about the dashboard (like version)

#### type DashboardWidget

```go
type DashboardWidget struct {
	Visualization string                      `json:"visualization,omitempty"`
	ID            int                         `json:"widget_id,omitempty"`
	AccountID     int                         `json:"account_id,omitempty"`
	Data          []DashboardWidgetData       `json:"data,omitempty"`
	Presentation  DashboardWidgetPresentation `json:"presentation,omitempty"`
	Layout        DashboardWidgetLayout       `json:"layout,omitempty"`
}
```

DashboardWidget represents a widget in a dashboard.

#### type DashboardWidgetData

```go
type DashboardWidgetData struct {
	NRQL          string                           `json:"nrql,omitempty"`
	Source        string                           `json:"source,omitempty"`
	Duration      int                              `json:"duration,omitempty"`
	EndTime       int                              `json:"end_time,omitempty"`
	EntityIds     []int                            `json:"entity_ids,omitempty"`
	CompareWith   []DashboardWidgetDataCompareWith `json:"compare_with,omitempty"`
	Metrics       []DashboardWidgetDataMetric      `json:"metrics,omitempty"`
	RawMetricName string                           `json:"raw_metric_name,omitempty"`
	Facet         string                           `json:"facet,omitempty"`
	OrderBy       string                           `json:"order_by,omitempty"`
	Limit         int                              `json:"limit,omitempty"`
}
```

DashboardWidgetData represents the data backing a dashboard widget.

#### type DashboardWidgetDataCompareWith

```go
type DashboardWidgetDataCompareWith struct {
	OffsetDuration string                                     `json:"offset_duration,omitempty"`
	Presentation   DashboardWidgetDataCompareWithPresentation `json:"presentation,omitempty"`
}
```

DashboardWidgetDataCompareWith represents the compare with configuration of the
widget.

#### type DashboardWidgetDataCompareWithPresentation

```go
type DashboardWidgetDataCompareWithPresentation struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}
```

DashboardWidgetDataCompareWithPresentation represents the compare with
presentation configuration of the widget.

#### type DashboardWidgetDataMetric

```go
type DashboardWidgetDataMetric struct {
	Name   string   `json:"name,omitempty"`
	Units  string   `json:"units,omitempty"`
	Scope  string   `json:"scope,omitempty"`
	Values []string `json:"values,omitempty"`
}
```

DashboardWidgetDataMetric represents the metrics data of the widget.

#### type DashboardWidgetLayout

```go
type DashboardWidgetLayout struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Row    int `json:"row"`
	Column int `json:"column"`
}
```

DashboardWidgetLayout represents the layout of a widget in a dashboard.

#### type DashboardWidgetPresentation

```go
type DashboardWidgetPresentation struct {
	Title                string                    `json:"title,omitempty"`
	Notes                string                    `json:"notes,omitempty"`
	DrilldownDashboardID int                       `json:"drilldown_dashboard_id,omitempty"`
	Threshold            *DashboardWidgetThreshold `json:"threshold,omitempty"`
}
```

DashboardWidgetPresentation represents the visual presentation of a dashboard
widget.

#### type DashboardWidgetThreshold

```go
type DashboardWidgetThreshold struct {
	Red    float64 `json:"red,omitempty"`
	Yellow float64 `json:"yellow,omitempty"`
}
```

DashboardWidgetThreshold represents the threshold configuration of a dashboard
widget.

#### type Dashboards

```go
type Dashboards struct {
}
```

Dashboards is used to communicate with the New Relic Dashboards product.

#### func  New

```go
func New(config config.Config) Dashboards
```
New is used to create a new Dashboards client instance.

#### func (*Dashboards) CreateDashboard

```go
func (dashboards *Dashboards) CreateDashboard(dashboard Dashboard) (*Dashboard, error)
```
CreateDashboard is used to create a New Relic dashboard.

#### func (*Dashboards) DeleteDashboard

```go
func (dashboards *Dashboards) DeleteDashboard(dashboardID int) (*Dashboard, error)
```
DeleteDashboard is used to delete a New Relic dashboard.

#### func (*Dashboards) GetDashboard

```go
func (dashboards *Dashboards) GetDashboard(dashboardID int) (*Dashboard, error)
```
GetDashboard is used to retrieve a single New Relic dashboard.

#### func (*Dashboards) ListDashboards

```go
func (dashboards *Dashboards) ListDashboards(params *ListDashboardsParams) ([]Dashboard, error)
```
ListDashboards is used to retrieve New Relic dashboards.

#### func (*Dashboards) UpdateDashboard

```go
func (dashboards *Dashboards) UpdateDashboard(dashboard Dashboard) (*Dashboard, error)
```
UpdateDashboard is used to update a New Relic dashboard.

#### type EditableType

```go
type EditableType string
```

EditableType represents an option for the dashboard's editable field.

#### type ListDashboardsParams

```go
type ListDashboardsParams struct {
	Category      string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	Page          int
	PerPage       int
	Sort          string
	Title         string
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
}
```

ListDashboardsParams represents a set of filters to be used when querying New
Relic dashboards.

#### type VisibilityType

```go
type VisibilityType string
```

VisibilityType represents an option for the dashboard's visibility field.
