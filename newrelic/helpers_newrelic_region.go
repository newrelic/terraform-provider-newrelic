package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// unsupportedRegionError returns a customer-friendly plan-time error for a resource
// that is not available in the given region. The message points to the resource's
// own documentation page where deprecation status and alternatives are documented.
func unsupportedRegionError(resourceName, region string) error {
	// Derive the docs URL slug from the resource name (strip "newrelic_" prefix).
	docSlug := strings.TrimPrefix(resourceName, "newrelic_")
	return fmt.Errorf(
		"%s is not available in the %s region. The underlying API this resource "+
			"depends on has been deprecated and is not provisioned in this region. "+
			"Please refer to the resource documentation for deprecation details "+
			"and guidance on supported alternatives: "+
			"https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/%s",
		resourceName, strings.ToUpper(region), docSlug,
	)
}

// errorIfUnsupportedRegion blocks resources that are unavailable in both the JP
// and GOV regions. Currently that is the legacy REST v2 Alert Channels family
// (newrelic_alert_channel, newrelic_alert_policy_channel) — confirmed absent on
// both regions. "FEDRAMP" is normalized to "GOV" in providerConfigure, so only
// JP and GOV are checked here.
func errorIfUnsupportedRegion(meta interface{}, resourceName string) error {
	providerConfig, ok := meta.(*ProviderConfig)
	if !ok || providerConfig == nil {
		return nil
	}
	r := providerConfig.Region
	if !strings.EqualFold(r, "JP") && !strings.EqualFold(r, "GOV") {
		return nil
	}
	return unsupportedRegionError(resourceName, r)
}

// errorIfGOVRegion blocks resources that are unavailable specifically in the GOV
// (FedRAMP) region but still work in JP and other regions. Use this for resources
// whose backing service is not provisioned in FedRAMP cells.
func errorIfGOVRegion(meta interface{}, resourceName string) error {
	providerConfig, ok := meta.(*ProviderConfig)
	if !ok || providerConfig == nil {
		return nil
	}
	r := providerConfig.Region
	if !strings.EqualFold(r, "GOV") {
		return nil
	}
	return unsupportedRegionError(resourceName, r)
}

// blockUnsupportedRegionDiff blocks a resource at plan time for both JP and GOV.
func blockUnsupportedRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return func(_ context.Context, _ *schema.ResourceDiff, meta interface{}) error {
		return errorIfUnsupportedRegion(meta, resourceName)
	}
}

// blockGOVRegionDiff blocks a resource at plan time for the GOV region only.
// Use this for resources confirmed unavailable on FedRAMP but working on JP.
func blockGOVRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return func(_ context.Context, _ *schema.ResourceDiff, meta interface{}) error {
		return errorIfGOVRegion(meta, resourceName)
	}
}

// blockJPRegionDiff is kept as a backward-compatible wrapper.
func blockJPRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return blockUnsupportedRegionDiff(resourceName)
}
