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

func TestNewRelicNotificationChannelWebhook_Basic(t *testing.T) {
	resourceName := "newrelic_notification_channel.webhook_test_foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationChannelConfigByType(rName, "WEBHOOK", "IINT", "b1e90a32-23b7-4028-b2c7-ffbdfe103852", `{
					key = "payload"
					value = "{\n\t\"id\": \"test\"\n}"
					label = "Payload Template"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testNewRelicNotificationChannelConfigByType(rName, "WEBHOOK", "IINT", "b1e90a32-23b7-4028-b2c7-ffbdfe103852", `{
					key = "payload"
					value = "{\n\t\"id\": \"test-update\"\n}"
					label = "Payload Template Update"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName),
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

func TestNewRelicNotificationChannelEmail_Basic(t *testing.T) {
	resourceName := "newrelic_notification_channel.email_test_foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationChannelConfigByType(rName, "EMAIL", "IINT", "0115e01f-5636-496e-947f-6ce0322d7c5d", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testNewRelicNotificationChannelConfigByType(rName, "EMAIL", "IINT", "0115e01f-5636-496e-947f-6ce0322d7c5d", `{
					key = "subject"
					value = "Update: {{ issueTitle }}"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName),
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

func testAccNewRelicNotificationChannelDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_notification_channel" {
			continue
		}

		var accountID int
		id := r.Primary.ID
		accountID = providerConfig.AccountID
		filters := ai.AiNotificationsChannelFilter{
			ID: id,
		}
		sorter := notifications.AiNotificationsChannelSorter{}

		_, err := client.Notifications.GetChannels(accountID, "", filters, sorter)
		if err == nil {
			return fmt.Errorf("notification channel still exists")
		}

	}
	return nil
}

func testNewRelicNotificationChannelConfigByType(name string, channelType string, product string, destinationId string, properties string) string {
	if properties == "" {
		return fmt.Sprintf(`
		resource "newrelic_notification_channel" "test_foo" {
			name = "%s"
			type = "%s"
			product = "%s"
			destination_id = "%s"
		}
	`, name, channelType, product, destinationId)
	}

	return fmt.Sprintf(`
		resource "newrelic_notification_channel" "test_foo" {
			name = "%s"
			type = "%s"
			product = "%s"
			destination_id = "%s"
			property %s
		}
	`, name, channelType, product, destinationId, properties)
}

func testAccCheckNewRelicNotificationChannelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no channel ID is set")
		}

		var accountID int
		id := rs.Primary.ID
		accountID = providerConfig.AccountID
		filters := ai.AiNotificationsChannelFilter{
			ID: id,
		}
		sorter := notifications.AiNotificationsChannelSorter{}

		found, err := client.Notifications.GetChannels(accountID, "", filters, sorter)
		if err != nil {
			return err
		}

		if string(found.Entities[0].ID) != rs.Primary.ID {
			return fmt.Errorf("channel not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}
