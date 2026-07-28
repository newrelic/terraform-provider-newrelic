package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// errorIfUnsupportedRegion returns an error when the provider is configured for
// a region that does not provision the legacy REST v2 Alert Channels API
// (`newrelic_alert_channel` and `newrelic_alert_policy_channel`). Currently
// that covers the JP and GOV (FedRAMP) regions - neither provisions the
// deprecated `alerts_channels.json` endpoint.
//
// The returned error is safe to bubble up directly from CustomizeDiff.
func errorIfUnsupportedRegion(meta interface{}, resourceName string) error {
	providerConfig, ok := meta.(*ProviderConfig)
	if !ok || providerConfig == nil {
		return nil
	}
	r := providerConfig.Region
	if !strings.EqualFold(r, "JP") &&
		!strings.EqualFold(r, "GOV") &&
		!strings.EqualFold(r, "FEDRAMP") {
		return nil
	}

	return fmt.Errorf(
		"%s cannot be used in the %s region: it is backed by the legacy REST v2 "+
			"Alert Channels API, which was globally deprecated in 2024 and is not "+
			"provisioned on this region. Migrate to the NerdGraph-based "+
			"`newrelic_notification_destination` and `newrelic_notification_channel` "+
			"resources (and `newrelic_workflow` in place of `newrelic_alert_policy_channel`); "+
			"see the getting-started guide at "+
			"https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started#add-a-notification-channel "+
			"for a walkthrough",
		resourceName, strings.ToUpper(r),
	)
}

// blockUnsupportedRegionDiff returns a CustomizeDiffFunc that fails a plan when
// the provider is configured for a region that does not support the legacy REST
// v2 Alert Channels API. Wire this into a `schema.Resource` via `CustomizeDiff`
// so the error surfaces at `terraform plan` time.
func blockUnsupportedRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return func(_ context.Context, _ *schema.ResourceDiff, meta interface{}) error {
		return errorIfUnsupportedRegion(meta, resourceName)
	}
}

// blockJPRegionDiff is kept as a thin wrapper around blockUnsupportedRegionDiff
// for backward compatibility with any callers that reference it by name.
func blockJPRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return blockUnsupportedRegionDiff(resourceName)
}
