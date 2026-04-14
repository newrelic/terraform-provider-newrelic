//go:build integration || WORKFLOW_INTEGRATIONS

package newrelic

import (
	"fmt"
	"strconv"
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
				Config: testNewRelicNotificationDestinationConfig(testAccountID, rName, authAttr, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, fmt.Sprintf("%s-updated", rName), authAttr, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
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

func TestNewRelicNotificationDestination_CustomHeadersAuth(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	authAttr := `
auth_custom_header {
		key = "testKey1"
		value = "testValue1"
	}
auth_custom_header {
		key = "testKey2"
		value = "testValue2"
	}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, rName, authAttr, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, fmt.Sprintf("%s-updated", rName), authAttr, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auth_custom_header.0.value",
					"auth_custom_header.1.value",
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
				Config: testNewRelicNotificationDestinationConfig(testAccountID, rName, authAttr, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, fmt.Sprintf("%s-updated", rName), authAttr, ""),
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

func TestNewRelicNotificationDestination_secureURL(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	urlAttr := `secure_url {
		prefix = "https://webbhook.site/"
		secure_suffix = "test"
	}`
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
				Config: testNewRelicNotificationDestinationConfig(testAccountID, rName, authAttr, urlAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, fmt.Sprintf("%s-updated", rName), authAttr, urlAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auth_custom_header.0.value",
					"auth_token.0.token",
					"auth_basic.0.password",
					"secure_url.0.secure_suffix",
				},
			},
		},
	})
}

func TestNewRelicNotificationDestination_secureURL_update(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	urlAttr := `secure_url {
		prefix = "https://webbhook.site/"
		secure_suffix = "test"
	}`
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
				Config: testNewRelicNotificationDestinationConfig(testAccountID, rName, authAttr, urlAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Update
			{
				Config: testNewRelicNotificationDestinationConfig(testAccountID, fmt.Sprintf("%s-updated", rName), authAttr, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auth_custom_header.0.value",
					"auth_token.0.token",
					"auth_basic.0.password",
					"secure_url.0.secure_suffix",
				},
			},
		},
	})
}

// TODO: Uncomment when organization environment variables are available in GitHub Actions
func TestNewRelicNotificationDestination_OrganizationScope(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)
	orgUUID := "fb33fea3-4d7e-4736-9701-acb59a634fdf"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create with ORGANIZATION scope (requires UUID)
			{
				Config: testNewRelicNotificationDestinationConfigWithScope(rName, "ORGANIZATION", orgUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.type", "ORGANIZATION"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.id", orgUUID),
				),
			},
			// Update name only (scope remains unchanged)
			{
				Config: testNewRelicNotificationDestinationConfigWithScope(fmt.Sprintf("%s-updated", rName), "ORGANIZATION", orgUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.type", "ORGANIZATION"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.id", orgUUID),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"scope.0.id",
				},
			},
		},
	})
}

func TestNewRelicNotificationDestination_AccountScope(t *testing.T) {
	resourceName := "newrelic_notification_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicNotificationDestinationDestroy,
		Steps: []resource.TestStep{
			// Test: Create with ACCOUNT scope (requires number)
			{
				Config: testNewRelicNotificationDestinationConfigWithScope(rName, "ACCOUNT", strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.type", "ACCOUNT"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.id", strconv.Itoa(testAccountID)),
				),
			},
			// Update name only (scope remains unchanged)
			{
				Config: testNewRelicNotificationDestinationConfigWithScope(fmt.Sprintf("%s-updated", rName), "ACCOUNT", strconv.Itoa(testAccountID)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotificationDestinationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.type", "ACCOUNT"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.id", strconv.Itoa(testAccountID)),
				),
			},
			// Import
			// Note: Import uses backward-compatible flow (sets account_id instead of scope)
			// so we ignore both scope and account_id differences
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"account_id",
					"scope.#",
					"scope.0.%",
					"scope.0.type",
					"scope.0.id",
				},
			},
		},
	})
}

func testNewRelicNotificationDestinationConfig(accountID int, name string, auth string, url string) string {
	var usedUrl string
	if url == "" {
		usedUrl = `property {
		key = "url"
		value = "https://webhook.site/"
	}`
	} else {
		usedUrl = url
	}

	sprintf := fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	account_id = %[1]d
	name = "%[2]s"
	type = "WEBHOOK"
	active = true

	%[4]s

	property {
		key = "source"
		value = "terraform"
		label = "terraform-integration-test"
	}

	%[3]s
}
`, accountID, name, auth, usedUrl)
	return sprintf
}

func testNewRelicNotificationDestinationConfigWithScope(name string, scopeType string, scopeID string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	name = "%[1]s"
	type = "WEBHOOK"
	active = true

	property {
		key = "url"
		value = "https://webhook.site/"
	}

	scope {
		type = "%[2]s"
		id   = "%[3]s"
	}
}
`, name, scopeType, scopeID)
}

func testAccNewRelicNotificationDestinationDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_notification_destination" {
			continue
		}

		id := r.Primary.ID
		filters := ai.AiNotificationsDestinationFilter{
			ID: id,
		}

		var resp *notifications.AiNotificationsDestinationsResponse
		var err error

		scopeType := r.Primary.Attributes["scope.0.type"]
		if scopeType == "ORGANIZATION" {
			resp, err = client.Notifications.GetDestinationsOrganization(nil, &filters, nil)
		} else {
			accountID := providerConfig.AccountID
			resp, err = client.Notifications.GetDestinationsAccount(accountID, nil, &filters, nil)
		}

		if err != nil {
			return err
		}

		if resp != nil && len(resp.Entities) > 0 {
			return fmt.Errorf("notification destination still exists")
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

		id := rs.Primary.ID
		filters := ai.AiNotificationsDestinationFilter{
			ID: id,
		}

		var found *notifications.AiNotificationsDestinationsResponse
		var err error

		scopeType := rs.Primary.Attributes["scope.0.type"]
		if scopeType == "ORGANIZATION" {
			found, err = client.Notifications.GetDestinationsOrganization(nil, &filters, nil)
		} else {
			accountID := providerConfig.AccountID
			found, err = client.Notifications.GetDestinationsAccount(accountID, nil, &filters, nil)
		}

		if err != nil {
			return err
		}

		if found == nil || len(found.Entities) == 0 || string(found.Entities[0].ID) != rs.Primary.ID {
			return fmt.Errorf("destination not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}
