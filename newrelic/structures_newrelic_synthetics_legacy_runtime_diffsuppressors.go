// naming inspired by newrelic/structures_newrelic_synthetics_all_monitors_validation_helpers.go in PR #2727
// this is not intended to be in different file, creating this in a different file for now to avoid merge conflicts
// after #2727 is merged

// !!!!! @pranav-new-relic DELETE THIS once this is no longer needed

package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func syntheticMonitorsLegacyRuntimeArgumentsDiffSuppressor(k, oldValue, newValue string, d *schema.ResourceData) bool {
	rawConfiguration := d.GetRawConfig()
	isRuntimeTypeNotSpecifiedInConfiguration := rawConfiguration.GetAttr(k).IsNull()
	isRuntimeTypeNullValue := newValue == ""

	if isRuntimeTypeNotSpecifiedInConfiguration && isRuntimeTypeNullValue {
		return true
	}

	return false
}
