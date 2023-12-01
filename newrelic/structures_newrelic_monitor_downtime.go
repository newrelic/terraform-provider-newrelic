package newrelic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	"golang.org/x/exp/slices"
)

// validate functions
func listValidOrdinalDayOfMonthValues() []string {
	return []string{
		string(synthetics.SyntheticsMonitorDowntimeDayOfMonthOrdinalTypes.FIRST),
		string(synthetics.SyntheticsMonitorDowntimeDayOfMonthOrdinalTypes.SECOND),
		string(synthetics.SyntheticsMonitorDowntimeDayOfMonthOrdinalTypes.THIRD),
		string(synthetics.SyntheticsMonitorDowntimeDayOfMonthOrdinalTypes.FOURTH),
		string(synthetics.SyntheticsMonitorDowntimeDayOfMonthOrdinalTypes.LAST),
	}
}

func listValidWeekDayValues() []string {
	return []string{
		string(synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.SUNDAY),
		string(synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.MONDAY),
		string(synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.TUESDAY),
		string(synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.WEDNESDAY),
		string(synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.THURSDAY),
		string(synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.FRIDAY),
		string(synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.SATURDAY),
	}
}

var syntheticsMonitorDowntimeMaintenanceDaysMap = map[string]synthetics.SyntheticsMonitorDowntimeWeekDays{
	"SUNDAY":    synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.SUNDAY,
	"MONDAY":    synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.MONDAY,
	"TUESDAY":   synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.TUESDAY,
	"WEDNESDAY": synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.WEDNESDAY,
	"THURSDAY":  synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.THURSDAY,
	"FRIDAY":    synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.FRIDAY,
	"SATURDAY":  synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.SATURDAY,
}

func listSyntheticsMonitorDowntimeValidMaintenanceDays() []string {
	keys := make([]string, 0, len(syntheticsMonitorDowntimeMaintenanceDaysMap))

	for k := range syntheticsMonitorDowntimeMaintenanceDaysMap {
		keys = append(keys, k)
	}

	return keys
}

// #################
// Validate functions
// #################

func validateMonitorDowntimeAttributes(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	var errorsList []string

	err := validateMonitorDowntimeEndRepeatStructure(d)
	if err != nil {
		errorsList = append(errorsList, err.Error())
	}

	err = validateMonitorDowntimeMaintenanceDaysStructure(d)
	if err != nil {
		errorsList = append(errorsList, err.Error())
	}

	err = validateMonitorDowntimeFrequencyStructure(d)
	if err != nil {
		errorsList = append(errorsList, err.Error())
	}

	err = validateMonitorDowntimeStartTimeEndTime(d)
	if err != nil {
		errorsList = append(errorsList, err.Error())
	}

	if len(errorsList) == 0 {
		return nil
	}

	errorsString := "the following validation errors have been identified: \n"

	for index, val := range errorsList {
		errorsString += fmt.Sprintf("(%d): %s\n", index+1, val)
	}

	return errors.New(errorsString)
}

func validateMonitorDowntimeStartTimeEndTime(d *schema.ResourceDiff) error {
	_, startTimeObtained := d.GetChange("start_time")
	_, endTimeObtained := d.GetChange("end_time")

	startTime, _ := time.Parse("2006-01-02T15:04:05", startTimeObtained.(string))
	endTime, _ := time.Parse("2006-01-02T15:04:05", endTimeObtained.(string))

	if endTime.Before(startTime) {
		return errors.New("`end_time` cannot be before `start_time`")
	}

	return nil
}

func validateMonitorDowntimeFrequencyStructure(d *schema.ResourceDiff) error {
	_, mode := d.GetChange("mode")
	_, frequencyObtained := d.GetChange("frequency")
	frequency := frequencyObtained.([]interface{})

	if mode != SyntheticsMonitorDowntimeModes.MONTHLY && len(frequency) > 0 {
		return errors.New("the argument `frequency` may only be used with the 'MONTHLY' mode")
	} else if mode == SyntheticsMonitorDowntimeModes.MONTHLY && len(frequency) == 0 {
		return errors.New("the argument `frequency` is mandatory to be specified with the 'MONTHLY' mode")
	}

	frequencyDaysOfMonth, frequencyDaysOfMonthOk := d.GetOkExists("frequency.0.days_of_month")
	if frequencyDaysOfMonthOk {
		for _, val := range frequencyDaysOfMonth.(*schema.Set).List() {
			if val.(int) < 1 || val.(int) > 31 {
				return errors.New("all `days_of_month` values need to be in the range of 1 and 31")
			}
		}
	}

	return nil
}
func validateMonitorDowntimeMaintenanceDaysStructure(d *schema.ResourceDiff) error {
	_, mode := d.GetChange("mode")
	_, maintenanceDaysObtained := d.GetChange("maintenance_days")
	maintenanceDays := maintenanceDaysObtained.(*schema.Set)

	if mode != SyntheticsMonitorDowntimeModes.WEEKLY && maintenanceDays.Len() > 0 {
		return errors.New("the argument `maintenance_days` may only be used with the 'WEEKLY' mode")
	} else if mode == SyntheticsMonitorDowntimeModes.WEEKLY && maintenanceDays.Len() == 0 {
		return errors.New("the argument `maintenance_days` is mandatory to be specified with the 'WEEKLY' mode")
	}

	listOfValidMaintenanceDays := listSyntheticsMonitorDowntimeValidMaintenanceDays()
	for _, val := range maintenanceDays.List() {
		isValidMaintenanceDay := false
		for _, day := range listOfValidMaintenanceDays {
			if day == val {
				isValidMaintenanceDay = true
			}
		}
		if isValidMaintenanceDay == false {
			return errors.New(fmt.Sprintf("%s is not an accepted value for maintenance_days; the acceptable list of values is %v", val, listOfValidMaintenanceDays))
		}
	}

	return nil
}

func validateMonitorDowntimeEndRepeatStructure(d *schema.ResourceDiff) error {
	_, mode := d.GetChange("mode")
	_, endRepeatObtained := d.GetChange("end_repeat")
	endRepeat := endRepeatObtained.([]interface{})
	// validModesWithEndRepeat := []string{"DAILY", "MONTHLY", "WEEKLY"}

	if len(endRepeat) != 0 && mode == SyntheticsMonitorDowntimeModes.ONE_TIME {
		return errors.New("the argument `end_repeat` may only be used with the modes `DAILY`, `MONTHLY` and `WEEKLY`")
	}

	return nil
}

func validateMonitorDowntimeTimeZone(val interface{}, key string) (warns []string, errs []error) {
	timezone := val.(string)
	_, err := time.LoadLocation(timezone)
	if err != nil {
		errs = append(errs, err)
	}

	return warns, errs
}

func validateMonitorDowntimeOnDate(val interface{}, key string) (warns []string, errs []error) {
	valueString := val.(string)
	_, err := time.Parse("2006-01-02", valueString)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid `on_date` %s: the attribute `on_date` needs to be in the format 'YYYY-MM-DD'", valueString))
	}
	return warns, errs
}

var SyntheticsMonitorDowntimeModes = struct {
	ONE_TIME string
	DAILY    string
	MONTHLY  string
	WEEKLY   string
}{
	ONE_TIME: "ONE_TIME",
	DAILY:    "DAILY",
	MONTHLY:  "MONTHLY",
	WEEKLY:   "WEEKLY",
}

// #################
// Classes whose objects would be used in create/update requests
// #################

type SyntheticsMonitorDowntimeCommonArgumentsInput struct {
	AccountID    int
	Name         string
	Mode         string
	StartTime    synthetics.NaiveDateTime
	EndTime      synthetics.NaiveDateTime
	Timezone     string
	MonitorGUIDs []synthetics.EntityGUID
}

type SyntheticsMonitorDowntimeOneTimeInput struct {
	SyntheticsMonitorDowntimeCommonArgumentsInput
}

type SyntheticsMonitorDowntimeDailyInput struct {
	SyntheticsMonitorDowntimeCommonArgumentsInput
	EndRepeat synthetics.SyntheticsDateWindowEndConfig
}

type SyntheticsMonitorDowntimeWeeklyInput struct {
	SyntheticsMonitorDowntimeDailyInput
	MaintenanceDays []synthetics.SyntheticsMonitorDowntimeWeekDays
}

type SyntheticsMonitorDowntimeMonthlyInput struct {
	SyntheticsMonitorDowntimeDailyInput
	Frequency synthetics.SyntheticsMonitorDowntimeMonthlyFrequency
}

// #################
// GET functions used to fetch values from the configuration
// #################

func getMonitorDowntimeValuesOfCommonArguments(d *schema.ResourceData) (*SyntheticsMonitorDowntimeCommonArgumentsInput, error) {
	commonArgumentsObject := &SyntheticsMonitorDowntimeCommonArgumentsInput{}

	accountID, err := getMonitorDowntimeAccountIDFromConfiguration(d)
	if err != nil {
		return nil, err
	}

	name, err := getMonitorDowntimeNameFromConfiguration(d)
	if err != nil {
		return nil, err
	}

	mode, err := getMonitorDowntimeModeFromConfiguration(d)
	if err != nil {
		return nil, err
	}

	startTime, err := getMonitorDowntimeStartTimeFromConfiguration(d)
	if err != nil {
		return nil, err
	}

	endTime, err := getMonitorDowntimeEndTimeFromConfiguration(d)
	if err != nil {
		return nil, err
	}

	timezone, err := getMonitorDowntimeTimezoneFromConfiguration(d)
	if err != nil {
		return nil, err
	}

	monitorGUIDs, err := getMonitorDowntimeMonitorGUIDsFromConfiguration(d)
	if err != nil {
		return nil, err
	}

	commonArgumentsObject.AccountID = accountID
	commonArgumentsObject.Name = name
	commonArgumentsObject.Mode = mode
	commonArgumentsObject.StartTime = startTime
	commonArgumentsObject.EndTime = endTime
	commonArgumentsObject.Timezone = timezone
	commonArgumentsObject.MonitorGUIDs = monitorGUIDs

	return commonArgumentsObject, nil
}

func getMonitorDowntimeAccountIDFromConfiguration(d *schema.ResourceData) (int, error) {
	val, ok := d.GetOk("account_id")
	if ok {
		if val.(string) == "" {
			return 0, errors.New(fmt.Sprintf("%s has value \"\"", `account_id`))
		} else {
			accountIdAsInteger, err := strconv.Atoi(val.(string))
			if err != nil {
				return 0, err
			}
			return accountIdAsInteger, nil
		}
	} else {
		accountIdAsInteger, err := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID"))
		if err != nil {
			return 0, err
		}
		return accountIdAsInteger, nil
	}
}

func getMonitorDowntimeNameFromConfiguration(d *schema.ResourceData) (string, error) {
	val, ok := d.GetOk("name")
	if ok {
		if val.(string) == "" {
			return "", errors.New(fmt.Sprintf("%s has value \"\"", `name`))
		} else {
			return val.(string), nil
		}
	} else {
		return "", errors.New(fmt.Sprintf(" value of argument %s not specified", `name`))
	}
}

func getMonitorDowntimeModeFromConfiguration(d *schema.ResourceData) (string, error) {
	val, ok := d.GetOk("mode")
	if ok {
		if val.(string) == "" {
			return "", errors.New(fmt.Sprintf("%s has value \"\"", `mode`))
		} else {
			return val.(string), nil
		}
	} else {
		return "", errors.New(fmt.Sprintf(" value of argument %s not specified", `mode`))
	}
}

func getMonitorDowntimeStartTimeFromConfiguration(d *schema.ResourceData) (synthetics.NaiveDateTime, error) {
	val, ok := d.GetOk("start_time")
	if ok {
		if val.(string) == "" {
			return "", errors.New(fmt.Sprintf("%s has value \"\"", `start_time`))
		} else {
			return synthetics.NaiveDateTime(val.(string)), nil
		}
	} else {
		return "", errors.New(fmt.Sprintf(" value of argument %s not specified", `start_time`))
	}
}

func getMonitorDowntimeEndTimeFromConfiguration(d *schema.ResourceData) (synthetics.NaiveDateTime, error) {
	val, ok := d.GetOk("end_time")
	if ok {
		if val.(string) == "" {
			return "", errors.New(fmt.Sprintf("%s has value \"\"", `end_time`))
		} else {
			return synthetics.NaiveDateTime(val.(string)), nil
		}
	} else {
		return "", errors.New(fmt.Sprintf(" value of argument %s not specified", `end_time`))
	}
}

func getMonitorDowntimeTimezoneFromConfiguration(d *schema.ResourceData) (string, error) {
	val, ok := d.GetOk("time_zone")
	if ok {
		if val.(string) == "" {
			return "", errors.New(fmt.Sprintf("%s has value \"\"", `time_zone`))
		} else {
			return val.(string), nil
		}
	} else {
		return "", errors.New(fmt.Sprintf(" value of argument %s not specified", `time_zone`))
	}
}

func getMonitorDowntimeMonitorGUIDsFromConfiguration(d *schema.ResourceData) ([]synthetics.EntityGUID, error) {
	val, ok := d.GetOk("monitor_guids")
	if ok {
		in := val.(*schema.Set).List()
		out := make([]synthetics.EntityGUID, len(in))
		for i := range in {
			out[i] = synthetics.EntityGUID(in[i].(string))
		}
		if len(out) == 0 {
			return []synthetics.EntityGUID{}, nil
		}
		return out, nil
	}
	return []synthetics.EntityGUID{}, nil
}

// #################
// GET functions used by create methods
// #################

func getMonitorDowntimeOneTimeValues(d *schema.ResourceData, commonArgumentsObject *SyntheticsMonitorDowntimeCommonArgumentsInput) (*SyntheticsMonitorDowntimeOneTimeInput, error) {
	return &SyntheticsMonitorDowntimeOneTimeInput{
		SyntheticsMonitorDowntimeCommonArgumentsInput: *commonArgumentsObject,
	}, nil
}

func getMonitorDowntimeDailyValues(d *schema.ResourceData, commonArgumentsObject *SyntheticsMonitorDowntimeCommonArgumentsInput) (*SyntheticsMonitorDowntimeDailyInput, error) {
	monitorDowntimeDailyInput := &SyntheticsMonitorDowntimeDailyInput{
		SyntheticsMonitorDowntimeCommonArgumentsInput: *commonArgumentsObject,
	}

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
		monitorDowntimeDailyInput.EndRepeat = endRepeatInput

	} else {
		monitorDowntimeDailyInput.EndRepeat = synthetics.SyntheticsDateWindowEndConfig{}
	}

	return monitorDowntimeDailyInput, nil
}

func getMonitorDowntimeWeeklyValues(d *schema.ResourceData, commonArgumentsObject *SyntheticsMonitorDowntimeCommonArgumentsInput) (*SyntheticsMonitorDowntimeWeeklyInput, error) {
	monitorDowntimeDailyInput, err := getMonitorDowntimeDailyValues(d, commonArgumentsObject)
	if err != nil {
		return nil, err
	}

	monitorDowntimeWeeklyInput := &SyntheticsMonitorDowntimeWeeklyInput{
		SyntheticsMonitorDowntimeDailyInput: *monitorDowntimeDailyInput,
	}

	// mandatory argument
	listOfMaintenanceDaysInConfiguration, err := getMaintenanceDaysList(d)
	if err != nil {
		return nil, err
	}
	maintenanceDays, err := convertSyntheticsMonitorDowntimeMaintenanceDays(listOfMaintenanceDaysInConfiguration)
	if err != nil {
		return nil, err
	}
	monitorDowntimeWeeklyInput.MaintenanceDays = maintenanceDays

	return monitorDowntimeWeeklyInput, nil
}

func getMonitorDowntimeMonthlyValues(d *schema.ResourceData, commonArgumentsObject *SyntheticsMonitorDowntimeCommonArgumentsInput) (*SyntheticsMonitorDowntimeMonthlyInput, error) {
	monitorDowntimeDailyInput, err := getMonitorDowntimeDailyValues(d, commonArgumentsObject)
	if err != nil {
		return nil, err
	}

	monitorDowntimeMonthlyInput := &SyntheticsMonitorDowntimeMonthlyInput{
		SyntheticsMonitorDowntimeDailyInput: *monitorDowntimeDailyInput,
	}

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
		monitorDowntimeMonthlyInput.Frequency = frequencyInput
	}

	return monitorDowntimeMonthlyInput, nil
}

// #################
// Methods which assist create methods
// #################

func (obj *SyntheticsMonitorDowntimeOneTimeInput) createMonitorDowntimeOneTime(ctx context.Context, client *newrelic.NewRelic) (string, error) {
	resp, err := client.Synthetics.SyntheticsCreateOnceMonitorDowntimeWithContext(
		ctx,
		obj.AccountID,
		obj.EndTime,
		obj.MonitorGUIDs,
		obj.Name,
		obj.StartTime,
		obj.Timezone,
	)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("encountered an API error while trying to create a monitor downtime: nil response returned")
	}

	return string(resp.GUID), nil
}

func (obj *SyntheticsMonitorDowntimeDailyInput) createMonitorDowntimeDaily(ctx context.Context, client *newrelic.NewRelic) (string, error) {
	resp, err := client.Synthetics.SyntheticsCreateDailyMonitorDowntimeWithContext(
		ctx,
		obj.AccountID,
		obj.EndRepeat,
		obj.EndTime,
		obj.MonitorGUIDs,
		obj.Name,
		obj.StartTime,
		obj.Timezone,
	)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("encountered an API error while trying to create a monitor downtime: nil response returned")
	}

	return string(resp.GUID), nil
}

func (obj *SyntheticsMonitorDowntimeWeeklyInput) createMonitorDowntimeWeekly(ctx context.Context, client *newrelic.NewRelic) (string, error) {
	resp, err := client.Synthetics.SyntheticsCreateWeeklyMonitorDowntimeWithContext(
		ctx,
		obj.AccountID,
		obj.EndRepeat,
		obj.EndTime,
		obj.MaintenanceDays,
		obj.MonitorGUIDs,
		obj.Name,
		obj.StartTime,
		obj.Timezone,
	)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("encountered an API error while trying to create a monitor downtime: nil response returned")
	}
	return string(resp.GUID), nil
}

func (obj *SyntheticsMonitorDowntimeMonthlyInput) createMonitorDowntimeMonthly(ctx context.Context, client *newrelic.NewRelic) (string, error) {
	resp, err := client.Synthetics.SyntheticsCreateMonthlyMonitorDowntimeWithContext(
		ctx,
		obj.AccountID,
		obj.EndRepeat,
		obj.EndTime,
		obj.Frequency,
		obj.MonitorGUIDs,
		obj.Name,
		obj.StartTime,
		obj.Timezone,
	)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("encountered an API error while trying to create a monitor downtime: nil response returned")
	}

	return string(resp.GUID), nil
}

// #################
// Methods which assist update methods
// #################

func (obj *SyntheticsMonitorDowntimeOneTimeInput) updateMonitorDowntimeOneTime(ctx context.Context, client *newrelic.NewRelic, guid synthetics.EntityGUID) (string, error) {
	resp, err := client.Synthetics.SyntheticsEditMonitorDowntimeWithContext(
		ctx,
		synthetics.SyntheticsMonitorDowntimeDailyConfig{},
		guid,
		obj.MonitorGUIDs,
		synthetics.SyntheticsMonitorDowntimeMonthlyConfig{},
		obj.Name,
		synthetics.SyntheticsMonitorDowntimeOnceConfig{
			EndTime:   obj.EndTime,
			StartTime: obj.StartTime,
			Timezone:  obj.Timezone,
		},
		synthetics.SyntheticsMonitorDowntimeWeeklyConfig{},
	)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("encountered an API error while trying to create a monitor downtime: nil response returned")
	}

	return string(resp.GUID), nil
}

func (obj *SyntheticsMonitorDowntimeDailyInput) updateMonitorDowntimeDaily(ctx context.Context, client *newrelic.NewRelic, guid synthetics.EntityGUID) (string, error) {
	resp, err := client.Synthetics.SyntheticsEditMonitorDowntimeWithContext(
		ctx,
		synthetics.SyntheticsMonitorDowntimeDailyConfig{
			EndTime:   obj.EndTime,
			StartTime: obj.StartTime,
			Timezone:  obj.Timezone,
			EndRepeat: obj.EndRepeat,
		},
		guid,
		obj.MonitorGUIDs,
		synthetics.SyntheticsMonitorDowntimeMonthlyConfig{},
		obj.Name,
		synthetics.SyntheticsMonitorDowntimeOnceConfig{},
		synthetics.SyntheticsMonitorDowntimeWeeklyConfig{},
	)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("encountered an API error while trying to create a monitor downtime: nil response returned")
	}

	return string(resp.GUID), nil
}

func (obj *SyntheticsMonitorDowntimeWeeklyInput) updateMonitorDowntimeWeekly(ctx context.Context, client *newrelic.NewRelic, guid synthetics.EntityGUID) (string, error) {
	resp, err := client.Synthetics.SyntheticsEditMonitorDowntimeWithContext(
		ctx,
		synthetics.SyntheticsMonitorDowntimeDailyConfig{},
		guid,
		obj.MonitorGUIDs,
		synthetics.SyntheticsMonitorDowntimeMonthlyConfig{},
		obj.Name,
		synthetics.SyntheticsMonitorDowntimeOnceConfig{},
		synthetics.SyntheticsMonitorDowntimeWeeklyConfig{
			EndTime:         obj.EndTime,
			StartTime:       obj.StartTime,
			Timezone:        obj.Timezone,
			EndRepeat:       obj.EndRepeat,
			MaintenanceDays: obj.MaintenanceDays,
		},
	)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("encountered an API error while trying to create a monitor downtime: nil response returned")
	}
	return string(resp.GUID), nil
}

func (obj *SyntheticsMonitorDowntimeMonthlyInput) updateMonitorDowntimeMonthly(ctx context.Context, client *newrelic.NewRelic, guid synthetics.EntityGUID) (string, error) {
	resp, err := client.Synthetics.SyntheticsEditMonitorDowntimeWithContext(
		ctx,
		synthetics.SyntheticsMonitorDowntimeDailyConfig{},
		guid,
		obj.MonitorGUIDs,
		synthetics.SyntheticsMonitorDowntimeMonthlyConfig{
			EndTime:   obj.EndTime,
			StartTime: obj.StartTime,
			Timezone:  obj.Timezone,
			EndRepeat: obj.EndRepeat,
			Frequency: obj.Frequency,
		},
		obj.Name,
		synthetics.SyntheticsMonitorDowntimeOnceConfig{},
		synthetics.SyntheticsMonitorDowntimeWeeklyConfig{},
	)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("encountered an API error while trying to create a monitor downtime: nil response returned")
	}

	return string(resp.GUID), nil
}

// #################
// Other Helper Functions
// #################

func getMaintenanceDaysList(d *schema.ResourceData) ([]string, error) {
	val, ok := d.GetOk("maintenance_days")
	if !ok {
		return nil, errors.New("`maintenance_days` not found in the configuration")
	}
	if ok {
		in := val.(*schema.Set).List()
		out := make([]string, len(in))
		for i := range in {
			out[i] = in[i].(string)
		}
		if len(out) == 0 {
			return nil, errors.New("invalid specification: empty list received in the argument 'maintenance_days'")
		} else {
			return out, nil
		}
	}
	return nil, nil
}

func getFrequencyDaysOfMonthList(daysOfMonth []interface{}) []int {
	out := make([]int, len(daysOfMonth))
	for i := range out {
		out[i] = daysOfMonth[i].(int)
	}
	return out
}

var syntheticsMonitorDowntimeMaintenanceDaysAliasesMap = map[string]synthetics.SyntheticsMonitorDowntimeWeekDays{
	"SU": synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.SUNDAY,
	"MO": synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.MONDAY,
	"TU": synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.TUESDAY,
	"WE": synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.WEDNESDAY,
	"TH": synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.THURSDAY,
	"FR": synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.FRIDAY,
	"SA": synthetics.SyntheticsMonitorDowntimeWeekDaysTypes.SATURDAY,
}

func convertSyntheticsMonitorDowntimeMaintenanceDays(maintenanceDays []string) ([]synthetics.SyntheticsMonitorDowntimeWeekDays, error) {
	maintenanceDaysTypeCasted := make([]synthetics.SyntheticsMonitorDowntimeWeekDays, 0, len(maintenanceDays))
	listOfValidMaintenanceDays := listSyntheticsMonitorDowntimeValidMaintenanceDays()

	for index, value := range maintenanceDays {
		isValidDay := slices.Contains(listOfValidMaintenanceDays, value)
		if !isValidDay {
			return nil, errors.New(fmt.Sprintf("expected maintenance_days[%d] to be one of %v, got %s", index, listOfValidMaintenanceDays, value))
		} else {
			maintenanceDaysTypeCasted = append(maintenanceDaysTypeCasted, syntheticsMonitorDowntimeMaintenanceDaysMap[value])
		}
	}
	return maintenanceDaysTypeCasted, nil
}

func getStringEntityTag(tags []entities.EntityTag, attributeName string) string {
	out := []string{}
	for _, t := range tags {
		if t.Key == attributeName {
			for _, v := range t.Values {
				out = append(out, string(v))
			}
		}
	}
	return out[0]
}

// #################
// Functions which help set values in the state
// #################

func setMonitorDowntimeFrequency(d *schema.ResourceData, tags []entities.EntityTag) {
	daysOfMonthTag := "monthDay"
	daysOfWeekTag := "specificWeekDay"

	var daysOfMonthValue []int
	var daysOfWeekValue string = ""

	// TODO: is there a better way to do this? (not initialise with -1 but something like nil which I am not able to do)

	for _, t := range tags {
		if t.Key == daysOfMonthTag {
			for _, v := range t.Values {
				value, _ := strconv.Atoi(v)
				daysOfMonthValue = append(daysOfMonthValue, value)
			}
		} else if t.Key == daysOfWeekTag {
			for _, v := range t.Values {
				daysOfWeekValue = v
			}
		}
	}

	// if neither monthDay or specificWeekDay is found in the tags, it means frequency is absent in the configuration
	if len(daysOfMonthValue) == 0 && daysOfWeekValue == "" {
		return
	}

	if len(daysOfMonthValue) != 0 {
		_ = d.Set("frequency", []map[string]interface{}{
			{
				"days_of_month": daysOfMonthValue,
			},
		})
	}

	if daysOfWeekValue != "" {
		value := strings.Split(daysOfWeekValue, ":")
		ordinalDayOfMonth := value[0]
		weekDay := value[1]

		_ = d.Set("frequency", []map[string]interface{}{
			{
				"days_of_week": []map[string]interface{}{
					{
						"ordinal_day_of_month": ordinalDayOfMonth,
						"week_day":             weekDay,
					},
				},
			},
		})
	}
}

func setMonitorDowntimeMaintenanceDays(d *schema.ResourceData, tags []entities.EntityTag) {
	maintenanceDaysTag := "weekDay"
	var maintenanceDaysValue []string
	for _, t := range tags {
		if t.Key == maintenanceDaysTag {
			for _, v := range t.Values {
				val, ok := syntheticsMonitorDowntimeMaintenanceDaysAliasesMap[v]
				if ok {
					maintenanceDaysValue = append(maintenanceDaysValue, string(val))
				}
			}
		}
	}
	_ = d.Set("maintenance_days", maintenanceDaysValue)
}

func setMonitorDowntimeEndRepeat(d *schema.ResourceData, tags []entities.EntityTag, timezone string) {
	onDateTag := "endRepeat"
	onRepeatTag := "occurrences"

	onDateValue := ""

	// TODO: is there a better way to do this? (not initialise with -1 but something like nil which I am not able to do)
	onRepeatValue := -1

	for _, t := range tags {
		if t.Key == onDateTag {
			for _, v := range t.Values {
				onDateValue = v
			}
		} else if t.Key == onRepeatTag {
			for _, v := range t.Values {
				onRepeatValue, _ = strconv.Atoi(v)
			}
		}
	}

	// if neither onDate or onRepeat is found in the tags, it means endRepeat is absent in the configuration
	if onDateValue == "" && onRepeatValue == -1 {
		return
	}

	if onDateValue != "" {
		onDateValueParsed, _ := strconv.ParseInt(onDateValue, 10, 64)
		dt := time.Unix(onDateValueParsed/1000, 0)
		loc, _ := time.LoadLocation(timezone)
		_ = d.Set("end_repeat", []map[string]interface{}{
			{
				"on_date": dt.In(loc).Format("2006-01-02"),
			},
		})
	}

	if onRepeatValue != -1 {
		_ = d.Set("end_repeat", []map[string]interface{}{
			{
				"on_repeat": onRepeatValue,
			},
		})
	}

}

func setMonitorDowntimeAccountId(tags []entities.EntityTag) string {
	return getStringEntityTag(tags, "accountId")
}

func setMonitorDowntimeMode(tags []entities.EntityTag) string {
	return getStringEntityTag(tags, "type")
}

func setMonitorDowntimeTimezone(tags []entities.EntityTag) string {
	return getStringEntityTag(tags, "timezone")
}

func setMonitorDowntimeStartTime(tags []entities.EntityTag) string {
	startTime := getStringEntityTag(tags, "startTime")
	startTimeIntParsed, _ := strconv.ParseInt(startTime, 10, 64)
	timezone := getStringEntityTag(tags, "timezone")
	dt := time.Unix(startTimeIntParsed/1000, 0)
	loc, _ := time.LoadLocation(timezone)
	return dt.In(loc).Format("2006-01-02T15:04:05")
}

func setMonitorDowntimeEndTime(tags []entities.EntityTag) string {
	endTime := getStringEntityTag(tags, "endTime")
	endTimeIntParsed, _ := strconv.ParseInt(endTime, 10, 64)
	timezone := getStringEntityTag(tags, "timezone")
	dt := time.Unix(endTimeIntParsed/1000, 0)
	loc, _ := time.LoadLocation(timezone)
	return dt.In(loc).Format("2006-01-02T15:04:05")
}

func setMonitorDowntimeMonitorGUIDs(relationships []entities.EntityRelationship, monitorDowntimeID common.EntityGUID) []string {
	var listOfRelatedMonitorGUIDs []string
	for _, rel := range relationships {
		source := rel.Source
		target := rel.Target
		if source.GUID == monitorDowntimeID && target.EntityType == "SYNTHETIC_MONITOR_ENTITY" {
			listOfRelatedMonitorGUIDs = append(listOfRelatedMonitorGUIDs, string(target.GUID))
		}
	}

	return listOfRelatedMonitorGUIDs
}
