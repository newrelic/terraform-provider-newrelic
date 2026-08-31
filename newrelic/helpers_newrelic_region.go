package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// unsupportedRegionError constructs a formal plan-time error directing the user
// to the resource's documentation page for deprecation details and migration guidance.
func unsupportedRegionError(resourceName, region string) error {
	docSlug := strings.TrimPrefix(resourceName, "newrelic_")
	return fmt.Errorf(
		"%s is not supported in the %s region. The API underlying this resource "+
			"has been deprecated and is not available within this region. "+
			"For information on the deprecation status and recommended migration paths, "+
			"refer to the resource documentation: "+
			"https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/%s",
		resourceName, strings.ToUpper(region), docSlug,
	)
}

// blockRegionsDiff returns a CustomizeDiffFunc that fails a plan when the provider
// is configured for any of the specified regions. Wire this into a schema.Resource
// via the CustomizeDiff field so the error surfaces at terraform plan time, before
// any state is written or API call is made.
//
// Note: "FEDRAMP" is normalised to "GOV" in providerConfigure, so passing "GOV"
// here covers both the GOV and FEDRAMP region aliases.
func blockRegionsDiff(resourceName string, regions ...string) schema.CustomizeDiffFunc {
	return func(_ context.Context, _ *schema.ResourceDiff, meta interface{}) error {
		providerConfig, ok := meta.(*ProviderConfig)
		if !ok || providerConfig == nil {
			return nil
		}
		r := providerConfig.Region
		for _, blocked := range regions {
			if strings.EqualFold(r, blocked) {
				return unsupportedRegionError(resourceName, r)
			}
		}
		return nil
	}
}

// blockUnsupportedRegionDiff blocks a resource at plan time for both JP and GOV.
// Use this for resources confirmed unavailable on both regions (e.g. the legacy
// REST v2 Alert Channels family).
func blockUnsupportedRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return blockRegionsDiff(resourceName, "JP", "GOV")
}

// blockGOVRegionDiff blocks a resource at plan time for the GOV region only.
// Use this for resources confirmed unavailable on FedRAMP but still working on JP.
func blockGOVRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return blockRegionsDiff(resourceName, "GOV")
}

// blockJPRegionDiff is retained as a backward-compatible alias for
// blockUnsupportedRegionDiff. New resources should use blockRegionsDiff directly
// with the explicit list of regions to block.
func blockJPRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return blockUnsupportedRegionDiff(resourceName)
}
