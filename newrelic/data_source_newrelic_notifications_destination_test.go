//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicNotificationDestinationDataSource_Id(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNotificationsDestinationDataSourceConfigById(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicNotificationDestination("data.newrelic_notification_destination.foo"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "WEBHOOK"),
				),
			},
		},
	})
}

func TestAccNewRelicNotificationDestinationDataSource_Name(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnvVars(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNotificationsDestinationDataSourceConfigByName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicNotificationDestination("data.newrelic_notification_destination.foo"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "WEBHOOK"),
				),
			},
		},
	})
}

func testAccNewRelicNotificationsDestinationDataSourceConfigById(name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	name = "%s"
	type = "WEBHOOK"
	active = true

	property {
		key = "url"
		value = "https://webhook.site/"
	}
}

data "newrelic_notification_destination" "foo" {
	id = newrelic_notification_destination.foo.id
}
`, name)
}

func testAccNewRelicNotificationsDestinationDataSourceConfigByName(name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	name = "%[1]s"
	type = "WEBHOOK"
	active = true

	property {
		key = "url"
		value = "https://webhook.site/"
	}
}
data "newrelic_notification_destination" "foo" {
	name = "%[1]s"
}
`, name)
}

func testAccNewRelicNotificationDestination(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		id := r.Primary.ID
		a := r.Primary.Attributes

		if id == "" {
			return fmt.Errorf("expected to get a notification destination id from New Relic")
		}

		if a["name"] == "" {
			return fmt.Errorf("expected to get a notification destination from New Relic")
		}

		return nil
	}
}
