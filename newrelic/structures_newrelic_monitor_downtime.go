package newrelic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
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

type SyntheticsMonitorDowntimeOneTimeCreateInput struct {
	AccountID    int
	Name         string
	StartTime    synthetics.NaiveDateTime
	EndTime      synthetics.NaiveDateTime
	Timezone     string
	MonitorGUIDs []synthetics.EntityGUID
}

type SyntheticsMonitorDowntimeDailyCreateInput struct {
	SyntheticsMonitorDowntimeOneTimeCreateInput
	EndRepeat synthetics.SyntheticsDateWindowEndConfig
}

type SyntheticsMonitorDowntimeWeeklyCreateInput struct {
	SyntheticsMonitorDowntimeDailyCreateInput
	MaintenanceDays []synthetics.SyntheticsMonitorDowntimeWeekDays
}

type SyntheticsMonitorDowntimeMonthlyCreateInput struct {
	SyntheticsMonitorDowntimeDailyCreateInput
	Frequency synthetics.SyntheticsMonitorDowntimeMonthlyFrequency
}
