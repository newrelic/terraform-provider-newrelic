package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

// termSchema returns the schema used for a critical or warning term priority.
func termSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Maps to `thresholdDuration` in NerdGraph and values are in seconds, not minutes.
			// Validation is different in NerdGraph - Value must be within 120-3600 seconds (2-60 minutes) and a multiple of 60 for BASELINE conditions.
			// Convert to seconds when using NerdGraph
			"duration": {
				Deprecated:   "use `threshold_duration` attribute instead",
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "In minutes, must be in the range of 1 to 120 (inclusive).",
				ValidateFunc: validation.IntBetween(1, 120),
			},
			// Value must be uppercase when using NerdGraph
			"operator": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "equals",
				Description:  "One of (above, above_or_equals, below, below_or_equals, equals, not_equals). Defaults to 'equals'.",
				ValidateFunc: validation.StringInSlice([]string{"above", "above_or_equals", "below", "below_or_equals", "equals", "not_equals"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"threshold": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "For baseline conditions must be in range [1, 1000].",
			},
			// Does not exist in NerdGraph. Equivalent to `threshold_occurrences`,
			// but with different wording.
			"time_function": {
				Deprecated:   "use `threshold_occurrences` attribute instead",
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Valid values are: 'all' or 'any'",
				ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
			},
			// NerdGraph only. Equivalent to `time_function`,
			// but with slightly different wording.
			// i.e. `any` (old) vs `AT_LEAST_ONCE` (new)
			"threshold_occurrences": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The criteria for how many data points must be in violation for the specified threshold duration. Valid values are: 'ALL' or 'AT_LEAST_ONCE' (case insensitive).",
				ValidateFunc: validation.StringInSlice([]string{"ALL", "AT_LEAST_ONCE"}, true),
				StateFunc: func(v interface{}) string {
					// Always store lowercase to prevent state drift
					return strings.ToLower(v.(string))
				},
			},
			// NerdGraph only. Equivalent to `duration`, but in seconds
			"threshold_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The duration, in seconds, that the threshold must violate in order to create an incident. Value must be a multiple of the 'aggregation_window' (which has a default of 60 seconds). Value must be within 120-86400 seconds for baseline conditions, and within 60-86400 seconds for static conditions",
			},
			"prediction": {
				Type:        schema.TypeSet,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				Description: "BETA PREVIEW: the `prediction` field is in limited release and only enabled for preview on a per-account basis. - Use `prediction` to open alerts when your static threshold is predicted to be reached in the future. The `prediction` field is only available for static conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"predict_by": {
							Type:     schema.TypeInt,
							Optional: true,
							// API default 3600 seconds = 1 hour
							Default:     3600,
							Description: "BETA PREVIEW: the `predict_by` field is in limited release and only enabled for preview on a per-account basis. - The duration, in seconds, that the prediction should look into the future.",
						},
						"prefer_prediction_violation": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "BETA PREVIEW: the `prefer_prediction_violation` field is in limited release and only enabled for preview on a per-account basis. - If a prediction incident is open when a term's static threshold is breached by the actual signal, default behavior is to close the prediction incident and open a static incident. Setting `prefer_prediction_violation` to `true` overrides this behavior leaving the prediction incident open and preventing a static incident from opening.",
						},
					},
				},
			},
			"disable_health_status_reporting": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Violations will not change system health status for this term.",
			},
		},
	}
}

func termSchemaDeprecated() *schema.Resource {
	rec := termSchema()

	prioritySchema := &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "critical",
		Description:  "One of (critical, warning). Defaults to 'critical'. At least one condition term must have priority set to 'critical'.",
		ValidateFunc: validation.StringInSlice([]string{"critical", "warning"}, false),
	}

	rec.Schema["priority"] = prioritySchema

	return rec
}

func resourceNewRelicNrqlAlertCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNrqlAlertConditionCreate,
		ReadContext:   resourceNewRelicNrqlAlertConditionRead,
		UpdateContext: resourceNewRelicNrqlAlertConditionUpdate,
		DeleteContext: resourceNewRelicNrqlAlertConditionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateWithMetadata(2, "type"),
		},
		CustomizeDiff: validateNrqlConditionAttributes,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the policy where this condition should be used.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the condition.",
			},
			"runbook_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Runbook URL to display in notifications.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not to enable the alert condition.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "static",
				Description:  "The type of NRQL alert condition to create. Valid values are: 'static', 'baseline', 'outlier'.",
				ValidateFunc: validation.StringInSlice([]string{"static", "baseline", "outlier"}, false),
			},
			"nrql": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "A NRQL query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:     schema.TypeString,
							Required: true,
						},
						"data_account_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The New Relic account ID to use as the basis for the NRQL alert condition's `query`; will default to `account_id` if unspecified.",
						},
						"since_value": {
							Deprecated:    "use `aggregation_method` attribute instead",
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "NRQL queries are evaluated in one-minute time windows. The start time depends on the value you provide in the NRQL condition's `since_value`.",
							ConflictsWith: []string{"nrql.0.evaluation_offset"},
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								valueString := val.(string)
								_, err := strconv.Atoi(valueString)
								if err != nil {
									errs = append(errs, fmt.Errorf("error converting string to int: %#v", err))
								}
								return
							},
						},
						// Equivalent to `since_value`.
						"evaluation_offset": {
							Deprecated:    "use `aggregation_method` attribute instead",
							Type:          schema.TypeInt,
							Optional:      true,
							Description:   "NRQL queries are evaluated in one-minute time windows. The start time depends on the value you provide in the NRQL condition's `evaluation_offset`.",
							ConflictsWith: []string{"nrql.0.since_value"},
						},
					},
				},
			},
			"term": {
				Type:          schema.TypeSet,
				MinItems:      1,
				MaxItems:      2,
				Optional:      true,
				Description:   "A set of terms for this condition. Max 2 terms allowed - at least one 1 critical term and 1 optional warning term.",
				Elem:          termSchemaDeprecated(),
				ConflictsWith: []string{"critical", "warning"},
				Deprecated:    "use `critical` and `warning` attributes instead",
			},
			"critical": {
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				Elem:          termSchema(),
				Description:   "A condition term with priority set to critical.",
				ConflictsWith: []string{"term"},
			},
			"warning": {
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				Elem:          termSchema(),
				Description:   "A condition term with priority set to warning.",
				ConflictsWith: []string{"term"},
			},
			"violation_time_limit_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  violationTimeLimitSecondsDefault,
				// Default value added as expected by the NerdGraph API to prevent discrepancies with `terraform plan`
				// Reference : https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/alert-violations/how-alert-condition-violations-are-closed/#time-limit
				Description:   "Sets a time limit, in seconds, that will automatically force-close a long-lasting incident after the time limit you select.  Must be in the range of 300 to 2592000 (inclusive)",
				ConflictsWith: []string{"violation_time_limit"},
				ValidateFunc:  validation.IntBetween(300, violationTimeLimitSecondsMax),
			},

			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID for managing your NRQL alert conditions.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the NRQL alert condition.",
			},
			"title_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This field allows you to create a custom title to be used when incidents are opened by the condition. Setting this field will override the default title. Must be Handlebars format.",
			},
			"violation_time_limit": {
				Type:          schema.TypeString,
				Deprecated:    "use `violation_time_limit_seconds` attribute instead",
				Optional:      true,
				Computed:      true,
				Description:   "Sets a time limit, in hours, that will automatically force-close a long-lasting incident after the time limit you select. Possible values are 'ONE_HOUR', 'TWO_HOURS', 'FOUR_HOURS', 'EIGHT_HOURS', 'TWELVE_HOURS', 'TWENTY_FOUR_HOURS', 'THIRTY_DAYS' (case insensitive).",
				ConflictsWith: []string{"violation_time_limit_seconds"},
				ValidateFunc:  validation.StringInSlice([]string{"ONE_HOUR", "TWO_HOURS", "FOUR_HOURS", "EIGHT_HOURS", "TWELVE_HOURS", "TWENTY_FOUR_HOURS", "THIRTY_DAYS"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
			},
			"open_violation_on_expiration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to create a new incident to capture that the signal expired.",
			},
			"close_violations_on_expiration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to close all open incidents when the signal expires.",
			},
			"aggregation_window": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The duration of the time window used to evaluate the NRQL query, in seconds.",
			},
			"slide_by": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The duration of overlapping time windows used to smooth the chart line, in seconds. Must be a factor of `aggregation_window` and less than the aggregation window. If `aggregation_window` is less than or equal to 3600 seconds, it should be greater or equal to 30 seconds. If `aggregation_window` is greater than 3600 seconds but less than 7200 seconds, it should be greater or equal to `aggregation_window / 120`.  If `aggregation_window` is greater than 7200 seconds, it should be greater or equal to `aggregation_window / 24",
			},
			"expiration_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The amount of time (in seconds) to wait before considering the signal expired.  Must be in the range of 30 to 172800 (inclusive)",
				ValidateFunc: validation.IntBetween(30, 172800),
			},
			"ignore_on_expected_termination": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to ignore expected termination of a signal when considering whether to create a loss of signal incident",
			},
			"fill_option": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Which strategy to use when filling gaps in the signal. If static, the 'fill value' will be used for filling gaps in the signal. Valid values are: 'NONE', 'LAST_VALUE', or 'STATIC' (case insensitive).",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "LAST_VALUE", "STATIC"}, true),
				StateFunc: func(v interface{}) string {
					// Always store lowercase to prevent state drift
					return strings.ToLower(v.(string))
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Assume that empty string and 'none' are the same for diff purposes due to API defaults
					return (old == "" || old == "none") && (new == "" || new == "none")
				},
			},
			"fill_value": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Description:  "If using the 'static' fill option, this value will be used for filling gaps in the signal.",
				RequiredWith: []string{"fill_option"},
			},
			"aggregation_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CADENCE", "EVENT_FLOW", "EVENT_TIMER"}, true),
				Description:  "The method that determines when we consider an aggregation window to be complete so that we can evaluate the signal for incidents. Default is EVENT_FLOW.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, sinceValueExists := d.GetOk("nrql.0.since_value")
					if sinceValueExists {
						return false
					}
					// If a value is not provided and the condition uses the default value, don't show a diff
					return (strings.EqualFold(old, "event_flow") && new == "") || strings.EqualFold(old, new)
				},
			},
			"aggregation_delay": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "How long we wait for data that belongs in each aggregation window. Depending on your data, a longer delay may increase accuracy but delay notifications. Use aggregationDelay with the EVENT_FLOW and CADENCE aggregation methods.",
				RequiredWith: []string{"aggregation_method"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, sinceValueExists := d.GetOk("nrql.0.since_value")
					if sinceValueExists {
						return false
					}
					// If a value is not provided and the condition uses the default value, don't show a diff
					oldInt, _ := strconv.ParseInt(old, 0, 8)
					newInt, _ := strconv.ParseInt(new, 0, 8)
					aggregationMethod := strings.ToLower(d.Get("aggregation_method").(string))
					return oldInt == 120 && newInt == 0 && (aggregationMethod == "event_flow" || aggregationMethod == "cadence")
				},
			},
			"evaluation_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "How long we wait until the signal starts evaluating. The maximum delay is 7200 seconds (120 minutes)",
				ValidateFunc: validation.IntBetween(1, 7200),
			},
			"aggregation_timer": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "How long we wait after each data point arrives to make sure we've processed the whole batch. Use aggregationTimer with the EVENT_TIMER aggregation method.",
				RequiredWith: []string{"aggregation_method"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, sinceValueExists := d.GetOk("nrql.0.since_value")
					if sinceValueExists {
						return false
					}
					// If a value is not provided and the condition uses the default value, don't show a diff
					oldInt, _ := strconv.ParseInt(old, 0, 8)
					newInt, _ := strconv.ParseInt(new, 0, 8)
					aggregationMethod := strings.ToLower(d.Get("aggregation_method").(string))
					return oldInt == 60 && newInt == 0 && aggregationMethod == "event_timer"
				},
			},
			"entity_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the NRQL Condition in New Relic.",
			},
			// Baseline ONLY
			"baseline_direction": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The baseline direction of a baseline NRQL alert condition. Valid values are: 'LOWER_ONLY', 'UPPER_AND_LOWER', 'UPPER_ONLY' (case insensitive).",
				ValidateFunc: validation.StringInSlice([]string{"LOWER_ONLY", "UPPER_AND_LOWER", "UPPER_ONLY"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
			},
			"signal_seasonality": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Seasonality under which a condition's signal(s) are evaluated. Valid values are: 'NEW_RELIC_CALCULATION', 'HOURLY', 'DAILY', 'WEEKLY', or 'NONE'. To have New Relic calculate seasonality automatically, set to 'NEW_RELIC_CALCULATION' (default). To turn off seasonality completely, set to 'NONE'.",
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(alerts.NrqlSignalSeasonalities.NewRelicCalculation),
						string(alerts.NrqlSignalSeasonalities.Hourly),
						string(alerts.NrqlSignalSeasonalities.Daily),
						string(alerts.NrqlSignalSeasonalities.Weekly),
						string(alerts.NrqlSignalSeasonalities.None),
					},
					true,
				),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If a value is not provided and the condition uses the default value, don't show a diff. Also case insensitive.
					return (strings.EqualFold(old, string(alerts.NrqlSignalSeasonalities.NewRelicCalculation)) && new == "") || strings.EqualFold(old, new)
				},
			},
			// Outlier ONLY
			"outlier_configuration": {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				Description: "BETA PREVIEW: the `outlier_configuration` field is in limited release and only enabled for preview on a per-account basis. - Defines parameters controlling outlier detection for an `outlier` NRQL condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dbscan": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "BETA PREVIEW: the `dbscan` field is in limited release and only enabled for preview on a per-account basis. - Container for DBSCAN settings used to cluster data points and classify noise as outliers. Requires `epsilon` and `minimum_points`; optional `evaluation_group_facet` partitions data before analysis.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"epsilon": {
										Type:         schema.TypeFloat,
										Required:     true,
										Description:  "BETA PREVIEW: the `epsilon` field is in limited release and only enabled for preview on a per-account basis. - Radius (distance threshold) for DBSCAN in the units of the query result. Smaller values tighten clusters; larger values broaden them. Must be > 0.",
										ValidateFunc: validation.FloatAtLeast(0.0000001),
									},
									"minimum_points": {
										Type:         schema.TypeInt,
										Required:     true,
										Description:  "BETA PREVIEW: the `minimum_points` field is in limited release and only enabled for preview on a per-account basis. - Minimum number of neighboring points needed to form a cluster. Must be >= 1.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"evaluation_group_facet": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "BETA PREVIEW: the `evaluation_group_facet` field is in limited release and only enabled for preview on a per-account basis. - Optional NRQL facet attribute used to segment data into groups (e.g. `host`, `region`) before running outlier detection. Omit to evaluate all results together.",
									},
								},
							},
						},
					},
				},
			},
			"target_entity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BETA PREVIEW: the `target_entity` field is in limited release and only enabled for preview on a per-account basis. - The GUID of the entity explicitly targeted by the condition. Issues triggered by this condition will affect the health status of this entity instead of having the affected entity detected automatically",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func resourceNewRelicNrqlAlertConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)
	policyID := strconv.Itoa(d.Get("policy_id").(int))

	conditionInput, err := expandNrqlAlertConditionCreateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic NRQL alert condition %s via NerdGraph API", conditionInput.Name)

	var condition *alerts.NrqlAlertCondition

	switch d.Get("type").(string) {
	case "baseline":
		condition, err = client.Alerts.CreateNrqlConditionBaselineMutationWithContext(ctx, accountID, policyID, *conditionInput)
	case "outlier":
		condition, err = client.Alerts.CreateNrqlConditionOutlierMutationWithContext(ctx, accountID, policyID, *conditionInput)
	case "static":
		condition, err = client.Alerts.CreateNrqlConditionStaticMutationWithContext(ctx, accountID, policyID, *conditionInput)
	}

	var diags diag.Diagnostics

	if graphQLError, ok := err.(*alerts.GraphQLErrorResponse); ok {
		for _, e := range graphQLError.Errors {
			var message string = e.Message
			var errorClass string = e.Extensions.ErrorClass
			var validationErrors = e.Extensions.ValidationErrors

			if len(validationErrors) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  message + ": " + errorClass,
				})
			} else {
				for _, validationError := range validationErrors {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  message + ": " + errorClass,
						Detail:   validationError.Name + ": " + validationError.Reason,
					})
				}
			}
		}

		return diags
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if condition == nil {
		return diag.Errorf("error creating nrql alert condition: response was nil")
	}

	conditionID, err := strconv.Atoi(condition.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		condition, err = client.Alerts.GetNrqlConditionQueryWithContext(ctx, accountID, strconv.Itoa(conditionID))
		if err != nil {
			if _, ok := err.(*errors.NotFound); ok {
				return resource.RetryableError(fmt.Errorf("nrql condition was not created"))
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	d.SetId(serializeIDs([]int{d.Get("policy_id").(int), conditionID})) // set to correct ID

	return diag.FromErr(flattenNrqlAlertCondition(accountID, condition, d))
}

func resourceNewRelicNrqlAlertConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic NRQL alert condition %s", d.Id())

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := ids[0]
	conditionID := ids[1]

	_, err = client.Alerts.QueryPolicyWithContext(ctx, accountID, strconv.Itoa(policyID))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	nrqlCondition, err := client.Alerts.GetNrqlConditionQueryWithContext(ctx, accountID, strconv.Itoa(conditionID))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenNrqlAlertCondition(accountID, nrqlCondition, d))
}

func resourceNewRelicNrqlAlertConditionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	conditionID := strconv.Itoa(ids[1])

	conditionInput, err := expandNrqlAlertConditionUpdateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	switch d.Get("type").(string) {
	case "baseline":
		_, err = client.Alerts.UpdateNrqlConditionBaselineMutationWithContext(ctx, accountID, conditionID, *conditionInput)
	case "outlier":
		_, err = client.Alerts.UpdateNrqlConditionOutlierMutationWithContext(ctx, accountID, conditionID, *conditionInput)
	case "static":
		_, err = client.Alerts.UpdateNrqlConditionStaticMutationWithContext(ctx, accountID, conditionID, *conditionInput)
	}

	var diags diag.Diagnostics

	if graphQLError, ok := err.(*alerts.GraphQLErrorResponse); ok {
		for _, e := range graphQLError.Errors {
			var message string = e.Message
			var errorClass string = e.Extensions.ErrorClass
			var validationErrors = e.Extensions.ValidationErrors

			if len(validationErrors) == 0 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  message + ": " + errorClass,
				})
			} else {
				for _, validationError := range validationErrors {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  message + ": " + errorClass,
						Detail:   validationError.Name + ": " + validationError.Reason,
					})
				}
			}
		}

		return diags
	}

	return resourceNewRelicNrqlAlertConditionRead(ctx, d, meta)
}

func resourceNewRelicNrqlAlertConditionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	conditionID := strconv.Itoa(ids[1])

	log.Printf("[INFO] Deleting New Relic NRQL alert condition %v", conditionID)

	_, err = client.Alerts.DeleteNrqlConditionMutationWithContext(ctx, accountID, conditionID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
