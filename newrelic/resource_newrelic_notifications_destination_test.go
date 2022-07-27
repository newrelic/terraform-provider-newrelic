//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestNewRelicNotificationDestinationWebhook_Basic(t *testing.T) {
	resourceName := "newrelic_notification_destination.webhook_test_foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationDestinationConfigByType(rName, "WEBHOOK", `{
					type = "BASIC"
					user = "test-user"
					password = "pass123"
				}`, `{
					key = "url"
					value = "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testNewRelicNotificationDestinationConfigByType(rName, "WEBHOOK", `{
					type = "BASIC"
					user = "test-user-update"
					password = "pass12345678-update"
				}`, `{
					key = "url"
					value = "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func TestNewRelicNotificationDestinationEmail_Basic(t *testing.T) {
	resourceName := "newrelic_notification_destination.email_test_foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationDestinationConfigByType(rName, "EMAIL", "", `{
					key = "email"
					value = "email_test@gmail.com"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testNewRelicNotificationDestinationConfigByType(rName, "EMAIL", "", `{
					key = "email"
					value = "update_email_test@gmail.com"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func TestNewRelicNotificationDestinationPagerDuty_Basic(t *testing.T) {
	resourceName := "newrelic_notification_destination.pagerduty_test_foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationDestinationConfigByType(rName, "PAGERDUTY_SERVICE_INTEGRATION", `{
					type = "TOKEN"
					prefix = "Bearer"
					token  = "10567a689d984d03c021034b22a789e2"
				}`, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testNewRelicNotificationDestinationConfigByType(rName, "PAGERDUTY_SERVICE_INTEGRATION", `{
					type = "TOKEN"
					prefix = "Bearer"
					token  = "another-689d984d03c021034b22a789e2"
				}`, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func testAccNewRelicNotificationDestinationDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_notification_destination" {
			continue
		}

		var accountID int
		id := r.Primary.ID
		accountID = providerConfig.AccountID
		filters := ai.AiNotificationsDestinationFilter{
			ID: id,
		}
		sorter := notifications.AiNotificationsDestinationSorter{}

		_, err := client.Notifications.GetDestinations(accountID, "", filters, sorter)
		if err == nil {
			return fmt.Errorf("notification destination still exists")
		}

	}
	return nil
}

func testNewRelicNotificationDestinationConfigByType(name string, channelType string, auth string, properties string) string {
	if auth == "" {
		return fmt.Sprintf(`
		resource "newrelic_notification_destination" "test_foo" {
			name = "%s"
			type = "%s"
			properties %s
		}
	`, name, channelType, properties)
	}

	if properties == "" {
		return fmt.Sprintf(`
		resource "newrelic_notification_destination" "test_foo" {
			name = "%s"
			type = "%s"
			auth = "%s"
		}
	`, name, channelType, auth)
	}

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
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no destination ID is set")
		}

		var accountID int
		id := rs.Primary.ID
		accountID = providerConfig.AccountID
		filters := ai.AiNotificationsDestinationFilter{
			ID: id,
		}
		sorter := notifications.AiNotificationsDestinationSorter{}

		found, err := client.Notifications.GetDestinations(accountID, "", filters, sorter)
		if err != nil {
			return err
		}

		if string(found.Entities[0].ID) != rs.Primary.ID {
			return fmt.Errorf("destination not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}
