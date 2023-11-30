package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	"golang.org/x/exp/maps"
)

func resourceNewRelicMonitorDowntime() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicMonitorDowntimeCreate,
		ReadContext:   resourceNewRelicMonitorDowntimeRead,
		UpdateContext: resourceNewRelicMonitorDowntimeUpdate,
		DeleteContext: resourceNewRelicMonitorDowntimeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "A name to identify the Monitor Downtime to be created.",
				Required:    true,
			},
			"monitor_guids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of GUIDs of monitors, to which the created Monitor Downtime shall be applied.",
				// ValidateFunc: validation included in validateMonitorDowntimeMonitorGUIDs as this is a set and is unsupported by the "validation" package
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The ID of the New Relic account in which the Monitor Downtime shall be created. Defaults to NEW_RELIC_ACCOUNT_ID if not specified.",
				Optional:    true,
				Default:     os.Getenv("NEW_RELIC_ACCOUNT_ID"),
			},
			"start_time": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A datetime stamp signifying the start of the Monitor Downtime.",
				ValidateFunc: validateNaiveDateTime,
			},
			"end_time": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A datetime stamp signifying the end of the Monitor Downtime.",
				ValidateFunc: validateNaiveDateTime,
			},
			"time_zone": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The timezone that applies to the Monitor Downtime schedule.",
				ValidateFunc: validateMonitorDowntimeTimeZone,
			},
			// used with daily, weekly and monthly monitor downtime
			"end_repeat": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "A specification of when the Monitor Downtime should end its repeat cycle, by number of occurrences or date.",
				// ValidateFunc: validation included in validateMonitorDowntimeEndRepeatStructure as this is a set; lists and sets are not supported by the "validation" package
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_date": {
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"end_repeat.0.on_date", "end_repeat.0.on_repeat"},
							Description:  "A date, on which the Monitor Downtime's repeat cycle is expected to end.",
							ValidateFunc: validateMonitorDowntimeOnDate,
						},
						"on_repeat": {
							Type:         schema.TypeInt,
							Optional:     true,
							ExactlyOneOf: []string{"end_repeat.0.on_date", "end_repeat.0.on_repeat"},
							Description:  "Number of repetitions after which the Monitor Downtime's repeat cycle is expected to end.",
						},
					},
				},
			},
			// used with weekly monitor downtime
			"maintenance_days": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of maintenance days to be included with the created weekly Monitor Downtime.",
				// ValidateFunc: validation included in validateMonitorDowntimeMaintenanceDaysStructure as this is a set; lists and sets are not supported by the "validation" package
			},
			// used with monthly monitor downtime
			"frequency": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "Configuration options for which days of the month a monitor downtime will occur",
				// ValidateFunc: validation included in validateMonitorDowntimeFrequencyStructure to use this argument only with "MONTHLY" mode
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days_of_month": {
							Type:         schema.TypeSet,
							Elem:         &schema.Schema{Type: schema.TypeInt},
							Optional:     true,
							ExactlyOneOf: []string{"frequency.0.days_of_month", "frequency.0.days_of_week"},
							Description:  "A numerical list of days of a month on which the Monitor Downtime is scheduled to run.",
							// ValidateFunc: validation included in validateMonitorDowntimeFrequencyStructure as this is a set; lists and sets are not supported by the "validation" package
						},
						"days_of_week": {
							Type:         schema.TypeList,
							MinItems:     1,
							MaxItems:     1,
							Optional:     true,
							ExactlyOneOf: []string{"frequency.0.days_of_month", "frequency.0.days_of_week"},
							Description:  "A list of days of the week on which the Monitor Downtime is scheduled to run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ordinal_day_of_month": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "An occurrence of the day selected within the month.",
										ValidateFunc: validation.StringInSlice(listValidOrdinalDayOfMonthValues(), false),
									},
									"week_day": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The day of the week on which the Monitor Downtime would run.",
										ValidateFunc: validation.StringInSlice(listValidWeekDayValues(), false),
									},
								},
							},
						},
					},
				},
			},
			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An identifier of the type of Monitor Downtime to be created.",
				ValidateFunc: validation.StringInSlice([]string{
					SyntheticsMonitorDowntimeModes.ONE_TIME,
					SyntheticsMonitorDowntimeModes.DAILY,
					SyntheticsMonitorDowntimeModes.MONTHLY,
					SyntheticsMonitorDowntimeModes.WEEKLY,
				}, false),
				ForceNew: true,
			},
		},
		CustomizeDiff: validateMonitorDowntimeAttributes,
	}
}

var requiredArgumentsList = []string{
	"account_id",
	"name",
	"mode",
	"start_time",
	"end_time",
	"time_zone",
}

func getValuesOfMonthlyMonitorDowntimeArguments(d *schema.ResourceData) (map[string]interface{}, error) {
	monthlyMonitorDowntimeArgumentsMap := make(map[string]interface{})

	dailyMonitorDowntimeArgumentsMap, err := getValuesOfDailyMonitorDowntimeArguments(d)
	if err != nil {
		return nil, err
	}

	maps.Copy(monthlyMonitorDowntimeArgumentsMap, dailyMonitorDowntimeArgumentsMap)

	_, ok := d.GetOk("frequency")
	if !ok {
		return nil, errors.New("`frequency` is a required argument with monthly monitor downtime")
	} else {
		var frequencyInput synthetics.SyntheticsMonitorDowntimeMonthlyFrequency
		daysOfMonth, daysOfMonthOk := d.GetOk("frequency.0.days_of_month")
		_, daysOfWeekOk := d.GetOk("frequency.0.days_of_week")
		if !daysOfMonthOk && !daysOfWeekOk {
			return nil, errors.New("the block `frequency` requires one of `days_of_month` or `days_of_week` to be specified")
		} else if daysOfMonthOk && daysOfWeekOk {
			return nil, errors.New("the block `frequency` requires one of `days_of_month` or `days_of_week` to be specified but not both")
		} else if daysOfMonthOk && !daysOfWeekOk {
			frequencyInput.DaysOfMonth = getFrequencyDaysOfMonthList(daysOfMonth.(*schema.Set).List())
		} else {
			var daysOfWeekInput synthetics.SyntheticsDaysOfWeek
			ordinalDayOfMonth, ordinalDayOfMonthOk := d.GetOk("frequency.0.days_of_week.0.ordinal_day_of_month")
			weekDay, weekDayOk := d.GetOk("frequency.0.days_of_week.0.week_day")
			if !ordinalDayOfMonthOk && !weekDayOk {
				return nil, errors.New("the block `days_of_week` requires specifying both `ordinal_day_of_month` and `week_day`")
			}
			daysOfWeekInput.WeekDay = synthetics.SyntheticsMonitorDowntimeWeekDays(weekDay.(string))
			daysOfWeekInput.OrdinalDayOfMonth = synthetics.SyntheticsMonitorDowntimeDayOfMonthOrdinal(ordinalDayOfMonth.(string))
			frequencyInput.DaysOfWeek = &daysOfWeekInput
		}
		monthlyMonitorDowntimeArgumentsMap["frequency"] = frequencyInput
	}

	return monthlyMonitorDowntimeArgumentsMap, nil
}

func getValuesOfWeeklyMonitorDowntimeArguments(d *schema.ResourceData) (map[string]interface{}, error) {
	weeklyMonitorDowntimeArgumentsMap := make(map[string]interface{})

	dailyMonitorDowntimeArgumentsMap, err := getValuesOfDailyMonitorDowntimeArguments(d)
	if err != nil {
		return nil, err
	}

	maps.Copy(weeklyMonitorDowntimeArgumentsMap, dailyMonitorDowntimeArgumentsMap)

	// mandatory argument
	listOfMaintenanceDaysInConfiguration, err := getMaintenanceDaysList(d)
	if err != nil {
		return nil, err
	}
	maintenanceDays, err := convertSyntheticsMonitorDowntimeMaintenanceDays(listOfMaintenanceDaysInConfiguration)
	if err != nil {
		return nil, err
	}
	weeklyMonitorDowntimeArgumentsMap["maintenance_days"] = maintenanceDays

	return weeklyMonitorDowntimeArgumentsMap, nil
}

func getValuesOfDailyMonitorDowntimeArguments(d *schema.ResourceData) (map[string]interface{}, error) {
	dailyMonitorDowntimeArgumentsMap := make(map[string]interface{})

	monitorGUIDs, err := getMonitorGUIDs(d)
	if err != nil {
		return nil, err
	}

	dailyMonitorDowntimeArgumentsMap["monitor_guids"] = monitorGUIDs

	_, ok := d.GetOk("end_repeat")
	if ok {
		// endRepeatStruct := endRepeat.(map[string]interface{})
		var endRepeatInput synthetics.SyntheticsDateWindowEndConfig
		onDate, onDateOk := d.GetOk("end_repeat.0.on_date")
		onRepeat, onRepeatOk := d.GetOk("end_repeat.0.on_repeat")

		if !onDateOk && !onRepeatOk {
			return nil, errors.New("the block `end_repeat` requires one of `on_date` or `on_repeat` to be specified")
		} else if onDateOk && onRepeatOk {
			return nil, errors.New("the block `end_repeat` requires only one of `on_date` or `on_repeat` to be specified, both cannot be specified")
		}

		endRepeatInput.OnDate = synthetics.Date(onDate.(string))
		endRepeatInput.OnRepeat = onRepeat.(int)
		dailyMonitorDowntimeArgumentsMap["end_repeat"] = endRepeatInput

	} else {
		dailyMonitorDowntimeArgumentsMap["end_repeat"] = synthetics.SyntheticsDateWindowEndConfig{}
	}

	return dailyMonitorDowntimeArgumentsMap, nil

}

func getMonitorGUIDs(d *schema.ResourceData) ([]synthetics.EntityGUID, error) {
	val, ok := d.GetOk("monitor_guids")
	if ok {
		in := val.(*schema.Set).List()
		out := make([]synthetics.EntityGUID, len(in))
		for i := range in {
			out[i] = synthetics.EntityGUID(in[i].(string))
		}
		if len(out) == 0 {
			return nil, errors.New("invalid specification of monitor GUIDs: empty list received in the argument 'monitor_guids'")
		} else {
			return out, nil
		}
	}
	return nil, nil
}

func getValuesOfRequiredArguments(d *schema.ResourceData) (map[string]string, error) {
	requiredArgumentsMap := make(map[string]string)
	for _, requiredAttribute := range requiredArgumentsList {
		val, ok := d.GetOk(requiredAttribute)
		switch requiredAttribute {
		case "account_id":
			if ok {
				if val.(string) == "" {
					return nil, errors.New(fmt.Sprintf("%s has value \"\"", requiredAttribute))
				} else {
					requiredArgumentsMap[requiredAttribute] = val.(string)
				}
			} else {
				requiredArgumentsMap[requiredAttribute] = os.Getenv("NEW_RELIC_ACCOUNT_ID")
			}
			break
		default:
			if ok {
				if val.(string) == "" {
					return nil, errors.New(fmt.Sprintf("%s has value \"\"", requiredAttribute))
				} else {
					requiredArgumentsMap[requiredAttribute] = val.(string)
				}
			} else {
				return nil, errors.New(fmt.Sprintf(" value of argument %s not specified", requiredAttribute))
			}
		}
	}
	return requiredArgumentsMap, nil
}

func resourceNewRelicMonitorDowntimeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	requiredArgumentsMap, err := getValuesOfRequiredArguments(d)
	if err != nil {
		return diag.FromErr(err)
	}
	accountIdAsInteger, err := strconv.Atoi(requiredArgumentsMap["account_id"])
	if err != nil {
		return diag.FromErr(err)
	}

	switch requiredArgumentsMap["mode"] {
	case "ONE_TIME":
		monitorGUIDs, err := getMonitorGUIDs(d)
		if err != nil {
			return diag.FromErr(err)
		}
		resp, err := client.Synthetics.SyntheticsCreateOnceMonitorDowntimeWithContext(
			ctx,
			accountIdAsInteger,
			synthetics.NaiveDateTime(requiredArgumentsMap["end_time"]),
			monitorGUIDs,
			requiredArgumentsMap["name"],
			synthetics.NaiveDateTime(requiredArgumentsMap["start_time"]),
			requiredArgumentsMap["time_zone"],
		)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if resp == nil {
			d.SetId("")
			return diag.FromErr(errors.New("encountered an API error while trying to create a monitor downtime: nil response returned"))
		}
		d.SetId(string(resp.GUID))
		break
	case "DAILY":
		conditionalAttributesMap, err := getValuesOfDailyMonitorDowntimeArguments(d)
		if err != nil {
			return diag.FromErr(err)
		}
		resp, err := client.Synthetics.SyntheticsCreateDailyMonitorDowntimeWithContext(
			ctx,
			accountIdAsInteger,
			conditionalAttributesMap["end_repeat"].(synthetics.SyntheticsDateWindowEndConfig),
			synthetics.NaiveDateTime(requiredArgumentsMap["end_time"]),
			conditionalAttributesMap["monitor_guids"].([]synthetics.EntityGUID),
			requiredArgumentsMap["name"],
			synthetics.NaiveDateTime(requiredArgumentsMap["start_time"]),
			requiredArgumentsMap["time_zone"],
		)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if resp == nil {
			d.SetId("")
			return diag.FromErr(errors.New("encountered an API error while trying to create a monitor downtime: nil response returned"))
		}
		d.SetId(string(resp.GUID))
		break
	case "WEEKLY":
		err := validateMonitorDowntimeMaintenanceDays(d)
		if err != nil {
			return diag.FromErr(err)
		}
		conditionalAttributesMap, err := getValuesOfWeeklyMonitorDowntimeArguments(d)
		if err != nil {
			return diag.FromErr(err)
		}
		resp, err := client.Synthetics.SyntheticsCreateWeeklyMonitorDowntimeWithContext(
			ctx,
			accountIdAsInteger,
			conditionalAttributesMap["end_repeat"].(synthetics.SyntheticsDateWindowEndConfig),
			synthetics.NaiveDateTime(requiredArgumentsMap["end_time"]),
			conditionalAttributesMap["maintenance_days"].([]synthetics.SyntheticsMonitorDowntimeWeekDays),
			conditionalAttributesMap["monitor_guids"].([]synthetics.EntityGUID),
			requiredArgumentsMap["name"],
			synthetics.NaiveDateTime(requiredArgumentsMap["start_time"]),
			requiredArgumentsMap["time_zone"],
		)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if resp == nil {
			d.SetId("")
			return diag.FromErr(errors.New("encountered an API error while trying to create a monitor downtime: nil response returned"))
		}
		d.SetId(string(resp.GUID))
		break
	case "MONTHLY":
		conditionalAttributesMap, err := getValuesOfMonthlyMonitorDowntimeArguments(d)
		if err != nil {
			return diag.FromErr(err)
		}
		resp, err := client.Synthetics.SyntheticsCreateMonthlyMonitorDowntimeWithContext(
			ctx,
			accountIdAsInteger,
			conditionalAttributesMap["end_repeat"].(synthetics.SyntheticsDateWindowEndConfig),
			synthetics.NaiveDateTime(requiredArgumentsMap["end_time"]),
			conditionalAttributesMap["frequency"].(synthetics.SyntheticsMonitorDowntimeMonthlyFrequency),
			conditionalAttributesMap["monitor_guids"].([]synthetics.EntityGUID),
			requiredArgumentsMap["name"],
			synthetics.NaiveDateTime(requiredArgumentsMap["start_time"]),
			requiredArgumentsMap["time_zone"],
		)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if resp == nil {
			d.SetId("")
			return diag.FromErr(errors.New("encountered an API error while trying to create a monitor downtime: nil response returned"))
		}
		d.SetId(string(resp.GUID))
		break
	default:
		return diag.FromErr(errors.New("invalid mode of operation: 'mode' can be 'ONE_TIME', 'DAILY', 'WEEKLY' or 'MONTHLY'"))
	}

	return resourceNewRelicMonitorDowntimeRead(ctx, d, meta)
}

func resourceNewRelicMonitorDowntimeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	// accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Synthetics Monitor Downtime %s", d.Id())

	// *** THIS WORKS TOO ***
	//time.Sleep(5 * time.Second)
	//resp, err := client.Entities.GetEntitySearchByQueryWithContext(ctx, entities.EntitySearchOptions{}, fmt.Sprintf("id = '%s'", d.Id()), []entities.EntitySearchSortCriteria{})
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//x := resp.Results.Entities
	//for _, val := range x {
	//	entity := val.(*entities.GenericEntityOutline)
	//	tags := entity.GetTags()
	//	_ = d.Set("name", entity.GetName())
	//	_ = d.Set("account_id", monitorDowntimeAttributeReaderMap["account_id"].(func([]entities.EntityTag) string)(tags))
	//	_ = d.Set("mode", monitorDowntimeAttributeReaderMap["mode"].(func([]entities.EntityTag) string)(tags))
	//	_ = d.Set("start_time", monitorDowntimeAttributeReaderMap["start_time"].(func([]entities.EntityTag) string)(tags))
	//	_ = d.Set("end_time", monitorDowntimeAttributeReaderMap["end_time"].(func([]entities.EntityTag) string)(tags))
	//	_ = d.Set("time_zone", monitorDowntimeAttributeReaderMap["time_zone"].(func([]entities.EntityTag) string)(tags))
	//}

	var tags []entities.EntityTag
	var entity *entities.GenericEntity

	// retry mechanism since the entity query "immediately" does NOT return all tags, and returns only three
	retryErr := resource.RetryContext(context.Background(), 30*time.Second, func() *resource.RetryError {
		resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))
		if err != nil {
			return resource.RetryableError(err)
		}
		entity = (*resp).(*entities.GenericEntity)
		tags = entity.GetTags()
		if len(tags) < 4 {
			return resource.RetryableError(fmt.Errorf("enough tags not found. retrying"))
		}
		return nil
	})

	if retryErr != nil {
		log.Fatalf("Unable to find application entity: %s", retryErr)
	}

	mode := monitorDowntimeAttributeReaderMap["mode"].(func([]entities.EntityTag) string)(tags)
	timezone := monitorDowntimeAttributeReaderMap["time_zone"].(func([]entities.EntityTag) string)(tags)
	_ = d.Set("name", entity.GetName())
	_ = d.Set("monitor_guids", monitorDowntimeAttributeReaderMap["monitor_guids"].(func([]entities.EntityRelationship, common.EntityGUID) []string)(entity.GetRelationships(), common.EntityGUID(d.Id())))
	_ = d.Set("account_id", monitorDowntimeAttributeReaderMap["account_id"].(func([]entities.EntityTag) string)(tags))
	_ = d.Set("mode", mode)
	_ = d.Set("start_time", monitorDowntimeAttributeReaderMap["start_time"].(func([]entities.EntityTag) string)(tags))
	_ = d.Set("end_time", monitorDowntimeAttributeReaderMap["end_time"].(func([]entities.EntityTag) string)(tags))
	_ = d.Set("time_zone", timezone)

	if mode != "ONE_TIME" {
		setMonitorDowntimeEndRepeat(d, tags, timezone)
	}

	if mode == "WEEKLY" {
		setMonitorDowntimeMaintenanceDays(d, tags)
	}

	if mode == "MONTHLY" {
		setMonitorDowntimeFrequency(d, tags)
	}
	return nil

}

func resourceNewRelicMonitorDowntimeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	requiredArgumentsMap, err := getValuesOfRequiredArguments(d)
	if err != nil {
		return diag.FromErr(err)
	}

	switch requiredArgumentsMap["mode"] {
	case "ONE_TIME":
		monitorGUIDs, err := getMonitorGUIDs(d)
		if err != nil {
			return diag.FromErr(err)
		}

		retryErr := resource.RetryContext(context.Background(), 30*time.Second, func() *resource.RetryError {
			resp, err := client.Synthetics.SyntheticsEditMonitorDowntimeWithContext(
				ctx,
				synthetics.SyntheticsMonitorDowntimeDailyConfig{},
				synthetics.EntityGUID(d.Id()),
				monitorGUIDs,
				synthetics.SyntheticsMonitorDowntimeMonthlyConfig{},
				requiredArgumentsMap["name"],
				synthetics.SyntheticsMonitorDowntimeOnceConfig{
					EndTime:   synthetics.NaiveDateTime(requiredArgumentsMap["end_time"]),
					StartTime: synthetics.NaiveDateTime(requiredArgumentsMap["start_time"]),
					Timezone:  requiredArgumentsMap["time_zone"],
				},
				synthetics.SyntheticsMonitorDowntimeWeeklyConfig{},
			)
			if err != nil {
				if err.Error() == "An error occurred resolving this field" {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)

			}

			if resp == nil {
				return resource.NonRetryableError(errors.New("encountered an API error while trying to create a monitor downtime: nil response returned"))
			}
			return nil
		})

		if retryErr != nil {
			log.Fatalf("Unable to find application entity: %s", retryErr)
		}

		if retryErr != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		break
	case "DAILY":
		conditionalAttributesMap, err := getValuesOfDailyMonitorDowntimeArguments(d)
		// TBD
		x := conditionalAttributesMap["end_repeat"].(synthetics.SyntheticsDateWindowEndConfig)

		if err != nil {
			return diag.FromErr(err)
		}
		resp, err := client.Synthetics.SyntheticsEditMonitorDowntimeWithContext(
			ctx,
			synthetics.SyntheticsMonitorDowntimeDailyConfig{},
			synthetics.EntityGUID(d.Id()),
			conditionalAttributesMap["monitor_guids"].([]synthetics.EntityGUID),
			synthetics.SyntheticsMonitorDowntimeMonthlyConfig{
				EndTime:   synthetics.NaiveDateTime(requiredArgumentsMap["end_time"]),
				StartTime: synthetics.NaiveDateTime(requiredArgumentsMap["start_time"]),
				Timezone:  requiredArgumentsMap["time_zone"],
				EndRepeat: x,
			},
			requiredArgumentsMap["name"],
			synthetics.SyntheticsMonitorDowntimeOnceConfig{},
			synthetics.SyntheticsMonitorDowntimeWeeklyConfig{},
		)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if resp == nil {
			d.SetId("")
			return diag.FromErr(errors.New("encountered an API error while trying to create a monitor downtime: nil response returned"))
		}
		d.SetId(string(resp.GUID))
		break
	case "WEEKLY":
		err := validateMonitorDowntimeMaintenanceDays(d)
		if err != nil {
			return diag.FromErr(err)
		}
		conditionalAttributesMap, err := getValuesOfWeeklyMonitorDowntimeArguments(d)
		if err != nil {
			return diag.FromErr(err)
		}

		x := conditionalAttributesMap["end_repeat"].(synthetics.SyntheticsDateWindowEndConfig)
		y := conditionalAttributesMap["maintenance_days"].([]synthetics.SyntheticsMonitorDowntimeWeekDays)

		resp, err := client.Synthetics.SyntheticsEditMonitorDowntimeWithContext(
			ctx,
			synthetics.SyntheticsMonitorDowntimeDailyConfig{},
			synthetics.EntityGUID(d.Id()),
			conditionalAttributesMap["monitor_guids"].([]synthetics.EntityGUID),
			synthetics.SyntheticsMonitorDowntimeMonthlyConfig{},
			requiredArgumentsMap["name"],
			synthetics.SyntheticsMonitorDowntimeOnceConfig{},
			synthetics.SyntheticsMonitorDowntimeWeeklyConfig{
				EndTime:         synthetics.NaiveDateTime(requiredArgumentsMap["end_time"]),
				StartTime:       synthetics.NaiveDateTime(requiredArgumentsMap["start_time"]),
				Timezone:        requiredArgumentsMap["time_zone"],
				EndRepeat:       x,
				MaintenanceDays: y,
			},
		)

		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if resp == nil {
			d.SetId("")
			return diag.FromErr(errors.New("encountered an API error while trying to create a monitor downtime: nil response returned"))
		}
		d.SetId(string(resp.GUID))
		break
	case "MONTHLY":
		conditionalAttributesMap, err := getValuesOfMonthlyMonitorDowntimeArguments(d)
		if err != nil {
			return diag.FromErr(err)
		}
		x := conditionalAttributesMap["end_repeat"].(synthetics.SyntheticsDateWindowEndConfig)
		y := conditionalAttributesMap["frequency"].(synthetics.SyntheticsMonitorDowntimeMonthlyFrequency)
		resp, err := client.Synthetics.SyntheticsEditMonitorDowntimeWithContext(
			ctx,
			synthetics.SyntheticsMonitorDowntimeDailyConfig{},
			synthetics.EntityGUID(d.Id()),
			conditionalAttributesMap["monitor_guids"].([]synthetics.EntityGUID),
			synthetics.SyntheticsMonitorDowntimeMonthlyConfig{
				EndTime:   synthetics.NaiveDateTime(requiredArgumentsMap["end_time"]),
				StartTime: synthetics.NaiveDateTime(requiredArgumentsMap["start_time"]),
				Timezone:  requiredArgumentsMap["time_zone"],
				EndRepeat: x,
				Frequency: y,
			},
			requiredArgumentsMap["name"],
			synthetics.SyntheticsMonitorDowntimeOnceConfig{},
			synthetics.SyntheticsMonitorDowntimeWeeklyConfig{},
		)

		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if resp == nil {
			d.SetId("")
			return diag.FromErr(errors.New("encountered an API error while trying to create a monitor downtime: nil response returned"))
		}
		d.SetId(string(resp.GUID))
		break
	default:
		return diag.FromErr(errors.New("invalid mode of operation: 'mode' can be 'ONE_TIME', 'DAILY', 'WEEKLY' or 'MONTHLY'"))
	}

	return resourceNewRelicMonitorDowntimeRead(ctx, d, meta)
}

func resourceNewRelicMonitorDowntimeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	resp, err := client.Synthetics.SyntheticsDeleteMonitorDowntimeWithContext(ctx, synthetics.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.FromErr(errors.New("encountered an API error while trying to delete the monitor downtime: nil response returned"))
	}
	return nil
}
