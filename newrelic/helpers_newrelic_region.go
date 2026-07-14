package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// errorIfJPRegion returns an error when the provider is configured for the JP
// region. It is used by resources whose upstream API is not (and will not be)
// available in JP - currently the legacy REST v2 Alert Channels endpoints
// (`newrelic_alert_channel` and `newrelic_alert_policy_channel`).
//
// The legacy `alerts_channels.json` REST endpoint has been deprecated globally
// since 2024 (see https://docs.newrelic.com/docs/alerts/scale-automate/rest-api/rest-api-calls-alerts/)
// and the Alerts platform team has confirmed it will not be provisioned on JP
// (see internal Slack thread in #help-alerts, Jul 2026). Attempting to Create
// these resources against JP produces an opaque "resource not found" error
// from the API - this helper surfaces a clearer, actionable message at plan
// time (via CustomizeDiff) instead.
//
// The returned error is safe to bubble up directly from CustomizeDiff.
func errorIfJPRegion(meta interface{}, resourceName string) error {
	providerConfig, ok := meta.(*ProviderConfig)
	if !ok || providerConfig == nil {
		return nil
	}
	if !strings.EqualFold(providerConfig.Region, "JP") {
		return nil
	}

	return fmt.Errorf(
		"%s cannot be used in the JP region: it is backed by the legacy REST v2 "+
			"Alert Channels API, which was globally deprecated in 2024 and is not "+
			"provisioned on JP. Migrate to the NerdGraph-based "+
			"`newrelic_notification_destination` and `newrelic_notification_channel` "+
			"resources (and `newrelic_workflow` in place of `newrelic_alert_policy_channel`); "+
			"see the getting-started guide at "+
			"https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/getting_started#add-a-notification-channel "+
			"for a walkthrough",
		resourceName,
	)
}

// blockJPRegionDiff returns a CustomizeDiffFunc that fails a plan when the
// provider is configured for the JP region. Wire this into a `schema.Resource`
// via the `CustomizeDiff` field so the error surfaces at `terraform plan`
// time - well before any state is written or any API call is made.
func blockJPRegionDiff(resourceName string) schema.CustomizeDiffFunc {
	return func(_ context.Context, _ *schema.ResourceDiff, meta interface{}) error {
		return errorIfJPRegion(meta, resourceName)
	}
}
