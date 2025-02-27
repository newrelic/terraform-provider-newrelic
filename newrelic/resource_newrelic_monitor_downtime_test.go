//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	tftest "github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

var resourceName = fmt.Sprintf("tf-test-%s", acctest.RandString(5))
var testMonitorGUIDsForCreate = testHelpersReturnTestSyntheticMonitorGUIDsForDowntime("create")
var testMonitorGUIDsForUpdate = testHelpersReturnTestSyntheticMonitorGUIDsForDowntime("update")

func TestAccNewRelicMonitorDowntime_Once(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("%s-once", resourceName)
	testAccountID, testAPIKey, _ := testHelpersFetchAPIKeyAndAccountID()

	terraformOptions := testHelpersSetTerraformOptions(
		t,
		"./integration_test_configuration/resource_newrelic_monitor_downtime",
		map[string]interface{}{
			"NEW_RELIC_ACCOUNT_ID": testAccountID,
			"NEW_RELIC_API_KEY":    testAPIKey,
			"NEW_RELIC_REGION":     "US",
			"name":                 rName,
			"mode":                 SyntheticsMonitorDowntimeModes.OneTime,
			"monitor_guids":        testMonitorGUIDsForCreate,
			"start_time":           generateRandomStartTime(),
			"end_time":             generateRandomEndTime(),
			"time_zone":            generateRandomTimeZone(),
		})

	defer tftest.Destroy(t, terraformOptions)

	// create, and wait for terraform plan to not be empty
	tftest.InitAndApply(t, terraformOptions)
	testHelpersWaitUntilAndConfirmEmptyPlan(t, terraformOptions)

	// update, and wait for terraform plan to not be empty
	terraformOptions.Vars["monitor_guids"] = testMonitorGUIDsForUpdate
	tftest.Apply(t, terraformOptions)
	testHelpersWaitUntilAndConfirmEmptyPlan(t, terraformOptions)
}

func TestAccNewRelicMonitorDowntime_Daily_Revamp(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("%s-daily", resourceName)
	testAccountID, testAPIKey, _ := testHelpersFetchAPIKeyAndAccountID()

	terraformOptions := testHelpersSetTerraformOptions(
		t,
		testHelpersSetIntegrationTestTerraformConfigDirectory("resource_newrelic_monitor_downtime"),
		map[string]interface{}{
			"NEW_RELIC_ACCOUNT_ID": testAccountID,
			"NEW_RELIC_API_KEY":    testAPIKey,
			"NEW_RELIC_REGION":     "US",
			"name":                 rName,
			"mode":                 SyntheticsMonitorDowntimeModes.DAILY,
			"monitor_guids":        testMonitorGUIDsForCreate,
			"start_time":           generateRandomStartTime(),
			"end_time":             generateRandomEndTime(),
			"time_zone":            generateRandomTimeZone(),
			"include_end_repeat":   true,
			"end_repeat_on_repeat": 3,
		})

	defer tftest.Destroy(t, terraformOptions)

	// create, and wait for terraform plan to not be empty
	tftest.InitAndApply(t, terraformOptions)
	testHelpersWaitUntilAndConfirmEmptyPlan(t, terraformOptions)

	// update, and wait for terraform plan to not be empty
	terraformOptions.Vars["monitor_guids"] = testMonitorGUIDsForUpdate
	terraformOptions.Vars["end_repeat_on_repeat"] = -1
	terraformOptions.Vars["end_repeat_on_date"] = generateRandomEndRepeatDate()

	tftest.Apply(t, terraformOptions)
	testHelpersWaitUntilAndConfirmEmptyPlan(t, terraformOptions)
}

func TestAccNewRelicMonitorDowntime_Weekly(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("%s-weekly", resourceName)
	testAccountID, testAPIKey, _ := testHelpersFetchAPIKeyAndAccountID()

	terraformOptions := testHelpersSetTerraformOptions(
		t,
		"./integration_test_configuration/resource_newrelic_monitor_downtime",
		map[string]interface{}{
			"NEW_RELIC_ACCOUNT_ID": testAccountID,
			"NEW_RELIC_API_KEY":    testAPIKey,
			"NEW_RELIC_REGION":     "US",
			"name":                 rName,
			"mode":                 SyntheticsMonitorDowntimeModes.WEEKLY,
			"monitor_guids":        testMonitorGUIDsForCreate,
			"start_time":           generateRandomStartTime(),
			"end_time":             generateRandomEndTime(),
			"time_zone":            generateRandomTimeZone(),
			"maintenance_days":     generateRandomMaintenanceDays(),
		})

	defer tftest.Destroy(t, terraformOptions)

	// create, and wait for terraform plan to not be empty
	tftest.InitAndApply(t, terraformOptions)
	testHelpersWaitUntilAndConfirmEmptyPlan(t, terraformOptions)

	// update, and wait for terraform plan to not be empty
	terraformOptions.Vars["monitor_guids"] = testMonitorGUIDsForUpdate
	terraformOptions.Vars["maintenance_days"] = generateRandomMaintenanceDays()
	terraformOptions.Vars["end_repeat_on_repeat"] = -1
	terraformOptions.Vars["end_repeat_on_date"] = generateRandomEndRepeatDate()

	tftest.Apply(t, terraformOptions)
	testHelpersWaitUntilAndConfirmEmptyPlan(t, terraformOptions)
}

func TestAccNewRelicMonitorDowntime_Monthly_Revamped(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("%s-weekly", resourceName)
	testAccountID, testAPIKey, _ := testHelpersFetchAPIKeyAndAccountID()

	terraformOptions := testHelpersSetTerraformOptions(
		t,
		"./integration_test_configuration/resource_newrelic_monitor_downtime",
		map[string]interface{}{
			"NEW_RELIC_ACCOUNT_ID":    testAccountID,
			"NEW_RELIC_API_KEY":       testAPIKey,
			"NEW_RELIC_REGION":        "US",
			"name":                    rName,
			"mode":                    SyntheticsMonitorDowntimeModes.MONTHLY,
			"monitor_guids":           testMonitorGUIDsForCreate,
			"start_time":              generateRandomStartTime(),
			"end_time":                generateRandomEndTime(),
			"time_zone":               generateRandomTimeZone(),
			"include_frequency":       true,
			"frequency_days_of_month": []int{5, 15, 25},
			"include_end_repeat":      true,
			"end_repeat_on_date":      generateRandomEndRepeatDate(),
		})

	defer tftest.Destroy(t, terraformOptions)

	// create, and wait for terraform plan to not be empty
	tftest.InitAndApply(t, terraformOptions)
	testHelpersWaitUntilAndConfirmEmptyPlan(t, terraformOptions)

	// update, and wait for terraform plan to not be empty
	terraformOptions.Vars["monitor_guids"] = testMonitorGUIDsForUpdate
	terraformOptions.Vars["include_end_repeat"] = false
	terraformOptions.Vars["frequency_days_of_month"] = []int{}
	terraformOptions.Vars["days_of_week_ordinal_day_of_month"] = "SECOND"
	terraformOptions.Vars["days_of_week_week_day"] = "SATURDAY"

	tftest.Apply(t, terraformOptions)
	testHelpersWaitUntilAndConfirmEmptyPlan(t, terraformOptions)
}

// helpers for all the tests written above
var fewValidTimeZones = []string{
	"Asia/Kolkata",
	"America/Los_Angeles",
	"Europe/Madrid",
	"Asia/Tokyo",
	"America/Vancouver",
	"Asia/Tel_Aviv",
	"Europe/Dublin",
	"Asia/Tashkent",
	"Europe/London",
	"Asia/Riyadh",
	"America/Chicago",
	"Australia/Sydney",
}

// monitors with the below GUIDs belong to the v2 Integration Tests Account
var monitorGUIDsAsString = convertStringListToString([]string{
	"MzgwNjUyNnxTWU5USHxNT05JVE9SfDQ3ZWI5YmYzLWRiOTEtNDljYy04MzM2LTBhZWJhNTE5MzhiOQ",
})
var monitorGUIDsUpdatedAsString = convertStringListToString([]string{
	"MzgwNjUyNnxTWU5USHxNT05JVE9SfDkwODEwNTRhLWRhYTAtNGI0Mi05YmIwLTY3M2M1MDI2ZWYyOA",
	"MzgwNjUyNnxTWU5USHxNT05JVE9SfGEyNWRmOTIwLTcxYjUtNDlmYy1iZTgzLTBhOGE0NjdiYWNhMg",
	"MzgwNjUyNnxTWU5USHxNT05JVE9SfDAzNjQ0ZDNlLTg0YzMtNDQyMC1hYjM4LTc0ZjBjODI4NTk3ZA",
})

func convertStringListToString(list []string) string {
	return fmt.Sprintf("[\"%s\"]", strings.Join(list, "\", \""))
}
func generateRandomTimeZone() string {
	rand.Seed(time.Now().Unix())
	return fewValidTimeZones[rand.Intn(len(fewValidTimeZones))]
}

func generateRandomStartTime() string {
	rand.Seed(time.Now().Unix())

	now := time.Now()
	hourLater := now.Add(time.Hour * 2)

	return hourLater.Format("2006-01-02T15:04:05")
}

func generateRandomEndTime() string {
	rand.Seed(time.Now().Unix())

	now := time.Now()

	// "5 +" to make sure end_time exceeds start_time by a minimum of 5 days
	randomDays := 5 + rand.Intn(25)
	daysLater := now.AddDate(0, 0, randomDays)

	return daysLater.Format("2006-01-02T15:04:05")
}

func generateRandomEndRepeatDate() string {
	rand.Seed(time.Now().Unix())

	now := time.Now()

	// "31 +" so that end_repeat > on_date can succeed the date in endTime by 30 days - endRepeat needs to be after endTime
	randomDays := 31 + rand.Intn(30)
	daysLater := now.AddDate(0, 0, randomDays)

	return daysLater.Format("2006-01-02")
}

func generateRandomMaintenanceDays() []string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	originalList := listSyntheticsMonitorDowntimeValidMaintenanceDays()

	// minimum 1, maximum 3
	randomLength := 1 + r.Intn(3)
	newList := make([]string, randomLength)

	for i := range newList {
		randomIndex := r.Intn(len(originalList))
		newList[i] = originalList[randomIndex]
	}

	return newList
}
