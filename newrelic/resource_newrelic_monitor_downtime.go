package newrelic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/accountmanagement"
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
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of GUIDs of monitors, to which the created Monitor Downtime shall be applied.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The ID of the New Relic account in which the Monitor Downtime shall be created. Defaults to NEW_RELIC_ACCOUNT_ID if not specified.",
				Optional:    true,
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The timezone that applies to the Monitor Downtime schedule.",
			},
			// used with daily, weekly and monthly monitor downtime
			"end_repeat": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "A specification of when the Monitor Downtime should end its repeat cycle, by number of occurrences or date.",
				// TODO: define validation to not use this with createOnce monitor downtime and keep this optional with other three types
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_date": {
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"on_date", "on_repeat"},
							Description:  "A date, on which the Monitor Downtime's repeat cycle is expected to end.",
							// TODO: define date validation here (possibly YYYY-MM-DD), didn't do it yet as the mutation is broken
						},
						"on_repeat": {
							Type:         schema.TypeInt,
							Optional:     true,
							ExactlyOneOf: []string{"on_date", "on_repeat"},
							Description:  "Number of repetitions after which the Monitor Downtime's repeat cycle is expected to end.",
						},
					},
				},
			},
			// used with weekly monitor downtime
			"maintenance_days": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of maintenance days to be included with the created weekly Monitor Downtime.",
				// TODO: define validation in such a way that this is used only with weekly monitor downtime
				// TODO: (reqd. only for weekly and not allowed for the rest)
				// TODO: in that function, include an "if" to check if it is "MONDAY", "TUESDAY", ... "SUNDAY"
				// TODO: !! Also check if this works as a list or as a single string; NG is not clear !!
			},
			// used with monthly monitor downtime
			"frequency": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "Configuration options for which days of the month a monitor downtime will occur",
				// TODO: define validation to use this only with monthly monitor downtime
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days_of_month": {
							Type:         schema.TypeList,
							Elem:         &schema.Schema{Type: schema.TypeInt},
							Optional:     true,
							ExactlyOneOf: []string{"days_of_month", "days_of_week"},
							Description:  "A numerical list of days of a month on which the Monitor Downtime is scheduled to run.",
							// TODO: define validation to have these values between 1 and 31
						},
						"days_of_week": {
							Type:         schema.TypeList,
							MinItems:     1,
							MaxItems:     1,
							Optional:     true,
							ExactlyOneOf: []string{"days_of_month", "days_of_week"},
							Description:  "A list of days of the week on which the Monitor Downtime is scheduled to run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ordinal_day_of_month": {
										Type:         schema.TypeString,
										Required:     true,
										ExactlyOneOf: []string{"on_date", "on_repeat"},
										Description:  "An occurrence of the day selected within the month.",
										// TODO: define this to belong to ["FIRST", "SECOND", "THIRD", "FOURTH", "LAST"]
									},
									"week_day": {
										Type:         schema.TypeInt,
										Required:     true,
										ExactlyOneOf: []string{"on_date", "on_repeat"},
										Description:  "The day of the week on which the Monitor Downtime would run.",
										// TODO: define this to belong to ["MONDAY", "TUESDAY", ... "SUNDAY"]
									},
								},
							},
						},
					},
				},
			},
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "An identifier of the type of Monitor Downtime to be created.",
				ValidateFunc: validation.StringInSlice([]string{"ONCE", "DAILY", "MONTHLY", "WEEKLY"}, false),
			},
		},
	}
}

var requiredArgumentsList = []string{
	"account_id",
	"name",
	"mode",
	"start_time",
	"end_time",
	"timezone",
}

func getValuesOfMonthlyMonitorDowntimeArguments(d *schema.ResourceData) (map[string]interface{}, error) {
	var monthlyMonitorDowntimeArgumentsMap map[string]interface{}

	dailyMonitorDowntimeArgumentsMap, err := getValuesOfDailyMonitorDowntimeArguments(d)
	if err != nil {
		return nil, err
	}

	maps.Copy(monthlyMonitorDowntimeArgumentsMap, dailyMonitorDowntimeArgumentsMap)

	frequency, ok := d.GetOk("frequency")
	if !ok {
		return nil, errors.New("`frequency` is a required argument with monthly monitor downtime")
	} else {
		frequencyStruct := frequency.(map[string]interface{})
		var frequencyInput synthetics.SyntheticsMonitorDowntimeMonthlyFrequency
		daysOfMonth, daysOfMonthOk := frequencyStruct["days_of_month"]
		daysOfWeek, daysOfWeekOk := frequencyStruct["days_of_week"]
		if !daysOfMonthOk && !daysOfWeekOk {
			return nil, errors.New("the block `frequency` requires one of `days_of_month` or `days_of_week` to be specified")
		} else if daysOfMonthOk && daysOfWeekOk {
			return nil, errors.New("the block `frequency` requires one of `days_of_month` or `days_of_week` to be specified but not both")
		} else if daysOfMonthOk && !daysOfWeekOk {
			frequencyInput.DaysOfMonth = daysOfMonth.([]int)
		} else {
			daysOfWeekStruct := daysOfWeek.(map[string]interface{})
			var daysOfWeekInput synthetics.SyntheticsDaysOfWeek
			ordinalDayOfMonth, ordinalDayOfMonthOk := daysOfWeekStruct["ordinal_day_of_month"]
			weekDay, weekDayOk := daysOfWeekStruct["week_day"]
			if !ordinalDayOfMonthOk && !weekDayOk {
				return nil, errors.New("the block `days_of_week` requires specifying both `ordinal_day_of_month` and `week_day`")
			}
			daysOfWeekInput.WeekDay = synthetics.SyntheticsMonitorDowntimeWeekDays(weekDay.(string))
			daysOfWeekInput.OrdinalDayOfMonth = synthetics.SyntheticsMonitorDowntimeDayOfMonthOrdinal(ordinalDayOfMonth.(string))
			frequencyInput.DaysOfWeek = daysOfWeekInput
		}
		monthlyMonitorDowntimeArgumentsMap["frequency"] = frequencyInput
	}

	return monthlyMonitorDowntimeArgumentsMap, nil
}

func getValuesOfWeeklyMonitorDowntimeArguments(d *schema.ResourceData) (map[string]interface{}, error) {
	var weeklyMonitorDowntimeArgumentsMap map[string]interface{}

	dailyMonitorDowntimeArgumentsMap, err := getValuesOfDailyMonitorDowntimeArguments(d)
	if err != nil {
		return nil, err
	}

	maps.Copy(weeklyMonitorDowntimeArgumentsMap, dailyMonitorDowntimeArgumentsMap)

	// mandatory argument
	maintenanceDays, ok := d.GetOk("maintenance_days")
	if !ok {
		return nil, errors.New("`maintenance_days` is a required argument with weekly monitor downtime")
	} else {
		weeklyMonitorDowntimeArgumentsMap["maintenance_days"] = maintenanceDays.([]string)
	}
	return weeklyMonitorDowntimeArgumentsMap, nil
}

func getValuesOfDailyMonitorDowntimeArguments(d *schema.ResourceData) (map[string]interface{}, error) {
	var dailyMonitorDowntimeArgumentsMap map[string]interface{}

	monitorGUIDs, err := getMonitorGUIDs(d)
	if err != nil {
		return nil, err
	}

	dailyMonitorDowntimeArgumentsMap["monitor_guids"] = monitorGUIDs

	endRepeat, ok := d.GetOk("end_repeat")
	if ok {
		endRepeatStruct := endRepeat.(map[string]interface{})
		var endRepeatInput synthetics.SyntheticsDateWindowEndConfig
		onDate, onDateOk := endRepeatStruct["on_date"]
		onRepeat, onRepeatOk := endRepeatStruct["on_repeat"]

		if !onDateOk && !onRepeatOk {
			return nil, errors.New("the block `end_repeat` requires one of `on_date` or `on_repeat` to be specified")
		} else if onDateOk && onRepeatOk {
			return nil, errors.New("the block `end_repeat` requires only one of `on_date` or `on_repeat` to be specified, both cannot be specified")
		}

		endRepeatInput.OnDate = onDate.(synthetics.Date)
		endRepeatInput.OnRepeat = onRepeat.(int)
		dailyMonitorDowntimeArgumentsMap["end_repeat"] = endRepeatInput

	} else {
		dailyMonitorDowntimeArgumentsMap["end_repeat"] = nil
	}

	return dailyMonitorDowntimeArgumentsMap, nil

}

func getMonitorGUIDs(d *schema.ResourceData) ([]string, error) {
	val, ok := d.GetOk("monitor_guids")
	if ok {
		if val.([]string) == nil || len(val.([]string)) == 0 {
			return nil, errors.New("invalid specification of monitor GUIDs: empty list received in the argument 'monitor_guids'")
		} else {
			return val.([]string), nil
		}
	}
	return nil, nil
}

func getValuesOfRequiredArguments(d *schema.ResourceData) (map[string]string, error) {
	var requiredArgumentsMap map[string]string
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
	// TODO: WRITE THE CREATE METHOD

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	requiredArgumentsMap, err := getValuesOfRequiredArguments(d)
	if err != nil {
		return diag.FromErr(err)
	}

	switch requiredArgumentsMap["mode"] {
	case "ONCE":
		break
	case "DAILY":
		break
	case "WEEKLY":
		break
	case "MONTHLY":
		break
	default:
		return diag.FromErr(errors.New("invalid mode of operation: 'mode' can be 'ONCE', 'DAILY', 'WEEKLY' or 'MONTHLY'"))

	}

	createAccountInput := accountmanagement.AccountManagementCreateInput{
		Name:       d.Get("name").(string),
		RegionCode: d.Get("region").(string),
	}
	created, err := client.AccountManagement.AccountManagementCreateAccount(createAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: Account not created. Please check the input details")
	}
	accountID := created.ManagedAccount.ID

	d.SetId(strconv.Itoa(accountID))
	return resourceNewRelicAccountRead(ctx, d, meta)
}

func resourceNewRelicMonitorDowntimeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: WRITE THE READ METHOD

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		account, err := getCreatedAccountByID(client, d.Id())
		//		fmt.Println("read", account.ID, err)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if account == nil {
			return resource.RetryableError(fmt.Errorf("account not found"))
		}
		_ = d.Set("region", account.RegionCode)
		_ = d.Set("name", account.Name)

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}
	return nil
}

func resourceNewRelicMonitorDowntimeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: WRITE THE UPDATE METHOD

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	updateAccountInput := accountmanagement.AccountManagementUpdateInput{
		Name: d.Get("name").(string),
		ID:   accountID,
	}
	updated, err := client.AccountManagement.AccountManagementUpdateAccount(updateAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		return diag.Errorf("err: Account not Updated. Please check the input details")
	}

	return resourceNewRelicAccountRead(ctx, d, meta)
}

func resourceNewRelicMonitorDowntimeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: WRITE THE DELETE METHOD

	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Account cannot be deleted via Terraform. https://docs.newrelic.com/docs/apis/nerdgraph/examples/manage-accounts-nerdgraph/#delete",
	})
	return diags
}
