package newrelic

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicAlertChannel_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertChannelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists("newrelic_alert_channel.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.foo", "type", "email"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.foo", "configuration.recipients", "terraform-acctest+foo@hashicorp.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.foo", "configuration.include_json_attachment", "1"),
					resource.TestCheckNoResourceAttr(
						"newrelic_alert_channel.foo", "headers"),
				),
			},
			{
				Config: testAccCheckNewRelicAlertChannelConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists("newrelic_alert_channel.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.foo", "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.foo", "type", "email"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.foo", "configuration.recipients", "terraform-acctest+bar@hashicorp.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.foo", "configuration.include_json_attachment", "0"),
					resource.TestCheckNoResourceAttr(
						"newrelic_alert_channel.foo", "headers"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertChannel_import(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertChannelConfig(rName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicAlertChannel_Webhook_withoutHeaders(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertChannelConfigWebhook_withoutHeaders(rName),

				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists("newrelic_alert_channel.channel_without_headers"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.channel_without_headers", "name", fmt.Sprintf("tf-test-webhook-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.channel_without_headers", "type", "webhook"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.channel_without_headers", "configuration.base_url", "http://test.com"),
					resource.TestCheckNoResourceAttr(
						"newrelic_alert_channel.channel_without_headers", "headers"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertChannel_Webhook_withHeaders(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertChannelConfigWebhook_withHeaders(rName),

				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists("newrelic_alert_channel.channel_with_headers"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.channel_with_headers", "name", fmt.Sprintf("tf-test-webhook-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.channel_with_headers", "type", "webhook"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.channel_with_headers", "configuration.base_url", "http://test.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.channel_with_headers", "headers.header1", "test"),
				),
			},
		},
	})
}

func testAccCheckNewRelicAlertChannelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_channel" {
			continue
		}

		id, err := strconv.ParseInt(r.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		_, err = client.GetAlertChannel(int(id))

		if err == nil {
			return fmt.Errorf("Alert channel still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicAlertChannelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No channel ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).Client

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		found, err := client.GetAlertChannel(int(id))
		if err != nil {
			return err
		}

		if strconv.Itoa(found.ID) != rs.Primary.ID {
			return fmt.Errorf("Channel not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckNewRelicAlertChannelConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_channel" "foo" {
  name = "tf-test-%s"
	type = "email"
	
	configuration = {
		recipients = "terraform-acctest+foo@hashicorp.com"
		include_json_attachment = "1"
	}
}
`, rName)
}

func testAccCheckNewRelicAlertChannelConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_channel" "foo" {
  name = "tf-test-updated-%s"
	type = "email"
	
	configuration = {
		recipients = "terraform-acctest+bar@hashicorp.com"
		include_json_attachment = "0"
	}
}
`, rName)
}

func testAccCheckNewRelicAlertChannelConfigWebhook_withoutHeaders(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_channel" "channel_without_headers" {
  name = "tf-test-webhook-%s"
  type = "webhook"

  configuration = {
    base_url = "http://test.com",
    auth_username = "username",
    auth_password = "password",
    payload_type = "application/json",
  }
}
`, rName)
}

func testAccCheckNewRelicAlertChannelConfigWebhook_withHeaders(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_channel" "channel_with_headers" {
  name = "tf-test-webhook-%s"
  type = "webhook"

  configuration = {
    base_url = "http://test.com",
    auth_username = "username",
    auth_password = "password",
    payload_type = "application/json",
  }

  headers {
    header1 = "test"
    header2 = "test2"
  }
}
`, rName)
}

func testAccCheckNewRelicAlertChannelConfigWebhook_withEmptyPayload(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_channel" "webhook_with_empty_payload" {
  name = "tf-test-webhook-%s"
  type = "webhook"

  configuration = {
    base_url = "http://test.com",
    auth_username = "username",
    auth_password = "password",
    payload_type = "application/json",
  }

  payload = {
  }
}
`, rName)
}

func testAccCheckNewRelicAlertChannelConfigWebhook_withPayload(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_channel" "webhook_with_payload" {
  name = "tf-test-webhook-%s"
  type = "webhook"

  configuration = {
    base_url = "http://test.com",
    auth_username = "username",
    auth_password = "password",
    payload_type = "application/json",
  }

  payload = {
    account_id = "test"
  }
}
`, rName)
}

func TestAccNewRelicAlertChannel_Webhook_withPayload(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccCheckNewRelicAlertChannelConfigWebhook_withPayload(rName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists("newrelic_alert_channel.webhook_with_payload"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.webhook_with_payload", "name", fmt.Sprintf("tf-test-webhook-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.webhook_with_payload", "type", "webhook"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.webhook_with_payload", "configuration.base_url", "http://test.com"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_channel.webhook_with_payload", "payload.account_id", "test"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertChannel_Webhook_withEmptyPayloadReturnsError(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile("expected payload not to be empty")
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckNewRelicAlertChannelConfigWebhook_withEmptyPayload(rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}
