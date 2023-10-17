//go:build integration
// +build integration

package newrelic

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestFlattenSchedule(t *testing.T) {
	t.Parallel()

	timestamp, _ := time.Parse(time.RFC3339, "2021-01-21T15:30:00+08:00")

	repeat := alerts.MutingRuleScheduleRepeat("WEEKLY")

	mockMutingRuleSchedule := alerts.MutingRuleSchedule{
		StartTime: &timestamp,
		EndTime:   &timestamp,
		TimeZone:  "America/Los_Angeles",
		Repeat:    &repeat,
		EndRepeat: &timestamp,
		WeeklyRepeatDays: &[]alerts.DayOfWeek{
			"MONDAY",
			"TUESDAY",
		},
	}

	mockScheduleConfig := map[string]interface{}{
		"start_time": "2021-01-21T15:30:00",
		"end_time":   "2021-01-21T15:30:00",
		"end_repeat": "2021-01-21T15:30:00",
		"time_zone":  "America/Los_Angeles",
		"repeat":     "WEEKLY",
		"weekly_repeat_days": []string{
			"MONDAY",
			"TUESDAY",
		},
	}

	result := flattenSchedule(&mockMutingRuleSchedule)

	require.Equal(t, []interface{}{mockScheduleConfig}, result)
}

func TestFlattenSchedule_EmptyDaysOfWeekWithWeeklyRepeat(t *testing.T) {
	// Flatten should send an empty slice for weekly_repeat_days if repeat is set to WEEKLY
	t.Parallel()

	timestamp, _ := time.Parse(time.RFC3339, "2021-01-21T15:30:00+08:00")

	repeat := alerts.MutingRuleScheduleRepeat("WEEKLY")

	mockMutingRuleSchedule := alerts.MutingRuleSchedule{
		StartTime:        &timestamp,
		EndTime:          &timestamp,
		TimeZone:         "America/Los_Angeles",
		Repeat:           &repeat,
		EndRepeat:        &timestamp,
		WeeklyRepeatDays: &[]alerts.DayOfWeek{},
	}

	mockScheduleConfig := map[string]interface{}{
		"start_time":         "2021-01-21T15:30:00",
		"end_time":           "2021-01-21T15:30:00",
		"end_repeat":         "2021-01-21T15:30:00",
		"time_zone":          "America/Los_Angeles",
		"repeat":             "WEEKLY",
		"weekly_repeat_days": []string{},
	}

	result := flattenSchedule(&mockMutingRuleSchedule)

	require.Equal(t, []interface{}{mockScheduleConfig}, result)
}

func TestFlattenSchedule_NilWeeklyRepeatDaysWeeklyRepeat(t *testing.T) {
	// Flatten should not null out weekly_repeat_days if repeat is set to WEEKLY

	t.Parallel()

	timestamp, _ := time.Parse(time.RFC3339, "2021-01-21T15:30:00+08:00")

	repeat := alerts.MutingRuleScheduleRepeat("WEEKLY")

	mockMutingRuleSchedule := alerts.MutingRuleSchedule{
		StartTime:        &timestamp,
		EndTime:          &timestamp,
		TimeZone:         "America/Los_Angeles",
		Repeat:           &repeat,
		EndRepeat:        &timestamp,
		WeeklyRepeatDays: nil,
	}

	mockScheduleConfig := map[string]interface{}{
		"start_time":         "2021-01-21T15:30:00",
		"end_time":           "2021-01-21T15:30:00",
		"end_repeat":         "2021-01-21T15:30:00",
		"time_zone":          "America/Los_Angeles",
		"repeat":             "WEEKLY",
		"weekly_repeat_days": []string{},
	}

	result := flattenSchedule(&mockMutingRuleSchedule)

	require.Equal(t, []interface{}{mockScheduleConfig}, result)
}

func TestFlattenSchedule_NilWeeklyRepeatDaysDailyRepeat(t *testing.T) {
	// Flatten should null out weekly_repeat_days if repeat is set to DAILY

	t.Parallel()

	timestamp, _ := time.Parse(time.RFC3339, "2021-01-21T15:30:00+08:00")

	repeat := alerts.MutingRuleScheduleRepeat("DAILY")

	mockMutingRuleSchedule := alerts.MutingRuleSchedule{
		StartTime:        &timestamp,
		EndTime:          &timestamp,
		TimeZone:         "America/Los_Angeles",
		Repeat:           &repeat,
		EndRepeat:        &timestamp,
		WeeklyRepeatDays: nil,
	}

	mockScheduleConfig := map[string]interface{}{
		"start_time":         "2021-01-21T15:30:00",
		"end_time":           "2021-01-21T15:30:00",
		"end_repeat":         "2021-01-21T15:30:00",
		"time_zone":          "America/Los_Angeles",
		"repeat":             "DAILY",
		"weekly_repeat_days": nil,
	}

	result := flattenSchedule(&mockMutingRuleSchedule)

	require.Equal(t, []interface{}{mockScheduleConfig}, result)
}

func TestExpandScheduleUpdate_Basic(t *testing.T) {
	t.Parallel()
	ts, _ := time.Parse("2006-01-02T15:04:05", "2021-01-21T15:30:00")
	timestamp := alerts.NaiveDateTime{Time: ts}
	timeZone := "America/Los_Angeles"
	repeat := alerts.MutingRuleScheduleRepeatTypes.WEEKLY

	testSchema := &schema.Set{F: schema.HashString}
	testSchema.Add("MONDAY")
	testSchema.Add("TUESDAY")

	mockScheduleConfig := map[string]interface{}{
		"start_time":         "2021-01-21T15:30:00",
		"end_time":           "2021-01-21T15:30:00",
		"end_repeat":         "2021-01-21T15:30:00",
		"time_zone":          "America/Los_Angeles",
		"repeat":             "WEEKLY",
		"weekly_repeat_days": testSchema,
	}

	result, _ := expandMutingRuleUpdateSchedule(mockScheduleConfig)

	expected := alerts.MutingRuleScheduleUpdateInput{
		StartTime: &timestamp,
		EndTime:   &timestamp,
		TimeZone:  &timeZone,
		Repeat:    &repeat,
		EndRepeat: &timestamp,
		WeeklyRepeatDays: &[]alerts.DayOfWeek{
			"TUESDAY",
			"MONDAY",
		},
	}

	require.Equal(t, expected, result)

}

func TestExpandScheduleCreate_Basic(t *testing.T) {
	t.Parallel()
	ts, _ := time.Parse("2006-01-02T15:04:05", "2021-01-21T15:30:00")
	timestamp := alerts.NaiveDateTime{Time: ts}
	timeZone := "America/Los_Angeles"
	repeat := alerts.MutingRuleScheduleRepeatTypes.WEEKLY

	testSchema := &schema.Set{F: schema.HashString}
	testSchema.Add("MONDAY")
	testSchema.Add("TUESDAY")

	mockScheduleConfig := map[string]interface{}{
		"start_time":         "2021-01-21T15:30:00",
		"end_time":           "2021-01-21T15:30:00",
		"end_repeat":         "2021-01-21T15:30:00",
		"time_zone":          "America/Los_Angeles",
		"repeat":             "WEEKLY",
		"weekly_repeat_days": testSchema,
	}

	result, _ := expandMutingRuleCreateSchedule(mockScheduleConfig)

	expected := alerts.MutingRuleScheduleCreateInput{
		StartTime: &timestamp,
		EndTime:   &timestamp,
		TimeZone:  timeZone,
		Repeat:    &repeat,
		EndRepeat: &timestamp,
		WeeklyRepeatDays: &[]alerts.DayOfWeek{
			"TUESDAY",
			"MONDAY",
		},
	}

	require.Equal(t, expected, result)

}

func TestExpandScheduleUpdate_EmptyFields(t *testing.T) {
	// similar to Basic, but assert that empty ("") fields are converted to nil
	t.Parallel()
	ts, _ := time.Parse("2006-01-02T15:04:05", "2021-01-21T15:30:00")
	timestamp := alerts.NaiveDateTime{Time: ts}
	timeZone := "America/Los_Angeles"
	repeat := alerts.MutingRuleScheduleRepeatTypes.WEEKLY

	testSchema := &schema.Set{F: schema.HashString}
	testSchema.Add("MONDAY")
	testSchema.Add("TUESDAY")

	mockScheduleConfig := map[string]interface{}{
		"start_time":         "2021-01-21T15:30:00",
		"end_time":           "",
		"end_repeat":         "",
		"time_zone":          "America/Los_Angeles",
		"repeat":             "WEEKLY",
		"weekly_repeat_days": testSchema,
	}

	result, _ := expandMutingRuleUpdateSchedule(mockScheduleConfig)

	expected := alerts.MutingRuleScheduleUpdateInput{
		StartTime: &timestamp,
		EndTime:   nil,
		TimeZone:  &timeZone,
		Repeat:    &repeat,
		EndRepeat: nil,
		WeeklyRepeatDays: &[]alerts.DayOfWeek{
			"TUESDAY",
			"MONDAY",
		},
	}

	require.Equal(t, expected, result)
}

func TestExpandScheduleCreate_EmptyFields(t *testing.T) {
	// similar to Basic, but assert that empty ("") and omitted fields are left out

	t.Parallel()
	ts, _ := time.Parse("2006-01-02T15:04:05", "2021-01-21T15:30:00")
	timestamp := alerts.NaiveDateTime{Time: ts}
	timeZone := "America/Los_Angeles"

	mockScheduleConfig := map[string]interface{}{
		"start_time": "2021-01-21T15:30:00",
		"end_time":   "2021-01-21T15:30:00",
		"end_repeat": "",
		"time_zone":  "America/Los_Angeles",
	}

	result, _ := expandMutingRuleCreateSchedule(mockScheduleConfig)

	expected := alerts.MutingRuleScheduleCreateInput{
		StartTime: &timestamp,
		EndTime:   &timestamp,
		TimeZone:  timeZone,
	}

	require.Equal(t, expected, result)

}

func TestExpandScheduleCreate_EmptyWeeklyRepeat(t *testing.T) {
	// similar to Basic, but assert that we can pass through an explicit empty slice of days

	t.Parallel()
	ts, _ := time.Parse("2006-01-02T15:04:05", "2021-01-21T15:30:00")
	timestamp := alerts.NaiveDateTime{Time: ts}
	timeZone := "America/Los_Angeles"
	repeat := alerts.MutingRuleScheduleRepeatTypes.WEEKLY

	testSchema := &schema.Set{F: schema.HashString}

	mockScheduleConfig := map[string]interface{}{
		"start_time":         "2021-01-21T15:30:00",
		"time_zone":          "America/Los_Angeles",
		"repeat":             "WEEKLY",
		"weekly_repeat_days": testSchema,
	}

	result, _ := expandMutingRuleCreateSchedule(mockScheduleConfig)

	expected := alerts.MutingRuleScheduleCreateInput{
		StartTime:        &timestamp,
		EndTime:          nil,
		TimeZone:         timeZone,
		Repeat:           &repeat,
		EndRepeat:        nil,
		WeeklyRepeatDays: &[]alerts.DayOfWeek{},
	}

	require.Equal(t, expected, result)

}

func TestExpandScheduleUpdate_EmptyWeeklyRepeat(t *testing.T) {
	// similar to Basic, but assert that we can pass through an explicit empty slice of days

	t.Parallel()
	ts, _ := time.Parse("2006-01-02T15:04:05", "2021-01-21T15:30:00")
	timestamp := alerts.NaiveDateTime{Time: ts}
	timeZone := "America/Los_Angeles"
	repeat := alerts.MutingRuleScheduleRepeatTypes.WEEKLY

	testSchema := &schema.Set{F: schema.HashString}

	mockScheduleConfig := map[string]interface{}{
		"start_time":         "2021-01-21T15:30:00",
		"time_zone":          "America/Los_Angeles",
		"repeat":             "WEEKLY",
		"weekly_repeat_days": testSchema,
	}

	result, _ := expandMutingRuleUpdateSchedule(mockScheduleConfig)

	expected := alerts.MutingRuleScheduleUpdateInput{
		StartTime:        &timestamp,
		EndTime:          nil,
		TimeZone:         &timeZone,
		Repeat:           &repeat,
		EndRepeat:        nil,
		WeeklyRepeatDays: &[]alerts.DayOfWeek{},
	}

	require.Equal(t, expected, result)

}

func TestCreate_ErrorHandlingForGraphQLErrors(t *testing.T) {
	client := alerts.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})

	mutingRuleCreateInput := alerts.MutingRuleCreateInput{
		Name:        "Some Muting Rule",
		Description: "Some Description",
		Enabled:     true,
		Condition: alerts.MutingRuleConditionGroup{
			Conditions: []alerts.MutingRuleCondition{
				{
					Attribute: "conditionName",
					Operator:  "INVALID_OPERATOR",
					Values:    []string{"Some Value"},
				},
			},
			Operator: "AND",
		},
	}

	// in the below function, the first return value is not needed as the muting rule would
	// fail to be created since an invalid operator is specified in one of the conditions
	_, err := client.CreateMutingRule(testAccountID, mutingRuleCreateInput)

	if err != nil {
		diags := handleMutingRuleCreateUpdateError(err)
		require.Equal(t, diags.HasError(), true)
	}
}

func TestCreate_ErrorHandlingForNonGraphQLErrors(t *testing.T) {
	diags := handleMutingRuleCreateUpdateError(errors.New("This is a non-GraphQL custom error"))
	require.Equal(t, diags.HasError(), true)
}
