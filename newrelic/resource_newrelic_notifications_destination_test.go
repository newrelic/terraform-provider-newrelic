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
	// t.Skip("Skipping TestNewRelicNotificationDestinationWebhook_Basic. AWAITING FINAL IMPLEMENTATION!")

	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationDestinationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
		},
	})
}

func testNewRelicNotificationDestinationConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	name = "%s"
	type = "WEBHOOK"

	properties {
		key = "url"
		value = "https://webhook.site/"
	}

	auth = {
		type = "BASIC"
		user = "username"
		password = "password"
	}
}
`, name)
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

		resp, err := client.Notifications.GetDestinations(accountID, "", filters, sorter)

		fmt.Print("\n\n **************************** \n")
		fmt.Printf("\n DestinationDestroy:  %+v \n", *resp)
		fmt.Print("\n **************************** \n\n")

		if len(resp.Entities) > 0 {
			return fmt.Errorf("notification destination still exists")
		}

		if err != nil {
			return err
		}
	}
	return nil
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
