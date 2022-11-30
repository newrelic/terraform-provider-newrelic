//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestNewRelicNotificationDestination_BasicAuth(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	authAttr := `auth_basic {
		user = "username"
		password = "abc123"
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, rName, authAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, fmt.Sprintf("%s-updated", rName), authAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auth_token.0.token",
					"auth_basic.0.password",
				},
			},
		},
	})
}

func TestNewRelicNotificationDestination_TokenAuth(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	authAttr := `auth_token {
		prefix = "testprefix"
		token = "abc123"
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, rName, authAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, fmt.Sprintf("%s-updated", rName), authAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auth_token.0.token",
					"auth_basic.0.password",
				},
			},
		},
	})
}

func testNewRelicNotificationDestinationConfig(accountID int, name string, auth string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	account_id = %[1]d
	name = "%[2]s"
	type = "WEBHOOK"
	active = true

	property {
		key = "url"
		value = "https://webhook.site/"
	}

	property {
		key = "source"
		value = "terraform"
		label = "terraform-integration-test"
	}

	%[3]s
}
`, accountID, name, auth)
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

		// fmt.Print("\n\n **************************** \n")
		// fmt.Printf("\n DestinationDestroy:  %+v \n", toJSON(r.Primary.Attributes))
		// fmt.Print("\n **************************** \n\n")

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
