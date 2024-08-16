package newrelic

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	"golang.org/x/exp/slices"
)

// Returns the common schema attributes shared by all Synthetics monitor types.
func syntheticsMonitorCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Description: "ID of the newrelic account.",
			ForceNew:    true,
			Computed:    true,
			Optional:    true,
		},
		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique entity identifier of the monitor in New Relic.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The title of this monitor.",
		},
		"status": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The monitor status (ENABLED or DISABLED).",
			ValidateFunc: validateSyntheticMonitorStatus,
		},
		"tag": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The tags that will be associated with the monitor.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of the tag key",
					},
					"values": {
						Type:        schema.TypeList,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Required:    true,
						Description: "Values associated with the tag key",
					},
				},
			},
		},
		"period": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.",
			ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorPeriods(), false),
		},
		"period_in_minutes": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The interval in minutes at which this monitor should run.",
		},
	}
}

func syntheticsMonitorRuntimeOptions() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"runtime_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The runtime type that the monitor will run.",
		},
		"runtime_type_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The specific semver version of the runtime type.",
		},
	}
}

// NOTE: This can be a shared schema partial for other synthetics monitor resources
func syntheticsMonitorLocationsAsStringsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"locations_private": {
			Type:         schema.TypeSet,
			Elem:         &schema.Schema{Type: schema.TypeString},
			Description:  "List private location GUIDs for which the monitor will run.",
			Optional:     true,
			AtLeastOneOf: []string{"locations_public", "locations_private"},
		},
		"locations_public": {
			Type:         schema.TypeSet,
			Elem:         &schema.Schema{Type: schema.TypeString},
			Description:  "Publicly available location names in which the monitor will run.",
			Optional:     true,
			AtLeastOneOf: []string{"locations_public", "locations_private"},
		},
	}
}

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

// This has been added to provide support to backward compatibility (older tf version supports period in minutes)
var syntheticsMonitorPeriodInMinutesValueMap = map[synthetics.SyntheticsMonitorPeriod]int{

	synthetics.SyntheticsMonitorPeriodTypes.EVERY_MINUTE:     1,
	synthetics.SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES:  5,
	synthetics.SyntheticsMonitorPeriodTypes.EVERY_10_MINUTES: 10,
	synthetics.SyntheticsMonitorPeriodTypes.EVERY_15_MINUTES: 15,
	synthetics.SyntheticsMonitorPeriodTypes.EVERY_30_MINUTES: 30,
	synthetics.SyntheticsMonitorPeriodTypes.EVERY_HOUR:       60,
	synthetics.SyntheticsMonitorPeriodTypes.EVERY_6_HOURS:    360,
	synthetics.SyntheticsMonitorPeriodTypes.EVERY_12_HOURS:   720,
	synthetics.SyntheticsMonitorPeriodTypes.EVERY_DAY:        1440,
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

func expandSyntheticsPublicLocations(locations []interface{}) []string {
	locationsOut := make([]string, len(locations))

	for i, v := range locations {
		locationsOut[i] = v.(string)
	}
	return locationsOut
}

func expandSyntheticsBrowsers(browsers []interface{}) []synthetics.SyntheticsBrowser {
	browsersOut := make([]synthetics.SyntheticsBrowser, len(browsers))

	for i, v := range browsers {
		browsersOut[i] = synthetics.SyntheticsBrowser(v.(string))

	}
	return browsersOut
}

func expandSyntheticsDevices(devices []interface{}) []synthetics.SyntheticsDevice {
	devicesOut := make([]synthetics.SyntheticsDevice, len(devices))

	for i, v := range devices {
		devicesOut[i] = synthetics.SyntheticsDevice(v.(string))
	}
	return devicesOut
}

func expandSyntheticsPrivateLocations(locations []interface{}) []synthetics.SyntheticsPrivateLocationInput {
	locationsOut := make([]synthetics.SyntheticsPrivateLocationInput, len(locations))

	for i, v := range locations {
		pl := v.(map[string]interface{})
		locationsOut[i].GUID = pl["guid"].(string)
		if v, ok := pl["vse_password"]; ok {
			locationsOut[i].VsePassword = synthetics.SecureValue(v.(string))
		}
	}
	return locationsOut
}

func expandSyntheticsCustomHeaders(headers []interface{}) []synthetics.SyntheticsCustomHeaderInput {
	output := make([]synthetics.SyntheticsCustomHeaderInput, len(headers))

	for i, v := range headers {
		header := v.(map[string]interface{})
		expanded := synthetics.SyntheticsCustomHeaderInput{
			Name:  header["name"].(string),
			Value: header["value"].(string),
		}
		output[i] = expanded
	}
	return output
}

// validation function to validate monitor period
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

// function to list out valid Synthetic Monitor status values
func listValidSyntheticsMonitorStatuses() []string {
	return []string{
		string(synthetics.SyntheticsMonitorStatusTypes.DISABLED),
		string(synthetics.SyntheticsMonitorStatusTypes.ENABLED),
		// synthetics.SyntheticsMonitorStatusTypes.MUTED removed on February 29, 2024 in accordance with its EOL
	}
}

// validate function that validates the status of Synthetic Monitors
// recent addition: return an error if 'MUTED' status is used, as this has reached EOL
func validateSyntheticMonitorStatus(val interface{}, key string) (warns []string, errs []error) {
	monitorStatusInput := val.(string)
	listOfValidSyntheticMonitorStatuses := listValidSyntheticsMonitorStatuses()
	containsValidSyntheticMonitorStatus := slices.Contains(listOfValidSyntheticMonitorStatuses, monitorStatusInput)
	if !containsValidSyntheticMonitorStatus {
		// hard-coding "MUTED" instead of using synthetics.SyntheticsMonitorStatusTypes.MUTED as this has been removed from newrelic-client-go, owing to the EOL
		if strings.ToUpper(monitorStatusInput) == "MUTED" {
			errs = append(errs, fmt.Errorf("invalid monitor status 'MUTED' \n"+
				"As of February 29, 2024, Synthetic Monitors no longer support the `MUTED` status.\n"+
				"Version 3.33.0 of the New Relic Terraform Provider is released to coincide with the `MUTED` status end-of-life.\n"+
				"Consequently, the only valid values for `status` for all types of Synthetic Monitors are `ENABLED` and `DISABLED`.\n"+
				"If you have a Terraform configuration with Synthetic Monitor resources previously applied with the status `MUTED` and are hence, seeing this error now upon `terraform plan`,\n"+"please change the status of the monitor to one of the two aforementioned values to plan and apply your configuration, and opt for other alternatives to mute monitors.\n"+
				"For additional information on alternatives to the `MUTED` status of Synthetic Monitors that can be managed via Terraform,\n"+"please refer to the Synthetic Monitors MUTED Status EOL Guide in the documentation of the New Relic Terraform Provider.\n"+
				"https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/upcoming_synthetics_muted_status_eol_guide"))
		} else {
			errs = append(errs, fmt.Errorf("expected status to be one of %v, got %s", listOfValidSyntheticMonitorStatuses, monitorStatusInput))
		}
	}
	return warns, errs
}

type SyntheticsMonitorType string

// nolint:revive
var SyntheticsMonitorTypes = struct {
	SIMPLE         SyntheticsMonitorType
	BROWSER        SyntheticsMonitorType
	SCRIPT_API     SyntheticsMonitorType
	SCRIPT_BROWSER SyntheticsMonitorType
}{
	SIMPLE:         "SIMPLE",
	BROWSER:        "BROWSER",
	SCRIPT_API:     "SCRIPT_API",
	SCRIPT_BROWSER: "SCRIPT_BROWSER",
}

func listValidSyntheticsScriptMonitorTypes() []string {
	return []string{
		string(SyntheticsMonitorTypes.SCRIPT_API),
		string(SyntheticsMonitorTypes.SCRIPT_BROWSER),
	}
}

// TODO: Move to newrelic-client-go
type SyntheticsPublicLocation string

// TODO: Move to newrelic-client-go
//
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

func getRuntimeValuesFromEntityTags(tags []entities.EntityTag) (runtimeType string, runtimeTypeVersion string) {
	runtimeType = ""
	runtimeTypeVersion = ""

	for _, t := range tags {
		if t.Key == "legacyRuntime" {
			for _, v := range t.Values {
				if v == "true" {
					return "", ""
				}
			}
		}

		if t.Key == "runtimeType" {
			runtimeType = t.Values[0]
		}

		if t.Key == "runtimeTypeVersion" {
			runtimeTypeVersion = t.Values[0]
		}
	}

	return runtimeType, runtimeTypeVersion
}

func getMonitorID(monitorGUID string) (string, error) {
	decodedGUID, err := base64.RawStdEncoding.DecodeString(monitorGUID)
	if err != nil {
		return "", err
	}

	// Check if "|" character is present in decodedGUID
	if !strings.Contains(string(decodedGUID), "|") {
		return "", fmt.Errorf("invalid monitor GUID '%s'", monitorGUID)
	}

	splitGUID := strings.Split(string(decodedGUID), "|")

	if len(splitGUID) < 4 {
		return "", fmt.Errorf("invalid monitor GUID '%s'", monitorGUID)
	}

	monitorID := splitGUID[3]

	return monitorID, nil
}

// This map facilitates safely setting the schema attributes which
// are returned as part of the monitor's entity tags. We have a limited
// set of attributes we can set from tags due to technical limitations in the API.
// Note this is a caveat in how synthetics monitor data is provided
// from the entity API. This is not a common resource pattern.
var syntheticsMonitorTagKeyToSchemaAttrMap = map[string]string{
	"runtimeType":        "runtime_type",
	"runtimeTypeVersion": "runtime_type_version",
	"scriptLanguage":     "script_language",
	"deviceOrientation":  "device_orientation",
	"deviceType":         "device_type",
	"browsers":           "browsers",
	"devices":            "devices",
}

func getCertCheckMonitorValuesFromEntityTags(tags []entities.EntityTag) (domain string, daysUntilExpiration int) {
	domain = ""
	daysUntilExpiration = 0

	for _, tag := range tags {

		if tag.Key == "domain" {
			domain = tag.Values[0]
		}

		if tag.Key == "daysUntilExpiration" {
			// Parse string to integer
			days := tag.Values[0]
			daysUntilExpiration, _ = strconv.Atoi(days)
		}
	}

	return domain, daysUntilExpiration
}
