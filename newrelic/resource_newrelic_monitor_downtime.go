package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
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
				Type:        schema.TypeInt,
				Description: "The ID of the New Relic account in which the Monitor Downtime shall be created. Defaults to the `account_id` in the provider{} configuration if not specified.",
				Optional:    true,
				Computed:    true,
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
					SyntheticsMonitorDowntimeModes.OneTime,
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

func resourceNewRelicMonitorDowntimeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	commonArgumentsObject, err := getMonitorDowntimeValuesOfCommonArguments(d, providerConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	switch commonArgumentsObject.Mode {
	case SyntheticsMonitorDowntimeModes.OneTime:
		oneTimeCreateObject, err := getMonitorDowntimeOneTimeValues(d, commonArgumentsObject)
		if err != nil {
			return diag.FromErr(err)
		}

		guid, err := oneTimeCreateObject.createMonitorDowntimeOneTime(ctx, client)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		d.SetId(guid)
	case SyntheticsMonitorDowntimeModes.DAILY:
		dailyCreateObject, err := getMonitorDowntimeDailyValues(d, commonArgumentsObject)
		if err != nil {
			return diag.FromErr(err)
		}

		guid, err := dailyCreateObject.createMonitorDowntimeDaily(ctx, client)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		d.SetId(guid)
	case SyntheticsMonitorDowntimeModes.WEEKLY:
		weeklyCreateObject, err := getMonitorDowntimeWeeklyValues(d, commonArgumentsObject)
		if err != nil {
			return diag.FromErr(err)
		}

		guid, err := weeklyCreateObject.createMonitorDowntimeWeekly(ctx, client)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		d.SetId(guid)
	case SyntheticsMonitorDowntimeModes.MONTHLY:
		monthlyCreateObject, err := getMonitorDowntimeMonthlyValues(d, commonArgumentsObject)
		if err != nil {
			return diag.FromErr(err)
		}

		guid, err := monthlyCreateObject.createMonitorDowntimeMonthly(ctx, client)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		d.SetId(guid)
	default:
		return diag.FromErr(errors.New("invalid mode of operation: 'mode' can be 'ONE_TIME', 'DAILY', 'WEEKLY' or 'MONTHLY'"))
	}

	return resourceNewRelicMonitorDowntimeRead(ctx, d, meta)
}

func resourceNewRelicMonitorDowntimeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic Synthetics Monitor Downtime %s", d.Id())

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

	mode := setMonitorDowntimeMode(tags)
	timezone := setMonitorDowntimeTimezone(tags)
	_ = d.Set("name", entity.GetName())
	_ = d.Set("account_id", setMonitorDowntimeAccountID(tags))
	_ = d.Set("mode", mode)
	_ = d.Set("start_time", setMonitorDowntimeStartTime(tags))
	_ = d.Set("end_time", setMonitorDowntimeEndTime(tags))
	_ = d.Set("time_zone", timezone)

	if mode != SyntheticsMonitorDowntimeModes.OneTime {
		setMonitorDowntimeEndRepeat(d, tags, timezone)
	}

	if mode == SyntheticsMonitorDowntimeModes.WEEKLY {
		setMonitorDowntimeMaintenanceDays(d, tags)
	}

	if mode == SyntheticsMonitorDowntimeModes.MONTHLY {
		setMonitorDowntimeFrequency(d, tags)
	}
	return nil

}

func resourceNewRelicMonitorDowntimeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	commonArgumentsObject, err := getMonitorDowntimeValuesOfCommonArguments(d, providerConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	switch commonArgumentsObject.Mode {
	case SyntheticsMonitorDowntimeModes.OneTime:
		oneTimeUpdateObject, err := getMonitorDowntimeOneTimeValues(d, commonArgumentsObject)
		if err != nil {
			return diag.FromErr(err)
		}

		guid, err := oneTimeUpdateObject.updateMonitorDowntimeOneTime(ctx, client, synthetics.EntityGUID(d.Id()))
		if err != nil {
			// d.SetId("")
			return diag.FromErr(err)
		}

		d.SetId(guid)
	case SyntheticsMonitorDowntimeModes.DAILY:
		dailyUpdateObject, err := getMonitorDowntimeDailyValues(d, commonArgumentsObject)
		if err != nil {
			return diag.FromErr(err)
		}

		guid, err := dailyUpdateObject.updateMonitorDowntimeDaily(ctx, client, synthetics.EntityGUID(d.Id()))
		if err != nil {
			// d.SetId("")
			return diag.FromErr(err)
		}

		d.SetId(guid)
	case SyntheticsMonitorDowntimeModes.WEEKLY:
		weeklyUpdateObject, err := getMonitorDowntimeWeeklyValues(d, commonArgumentsObject)
		if err != nil {
			return diag.FromErr(err)
		}

		guid, err := weeklyUpdateObject.updateMonitorDowntimeWeekly(ctx, client, synthetics.EntityGUID(d.Id()))
		if err != nil {
			// d.SetId("")
			return diag.FromErr(err)
		}

		d.SetId(guid)
	case SyntheticsMonitorDowntimeModes.MONTHLY:
		monthlyUpdateObject, err := getMonitorDowntimeMonthlyValues(d, commonArgumentsObject)
		if err != nil {
			return diag.FromErr(err)
		}

		guid, err := monthlyUpdateObject.updateMonitorDowntimeMonthly(ctx, client, synthetics.EntityGUID(d.Id()))
		if err != nil {
			// d.SetId("")
			return diag.FromErr(err)
		}

		d.SetId(guid)
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
