package newrelic

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

func resourceNewRelicOneDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicOneDashboardCreate,
		ReadContext:   resourceNewRelicOneDashboardRead,
		UpdateContext: resourceNewRelicOneDashboardUpdate,
		DeleteContext: resourceNewRelicOneDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		//		CustomizeDiff: sortSampleAttributeList,

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			return ignoreLineWidgetOrderDiff(diff, meta)
		},

		//CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
		//	return alignWidgetsOrderWithTerraformConfig(diff)
		//},

		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dashboard's name.",
			},
			"page": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				Elem:        dashboardPageSchemaElem(),
			},
			// Optional
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where you want to create the dashboard.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The dashboard's description.",
			},
			"permissions": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "public_read_only",
				ValidateFunc: validation.StringInSlice([]string{"private", "public_read_only", "public_read_write"}, false),
				Description:  "Determines who can see or edit the dashboard. Valid values are private, public_read_only, public_read_write. Defaults to public_read_only.",
			},
			// Computed
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the dashboard in New Relic.",
			},
			"permalink": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the dashboard.",
			},
			"variable": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Dashboard-local variable definitions.",
				Elem:        dashboardVariableSchemaElem(),
			},
		},
	}
}

func dashboardVariableSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Default values for this variable.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_multi_selection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether this variable supports multiple selection or not. Only applies to variables of type NRQL or ENUM.",
			},
			"item": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of possible values for variables of type ENUM",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-friendly display string for this value.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A possible variable value",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The variable identifier.",
			},
			"nrql_query": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration for variables of type NRQL.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "New Relic account ID(s) to issue the query against.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NRQL formatted query.",
						},
					},
				},
			},
			"replacement_strategy": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Indicates the strategy to apply when replacing a variable in a NRQL query.",
				ValidateFunc: validation.StringInSlice([]string{"default", "identifier", "number", "string"}, false),
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human-friendly display string for this variable.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies the data type of the variable and where its possible values may come from.",
				ValidateFunc: validation.StringInSlice([]string{"enum", "nrql", "string"}, false),
			},
		},
	}
}

// dashboardPageElem returns the schema for a New Relic dashboard Page
func dashboardPageSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The dashboard page's description.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dashboard page's name.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the dashboard page in New Relic.",
			},

			// All the widget types below
			"widget_area": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An area widget.",
				Elem:        dashboardWidgetAreaSchemaElem(),
			},
			"widget_bar": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A bar widget.",
				Elem:        dashboardWidgetBarSchemaElem(),
			},
			"widget_billboard": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A billboard widget.",
				Elem:        dashboardWidgetBillboardSchemaElem(),
			},
			"widget_bullet": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A bullet widget.",
				Elem:        dashboardWidgetBulletSchemaElem(),
			},
			"widget_funnel": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A funnel widget.",
				Elem:        dashboardWidgetFunnelSchemaElem(),
			},
			"widget_heatmap": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A heatmap widget.",
				Elem:        dashboardWidgetHeatmapSchemaElem(),
			},
			"widget_histogram": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A histogram widget.",
				Elem:        dashboardWidgetHistogramSchemaElem(),
			},
			"widget_line": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A line widget.",
				Elem:        dashboardWidgetLineSchemaElem(),
				//DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
				//	return widgetOrderSuppressFunc(old, new, d)
				//},
				// DiffSuppressFunc: elementOrderDiffSuppressFunc,

			},
			"widget_markdown": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A markdown widget.",
				Elem:        dashboardWidgetMarkdownSchemaElem(),
			},
			"widget_pie": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A pie widget.",
				Elem:        dashboardWidgetPieSchemaElem(),
			},
			"widget_log_table": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A log table widget.",
				Elem:        dashboardWidgetLogTableSchemaElem(),
			},
			"widget_table": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A table widget.",
				Elem:        dashboardWidgetTableSchemaElem(),
			},
			"widget_json": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A JSON widget.",
				Elem:        dashboardWidgetJSONSchemaElem(),
			},
			"widget_stacked_bar": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A stacked bar widget.",
				Elem:        dashboardWidgetStackedBarSchemaElem(),
			},
		},
	}
}

func dashboardWidgetSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ID of the widget.",
		},
		"title": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A title for the widget.",
		},
		"column": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"height": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      3,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"row": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"width": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(1, 12),
		},
		"nrql_query": {
			Type:     schema.TypeList,
			Required: true,
			Elem:     dashboardWidgetNRQLQuerySchemaElem(),
		},
		"ignore_time_range": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"facet_show_other_series": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"legend_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"y_axis_left_min": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		"y_axis_left_max": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		"null_values": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     dashboardWidgetNullValuesSchemaElem(),
		},
		"units": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     dashboardWidgetUnitsSchemaElem(),
		},
		"colors": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     dashboardWidgetColorSchemaElem(),
		},
	}
}

func dashboardWidgetColorSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"series_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Color code",
						},
						"series_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Series name",
						},
					},
				},
			},
		},
	}
}

func dashboardWidgetUnitsSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"series_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unit name",
						},
						"series_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Series name",
						},
					},
				},
			},
		},
	}
}

func dashboardWidgetNullValuesSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"null_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"series_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"null_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Null value",
						},
						"series_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Series name",
						},
					},
				},
			},
		},
	}
}

// dashboardWidgetNRQLQuerySchemaElem defines a NRQL query for use on a dashboard
//
// see: newrelic/newrelic-client-go/pkg/entities/DashboardWidgetQuery
func dashboardWidgetNRQLQuerySchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The account id used for the NRQL query.",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The NRQL query.",
			},
		},
	}
}

func dashboardWidgetAreaSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetBarSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["linked_entity_guids"] = dashboardWidgetLinkedEntityGUIDsSchema()
	s["filter_current_dashboard"] = dashboardWidgetFilterCurrentDashboardSchema()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetBillboardSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["critical"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The critical threshold value.",
	}

	s["warning"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The warning threshold value.",
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetBulletSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["limit"] = &schema.Schema{
		Type:        schema.TypeFloat,
		Required:    true,
		Description: "The maximum value for the visualization",
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetFunnelSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetHeatmapSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["linked_entity_guids"] = dashboardWidgetLinkedEntityGUIDsSchema()
	s["filter_current_dashboard"] = dashboardWidgetFilterCurrentDashboardSchema()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetHistogramSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetLineSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["y_axis_left_zero"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies if the values on the graph to be rendered need to be fit to scale, or printed within the specified range.",
	}
	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetJSONSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetMarkdownSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	delete(s, "nrql_query") // No queries for Markdown

	s["text"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetStackedBarSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetPieSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["linked_entity_guids"] = dashboardWidgetLinkedEntityGUIDsSchema()
	s["filter_current_dashboard"] = dashboardWidgetFilterCurrentDashboardSchema()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetLogTableSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetTableSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["linked_entity_guids"] = dashboardWidgetLinkedEntityGUIDsSchema()
	s["filter_current_dashboard"] = dashboardWidgetFilterCurrentDashboardSchema()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetLinkedEntityGUIDsSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		Computed:    true,
		Description: "Related entities. Currently only supports Dashboard entities, but may allow other cases in the future.",
	}
}

func dashboardWidgetFilterCurrentDashboardSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Use this item to filter the current dashboard",
	}
}

//func getWidgetsFromTerraformConfig(d *schema.ResourceData) []interface{} {
//	configWidgets := d.Get("widget_line").([]interface{})
//	widgets := make([]interface{}, len(configWidgets))
//	for idx, widget := range configWidgets {
//		widgets[idx] = widget
//	}
//	return widgets
//}

func resourceNewRelicOneDashboardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	defaultInfo := map[string]interface{}{
		"account_id": accountID,
	}
	dashboard, err := expandDashboardInput(d, defaultInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	filterWidgets, err := findDashboardWidgetFilterCurrentDashboard(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic One dashboard: %s", dashboard.Name)

	created, err := client.Dashboards.DashboardCreateWithContext(ctx, accountID, *dashboard)
	if err != nil {
		return diag.FromErr(err)
	}
	guid := created.EntityResult.GUID
	if guid == "" {
		var errMessages string
		for _, e := range created.Errors {
			errMessages += "[" + string(e.Type) + ": " + e.Description + "]"
		}

		return diag.Errorf("err: newrelic_one_dashboard Create failed: %s", errMessages)
	}

	log.Printf("[INFO] New Dashboard GUID: %s", guid)

	d.SetId(string(guid))

	res := resourceNewRelicOneDashboardRead(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Number of widgets with filter_current_dashboard: %d", len(filterWidgets))
	if len(filterWidgets) > 0 {

		err = setDashboardWidgetFilterCurrentDashboardLinkedEntity(d, filterWidgets)
		if err != nil {
			return diag.FromErr(err)
		}

		dashboard, err := expandDashboardInput(d, defaultInfo)
		if err != nil {
			return diag.FromErr(err)
		}

		result, err := client.Dashboards.DashboardUpdateWithContext(ctx, *dashboard, guid)
		if err != nil {
			return diag.FromErr(err)
		}

		return diag.FromErr(flattenDashboardUpdateResult(result, d))

	}

	return res
}

// resourceNewRelicOneDashboardRead NerdGraph => Terraform reader
func resourceNewRelicOneDashboardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic One dashboard %s", d.Id())

	dashboard, err := client.Dashboards.GetDashboardEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenDashboardEntity(dashboard, d))
}

func resourceNewRelicOneDashboardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	defaultInfo := map[string]interface{}{
		"account_id": accountID,
	}

	filterWidgets, err := findDashboardWidgetFilterCurrentDashboard(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the linked guid entities, if the page guid value is empty it is set to nil
	err = setDashboardWidgetFilterCurrentDashboardLinkedEntity(d, filterWidgets)
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := expandDashboardInput(d, defaultInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating New Relic One dashboard '%s' (%s)", dashboard.Name, d.Id())

	updated, err := client.Dashboards.DashboardUpdateWithContext(ctx, *dashboard, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	guid := updated.EntityResult.GUID
	if guid == "" {
		var errMessages string
		for _, e := range updated.Errors {
			errMessages += "[" + string(e.Type) + ": " + e.Description + "]"
		}

		return diag.Errorf("err: newrelic_one_dashboard Update failed: %s", errMessages)
	}

	diagErr := resourceNewRelicOneDashboardRead(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Printf("[INFO] Number of widgets with filter_current_dashboard: %d", len(filterWidgets))
	// If there are widgets with filter_current_dashboard, we need to update the linked guid entities
	if len(filterWidgets) > 0 {
		err = setDashboardWidgetFilterCurrentDashboardLinkedEntity(d, filterWidgets)
		if err != nil {
			return diag.FromErr(err)
		}

		dashboard, err = expandDashboardInput(d, defaultInfo)
		if err != nil {
			return diag.FromErr(err)
		}

		updated, err = client.Dashboards.DashboardUpdateWithContext(ctx, *dashboard, guid)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// We have to use the Update Result, not a re-read of the entity as the changes take
	// some amount of time to be re-indexed
	return diag.FromErr(flattenDashboardUpdateResult(updated, d))
}

func resourceNewRelicOneDashboardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One dashboard %v", d.Id())

	if _, err := client.Dashboards.DashboardDeleteWithContext(ctx, common.EntityGUID(d.Id())); err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

//func convertToCanonical(input string) map[string]interface{} {
//	var data map[string]interface{}
//	if err := json.Unmarshal([]byte(input), &data); err != nil {
//		return nil
//	}
//
//	keys := make([]string, 0, len(data))
//	for key := range data {
//		keys = append(keys, key)
//	}
//	sort.Strings(keys)
//
//	canonicalMap := make(map[string]interface{})
//	for _, key := range keys {
//		canonicalMap[key] = data[key]
//	}
//
//	log.Println("ENTERED THIS STATEMENT")
//	log.Println(canonicalMap)
//	return canonicalMap
//}

//func widgetLineOrderSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
//	oldList := strings.Split(old, ",")
//	newList := strings.Split(new, ",")
//
//	sort.Strings(oldList)
//	sort.Strings(newList)
//
//	if reflect.DeepEqual(oldList, newList) {
//		return true
//	}
//
//	return false
//}

//func widgetOrderSuppressFunc(d *schema.ResourceData) bool {
//	oldValue, newValue := d.GetChange("widget_line")
//	oldWidgets, newWidgets := toWidgets(oldValue, newValue)
//
//	if len(oldWidgets) != len(newWidgets) {
//		return false
//	}
//
//	matched := 0
//	for _, oldWidget := range oldWidgets {
//		for _, newWidget := range newWidgets {
//			if reflect.DeepEqual(oldWidget, newWidget) {
//				matched++
//				break
//			}
//		}
//	}
//
//	return matched == len(oldWidgets)
//}

//func widgetOrderSuppressFunc(d *schema.ResourceData) bool {
//	oldValue, newValue := d.GetChange("widget_line")
//	oldWidgets, newWidgets := toWidgets(oldValue, newValue)
//
//	if len(oldWidgets) != len(newWidgets) {
//		return false
//	}
//
//	matched := 0
//	for i, oldWidget := range oldWidgets {
//		if reflect.DeepEqual(oldWidget, newWidgets[i]) {
//			matched++
//			continue
//		}
//
//		found := false
//		for _, newWidget := range newWidgets {
//			if reflect.DeepEqual(oldWidget, newWidget) {
//				found = true
//				break
//			}
//		}
//
//		if found {
//			matched++
//		} else {
//			return false
//		}
//	}
//
//	return matched == len(oldWidgets)
//}

//func widgetOrderSuppressFunc(d *schema.ResourceData) bool {
//	oldValue, newValue := d.GetChange("widget_line")
//	oldWidgets, newWidgets := toWidgets(oldValue, newValue)
//
//	if len(oldWidgets) != len(newWidgets) {
//		return false
//	}
//
//	indexMismatches := 0
//	matchedSet := make(map[int]bool)
//
//	for i, oldWidget := range oldWidgets {
//		if !reflect.DeepEqual(oldWidget, newWidgets[i]) {
//			// Found a mismatch in index, search for possible matches in newWidgets.
//			indexMismatches++
//			matchIndex := -1
//
//			for j, newWidget := range newWidgets {
//				if reflect.DeepEqual(oldWidget, newWidget) && !matchedSet[j] {
//					// Found a matching widget in a different position.
//					matchIndex = j
//					break
//				}
//			}
//
//			if matchIndex == -1 {
//				// No matching widget found in newWidgets, actual config change.
//				return false
//			} else {
//				// Matching widget found in a different position.
//				matchedSet[matchIndex] = true
//			}
//		}
//	}
//
//	return indexMismatches == len(oldWidgets) || indexMismatches == 0
//}
//
//func toWidgets(oldInterface, newInterface interface{}) ([]interface{}, []interface{}) {
//	//oldWidgetLines, newWidgetLines := oldInterface.([]interface{}), newInterface.([]interface{})
//	var oldWidgetLines, newWidgetLines []interface{}
//
//	if oldInterface != nil {
//		oldWidgetLines = oldInterface.([]interface{})
//	}
//
//	if newInterface != nil {
//		newWidgetLines = newInterface.([]interface{})
//	}
//	oldWidgets := make([]interface{}, len(oldWidgetLines))
//	newWidgets := make([]interface{}, len(newWidgetLines))
//
//	for idx, item := range oldWidgetLines {
//		widget := item.(map[string]interface{})
//		oldWidgets[idx] = widget
//	}
//
//	for idx, item := range newWidgetLines {
//		widget := item.(map[string]interface{})
//		newWidgets[idx] = widget
//	}
//
//	return oldWidgets, newWidgets
//}

//func widgetOrderSuppressFunc(old, new string, d *schema.ResourceData) bool {
//	oldValue, newValue := d.GetChange("widget_line")
//	var oldWidgets, newWidgets []interface{}
//	if oldValue != nil {
//		oldWidgets = oldValue.([]interface{})
//	}
//	if newValue != nil {
//		newWidgets = newValue.([]interface{})
//	}
//
//	if len(oldWidgets) != len(newWidgets) {
//		return false
//	}
//
//	// Sort lists without considering the order of elements
//	sort.SliceStable(oldWidgets, func(i, j int) bool {
//		return widgetComparison(oldWidgets[i], oldWidgets[j]) < 0
//	})
//
//	sort.SliceStable(newWidgets, func(i, j int) bool {
//		return widgetComparison(newWidgets[i], newWidgets[j]) < 0
//	})
//
//	for i, oldWidget := range oldWidgets {
//		newWidget := newWidgets[i]
//
//		if !reflect.DeepEqual(oldWidget, newWidget) {
//			return false
//		}
//	}
//
//	return true
//}
//
//func widgetComparison(a, b interface{}) int {
//	widgetA := a.(map[string]interface{})
//	widgetB := b.(map[string]interface{})
//
//	param1Comparison := strings.Compare(widgetA["row"].(string), widgetB["row"].(string))
//	if param1Comparison != 0 {
//		return param1Comparison
//	}
//
//	param2Comparison := strings.Compare(widgetA["column"].(string), widgetB["column"].(string))
//	if param2Comparison != 0 {
//		return param2Comparison
//	}
//
//	return 0
//}

//func elementOrderDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
//	log.Println("ENTERED HERE")
//	log.Println("OLD")
//	log.Println(old)
//	log.Println("NEW")
//	log.Println(new)
//	log.Printf("[DEBUG] Key: %q; Old value: %q; New value: %q", k, old, new) // Log the values of k, old, and new
//	return true
//}

//func sortSampleAttributeList(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
//	log.Println("REACHED HERE")
//	oldList, _ := d.GetChange("widget_line")
//	var oldListTyped []interface{}
//	if oldList != nil {
//		oldListTyped = oldList.([]interface{})
//	}
//	log.Println(oldListTyped)
//	var newList []interface{}
//	if d.Get("widget_line") != nil {
//		newList = d.Get("widget_line").([]interface{})
//	}
//	log.Println(newList)
//	// Sort `newList` according to the order of elements in `oldList`.
//
//	//sortedNewList := make([]interface{}, len(newList))
//	//// Implement your sorting logic based on attributes in the old list.
//	//// ...
//	//
//	//if !reflect.DeepEqual(sortedNewList, newList) {
//	//	err := d.SetNew("sample_attribute", sortedNewList)
//	//	if err != nil {
//	//		return fmt.Errorf("failed to set sorted sample_attribute list in the plan: %s", err)
//	//	}
//	//}
//
//	return nil
//}

func alignWidgetsOrderWithTerraformConfig(diff *schema.ResourceDiff) error {
	log.Println("REACHED HERE")
	if !diff.HasChange("widget_line") {
		log.Println("DID NOT REACH HERE")
		oldValue, newValue := diff.GetChange("line_widgets")

		// Sort the widgets values
		oldWidgets := oldValue.([]interface{})
		log.Println(oldWidgets)
		newWidgets := newValue.([]interface{})
		log.Println(newWidgets)

		x := *diff
		log.Println(x)
		attributeNames := diff.GetChangedKeysPrefix("")
		for _, key := range attributeNames {
			oldValue, newValue := diff.GetChange(key)
			log.Println("----------")
			log.Printf("- %s:\n", key)
			log.Printf("  Old Value: %v\n", oldValue)
			log.Printf("  New Value: %v\n", newValue)
			log.Println("----------")
		}

		log.Println(diff.Get("widget_line"))
		return nil
	}
	log.Println("ENTERED HERE")
	log.Println("ENTERED HERE")
	log.Println(diff.Get("widget_line"))

	//oldValue, newValue := diff.GetChange("widget_line")
	//oldWidgets := oldValue.([]interface{})
	//newWidgets := newValue.([]interface{})
	//
	//if len(oldWidgets) != len(newWidgets) {
	//	return nil
	//}
	//
	//sortedNewWidgets := make([]interface{}, len(newWidgets))
	//for idx, oldWidget := range oldWidgets {
	//	oldWidgetMap := oldWidget.(map[string]interface{})
	//
	//	found := false
	//	for idxNew, newWidget := range newWidgets {
	//		newWidgetMap := newWidget.(map[string]interface{})
	//
	//		if reflect.DeepEqual(oldWidgetMap, newWidgetMap) {
	//			sortedNewWidgets[idx] = newWidgets[idxNew]
	//			found = true
	//			break
	//		}
	//	}
	//
	//	if !found {
	//		return nil
	//	}
	//}
	//
	//for idx, widget := range sortedNewWidgets {
	//	prefixedKey := "widget_line." + strconv.Itoa(idx)
	//	for k, v := range widget.(map[string]interface{}) {
	//		err := diff.SetNew(prefixedKey+"."+k, v)
	//		if err != nil {
	//			return fmt.Errorf("[DEBUG] error setting widget_line attribute: %w", err)
	//		}
	//	}
	//}
	//
	return nil

}

//
//
//
//
//
//
//
//
//
//
//

// lineWidget represents the basic structure for each line_widget
type lineWidget struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Add other attributes from the line_widget schema
}

func ignoreLineWidgetOrderDiff(d *schema.ResourceDiff, meta interface{}) error {
	oldValue, newValue := d.GetChange("page")

	oldPages := oldValue.([]interface{})
	newPages := newValue.([]interface{})

	// Iterate through the 'page' attributes and sort the 'line_widgets'
	for pageIndex := range oldPages {
		log.Println("REACHED PAGE BLOCK IN IGNORE")
		oldWidgets, ok1 := oldPages[pageIndex].(map[string]interface{})["line_widgets"].([]interface{})
		newWidgets, ok2 := newPages[pageIndex].(map[string]interface{})["line_widgets"].([]interface{})

		if ok1 && ok2 {
			log.Println("REACHED WIDGET BLOCK IN IGNORE")
			log.Println(oldWidgets)
			log.Println(newWidgets)
			sortLineWidgets(oldWidgets)
			sortLineWidgets(newWidgets)

			pageKey := fmt.Sprintf("page.%d.line_widgets", pageIndex)
			if err := d.SetNewComputed(pageKey); err != nil {
				return err
			}
		}
	}

	return nil
}

func sortLineWidgets(widgets []interface{}) {
	// Convert each map element to a lineWidget struct and store them in a slice.
	var lineWidgets []lineWidget
	for _, widget := range widgets {
		widgetMap := widget.(map[string]interface{})
		lineWidgets = append(lineWidgets, lineWidget{
			ID:    widgetMap["id"].(string),
			Title: widgetMap["title"].(string),
			// Add other attributes from the line_widget schema
		})
	}

	// Sort line_widgets slice based on their ID (or any other suitable attribute)
	sort.Slice(lineWidgets, func(i, j int) bool {
		return lineWidgets[i].ID < lineWidgets[j].ID
	})

	// Set the sorted values back to the original widgets slice.
	for i, widget := range lineWidgets {
		widgets[i] = map[string]interface{}{
			"id":    widget.ID,
			"title": widget.Title,
			// Add other attributes from the line_widget schema
		}
	}
}
