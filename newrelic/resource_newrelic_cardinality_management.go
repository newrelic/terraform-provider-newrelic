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
	// cardinalityLimitName is the internal name used by the New Relic platform
	// to identify the per-metric cardinality limit.
	cardinalityLimitName = "Dimensional Metric per-metric cardinality ingested per day"

	cardinalityModeDefault          = "DEFAULT"
	cardinalityModePerMetric        = "PER_METRIC"
	cardinalityLimitPlatformDefault = 100000
)

func resourceNewRelicCardinalityManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCardinalityManagementCreate,
		ReadContext:   resourceNewRelicCardinalityManagementRead,
		UpdateContext: resourceNewRelicCardinalityManagementUpdate,
		DeleteContext: resourceNewRelicCardinalityManagementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceNewRelicCardinalityManagementDiff,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID. Defaults to the account ID configured on the provider.",
			},
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{cardinalityModeDefault, cardinalityModePerMetric}, false),
				Description:  "The management mode. Use 'DEFAULT' to set the account-wide limit that applies to all metrics, or 'PER_METRIC' to set individual limits per metric name.",
			},
			"cardinality_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The account-wide cardinality limit. Required when mode is 'DEFAULT'; must not be set when mode is 'PER_METRIC'.",
			},
			"metric": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "One or more per-metric cardinality overrides. Required when mode is 'PER_METRIC'; must not be set when mode is 'DEFAULT'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the metric to override.",
						},
						"cardinality_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The cardinality limit for this metric.",
						},
					},
				},
				// Hash by metric name so that updating only the limit for an existing
				// metric is treated as an in-place change rather than remove + add.
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return schema.HashString(m["name"].(string))
				},
			},
		},
	}
}

func resourceNewRelicCardinalityManagementDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	mode := d.Get("mode").(string)
	limit := d.Get("cardinality_limit").(int)
	metrics := d.Get("metric").(*schema.Set)

	switch mode {
	case cardinalityModeDefault:
		if limit == 0 {
			return fmt.Errorf("cardinality_limit is required when mode is %q", cardinalityModeDefault)
		}
		if metrics.Len() > 0 {
			return fmt.Errorf("metric blocks must not be set when mode is %q", cardinalityModeDefault)
		}
	case cardinalityModePerMetric:
		if metrics.Len() == 0 {
			return fmt.Errorf("at least one metric block is required when mode is %q", cardinalityModePerMetric)
		}
		if limit != 0 {
			return fmt.Errorf("cardinality_limit must not be set when mode is %q; set it inside each metric block instead", cardinalityModePerMetric)
		}
	}
	return nil
}

func resourceNewRelicCardinalityManagementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	mode := d.Get("mode").(string)

	switch mode {
	case cardinalityModeDefault:
		limit := d.Get("cardinality_limit").(int)
		log.Printf("[INFO] Setting account-wide cardinality limit for account %d to %d", accountID, limit)

		input := buildDefaultCardinalityInput(accountID, limit)
		if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(buildCardinalityManagementID(accountID, mode))
		if err := d.Set("account_id", accountID); err != nil {
			return diag.FromErr(err)
		}

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Account-wide cardinality limit set",
			Detail:   fmt.Sprintf("The account-wide cardinality limit has been set to %d. Please allow a few minutes for this change to appear in the New Relic UI.", limit),
		}}

	case cardinalityModePerMetric:
		metrics := d.Get("metric").(*schema.Set)
		log.Printf("[INFO] Setting per-metric cardinality limits for account %d (%d metric(s))", accountID, metrics.Len())

		for _, raw := range metrics.List() {
			m := raw.(map[string]interface{})
			name := m["name"].(string)
			limit := m["cardinality_limit"].(int)

			input := buildPerMetricCardinalityInput(accountID, name, limit)
			if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
				return diag.FromErr(err)
			}
		}

		d.SetId(buildCardinalityManagementID(accountID, mode))
		if err := d.Set("account_id", accountID); err != nil {
			return diag.FromErr(err)
		}

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Per-metric cardinality limits applied",
			Detail:   fmt.Sprintf("Cardinality limits have been applied for %d metric(s). Please allow a few minutes for these changes to appear in the New Relic UI.", metrics.Len()),
		}}
	}

	return nil
}

func resourceNewRelicCardinalityManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID, mode, err := parseCardinalityManagementID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mode", mode); err != nil {
		return diag.FromErr(err)
	}

	if mode == cardinalityModeDefault {
		log.Printf("[INFO] Reading account-wide cardinality limit for account %d", accountID)

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
		return nil
	}

	// PER_METRIC: metric limits are tied to metric activity on the platform and
	// are not independently queryable, so Terraform preserves the last applied
	// values in state. A plan will not detect changes made outside of Terraform.
	log.Printf("[INFO] Skipping live read for PER_METRIC cardinality management on account %d (state preserved)", accountID)
	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "Metric cardinality limits reflect the last values applied by Terraform",
		Detail: "The metric limits shown here are the last values applied by Terraform. " +
			"Changes may take a few minutes to appear in the New Relic UI once metric data flows through.",
	}}
}

func resourceNewRelicCardinalityManagementUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := d.Get("account_id").(int)
	mode := d.Get("mode").(string)

	switch mode {
	case cardinalityModeDefault:
		limit := d.Get("cardinality_limit").(int)
		log.Printf("[INFO] Updating account-wide cardinality limit for account %d to %d", accountID, limit)

		input := buildDefaultCardinalityInput(accountID, limit)
		if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
			return diag.FromErr(err)
		}

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Account-wide cardinality limit updated",
			Detail:   fmt.Sprintf("The account-wide cardinality limit has been updated to %d. Please allow a few minutes for this change to appear in the New Relic UI.", limit),
		}}

	case cardinalityModePerMetric:
		if !d.HasChange("metric") {
			return nil
		}

		oldRaw, newRaw := d.GetChange("metric")
		oldMetrics := metricSetToMap(oldRaw.(*schema.Set))
		newMetrics := metricSetToMap(newRaw.(*schema.Set))

		// Reset limits for any metrics that were removed from the configuration.
		for name := range oldMetrics {
			if _, stillPresent := newMetrics[name]; !stillPresent {
				log.Printf("[INFO] Resetting removed metric %q in account %d to platform default", name, accountID)
				input := buildResetMetricCardinalityInput(accountID, name)
				if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// Apply the current set of metric limits.
		for name, limit := range newMetrics {
			log.Printf("[INFO] Setting cardinality limit for metric %q in account %d to %d", name, accountID, limit)
			input := buildPerMetricCardinalityInput(accountID, name, limit)
			if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
				return diag.FromErr(err)
			}
		}

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Per-metric cardinality limits updated",
			Detail:   "The metric cardinality limits have been updated. Please allow a few minutes for these changes to appear in the New Relic UI.",
		}}
	}

	return nil
}

func resourceNewRelicCardinalityManagementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID, mode, err := parseCardinalityManagementID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	switch mode {
	case cardinalityModeDefault:
		log.Printf("[INFO] Resetting account-wide cardinality limit for account %d to platform default (%d)", accountID, cardinalityLimitPlatformDefault)

		input := buildDefaultCardinalityInput(accountID, cardinalityLimitPlatformDefault)
		if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
			return diag.FromErr(err)
		}

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Account-wide cardinality limit reset to platform default",
			Detail:   fmt.Sprintf("The account-wide cardinality limit has been reset to the New Relic platform default of %d. Please allow a few minutes for this change to appear in the New Relic UI.", cardinalityLimitPlatformDefault),
		}}

	case cardinalityModePerMetric:
		metrics := d.Get("metric").(*schema.Set)
		log.Printf("[INFO] Resetting %d per-metric cardinality limit(s) for account %d to platform default (%d)", metrics.Len(), accountID, cardinalityLimitPlatformDefault)

		for _, raw := range metrics.List() {
			m := raw.(map[string]interface{})
			name := m["name"].(string)
			input := buildResetMetricCardinalityInput(accountID, name)
			if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
				return diag.FromErr(err)
			}
		}

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Per-metric cardinality limits reset to platform default",
			Detail:   fmt.Sprintf("All per-metric cardinality overrides have been reset to the New Relic platform default of %d. Please allow a few minutes for these changes to appear in the New Relic UI.", cardinalityLimitPlatformDefault),
		}}
	}

	return nil
}

// buildDefaultCardinalityInput constructs the mutation input for the account-wide default limit.
func buildDefaultCardinalityInput(accountID, limit int) datamanagement.DataManagementAccountLimitInput {
	return datamanagement.DataManagementAccountLimitInput{
		Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
		OverrideValue:  limit,
		OverrideReason: fmt.Sprintf("Account-wide cardinality limit for account %d set to %d via Terraform", accountID, limit),
		Qualifier:      "",
	}
}

// buildPerMetricCardinalityInput constructs the mutation input for a single metric override.
func buildPerMetricCardinalityInput(accountID int, metricName string, limit int) datamanagement.DataManagementAccountLimitInput {
	return datamanagement.DataManagementAccountLimitInput{
		Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
		OverrideValue:  limit,
		OverrideReason: fmt.Sprintf("Cardinality limit for metric %q in account %d set to %d via Terraform", metricName, accountID, limit),
		Qualifier:      metricName,
	}
}

// buildResetMetricCardinalityInput constructs the mutation input to reset a single metric
// back to the New Relic platform default of 100,000.
func buildResetMetricCardinalityInput(accountID int, metricName string) datamanagement.DataManagementAccountLimitInput {
	return datamanagement.DataManagementAccountLimitInput{
		Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
		OverrideValue:  cardinalityLimitPlatformDefault,
		OverrideReason: fmt.Sprintf("Cardinality limit for metric %q in account %d reset to platform default (%d) via Terraform", metricName, accountID, cardinalityLimitPlatformDefault),
		Qualifier:      metricName,
	}
}

// metricSetToMap converts the TypeSet of metric blocks into a map of metric name → limit,
// making it easy to compare old and new sets when computing updates.
func metricSetToMap(s *schema.Set) map[string]int {
	result := make(map[string]int, s.Len())
	for _, raw := range s.List() {
		m := raw.(map[string]interface{})
		result[m["name"].(string)] = m["cardinality_limit"].(int)
	}
	return result
}

func buildCardinalityManagementID(accountID int, mode string) string {
	return fmt.Sprintf("%d:%s", accountID, mode)
}

func parseCardinalityManagementID(id string) (int, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid cardinality management resource ID %q: expected format \"<accountId>:DEFAULT\" or \"<accountId>:PER_METRIC\"", id)
	}
	accountID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid account ID in cardinality management resource ID %q: %w", id, err)
	}
	mode := parts[1]
	if mode != cardinalityModeDefault && mode != cardinalityModePerMetric {
		return 0, "", fmt.Errorf("invalid mode %q in cardinality management resource ID %q: expected DEFAULT or PER_METRIC", mode, id)
	}
	return accountID, mode, nil
}
