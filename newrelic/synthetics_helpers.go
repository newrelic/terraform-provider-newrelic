package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

type SyntheticsMonitorBase struct {
	Name   string
	Period synthetics.SyntheticsMonitorPeriod
	Status synthetics.SyntheticsMonitorStatus
	Tags   []synthetics.SyntheticsTag
}

// Handles setting simple string attributes in the schema. If the attribute/key is
// invalid or the value is not a correct type, an error will be returned.
func setSyntheticsMonitorAttributes(d *schema.ResourceData, attributes map[string]string) error {
	for key := range attributes {
		err := d.Set(key, attributes[key])
		if err != nil {
			return err
		}
	}

	return nil
}

func buildSyntheticsBrokenLinksMonitorCreateInput(d *schema.ResourceData) *synthetics.SyntheticsCreateBrokenLinksMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateBrokenLinksMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
	}

	if attr, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandStringSlice(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(attr.([]interface{}))
	}
	if v, ok := d.GetOk("uri"); ok {
		input.Uri = v.(string)
	}
	return &input
}

func buildSyntheticsBrokenLinksMonitorUpdateInput(d *schema.ResourceData) *synthetics.SyntheticsUpdateBrokenLinksMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateBrokenLinksMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
	}

	if attr, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandStringSlice(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(attr.([]interface{}))
	}
	if v, ok := d.GetOk("uri"); ok {
		input.Uri = v.(string)
	}

	return &input
}

// Builds an array of typed diagnostic errors based on the GraphQL `response.errors` array.
func buildCreateSyntheticsMonitorResponseErrors(errors []synthetics.SyntheticsMonitorCreateError) diag.Diagnostics {
	var diagErrors diag.Diagnostics
	for _, err := range errors {
		diagErrors = append(diagErrors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
		})
	}
	return diagErrors
}

// Builds an array of typed diagnostic errors based on the GraphQL `response.errors` array.
func buildUpdateSyntheticsMonitorResponseErrors(errors []synthetics.SyntheticsMonitorUpdateError) diag.Diagnostics {
	var diagErrors diag.Diagnostics
	for _, err := range errors {
		diagErrors = append(diagErrors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
		})
	}
	return diagErrors
}

func expandSyntheticsMonitorBase(d *schema.ResourceData) SyntheticsMonitorBase {
	inputBase := SyntheticsMonitorBase{}

	name := d.Get("name")
	inputBase.Name = name.(string)
	status := d.Get("status")
	inputBase.Status = synthetics.SyntheticsMonitorStatus(status.(string))
	period := d.Get("period")
	inputBase.Period = synthetics.SyntheticsMonitorPeriod(period.(string))

	if tags, ok := d.GetOk("tag"); ok {
		inputBase.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}

	if tags, ok := d.GetOk("tag"); ok {
		inputBase.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}

	if tags, ok := d.GetOk("tag"); ok {
		inputBase.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}
	return inputBase
}

func expandSyntheticsTags(tags []interface{}) []synthetics.SyntheticsTag {
	out := make([]synthetics.SyntheticsTag, len(tags))
	for i, v := range tags {
		tag := v.(map[string]interface{})
		expanded := synthetics.SyntheticsTag{
			Key:    tag["key"].(string),
			Values: expandSyntheticsTagValues(tag["values"].([]interface{})),
		}
		out[i] = expanded
	}
	return out
}

func expandSyntheticsTagValues(v []interface{}) []string {
	values := make([]string, len(v))
	for i, value := range v {
		values[i] = value.(string)
	}
	return values
}

func expandStringSlice(strings []interface{}) []string {
	out := make([]string, len(strings))

	for i, v := range strings {
		out[i] = v.(string)
	}
	return out
}

//validation function to validate monitor period
func listValidSyntheticsMonitorPeriods() []string {
	return []string{
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_MINUTE),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_10_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_15_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_30_MINUTES),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_HOUR),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_6_HOURS),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_12_HOURS),
		string(synthetics.SyntheticsMonitorPeriodTypes.EVERY_DAY),
	}
}

//validate func to validate monitor status
func listValidSyntheticsMonitorStatuses() []string {
	return []string{
		string(synthetics.SyntheticsMonitorStatusTypes.DISABLED),
		string(synthetics.SyntheticsMonitorStatusTypes.ENABLED),
		string(synthetics.SyntheticsMonitorStatusTypes.MUTED),
	}
}

func listValidSyntheticsMonitorPublicLocations() []string {
	return []string{
		string(syntheticsPublicLocations.AWS_US_EAST_1),
		string(syntheticsPublicLocations.AWS_US_EAST_2),
		string(syntheticsPublicLocations.AWS_US_WEST_1),
		string(syntheticsPublicLocations.AWS_US_WEST_2),
		string(syntheticsPublicLocations.AWS_CA_CENTRAL_1),
		string(syntheticsPublicLocations.AWS_EU_WEST_1),
		string(syntheticsPublicLocations.AWS_EU_WEST_2),
		string(syntheticsPublicLocations.AWS_EU_WEST_3),
		string(syntheticsPublicLocations.AWS_EU_CENTRAL_1),
		string(syntheticsPublicLocations.AWS_EU_SOUTH_1),
		string(syntheticsPublicLocations.AWS_EU_NORTH_1),
		string(syntheticsPublicLocations.AWS_SA_EAST_1),
		string(syntheticsPublicLocations.AWS_AF_SOUTH_1),
		string(syntheticsPublicLocations.AWS_AP_EAST_1),
		string(syntheticsPublicLocations.AWS_ME_SOUTH_1),
		string(syntheticsPublicLocations.AWS_AP_SOUTH_1),
		string(syntheticsPublicLocations.AWS_AP_NORTHEAST_2),
		string(syntheticsPublicLocations.AWS_AP_SOUTHEAST_1),
		string(syntheticsPublicLocations.AWS_AP_NORTHEAST_1),
		string(syntheticsPublicLocations.AWS_AP_SOUTHEAST_2),
	}
}

// TODO: Move to newrelic-client-go
type SyntheticsPrivateLocation string

// TODO: Move to newrelic-client-go
//nolint:revive
var syntheticsPublicLocations = struct {
	AWS_US_EAST_1      SyntheticsPrivateLocation
	AWS_US_EAST_2      SyntheticsPrivateLocation
	AWS_CA_CENTRAL_1   SyntheticsPrivateLocation
	AWS_US_WEST_1      SyntheticsPrivateLocation
	AWS_US_WEST_2      SyntheticsPrivateLocation
	AWS_EU_WEST_1      SyntheticsPrivateLocation
	AWS_EU_WEST_2      SyntheticsPrivateLocation
	AWS_EU_WEST_3      SyntheticsPrivateLocation
	AWS_EU_CENTRAL_1   SyntheticsPrivateLocation
	AWS_EU_SOUTH_1     SyntheticsPrivateLocation
	AWS_EU_NORTH_1     SyntheticsPrivateLocation
	AWS_SA_EAST_1      SyntheticsPrivateLocation
	AWS_AF_SOUTH_1     SyntheticsPrivateLocation
	AWS_AP_EAST_1      SyntheticsPrivateLocation
	AWS_ME_SOUTH_1     SyntheticsPrivateLocation
	AWS_AP_SOUTH_1     SyntheticsPrivateLocation
	AWS_AP_NORTHEAST_2 SyntheticsPrivateLocation
	AWS_AP_SOUTHEAST_1 SyntheticsPrivateLocation
	AWS_AP_NORTHEAST_1 SyntheticsPrivateLocation
	AWS_AP_SOUTHEAST_2 SyntheticsPrivateLocation
}{
	// US
	AWS_US_EAST_1:    "AWS_US_EAST_1",
	AWS_US_EAST_2:    "AWS_US_EAST_2",
	AWS_US_WEST_1:    "AWS_US_WEST_1",
	AWS_US_WEST_2:    "AWS_US_WEST_2",
	AWS_CA_CENTRAL_1: "AWS_CA_CENTRAL_1",

	// Europe
	AWS_EU_WEST_1:    "AWS_EU_WEST_1",
	AWS_EU_WEST_2:    "AWS_EU_WEST_2",
	AWS_EU_WEST_3:    "AWS_EU_WEST_3",
	AWS_EU_CENTRAL_1: "AWS_EU_CENTRAL_1",
	AWS_EU_SOUTH_1:   "AWS_EU_SOUTH_1",
	AWS_EU_NORTH_1:   "AWS_EU_NORTH_1",

	// South America
	AWS_SA_EAST_1: "AWS_SA_EAST_1",

	// Africa
	AWS_AF_SOUTH_1: "AWS_AF_SOUTH_1",

	// Asia
	AWS_AP_EAST_1:      "AWS_AP_EAST_1",
	AWS_ME_SOUTH_1:     "AWS_ME_SOUTH_1",
	AWS_AP_SOUTH_1:     "AWS_AP_SOUTH_1",
	AWS_AP_NORTHEAST_2: "AWS_AP_NORTHEAST_2",
	AWS_AP_SOUTHEAST_1: "AWS_AP_SOUTHEAST_1",
	AWS_AP_NORTHEAST_1: "AWS_AP_NORTHEAST_1",

	// Australia
	AWS_AP_SOUTHEAST_2: "AWS_AP_SOUTHEAST_2",
}
