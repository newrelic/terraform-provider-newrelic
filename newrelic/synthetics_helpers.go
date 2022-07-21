package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

var syntheticsMonitorPeriodValueMap = map[int]synthetics.SyntheticsMonitorPeriod{
	1:    synthetics.SyntheticsMonitorPeriodTypes.EVERY_MINUTE,
	5:    synthetics.SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES,
	10:   synthetics.SyntheticsMonitorPeriodTypes.EVERY_10_MINUTES,
	15:   synthetics.SyntheticsMonitorPeriodTypes.EVERY_15_MINUTES,
	30:   synthetics.SyntheticsMonitorPeriodTypes.EVERY_30_MINUTES,
	60:   synthetics.SyntheticsMonitorPeriodTypes.EVERY_HOUR,
	360:  synthetics.SyntheticsMonitorPeriodTypes.EVERY_6_HOURS,
	720:  synthetics.SyntheticsMonitorPeriodTypes.EVERY_12_HOURS,
	1440: synthetics.SyntheticsMonitorPeriodTypes.EVERY_DAY,
}

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

// TODO: Move to newrelic-client-go
type SyntheticsPublicLocation string

// TODO: Move to newrelic-client-go
//nolint:revive
var syntheticsPublicLocations = struct {
	US_EAST_1      SyntheticsPublicLocation
	US_EAST_2      SyntheticsPublicLocation
	CA_CENTRAL_1   SyntheticsPublicLocation
	US_WEST_1      SyntheticsPublicLocation
	US_WEST_2      SyntheticsPublicLocation
	EU_WEST_1      SyntheticsPublicLocation
	EU_WEST_2      SyntheticsPublicLocation
	EU_WEST_3      SyntheticsPublicLocation
	EU_CENTRAL_1   SyntheticsPublicLocation
	EU_SOUTH_1     SyntheticsPublicLocation
	EU_NORTH_1     SyntheticsPublicLocation
	SA_EAST_1      SyntheticsPublicLocation
	AF_SOUTH_1     SyntheticsPublicLocation
	AP_EAST_1      SyntheticsPublicLocation
	ME_SOUTH_1     SyntheticsPublicLocation
	AP_SOUTH_1     SyntheticsPublicLocation
	AP_NORTHEAST_2 SyntheticsPublicLocation
	AP_SOUTHEAST_1 SyntheticsPublicLocation
	AP_NORTHEAST_1 SyntheticsPublicLocation
	AP_SOUTHEAST_2 SyntheticsPublicLocation
}{
	// US
	US_EAST_1:    "US_EAST_1",
	US_EAST_2:    "US_EAST_2",
	US_WEST_1:    "US_WEST_1",
	US_WEST_2:    "US_WEST_2",
	CA_CENTRAL_1: "CA_CENTRAL_1",

	// Europe
	EU_WEST_1:    "EU_WEST_1",
	EU_WEST_2:    "EU_WEST_2",
	EU_WEST_3:    "EU_WEST_3",
	EU_CENTRAL_1: "EU_CENTRAL_1",
	EU_SOUTH_1:   "EU_SOUTH_1",
	EU_NORTH_1:   "EU_NORTH_1",

	// South America
	SA_EAST_1: "SA_EAST_1",

	// Africa
	AF_SOUTH_1: "AF_SOUTH_1",

	// Asia
	AP_EAST_1:      "AP_EAST_1",
	ME_SOUTH_1:     "ME_SOUTH_1",
	AP_SOUTH_1:     "AP_SOUTH_1",
	AP_NORTHEAST_2: "AP_NORTHEAST_2",
	AP_SOUTHEAST_1: "AP_SOUTHEAST_1",
	AP_NORTHEAST_1: "AP_NORTHEAST_1",

	// Australia
	AP_SOUTHEAST_2: "AP_SOUTHEAST_2",
}

var syntheticsPublicLocationsMap = map[string]SyntheticsPublicLocation{
	"Columbus, OH, USA":      syntheticsPublicLocations.US_EAST_2,
	"Montreal, Québec, CA":   syntheticsPublicLocations.CA_CENTRAL_1,
	"Portland, OR, USA":      syntheticsPublicLocations.US_WEST_2,
	"San Francisco, CA, USA": syntheticsPublicLocations.US_WEST_1,
	"Washington, DC, USA":    syntheticsPublicLocations.US_EAST_1,
	"São Paulo, BR":          syntheticsPublicLocations.SA_EAST_1,
	"Hong Kong, HK":          syntheticsPublicLocations.AP_EAST_1,
	"Manama, BH":             syntheticsPublicLocations.ME_SOUTH_1,
	"Mumbai, IN":             syntheticsPublicLocations.AP_SOUTH_1,
	"Seoul, KR":              syntheticsPublicLocations.AP_NORTHEAST_2,
	"Singapore, SG":          syntheticsPublicLocations.AP_SOUTHEAST_1,
	"Tokyo, JP":              syntheticsPublicLocations.AP_NORTHEAST_1,
	"Dublin, IE":             syntheticsPublicLocations.EU_WEST_1,
	"Frankfurt, DE":          syntheticsPublicLocations.EU_CENTRAL_1,
	"London, England, UK":    syntheticsPublicLocations.EU_WEST_2,
	"Milan, IT":              syntheticsPublicLocations.EU_SOUTH_1,
	"Paris, FR":              syntheticsPublicLocations.EU_WEST_3,
	"Stockholm, SE":          syntheticsPublicLocations.EU_NORTH_1,
	"Sydney, AU":             syntheticsPublicLocations.AP_SOUTHEAST_2,
	"Cape Town, ZA":          syntheticsPublicLocations.AF_SOUTH_1,
}

func getPublicLocationsFromEntityTags(tags []entities.EntityTag) []string {
	out := []string{}

	for _, t := range tags {
		if t.Key == "publicLocation" {
			for _, v := range t.Values {
				out = append(out, string(syntheticsPublicLocationsMap[v]))
			}
		}
	}

	return out
}
