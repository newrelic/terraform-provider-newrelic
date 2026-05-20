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
	cardinalityLimitName            = "Dimensional Metric per-metric cardinality ingested per day"
	cardinalityModeDefault          = "DEFAULT"
	cardinalityModePerMetric        = "PER_METRIC"
	cardinalityLimitPlatformDefault = 100000

	// cardinalityUILagNotice is appended to warnings on write operations. The
	// updated limit takes effect as metric data is received and processed by
	// the platform, so visibility in the New Relic UI follows the metric
	// ingestion cycle rather than being instantaneous.
	cardinalityUILagNotice = "The updated limit takes effect as metric data is received and processed by the platform, and will be visible in the New Relic UI once the next metric ingestion cycle completes."
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

	mode := d.Get("mode").(string)
	input := buildCardinalityLimitInput(accountID, d)

	log.Printf("[INFO] Creating New Relic account cardinality limit for account %d, mode %q, metric %q", accountID, mode, input.Qualifier)

	_, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildCardinalityLimitID(accountID, input.Qualifier))

	var diags diag.Diagnostics

	switch mode {
	case cardinalityModeDefault:
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Account-wide cardinality limit updated",
			Detail:   fmt.Sprintf("The account-wide cardinality limit has been set to %d. %s", input.OverrideValue, cardinalityUILagNotice),
		})
	case cardinalityModePerMetric:
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Per-metric cardinality limit applied",
			Detail: fmt.Sprintf(
				"The cardinality limit for metric %q has been set to %d. %s",
				input.Qualifier, input.OverrideValue, cardinalityUILagNotice,
			),
		})
	}

	return append(diags, resourceNewRelicAccountCardinalityLimitRead(ctx, d, meta)...)
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

	if err := d.Set("mode", cardinalityModePerMetric); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metric_name", metricName); err != nil {
		return diag.FromErr(err)
	}

	// PER_METRIC mode: the enforced limit is tied to metric ingestion on the
	// platform, so cardinality_limit is preserved from state rather than read
	// back on each plan.
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  "Per-metric cardinality limit reflects last applied value",
			Detail: fmt.Sprintf(
				"The enforced cardinality limit for metric %q is tied to metric ingestion on the platform. "+
					"The value in state is the last limit applied by Terraform. "+
					"Run 'terraform apply' to re-enforce the desired limit at any time.",
				metricName,
			),
		},
	}
}

func resourceNewRelicAccountCardinalityLimitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID, metricName, err := parseCardinalityLimitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if metricName == "" {
		// DEFAULT mode: the NerdGraph API has no delete operation, so we reset the
		// account-wide default back to the New Relic platform default of 100,000.
		log.Printf("[INFO] Resetting account-wide default cardinality limit for account %d to platform default (%d)", accountID, cardinalityLimitPlatformDefault)

		resetInput := datamanagement.DataManagementAccountLimitInput{
			Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
			OverrideValue:  cardinalityLimitPlatformDefault,
			OverrideReason: fmt.Sprintf("Default cardinality limit for account %d reset to platform default (%d) via Terraform destroy", accountID, cardinalityLimitPlatformDefault),
			Qualifier:      "",
		}

		if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, resetInput); err != nil {
			return diag.FromErr(err)
		}

		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Account-wide default cardinality limit reset to platform default",
				Detail: fmt.Sprintf(
					"The account-wide default cardinality limit has been reset to the New Relic platform default of %d. %s",
					cardinalityLimitPlatformDefault, cardinalityUILagNotice,
				),
			},
		}
	}

	// PER_METRIC mode: the NerdGraph API has no delete operation, so we reset this
	// metric's override to the current account-wide default, effectively removing
	// the per-metric exception.
	log.Printf("[INFO] Resetting per-metric cardinality limit for metric %q in account %d to the current account-wide default", metricName, accountID)

	defaultLimit, err := fetchDefaultCardinalityLimitValue(ctx, providerConfig, accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	resetInput := datamanagement.DataManagementAccountLimitInput{
		Limit:          datamanagement.DataManagementLimitLookupInput{Name: cardinalityLimitName},
		OverrideValue:  defaultLimit,
		OverrideReason: fmt.Sprintf("Cardinality limit for metric '%s' in account %d reset to account-wide default (%d) via Terraform destroy", metricName, accountID, defaultLimit),
		Qualifier:      metricName,
	}

	if _, err := client.DataManagement.DataManagementCreateAccountLimitWithContext(ctx, accountID, resetInput); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  "Per-metric cardinality limit reset to account-wide default",
			Detail: fmt.Sprintf(
				"The cardinality limit override for metric %q has been reset to the current account-wide default value of %d. %s",
				metricName, defaultLimit, cardinalityUILagNotice,
			),
		},
	}
}

// fetchDefaultCardinalityLimitValue retrieves the current account-wide default
// cardinality limit from the dataManagement limits API. Falls back to the
// platform default of 100,000 if the limit entry is not found.
func fetchDefaultCardinalityLimitValue(ctx context.Context, providerConfig *ProviderConfig, accountID int) (int, error) {
	limits, err := providerConfig.NewClient.DataManagement.GetLimitsWithContext(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch account-wide default cardinality limit: %w", err)
	}
	if limits != nil {
		for _, l := range *limits {
			if l.Name == cardinalityLimitName {
				return l.Value, nil
			}
		}
	}
	return cardinalityLimitPlatformDefault, nil
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
