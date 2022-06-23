//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestNewRelicNotificationDestination_Webhook(t *testing.T) {
	resourceName := "newrelic_alert_channel.test_foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationDestinationByType(rName, "webhook", `{
					type = "BASIC"
					user = "test-user"
					password = "pass123"
				}`, `{
					key = "url"
					value = "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccNewRelicNotificationDestinationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	accountID := 10867072

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_notification_destination" {
			continue
		}

		id := r.Primary.ID

		_, err := client.Notifications.GetDestination(accountID, notifications.UUID(id))
		if err == nil {
			return fmt.Errorf("notification destination still exists")
		}

	}
	return nil
}

func testNewRelicNotificationDestinationByType(name string, channelType string, auth string, properties string) string {
	return fmt.Sprintf(`
		resource "newrelic_notification_destination" "test_foo" {
			name = "%s"
			type = "%s"
			auth = "%s"
			properties %s
		}
	`, name, channelType, auth, properties)
}

func testAccCheckNewRelicNotificationDestinationExists(n string) resource.TestCheckFunc {
	accountID := 10867072

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no destination ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		id := rs.Primary.ID

		found, err := client.Notifications.GetDestination(accountID, notifications.UUID(id))
		if err != nil {
			return err
		}

		if string(found.ID) != rs.Primary.ID {
			return fmt.Errorf("destination not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}
