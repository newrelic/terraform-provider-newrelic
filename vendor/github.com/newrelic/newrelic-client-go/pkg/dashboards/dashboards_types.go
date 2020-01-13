package dashboards

import "time"

// Dashboard represents information about a New Relic dashboard.
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

// VisibilityType represents an option for the dashboard's visibility field.
type VisibilityType string

// EditableType represents an option for the dashboard's editable field.
type EditableType string

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

// DashboardMetadata represents metadata about the dashboard (like version)
type DashboardMetadata struct {
	Version int `json:"version"`
}

// DashboardWidget represents a widget in a dashboard.
type DashboardWidget struct {
	Visualization string                      `json:"visualization,omitempty"`
	ID            int                         `json:"widget_id,omitempty"`
	AccountID     int                         `json:"account_id,omitempty"`
	Data          []DashboardWidgetData       `json:"data,omitempty"`
	Presentation  DashboardWidgetPresentation `json:"presentation,omitempty"`
	Layout        DashboardWidgetLayout       `json:"layout,omitempty"`
}

// DashboardWidgetData represents the data backing a dashboard widget.
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

// DashboardWidgetDataCompareWith represents the compare with configuration of the widget.
type DashboardWidgetDataCompareWith struct {
	OffsetDuration string                                     `json:"offset_duration,omitempty"`
	Presentation   DashboardWidgetDataCompareWithPresentation `json:"presentation,omitempty"`
}

// DashboardWidgetDataCompareWithPresentation represents the compare with presentation configuration of the widget.
type DashboardWidgetDataCompareWithPresentation struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

// DashboardWidgetDataMetric represents the metrics data of the widget.
type DashboardWidgetDataMetric struct {
	Name   string   `json:"name,omitempty"`
	Units  string   `json:"units,omitempty"`
	Scope  string   `json:"scope,omitempty"`
	Values []string `json:"values,omitempty"`
}

// DashboardWidgetPresentation represents the visual presentation of a dashboard widget.
type DashboardWidgetPresentation struct {
	Title                string                    `json:"title,omitempty"`
	Notes                string                    `json:"notes,omitempty"`
	DrilldownDashboardID int                       `json:"drilldown_dashboard_id,omitempty"`
	Threshold            *DashboardWidgetThreshold `json:"threshold,omitempty"`
}

// DashboardWidgetThreshold represents the threshold configuration of a dashboard widget.
type DashboardWidgetThreshold struct {
	Red    float64 `json:"red,omitempty"`
	Yellow float64 `json:"yellow,omitempty"`
}

// DashboardWidgetLayout represents the layout of a widget in a dashboard.
type DashboardWidgetLayout struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Row    int `json:"row"`
	Column int `json:"column"`
}

// DashboardFilter represents the filter in a dashboard.
type DashboardFilter struct {
	EventTypes []string `json:"event_types,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
}
