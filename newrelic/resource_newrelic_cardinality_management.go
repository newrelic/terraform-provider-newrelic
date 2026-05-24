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
	cardinalityLimitName     = "Dimensional Metric per-metric cardinality ingested per day"
	cardinalityModeDefault   = "DEFAULT"
	cardinalityModePerMetric = "PER_METRIC"

	// cardinalityLimitPlatformDefault is the out-of-the-box limit New Relic applies
	// to every account before any overrides are configured.
	cardinalityLimitPlatformDefault = 100000

	// cardinalityUILagNotice is included in warnings after write operations to let
	// users know that cardinality limit changes are not always instant in the UI.
	// Limits take effect in the enforcement layer right away, but the New Relic UI
	// and NRDB consumption events may lag by a few minutes, especially if the
	// affected metrics are not actively ingesting data at the time of the change.
	cardinalityUILagNotice = "Changes take effect in the enforcement layer right away, but may take a few minutes " +
		"to appear in the New Relic UI — particularly if the affected metrics have not sent data recently."
)

func resourceNewRelicCardinalityManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCardinalityManagementCreate,
		ReadContext:   resourceNewRelicCardinalityManagementRead,
		UpdateContext: resourceNewRelicCardinalityManagementCreate,
		DeleteContext: resourceNewRelicCardinalityManagementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceNewRelicCardinalityManagementDiff,
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{cardinalityModeDefault, cardinalityModePerMetric}, false),
				Description: "The override mode. Use `DEFAULT` to set a single account-wide limit that applies to " +
					"all metrics, or `PER_METRIC` to set individual limits for one or more named metrics.",
			},
			"cardinality_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The account-wide cardinality limit — the maximum number of unique " +
					"dimension-value combinations allowed per metric per day. " +
					"Required when `mode` is `DEFAULT`; must not be set when `mode` is `PER_METRIC`.",
			},
			"metric": {
				Type:     schema.TypeList,
				Optional: true,
				Description: "One or more metrics to set individual cardinality limits for. " +
					"Required when `mode` is `PER_METRIC`; must not be set when `mode` is `DEFAULT`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The full name of the metric (e.g. `http.server.duration`).",
						},
						"cardinality_limit": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The maximum number of unique dimension-value combinations allowed per day for this metric.",
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicCardinalityManagementDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	mode := d.Get("mode").(string)
	metrics := d.Get("metric").([]interface{})
	topLevelLimit := d.Get("cardinality_limit").(int)

	switch mode {
	case cardinalityModeDefault:
		if len(metrics) > 0 {
			return fmt.Errorf("metric blocks must not be set when mode is %q — use cardinality_limit at the top level instead", cardinalityModeDefault)
		}
		if topLevelLimit == 0 {
			return fmt.Errorf("cardinality_limit is required when mode is %q", cardinalityModeDefault)
		}
	case cardinalityModePerMetric:
		if topLevelLimit != 0 {
			return fmt.Errorf("cardinality_limit must not be set at the top level when mode is %q — set the limit inside each metric block instead", cardinalityModePerMetric)
		}
		if len(metrics) == 0 {
			return fmt.Errorf("at least one metric block is required when mode is %q", cardinalityModePerMetric)
		}
	}
	return nil
}

func resourceNewRelicCardinalityManagementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := providerConfig.AccountID
	mode := d.Get("mode").(string)

	if mode == cardinalityModeDefault {
		limit := d.Get("cardinality_limit").(int)
		log.Printf("[INFO] Setting account-wide cardinality limit for account %d to %d", accountID, limit)

		input := datamanagement.DataManagementAccountLimitInput{
			Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
			OverrideValue:  limit,
			OverrideReason: fmt.Sprintf("Account-wide cardinality limit for account %d set to %d via Terraform", accountID, limit),
		}
		if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(buildCardinalityLimitID(accountID, cardinalityModeDefault))
		return resourceNewRelicCardinalityManagementRead(ctx, d, meta)
	}

	// PER_METRIC: apply one override per metric block.
	metrics := d.Get("metric").([]interface{})
	log.Printf("[INFO] Setting per-metric cardinality limits for account %d (%d metric(s))", accountID, len(metrics))

	for _, raw := range metrics {
		m := raw.(map[string]interface{})
		name := m["name"].(string)
		limit := m["cardinality_limit"].(int)

		input := datamanagement.DataManagementAccountLimitInput{
			Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
			OverrideValue:  limit,
			OverrideReason: fmt.Sprintf("Cardinality limit for metric %q in account %d set to %d via Terraform", name, accountID, limit),
			Qualifier:      name,
		}
		if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input); err != nil {
			return diag.Errorf("failed to set cardinality limit for metric %q: %s", name, err)
		}
	}

	d.SetId(buildCardinalityLimitID(accountID, cardinalityModePerMetric))

	// No Read call here — per-metric override values cannot be read back from the
	// API, so there is nothing to reconcile. Read will surface its own advisory
	// warning during the next plan or refresh.
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  "Metric cardinality limit override(s) applied",
			Detail: fmt.Sprintf(
				"Cardinality limit overrides have been set for %d metric(s) in account %d.\n\n%s",
				len(metrics), accountID, cardinalityUILagNotice,
			),
		},
	}
}

func resourceNewRelicCardinalityManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID, mode, err := parseCardinalityLimitID(d.Id())
	if err != nil {
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

	// PER_METRIC: override values are write-only in this mode — state is always
	// preserved from the last apply rather than reconciled from a live read.
	log.Printf("[INFO] Skipping live read for PER_METRIC cardinality limits in account %d — state reflects last apply", accountID)

	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  "Per-metric cardinality limit values reflect the last Terraform apply",
			Detail: "In PER_METRIC mode, cardinality limit values in state are indicative of the last configuration applied via Terraform — this is the expected behaviour for this mode.\n\n" +
				"If any of these limits have been adjusted outside of Terraform, run terraform apply to re-apply the desired values.",
		},
	}
}

func resourceNewRelicCardinalityManagementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID, mode, err := parseCardinalityLimitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if mode == cardinalityModeDefault {
		// The New Relic API has no delete operation for cardinality limit overrides,
		// so destroy resets the account-wide limit back to the platform default.
		log.Printf("[INFO] Resetting account-wide cardinality limit for account %d to platform default (%d)", accountID, cardinalityLimitPlatformDefault)

		resetInput := datamanagement.DataManagementAccountLimitInput{
			Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
			OverrideValue:  cardinalityLimitPlatformDefault,
			OverrideReason: fmt.Sprintf("Account-wide cardinality limit for account %d reset to platform default (%d) via Terraform destroy", accountID, cardinalityLimitPlatformDefault),
		}
		if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, resetInput); err != nil {
			return diag.FromErr(err)
		}

		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Account-wide cardinality limit reset to platform default",
				Detail: fmt.Sprintf(
					"The account-wide cardinality limit for account %d has been reset to the New Relic platform default of %d.\n\n"+
						"This value applies to all metrics in the account that do not have a per-metric override.\n\n%s",
					accountID, cardinalityLimitPlatformDefault, cardinalityUILagNotice,
				),
			},
		}
	}

	// PER_METRIC: reset each managed metric to the platform default.
	// The API has no delete operation, so this is the closest equivalent to removal.
	metrics := d.Get("metric").([]interface{})
	log.Printf("[INFO] Resetting per-metric cardinality limits to platform default for account %d (%d metric(s))", accountID, len(metrics))

	for _, raw := range metrics {
		m := raw.(map[string]interface{})
		name := m["name"].(string)

		resetInput := datamanagement.DataManagementAccountLimitInput{
			Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
			OverrideValue:  cardinalityLimitPlatformDefault,
			OverrideReason: fmt.Sprintf("Cardinality limit override for metric %q in account %d removed via Terraform; reset to platform default (%d)", name, accountID, cardinalityLimitPlatformDefault),
			Qualifier:      name,
		}
		if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, resetInput); err != nil {
			return diag.Errorf("failed to reset cardinality limit for metric %q: %s", name, err)
		}
	}

	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  "Per-metric cardinality limit overrides removed",
			Detail: fmt.Sprintf(
				"Cardinality limit overrides for %d metric(s) in account %d have been reset to the platform default of %d.\n\n%s",
				len(metrics), accountID, cardinalityLimitPlatformDefault, cardinalityUILagNotice,
			),
		},
	}
}

func buildCardinalityLimitID(accountID int, mode string) string {
	return fmt.Sprintf("%d:%s", accountID, mode)
}

func parseCardinalityLimitID(id string) (int, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("unexpected cardinality management ID format %q — expected \"<accountId>:<mode>\"", id)
	}
	accountID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid account ID in cardinality management ID %q: %w", id, err)
	}
	return accountID, parts[1], nil
}
