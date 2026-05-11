package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"


	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/datamanagement"
)

const (
	cardinalityLimitName        = "Dimensional Metric per-metric cardinality ingested per day"
	cardinalityModeDefault      = "DEFAULT"
	cardinalityModePerMetric    = "PER_METRIC"
)

func resourceNewRelicAccountCardinalityLimit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAccountCardinalityLimitCreate,
		ReadContext:   resourceNewRelicAccountCardinalityLimitRead,
		UpdateContext: resourceNewRelicAccountCardinalityLimitCreate,
		DeleteContext: resourceNewRelicAccountCardinalityLimitDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceNewRelicAccountCardinalityLimitDiff,
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{cardinalityModeDefault, cardinalityModePerMetric}, false),
				Description:  "The override mode. Use 'DEFAULT' to set the account-wide default limit for all metrics, or 'PER_METRIC' to override the limit for a single metric (requires metric_name).",
			},
			"metric_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the metric to override. Required when mode is 'PER_METRIC'. Must not be set when mode is 'DEFAULT'.",
			},
			"cardinality_limit": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The cardinality limit value (maximum unique dimension-value combinations allowed per day).",
			},
		},
	}
}

func resourceNewRelicAccountCardinalityLimitDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	mode := d.Get("mode").(string)
	metricName := d.Get("metric_name").(string)

	switch mode {
	case cardinalityModeDefault:
		if metricName != "" {
			return fmt.Errorf("metric_name must not be set when mode is %q", cardinalityModeDefault)
		}
	case cardinalityModePerMetric:
		if metricName == "" {
			return fmt.Errorf("metric_name is required when mode is %q", cardinalityModePerMetric)
		}
	}
	return nil
}

func resourceNewRelicAccountCardinalityLimitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := providerConfig.AccountID

	input := buildCardinalityLimitInput(accountID, d)

	log.Printf("[INFO] Creating New Relic account cardinality limit for account %d, mode %q, metric %q", accountID, d.Get("mode").(string), input.Qualifier)

	_, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildCardinalityLimitID(accountID, input.Qualifier))

	return resourceNewRelicAccountCardinalityLimitRead(ctx, d, meta)
}

func resourceNewRelicAccountCardinalityLimitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID, metricName, err := parseCardinalityLimitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading New Relic account cardinality limit for account %d, metric %q", accountID, metricName)

	if metricName == "" {
		if err := d.Set("mode", cardinalityModeDefault); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("metric_name", ""); err != nil {
			return diag.FromErr(err)
		}

		// DEFAULT mode: read the current value from the dataManagement limits query.
		limits, err := client.DataManagement.GetLimitsWithContext(ctx, accountID)
		if err != nil {
			return diag.FromErr(err)
		}
		if limits != nil {
			for _, l := range *limits {
				if l.Name == cardinalityLimitName {
					if err := d.Set("cardinality_limit", l.Value); err != nil {
						return diag.FromErr(err)
					}
					break
				}
			}
		}
	} else {
		if err := d.Set("mode", cardinalityModePerMetric); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("metric_name", metricName); err != nil {
			return diag.FromErr(err)
		}

		// PER_METRIC mode: no reliable synchronous read path exists for per-metric
		// override values. The dataManagement limits query does not expose qualifier,
		// and newrelic.resourceConsumption.limitValue in NRDB lags behind the mutation
		// API by the metric ingestion interval. cardinality_limit is kept from state.
	}

	return nil
}

func resourceNewRelicAccountCardinalityLimitDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// The NerdGraph API does not expose a delete mutation for account cardinality
	// limits. Destroying this resource removes it from Terraform state only; the
	// override remains in New Relic until changed externally or via a new apply.
	log.Printf("[INFO] Removing New Relic account cardinality limit %s from state (no API delete available)", d.Id())
	return nil
}

func buildCardinalityLimitInput(accountID int, d *schema.ResourceData) datamanagement.DataManagementAccountLimitInput {
	metricName := d.Get("metric_name").(string)
	limit := d.Get("cardinality_limit").(int)

	var reason string
	if metricName == "" {
		reason = fmt.Sprintf("Default cardinality limit for account %d set to %d via Terraform", accountID, limit)
	} else {
		reason = fmt.Sprintf("Cardinality limit for metric '%s' in account %d set to %d via Terraform", metricName, accountID, limit)
	}

	return datamanagement.DataManagementAccountLimitInput{
		Limit: datamanagement.DataManagementLimitLookupInput{
			Name: cardinalityLimitName,
		},
		OverrideValue:  limit,
		OverrideReason: reason,
		Qualifier:      metricName,
	}
}

func buildCardinalityLimitID(accountID int, metricName string) string {
	return fmt.Sprintf("%d:%s", accountID, metricName)
}

func parseCardinalityLimitID(id string) (int, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid cardinality limit ID %q: expected format \"<accountId>:<metricName>\"", id)
	}
	accountID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid account ID in cardinality limit ID %q: %w", id, err)
	}
	return accountID, parts[1], nil
}
