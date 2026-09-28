package newrelic

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

func expandChangeTrackingDeployment(d *schema.ResourceData) (changetracking.ChangeTrackingDataHandlingRules, changetracking.ChangeTrackingDeploymentInput, error) {
	dataHandlingRules := expandChangeTrackingDataHandlingRules(d)
	deployment := expandChangeTrackingDeploymentInput(d)
	return dataHandlingRules, deployment, nil
}

func expandChangeTrackingDataHandlingRules(d *schema.ResourceData) changetracking.ChangeTrackingDataHandlingRules {
	rules := changetracking.ChangeTrackingDataHandlingRules{}

	if v, ok := d.GetOk("validation_flags"); ok {
		flags := v.([]interface{})
		validationFlags := make([]changetracking.ChangeTrackingValidationFlag, len(flags))
		for i, f := range flags {
			validationFlags[i] = changetracking.ChangeTrackingValidationFlag(f.(string))
		}
		rules.ValidationFlags = validationFlags
	}

	return rules
}

func expandChangeTrackingDeploymentInput(d *schema.ResourceData) changetracking.ChangeTrackingDeploymentInput {
	deployment := changetracking.ChangeTrackingDeploymentInput{
		Version:    d.Get("version").(string),
		EntityGUID: common.EntityGUID(d.Get("entity_guid").(string)),
	}

	if v, ok := d.GetOk("changelog"); ok {
		deployment.Changelog = v.(string)
	}

	if v, ok := d.GetOk("commit"); ok {
		deployment.Commit = v.(string)
	}

	if v, ok := d.GetOk("deep_link"); ok {
		deployment.DeepLink = v.(string)
	}

	if v, ok := d.GetOk("deployment_type"); ok {
		deployment.DeploymentType = changetracking.ChangeTrackingDeploymentType(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		deployment.Description = v.(string)
	}

	if v, ok := d.GetOk("group_id"); ok {
		deployment.GroupId = v.(string)
	}

	if v, ok := d.GetOk("timestamp"); ok {
		deployment.Timestamp = nrtime.EpochMilliseconds(time.UnixMilli(int64(v.(int))))
	}

	if v, ok := d.GetOk("user"); ok {
		deployment.User = v.(string)
	}

	return deployment
}

// No flatten functions: this API is write-only.
