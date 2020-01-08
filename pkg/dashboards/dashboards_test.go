// +build unit

package dashboards

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

func NewTestDashboards(handler http.Handler) Dashboards {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

var (
	testBillboardWidgetLayout = DashboardWidgetLayout{
		Width:  1,
		Height: 1,
		Row:    1,
		Column: 1,
	}

	testBillboardWidgetPresentation = DashboardWidgetPresentation{
		Title: "95th Percentile Load Time (ms)",
		Notes: "",
		Threshold: &DashboardWidgetThreshold{
			Red:    100,
			Yellow: 50,
		},
	}

	testBillboardWidgetData = []DashboardWidgetData{
		{
			NRQL: "SELECT percentile(duration, 95) FROM SyntheticCheck FACET monitorName since 7 days ago",
		},
	}

	testMarkdownWidgetLayout = DashboardWidgetLayout{
		Width:  1,
		Height: 1,
		Row:    1,
		Column: 2,
	}

	testMarkdownWidgetPresentation = DashboardWidgetPresentation{
		Title: "Links",
		Notes: "",
	}

	testMarkdownWidgetData = []DashboardWidgetData{
		{
			Source: "[test link](https://test.com)",
		},
	}

	testMetricsWidgetLayout = DashboardWidgetLayout{
		Width:  1,
		Height: 1,
		Row:    1,
		Column: 3,
	}

	testMetricsWidgetPresentation = DashboardWidgetPresentation{
		Title: "Links",
		Notes: "",
	}

	testMetricsWidgetData = []DashboardWidgetData{
		{
			Duration: 1800000,
			EndTime:  1800000000000,
			EntityIds: []int{
				1234,
			},
			CompareWith: []DashboardWidgetDataCompareWith{
				{
					OffsetDuration: "P7D",
					Presentation: DashboardWidgetDataCompareWithPresentation{
						Name:  "Last week",
						Color: "#b1b6ba",
					},
				},
				{
					OffsetDuration: "P1D",
					Presentation: DashboardWidgetDataCompareWithPresentation{
						Name:  "Yesterday",
						Color: "#77add4",
					},
				},
			},
			Metrics: []DashboardWidgetDataMetric{
				{
					Name:  "CPU/System/Utilization",
					Units: "",
					Scope: "",
					Values: []string{
						"percent",
					},
				},
			},
			RawMetricName: "CPU/System/Utilization",
			Facet:         "host",
			OrderBy:       "score",
			Limit:         10,
		},
	}

	testDashboardWidgets = []DashboardWidget{
		{
			Visualization: "billboard",
			AccountID:     1,
			Data:          testBillboardWidgetData,
			Presentation:  testBillboardWidgetPresentation,
			Layout:        testBillboardWidgetLayout,
		},
		{
			Visualization: "markdown",
			AccountID:     1,
			Data:          testMarkdownWidgetData,
			Presentation:  testMarkdownWidgetPresentation,
			Layout:        testMarkdownWidgetLayout,
		},
		{
			Visualization: "metric_line_chart",
			AccountID:     1,
			Data:          testMetricsWidgetData,
			Presentation:  testMetricsWidgetPresentation,
			Layout:        testMetricsWidgetLayout,
		},
	}

	testDashboardMetadata = DashboardMetadata{
		Version: 1,
	}

	testDashboardFilter = DashboardFilter{}

	testCreatedAt, _ = time.Parse(time.RFC3339, "2016-02-20T01:57:58Z")
	testUpdatedAt, _ = time.Parse(time.RFC3339, "2016-09-27T22:59:21Z")

	testDashboard = Dashboard{
		ID:         1234,
		Title:      "Test",
		Icon:       "bar-chart",
		Widgets:    testDashboardWidgets,
		Metadata:   testDashboardMetadata,
		Filter:     testDashboardFilter,
		Visibility: "all",
		Editable:   "editable_by_all",
		UIURL:      "https://insights.newrelic.com/accounts/1136088/dashboards/129507",
		APIURL:     "https://api.newrelic.com/v2/dashboards/129507",
		OwnerEmail: "foo@bar.com",
		CreatedAt:  testCreatedAt,
		UpdatedAt:  testUpdatedAt,
	}
	testDashboardJson = `
	{
		"id":1234,
		"title":"Test",
		"icon":"bar-chart",
		"created_at":"2016-02-20T01:57:58Z",
		"updated_at":"2016-09-27T22:59:21Z",
		"visibility":"all",
		"editable":"editable_by_all",
		"ui_url":"https://insights.newrelic.com/accounts/1136088/dashboards/129507",
		"api_url":"https://api.newrelic.com/v2/dashboards/129507",
		"owner_email":"foo@bar.com",
		"metadata":{
			"version":1
		},
		"filter":null,
		"widgets":[
			{
				"visualization":"billboard",
				"account_id":1,
				"data":[
					{
						"nrql":"SELECT percentile(duration, 95) FROM SyntheticCheck FACET monitorName since 7 days ago"
					}
				],
				"presentation":{
					"title":"95th Percentile Load Time (ms)",
					"notes":null,
					"drilldown_dashboard_id":null,
					"threshold":{
						"red":100,
						"yellow":50
					}
				},
				"layout":{
					"width":1,
					"height":1,
					"row":1,
					"column":1
				}
			},
			{
				"visualization":"markdown",
				"account_id":1,
				"data":[
					{
						"source":"[test link](https://test.com)"
					}
				],
				"presentation":{
					"title":"Links",
					"notes":null,
					"drilldown_dashboard_id":null
				},
				"layout":{
					"width":1,
					"height":1,
					"row":1,
					"column":2
				}
			},
			{
				"visualization":"metric_line_chart",
				"account_id":1,
				"data":[
					{
						"duration":1800000,
						"end_time":1800000000000,
						"entity_ids":[
							1234
						],
						"compare_with":[
							{
								"offset_duration": "P7D",
								"presentation": {
									"name": "Last week",
									"color": "#b1b6ba"
								}
							},
							{
								"offset_duration": "P1D",
								"presentation": {
									"name": "Yesterday",
									"color": "#77add4"
								}
							}
						  ],
						"metrics":[
							{
								"name":"CPU/System/Utilization",
								"units":null,
								"scope":"",
								"values":[
									"percent"
								]
							}
						],
						"order_by":"score",
						"limit":10,
						"facet":"host",
						"raw_metric_name":"CPU/System/Utilization"
					}
				],
				"presentation":{
					"title":"Links",
					"notes":null
				},
				"layout":{
					"width":1,
					"height":1,
					"row":1,
					"column":3
				}
			}
		]
	}`
)

func TestListDashboards(t *testing.T) {
	t.Parallel()
	client := NewTestDashboards(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "dashboard/json")
		_, err := w.Write([]byte(fmt.Sprintf(`
		{
			"dashboards": [%s]
		}
		`, testDashboardJson)))

		assert.NoError(t, err)
	}))

	expected := []*Dashboard{&testDashboard}

	actual, err := client.ListDashboards(nil)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestListDashboardsWithParams(t *testing.T) {
	t.Parallel()
	expectedCategory := "category"
	expectedTime := time.Now()
	expectedPage := 2
	expectedPerPage := 10
	expectedSort := "sort"
	expectedTitle := "title"

	client := NewTestDashboards(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		assert.Equal(t, expectedCategory, values.Get("filter[category]"))
		assert.Equal(t, expectedTitle, values.Get("filter[title]"))
		assert.Equal(t, expectedTime.Format(time.RFC3339), values.Get("filter[created_before]"))
		assert.Equal(t, expectedTime.Format(time.RFC3339), values.Get("filter[created_after]"))
		assert.Equal(t, expectedTime.Format(time.RFC3339), values.Get("filter[updated_before]"))
		assert.Equal(t, expectedTime.Format(time.RFC3339), values.Get("filter[updated_after]"))
		assert.Equal(t, expectedSort, values.Get("sort"))
		assert.Equal(t, strconv.Itoa(expectedPage), values.Get("page"))
		assert.Equal(t, strconv.Itoa(expectedPerPage), values.Get("per_page"))

		w.Header().Set("Content-Type", "dashboard/json")
		_, err := w.Write([]byte(`{"dashboards":[]}`))

		assert.NoError(t, err)
	}))

	params := ListDashboardsParams{
		Category:      expectedCategory,
		CreatedAfter:  &expectedTime,
		CreatedBefore: &expectedTime,
		Page:          expectedPage,
		PerPage:       expectedPerPage,
		Sort:          expectedSort,
		Title:         expectedTitle,
		UpdatedAfter:  &expectedTime,
		UpdatedBefore: &expectedTime,
	}

	_, err := client.ListDashboards(&params)

	assert.NoError(t, err)
}

func TestGetDashboard(t *testing.T) {
	t.Parallel()

	client := NewTestDashboards(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "dashboard/json")
		_, err := w.Write([]byte(fmt.Sprintf(`
		{
			"dashboard": %s
		}
		`, testDashboardJson)))

		assert.NoError(t, err)
	}))

	actual, err := client.GetDashboard(testDashboard.ID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testDashboard, actual)
}

func TestCreateDashboard(t *testing.T) {
	t.Parallel()

	client := NewTestDashboards(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "dashboard/json")
		_, err := w.Write([]byte(fmt.Sprintf(`
		{
			"dashboard": %s
		}
		`, testDashboardJson)))

		assert.NoError(t, err)
	}))

	actual, err := client.CreateDashboard(testDashboard)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testDashboard, actual)
}

func TestUpdateDashboard(t *testing.T) {
	t.Parallel()

	client := NewTestDashboards(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "dashboard/json")
		_, err := w.Write([]byte(fmt.Sprintf(`
		{
			"dashboard": %s
		}
		`, testDashboardJson)))

		assert.NoError(t, err)
	}))

	actual, err := client.UpdateDashboard(testDashboard)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testDashboard, actual)
}

func TestDeleteDashboard(t *testing.T) {
	t.Parallel()

	client := NewTestDashboards(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "dashboard/json")
		_, err := w.Write([]byte(fmt.Sprintf(`
		{
			"dashboard": %s
		}
		`, testDashboardJson)))

		assert.NoError(t, err)
	}))

	actual, err := client.DeleteDashboard(testDashboard.ID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testDashboard, actual)
}
