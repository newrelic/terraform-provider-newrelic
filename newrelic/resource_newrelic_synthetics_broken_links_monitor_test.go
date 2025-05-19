//go:build integration || SYNTHETICS
// +build integration SYNTHETICS

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func TestAccNewRelicSyntheticsBrokenLinksMonitor(t *testing.T) {
	resourceName := "newrelic_synthetics_broken_links_monitor.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitorConfig(
					rName,
					SyntheticsNodeRuntimeType,
					SyntheticsNodeNewRuntimeTypeVersion,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorEntityExists(resourceName),
				),
			},
			// Update
			{
				Config: testAccNewRelicSyntheticsBrokenLinksMonitorConfig(
					fmt.Sprintf("%s-updated", rName),
					SyntheticsNodeRuntimeType,
					SyntheticsNodeNewRuntimeTypeVersion,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsMonitorEntityExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"locations_private",
					"tag",
				},
			},
		},
	})
}

func testAccNewRelicSyntheticsBrokenLinksMonitorConfig(
	name string,
	runtimeType string,
	runtimeTypeVersion string,
) string {
	return fmt.Sprintf(`
		resource "newrelic_synthetics_broken_links_monitor" "foo" {
		  name             = "%[1]s"
		  period           = "EVERY_HOUR"
		  status           = "ENABLED"
		  locations_public = ["AP_SOUTH_1"]
		  uri              = "https://www.google.com"
		  tag {
			key    = "tf-test"
			values = ["tf-acc-test"]
		  }
		  %[2]s
		  %[3]s
}`,
		name,
		testConfigurationStringBuilder("runtime_type", runtimeType),
		testConfigurationStringBuilder("runtime_type_version", runtimeTypeVersion),
	)
}

func testAccCheckNewRelicSyntheticsMonitorEntityExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics monitor ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(60 * time.Second)

		result, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}
		if string((*result).GetGUID()) != rs.Primary.ID {
			return fmt.Errorf("the monitor is not found %v - %v", (*result).GetGUID(), rs.Primary.ID)
		}

		if rs.Primary.Attributes["monitor_id"] != string((*result).(*entities.SyntheticMonitorEntity).MonitorId) {
			return fmt.Errorf("the monitor id doesnot match, expected: %v", rs.Primary.Attributes["monitor_id"])
		}

		if rs.Primary.Attributes["runtime_type"] != "" && rs.Primary.Attributes["runtime_type_version"] != "" {
			runtimeTagsExist := false
			tags := (*result).GetTags()
			for _, t := range tags {
				if t.Key == "runtimeType" || t.Key == "runtimeTypeVersion" {
					runtimeTagsExist = true
				}
			}

			if runtimeTagsExist == false {
				return fmt.Errorf("runtimeType and runtimeTypeVersion not found in the entity fetched")
			}
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsMonitorResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_broken_links_monitor" {
			continue
		}

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(60 * time.Second)

		found, _ := client.Entities.GetEntity(common.EntityGUID(r.Primary.ID))
		if (*found) != nil {
			return fmt.Errorf("synthetics monitor still exists")
		}
	}
	return nil
}
