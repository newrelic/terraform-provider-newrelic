package newrelic

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func resourceNewRelicDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicDashboardCreate,
		Read:   resourceNewRelicDashboardRead,
		Update: resourceNewRelicDashboardUpdate,
		Delete: resourceNewRelicDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"icon": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "bar-chart",
			},
			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: validation.StringInSlice([]string{"owner", "all"}, false),
			},
			"dashboard_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"editable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "editable_by_all",
				ValidateFunc: validation.StringInSlice([]string{"read_only", "editable_by_owner", "editable_by_all", "all"}, false),
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_types": &schema.Schema{
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
							Set:      schema.HashString,
						},
						"attributes": &schema.Schema{
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Set:      schema.HashString,
						},
					},
				},
			},
			"widget": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 60,
				Set:      resourceNewRelicDashboardWidgetsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"visualization": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"width": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"height": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"row": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"column": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"notes": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						// TODO: Move this to a set/map?
						"nrql": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicDashboardWidgetsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	row := m["row"].(int)
	column := m["column"].(int)
	width := m["width"].(int)
	height := m["height"].(int)
	nrql := m["nrql"].(string)
	title := m["title"].(string)
	notes := m["notes"].(string)
	viz := m["visualization"].(string)

	buf.WriteString(fmt.Sprintf("%d-%d-%d-%d-%s-%s-%s-%s",
		row, column, width, height, nrql, title, viz, notes))

	return hashcode.String(buf.String())
}

// Assemble the *newrelic.Dashboard variable.
//
// Used by the newrelic_dashboard Create and Update functions.
func expandDashboard(d *schema.ResourceData) *newrelic.Dashboard {
	metadata := newrelic.DashboardMetadata{
		Version: 1,
	}

	// TODO: Some of these should be terraform defaults and validated
	dashboard := newrelic.Dashboard{
		Title:      d.Get("title").(string),
		Metadata:   metadata,
		Icon:       d.Get("icon").(string),
		Visibility: d.Get("visibility").(string),
		Editable:   d.Get("editable").(string),
	}

	if f, ok := d.GetOk("filter"); ok {
		filter := f.([]interface{})[0].(map[string]interface{})
		dashboardFilter := newrelic.DashboardFilter{}

		if v, ok := filter["attributes"]; ok {
			attributes := v.(*schema.Set).List()
			vs := make([]string, 0, len(attributes))
			for _, a := range attributes {
				vs = append(vs, a.(string))
			}

			dashboardFilter.Attributes = vs
		}

		if v, ok := filter["event_types"]; ok {
			eventTypes := v.(*schema.Set).List()
			vs := make([]string, 0, len(eventTypes))
			for _, e := range eventTypes {
				vs = append(vs, e.(string))
			}
			dashboardFilter.EventTypes = vs
		}
		dashboard.Filter = dashboardFilter
	}

	log.Printf("[INFO] widget schema: %+v\n", d.Get("widget"))

	if w, ok := d.GetOk("widget"); ok {
		widgets := w.(*schema.Set).List()
		for _, widget := range widgets {
			w := widget.(map[string]interface{})

			widgetPresentation := newrelic.DashboardWidgetPresentation{
				Title: w["title"].(string),
				Notes: w["notes"].(string),
			}

			widgetLayout := newrelic.DashboardWidgetLayout{
				Row:    w["row"].(int),
				Column: w["column"].(int),
				Width:  w["width"].(int),
				Height: w["height"].(int),
			}

			// TODO: Support non-NRQL Widgets
			widgetData := []newrelic.DashboardWidgetData{
				{
					NRQL: w["nrql"].(string),
				},
			}

			dashboard.Widgets = append(dashboard.Widgets, newrelic.DashboardWidget{
				Visualization: w["visualization"].(string),
				Layout:        widgetLayout,
				Presentation:  widgetPresentation,
				Data:          widgetData,
			})
		}
	}

	return &dashboard
}

// Unpack the *newrelic.Dashboard variable and set resource data.
//
// Used by the newrelic_dashboard Read function (resourceNewRelicDashboardRead)
func flattenDashboard(dashboard *newrelic.Dashboard, d *schema.ResourceData) error {
	d.Set("title", dashboard.Title)
	d.Set("icon", dashboard.Icon)
	d.Set("visibility", dashboard.Visibility)
	d.Set("editable", dashboard.Editable)
	d.Set("dashboard_url", dashboard.UIURL)
	if filterErr := d.Set("filter", flattenFilter(&dashboard.Filter)); filterErr != nil {
		return filterErr
	}
	widgetSet := schema.Set{
		F: resourceNewRelicDashboardWidgetsHash,
	}
	for _, widget := range dashboard.Widgets {
		values := map[string]interface{}{}
		values["visualization"] = widget.Visualization
		values["title"] = widget.Presentation.Title
		values["notes"] = widget.Presentation.Notes
		values["row"] = widget.Layout.Row
		values["column"] = widget.Layout.Column
		values["width"] = widget.Layout.Width
		values["height"] = widget.Layout.Height

		// TODO: Support non-NRQL Widgets
		if len(widget.Data) > 0 {
			values["nrql"] = widget.Data[0].NRQL
		}
		widgetSet.Add(values)
	}
	err := d.Set("widget", &widgetSet)
	if err != nil {
		return err
	}

	return nil
}

func resourceNewRelicDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	dashboard := expandDashboard(d)
	log.Printf("[INFO] Creating New Relic dashboard: %s", dashboard.Title)

	dashboard, err := client.CreateDashboard(*dashboard)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(dashboard.ID))

	return resourceNewRelicDashboardRead(d, meta)
}

func resourceNewRelicDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic dashboard %s", d.Id())

	dashboardID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	dashboard, err := client.GetDashboard(dashboardID)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenDashboard(dashboard, d)
}

func resourceNewRelicDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	dashboard := expandDashboard(d)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	dashboard.ID = id
	log.Printf("[INFO] Updating New Relic dashboard %d", id)

	_, err = client.UpdateDashboard(*dashboard)
	if err != nil {
		return err
	}

	return resourceNewRelicDashboardRead(d, meta)
}

func resourceNewRelicDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting New Relic dashboard %v", id)

	if err := client.DeleteDashboard(id); err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId("")

	return nil
}

func flattenFilter(f *newrelic.DashboardFilter) []interface{} {
	if f == nil {
		return nil
	}
	filterResult := make(map[string]interface{})

	attributesList := make([]interface{}, 0, len(f.Attributes))
	for _, v := range f.Attributes {
		attributesList = append(attributesList, v)
	}

	eventTypesList := make([]interface{}, 0, len(f.EventTypes))
	for _, v := range f.EventTypes {
		eventTypesList = append(eventTypesList, v)
	}

	filterResult["attributes"] = schema.NewSet(schema.HashString, attributesList)
	filterResult["event_types"] = schema.NewSet(schema.HashString, eventTypesList)
	return []interface{}{filterResult}
}
