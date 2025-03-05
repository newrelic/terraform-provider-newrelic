//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestNewRelicNotificationChannel_Webhook(t *testing.T) {
	resourceName := "newrelic_notification_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	channelPropsAttr := `property {
		key = "payload"
		value = "{\n\t\"id\": \"test\"\n}"
		label = "Payload Template"
	}

	property {
		key = "url"
		value = "https://webhook.site/"
	}
	`
	destinationPropsAttr := `property {
		key = "url"
		value = "https://webhook.site/"
	}
	`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testNewRelicNotificationChannelConfig(
					testAccountID,
					rName,
					string(notifications.AiNotificationsChannelTypeTypes.WEBHOOK),
					channelPropsAttr,
					destinationPropsAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationChannelConfig(
					testAccountID,
					fmt.Sprintf("%s-updated", rName),
					string(notifications.AiNotificationsChannelTypeTypes.WEBHOOK),
					channelPropsAttr,
					destinationPropsAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestNewRelicNotificationChannel_WebhookPropertyError(t *testing.T) {
	t.Skipf("Skipping this test until we are sure on the property block that is expected to throw an error with a Webhook")
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	channelPropsAttr := `property {
		key = "payload"
		value = "{\n\t\"id\": \"test\"\n}"
		label = "Payload Template"
	}

	# Test error for missing property key = url
	`
	destinationPropsAttr := `
	`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testNewRelicNotificationChannelConfig(
					testAccountID,
					rName,
					string(notifications.AiNotificationsChannelTypeTypes.WEBHOOK),
					channelPropsAttr,
					destinationPropsAttr,
				),
				ExpectError: regexp.MustCompile(`Missing mandatory field: "Domain"`),
			},
		},
	})
}

func TestNewRelicNotificationChannel_Email(t *testing.T) {
	resourceName := "newrelic_notification_channel.foo"
	rand := acctest.RandString(6)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	channelPropsAttr := `property {
		key = "email"
		value = "no-reply+terraformtest@newrelic.com"
	}`
	destinationPropsAttr := `property {
		key = "subject"
		value = "some subject"
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testNewRelicNotificationChannelConfig(
					testAccountID,
					rName,
					string(notifications.AiNotificationsChannelTypeTypes.EMAIL),
					channelPropsAttr,
					destinationPropsAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationChannelConfig(
					testAccountID,
					fmt.Sprintf("%s-updated", rName),
					string(notifications.AiNotificationsChannelTypeTypes.EMAIL),
					channelPropsAttr,
					destinationPropsAttr,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationChannelExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestNewRelicNotificationChannel_EmailPropertyError(t *testing.T) {
	rand := acctest.RandString(6)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	propsAttr := `property {
		key = "invalid-for-email"
		value = "email-error-test"
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Test error scenario
			{
				Config: testNewRelicNotificationChannelConfig(
					testAccountID,
					rName,
					string(notifications.AiNotificationsChannelTypeTypes.EMAIL),
					propsAttr,
					propsAttr,
				),
				ExpectError: regexp.MustCompile(`Missing mandatory field: "Email"`),
			},
		},
	})
}

func TestNewRelicNotificationChannel_ImportChannel_WrongId(t *testing.T) {
	resourceName := "newrelic_notification_channel.foo"
	rand := acctest.RandString(6)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	channelPropsAttr := `property {
		key = "email"
		value = "no-reply+terraformtest@newrelic.com"
	}`
	destinationPropsAttr := `property {
		key = "subject"
		value = "some subject"
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationChannelDestroy,
		Steps: []resource.TestStep{
			// Import
			{
				Config: testNewRelicNotificationChannelConfig(
					testAccountID,
					rName,
					string(notifications.AiNotificationsChannelTypeTypes.EMAIL),
					channelPropsAttr,
					destinationPropsAttr,
				),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "d20a6505-dbe3-4484-82ef-7723cb17a7d2",
				ExpectError:       regexp.MustCompile("Error: Cannot import non-existent remote object"),
			},
		},
	})
}

func testNewRelicNotificationChannelConfig(accountID int, name string, notificationType string, channelProps string, destinationProps string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	account_id = %[1]d
	name = "destination-%[2]s"
	type = "%[3]s"

	auth_basic {
		user = "username"
		password = "password"
	}

	%[4]s
}

resource "newrelic_notification_channel" "foo" {
	account_id = newrelic_notification_destination.foo.account_id
	name = "%[2]s"
	type = "%[3]s"
	product = "IINT"
	destination_id = newrelic_notification_destination.foo.id

	%[4]s
}
`, accountID, name, notificationType, channelProps, destinationProps)
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

		if len(found.Entities) == 0 || string(found.Entities[0].ID) != rs.Primary.ID {
			return fmt.Errorf("channel not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}
