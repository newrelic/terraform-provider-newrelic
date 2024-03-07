//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

var resourceName = fmt.Sprintf("tf-test-%s", acctest.RandString(5))

// TestAccNewRelicMonitorDowntime_Once tests create, update operations of a one time monitor downtime
func TestAccNewRelicMonitorDowntime_Once(t *testing.T) {
	rName := fmt.Sprintf("%s-once", resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccCheckNewRelicMonitorDowntime_OnceConfiguration(rName, SyntheticsMonitorDowntimeModes.OneTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Update
			{
				Config: testAccCheckNewRelicMonitorDowntime_OnceConfigurationUpdated(rName, SyntheticsMonitorDowntimeModes.OneTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Import
			{
				ResourceName: "newrelic_monitor_downtime.foo",
				ImportState:  true,
			},
		},
	})
}

// TestAccNewRelicMonitorDowntime_Daily tests create, update operations of a daily monitor downtime
func TestAccNewRelicMonitorDowntime_Daily(t *testing.T) {
	rName := fmt.Sprintf("%s-daily", resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccCheckNewRelicMonitorDowntime_DailyConfiguration(rName, SyntheticsMonitorDowntimeModes.DAILY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Update
			{
				Config: testAccCheckNewRelicMonitorDowntime_DailyConfigurationUpdated(rName, SyntheticsMonitorDowntimeModes.DAILY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Import
			{
				ResourceName: "newrelic_monitor_downtime.foo",
				ImportState:  true,
			},
		},
	})
}

// TestAccNewRelicMonitorDowntime_Weekly tests create, update operations of a weekly monitor downtime
func TestAccNewRelicMonitorDowntime_Weekly(t *testing.T) {
	rName := fmt.Sprintf("%s-weekly", resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccCheckNewRelicMonitorDowntime_WeeklyConfiguration(rName, SyntheticsMonitorDowntimeModes.WEEKLY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Update
			{
				Config: testAccCheckNewRelicMonitorDowntime_WeeklyConfigurationUpdated(rName, SyntheticsMonitorDowntimeModes.WEEKLY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Import
			{
				ResourceName: "newrelic_monitor_downtime.foo",
				ImportState:  true,
			},
		},
	})
}

// TestAccNewRelicMonitorDowntime_Monthly tests create, update operations of a monthly monitor downtime
func TestAccNewRelicMonitorDowntime_Monthly(t *testing.T) {
	rName := fmt.Sprintf("%s-monthly", resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccCheckNewRelicMonitorDowntime_MonthlyConfiguration(rName, SyntheticsMonitorDowntimeModes.MONTHLY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Update
			{
				Config: testAccCheckNewRelicMonitorDowntime_MonthlyConfigurationUpdated(rName, SyntheticsMonitorDowntimeModes.MONTHLY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Import
			{
				ResourceName: "newrelic_monitor_downtime.foo",
				ImportState:  true,
			},
		},
	})
}

// TestAccNewRelicMonitorDowntime_MultiMode tests creating a monthly downtime and updating it as a daily downtime
func TestAccNewRelicMonitorDowntime_MultiMode(t *testing.T) {
	rName := fmt.Sprintf("%s-multimode", resourceName)
	resource.ParallelTest(t, resource.TestCase{
		//PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccCheckNewRelicMonitorDowntime_MonthlyConfiguration(rName, SyntheticsMonitorDowntimeModes.MONTHLY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Update
			{
				Config: testAccCheckNewRelicMonitorDowntime_DailyConfigurationUpdated(rName, SyntheticsMonitorDowntimeModes.DAILY),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicMonitorDowntimeExists("newrelic_monitor_downtime.foo"),
				),
			},
			// Import
			{
				ResourceName: "newrelic_monitor_downtime.foo",
				ImportState:  true,
			},
		},
	})
}

// TestAccNewRelicMonitorDowntime_Once_IncorrectConfig tests the creation of a once time monitor downtime with incorrect configuration
func TestAccNewRelicMonitorDowntime_Once_IncorrectConfig(t *testing.T) {
	rName := fmt.Sprintf("%s-once-incorrect", resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckNewRelicMonitorDowntime_DailyConfiguration(rName, SyntheticsMonitorDowntimeModes.OneTime),
				ExpectError: regexp.MustCompile(`validation errors`),
			},
		},
	})
}

// TestAccNewRelicMonitorDowntime_Weekly_IncorrectConfig tests the creation of a weekly monitor downtime with incorrect configuration
func TestAccNewRelicMonitorDowntime_Weekly_IncorrectConfig(t *testing.T) {
	rName := fmt.Sprintf("%s-weekly-incorrect", resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckNewRelicMonitorDowntime_MonthlyConfiguration(rName, SyntheticsMonitorDowntimeModes.WEEKLY),
				ExpectError: regexp.MustCompile(`validation errors`),
			},
		},
	})
}

// TestAccNewRelicMonitorDowntime_Monthly_IncorrectConfig tests the creation of a monthly monitor downtime with incorrect configuration
func TestAccNewRelicMonitorDowntime_Monthly_IncorrectConfig(t *testing.T) {
	rName := fmt.Sprintf("%s-monthly-incorrect", resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckNewRelicMonitorDowntime_WeeklyConfiguration(rName, SyntheticsMonitorDowntimeModes.MONTHLY),
				ExpectError: regexp.MustCompile(`validation errors`),
			},
		},
	})
}

func testAccCheckNewRelicMonitorDowntimeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no monitor downtime ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		// discarding this sleep duration, as this seems to be breaking integration tests
		// time.Sleep(20 * time.Second)

		found, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		x := (*found).(*entities.GenericEntity)
		if x.GUID != common.EntityGUID(rs.Primary.ID) {
			return fmt.Errorf("monitor downtime not found")
		}

		return nil
	}
}

func testAccCheckNewRelicMonitorDowntime_BaseConfiguration(name string, monitorGUIDsConverted string, mode string) string {
	return `
	  	name = "` + name + `"
      	monitor_guids = ` + monitorGUIDsConverted + `
  		mode       = "` + mode + `"
  		start_time = "` + generateRandomStartTime() + `"
  		end_time   = "` + generateRandomEndTime() + `"
  		time_zone  = "` + generateRandomTimeZone() + `"
	`
}
func testAccCheckNewRelicMonitorDowntime_OnceConfiguration(name string, mode string) string {
	return `
		resource "newrelic_monitor_downtime" "foo" {
  		` + testAccCheckNewRelicMonitorDowntime_BaseConfiguration(name, monitorGUIDsAsString, mode) + `
		}
	`
}

func testAccCheckNewRelicMonitorDowntime_OnceConfigurationUpdated(name string, mode string) string {
	return `
		resource "newrelic_monitor_downtime" "foo" {
  		` + testAccCheckNewRelicMonitorDowntime_BaseConfiguration(fmt.Sprintf("%s-updated", name), monitorGUIDsUpdatedAsString, mode) + `
		}
	`
}

func testAccCheckNewRelicMonitorDowntime_DailyConfiguration(name string, mode string) string {
	return `
		resource "newrelic_monitor_downtime" "foo" {
  		` + testAccCheckNewRelicMonitorDowntime_BaseConfiguration(name, monitorGUIDsAsString, mode) + `
			end_repeat {
				on_repeat = 3
			}
		}
	`
}

func testAccCheckNewRelicMonitorDowntime_DailyConfigurationUpdated(name string, mode string) string {
	return `
		resource "newrelic_monitor_downtime" "foo" {
  		` + testAccCheckNewRelicMonitorDowntime_BaseConfiguration(fmt.Sprintf("%s-updated", name), monitorGUIDsUpdatedAsString, mode) + `
			end_repeat {
				on_date = "` + generateRandomEndRepeatDate() + `"
			}
		}
	`
}

func testAccCheckNewRelicMonitorDowntime_WeeklyConfiguration(name string, mode string) string {
	return `
		resource "newrelic_monitor_downtime" "foo" {
  		` + testAccCheckNewRelicMonitorDowntime_BaseConfiguration(name, monitorGUIDsAsString, mode) + `
			maintenance_days = ` + generateRandomMaintenanceDays() + `
		}
	`
}

func testAccCheckNewRelicMonitorDowntime_WeeklyConfigurationUpdated(name string, mode string) string {
	return `
		resource "newrelic_monitor_downtime" "foo" {
  		` + testAccCheckNewRelicMonitorDowntime_BaseConfiguration(fmt.Sprintf("%s-updated", name), monitorGUIDsUpdatedAsString, mode) + `
			maintenance_days = ` + generateRandomMaintenanceDays() + `
			end_repeat {
				on_date = "` + generateRandomEndRepeatDate() + `"
			}
		}
	`
}

func testAccCheckNewRelicMonitorDowntime_MonthlyConfiguration(name string, mode string) string {
	return `
		resource "newrelic_monitor_downtime" "foo" {
  		` + testAccCheckNewRelicMonitorDowntime_BaseConfiguration(name, monitorGUIDsAsString, mode) + `
			end_repeat {
				on_date = "` + generateRandomEndRepeatDate() + `"
			}
			frequency {
				days_of_month = [5, 15, 25]
			}
		}
	`
}

func testAccCheckNewRelicMonitorDowntime_MonthlyConfigurationUpdated(name string, mode string) string {
	return `
		resource "newrelic_monitor_downtime" "foo" {
  		` + testAccCheckNewRelicMonitorDowntime_BaseConfiguration(fmt.Sprintf("%s-updated", name), monitorGUIDsUpdatedAsString, mode) + `
			frequency {
				days_of_week {
					ordinal_day_of_month = "SECOND"
					week_day = "SATURDAY"
				}
			}
		}
	`
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
	"MzgwNjUyNnxTWU5USHxNT05JVE9SfGFmZmM0MTRiLTVhNmMtNGI5NS1iMzYwLThhNmQ2ZTkzOTM3Yw",
})
var monitorGUIDsUpdatedAsString = convertStringListToString([]string{
	"MzgwNjUyNnxTWU5USHxNT05JVE9SfDhkYmMyYmIwLTQwZjgtNDA5NC05OTA1LTdhZGE2ZGViMmEwNg",
	"MzgwNjUyNnxTWU5USHxNT05JVE9SfDViZDNmYTk4LTA2NjgtNGQ1Yy05ODU2LTk3MzlmNWViY2JlNg",
	"MzgwNjUyNnxTWU5USHxNT05JVE9SfGFmZmM0MTRiLTVhNmMtNGI5NS1iMzYwLThhNmQ2ZTkzOTM3Yw",
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

func generateRandomMaintenanceDays() string {
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

	return convertStringListToString(newList)
}
