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
	resourceName := "newrelic_notification_channel.test_foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	destinationId := "4756c466-c29f-4f89-9cb4-382cabfcef61"
	var channelID *string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationChannelConfigByType(rName, "WEBHOOK", "IINT", destinationId, `{
					key = "payload"
					value = "{\n\t\"id\": \"test\"\n}"
					label = "Payload Template"
				}`, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName, channelID),
				),
			},
			// Test: Update
			{
				Config: testNewRelicNotificationChannelConfigByType(rName, "WEBHOOK", "IINT", destinationId, `{
					key = "payload"
					value = "{\n\t\"id\": \"test-update\"\n}"
					label = "Payload Template Update"
				}`, channelID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName, channelID),
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
	resourceName := "newrelic_notification_channel.test_foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	destinationID := "d112de81-46be-4b52-959d-945448a64cc1"
	var channelID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationChannelConfigByType(rName, "EMAIL", "IINT", destinationID, "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName, channelID),
				),
			},
			// Test: Update
			{
				Config: testNewRelicNotificationChannelConfigByType(rName, "EMAIL", "IINT", destinationID, `{
					key = "subject"
					value = "Update: {{ issueTitle }}"
				}`, channelID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName, channelID),
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

		channelsResponse, _ := client.Notifications.GetChannels(accountID, "", filters, sorter)
		if len(channelsResponse.Entities) != 0 {
			return fmt.Errorf("notification channel still exists")
		}

	}
	return nil
}

func testNewRelicNotificationChannelConfigByType(name string, channelType string, product string, destinationId string, properties string, channelID string) string {
	if channelID != "" {
		if properties == "" {
			return fmt.Sprintf(`
				resource "newrelic_notification_channel" "test_foo" {
					id = "%s"
					name = "%s"
					type = "%s"
					product = "%s"
					destination_id = "%s"
				}
			`, channelID, name, channelType, product, destinationId)
		} else {
			return fmt.Sprintf(`
				resource "newrelic_notification_channel" "test_foo" {
					id = "%s"
					name = "%s"
					type = "%s"
					product = "%s"
					destination_id = "%s"
					properties %s
				}
			`, channelID, name, channelType, product, destinationId, properties)
		}
	}

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
			properties %s
		}
	`, name, channelType, product, destinationId, properties)
}

func testAccCheckNewRelicNotificationChannelExists(n string, channelID string) resource.TestCheckFunc {
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

		channelID = id

		return nil
	}
}
