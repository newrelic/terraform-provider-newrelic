// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicAlertChannel_Basic(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	rNameUpdated := fmt.Sprintf("tf-test-updated-%s", rand)
	rNameDeprecatedUpdated := fmt.Sprintf("tf-test-deprecated-updated-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Update (Migrate configuration)
			{
				Config: testAccNewRelicAlertChannelConfigUpdated(rNameDeprecatedUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameDeprecatedUpdated),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "config.0.recipients", "terraform-acctest+bar@hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "config.0.include_json_attachment", "0"),
				),
			},
			// Test: Create
			{
				Config: testAccNewRelicAlertChannelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "config.0.recipients", "terraform-acctest+foo@hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "config.0.include_json_attachment", "1"),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertChannelConfigUpdated(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "config.0.recipients", "terraform-acctest+bar@hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "config.0.include_json_attachment", "0"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicAlertChannel_Webhook(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "webhook", `{
					base_url = "http://www.test.com"
					payload_type = "application/json"
					payload = {
						"test": "value"
					}
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccNewRelicAlertChannel_WebhookPayloadHeaderStringConflicts(t *testing.T) {
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Headers attribute conflict
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "webhook", `{
					base_url = "http://www.test.com"
					payload_type = "application/json"
					payload = {
						test = "value"
					}
					headers_string = "{ \"test\": { \"some\": \"value\" } }"
					headers = {
						test = "value"
					}
				}`),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
			// Test: Payload attribute conflict
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "webhook", `{
					base_url = "http://www.test.com"
					payload_type = "application/json"
					payload_string = "{ \"test\": { \"some\": \"value\" } }"
					payload = {
						test = "value"
					}
				}`),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
		},
	})
}

func TestAccNewRelicAlertChannel_WebhookPayloadHeaderString(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "webhook", `{
					base_url = "http://www.test.com"
					payload_type = "application/json"
					payload_string = "{ \"test\": { \"some\": \"value\" } }"
					headers_string = "{ \"test\": { \"some\": \"value\" } }"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicAlertChannel_Slack(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "slack", `{
					url = "https://example.slack.com"
					channel = "example-channel"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "slack"),
					resource.TestCheckResourceAttr(resourceName, "config.0.url", "https://example.slack.com"),
					resource.TestCheckResourceAttr(resourceName, "config.0.channel", "example-channel"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config.0.url", // ignore sensitive data that's not returned from the API
				},
			},
		},
	})
}

func TestAccNewRelicAlertChannel_PagerDuty(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "pagerduty", `{
					service_key = "abc123"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "pagerduty"),
					resource.TestCheckResourceAttr(resourceName, "config.0.service_key", "abc123"),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "pagerduty", `{
					service_key = "abc321"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "pagerduty"),
					resource.TestCheckResourceAttr(resourceName, "config.0.service_key", "abc321"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config.0.service_key", // ignore sensitive data that's not returned from the API
				},
			},
		},
	})
}

func TestAccNewRelicAlertChannel_OpsGenie(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "opsgenie", `{
					api_key = "abc123"
					teams = "example-team"
					tags = "tag1"
					recipients = "example@somedomain.com"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "opsgenie"),
					resource.TestCheckResourceAttr(resourceName, "config.0.teams", "example-team"),
					resource.TestCheckResourceAttr(resourceName, "config.0.tags", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.recipients", "example@somedomain.com"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config.0.api_key", // ignore sensitive data that's not returned from the API
				},
			},
		},
	})
}

func TestAccNewRelicAlertChannel_VictorOps(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "victorops", `{
					key = "abc123"
					route_key = "/example-route"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "victorops"),
					resource.TestCheckResourceAttr(resourceName, "config.0.route_key", "/example-route"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config.0.key", // ignore sensitive data that's not returned from the API
				},
			},
		},
	})
}

func TestAccNewRelicAlertChannel_WebhookPayloadValidation(t *testing.T) {
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	expectedErrorMsg, _ := regexp.Compile(`payload_type is required when using payload`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create (expect error)
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "webhook", `{
					base_url = "https://example.com"
					payload = {
						"key" = "value"
					}
				}`),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func TestAccNewRelicAlertChannel_ResourceNotFound(t *testing.T) {
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertChannelConfig(rName),
			},
			{
				PreConfig: testAccDeleteAlertChannel(rName),
				Config:    testAccNewRelicAlertChannelConfig(rName),
			},
		},
	})
}

func testAccNewRelicAlertChannelConfig(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_alert_channel" "foo" {
		name = "%s"
		type = "email"

		config {
			recipients = "terraform-acctest+foo@hashicorp.com"
			include_json_attachment = "1"
		}
	}
`, name)
}

func testAccNewRelicAlertChannelConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_alert_channel" "foo" {
		name = "%s"
		type = "email"

		config {
			recipients = "terraform-acctest+bar@hashicorp.com"
			include_json_attachment = "0"
		}
	}
`, name)
}

func testAccNewRelicAlertChannelConfigByType(name string, channelType string, configuration string) string {
	return fmt.Sprintf(`
		resource "newrelic_alert_channel" "foo" {
			name = "%s"
			type = "%s"

			config %s
		}
	`, name, channelType, configuration)
}

func testAccCheckNewRelicAlertChannelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_channel" {
			continue
		}

		id, err := strconv.ParseInt(r.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		_, err = client.Alerts.GetChannel(int(id))

		if err == nil {
			return fmt.Errorf("alert channel still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicAlertChannelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no channel ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		found, err := client.Alerts.GetChannel(int(id))
		if err != nil {
			return err
		}

		if strconv.Itoa(found.ID) != rs.Primary.ID {
			return fmt.Errorf("channel not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccDeleteAlertChannel(name string) func() {
	return func() {
		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		alertChannels, _ := client.Alerts.ListChannels()

		for _, d := range alertChannels {
			if d.Name == name {
				_, _ = client.Alerts.DeleteChannel(d.ID)
				break
			}
		}
	}
}
