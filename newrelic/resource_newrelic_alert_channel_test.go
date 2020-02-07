package newrelic

import (
	"fmt"
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
	rNameDeprecated := fmt.Sprintf("tf-test-deprecated-%s", rand)
	rNameDeprecatedUpdated := fmt.Sprintf("tf-test-deprecated-updated-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create (Deprecated)
			{
				Config: testAccNewRelicAlertChannelDeprecatedConfig(rNameDeprecated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameDeprecated),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "configuration.recipients", "terraform-acctest+foo@hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "configuration.include_json_attachment", "1"),
				),
				// The deprecated configuration attribute requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
			},
			// Test: Update (Deprecated)
			{
				Config: testAccNewRelicAlertChannelDeprecatedConfigUpdated(rNameDeprecatedUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameDeprecatedUpdated),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "configuration.recipients", "terraform-acctest+bar@hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "configuration.include_json_attachment", "0"),
				),
				// The deprecated configuration attribute requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
			},
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

func TestAccNewRelicAlertChannel_Slack(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	rNameDeprecated := fmt.Sprintf("tf-test-deprecated-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create (Deprecated)
			{
				Config: testAccNewRelicAlertChannelDeprecatedConfigByType(rNameDeprecated, "slack", `{
					url = "https://example.slack.com"
					channel = "example-channel"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameDeprecated),
					resource.TestCheckResourceAttr(resourceName, "type", "slack"),
					resource.TestCheckResourceAttr(resourceName, "configuration.url", "https://example.slack.com"),
					resource.TestCheckResourceAttr(resourceName, "configuration.channel", "example-channel"),
				),
				// The deprecated configuration attribute requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
			},
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
				// The config block requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
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

func TestAccNewRelicAlertChannel_PagerDuty(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	rNameDeprecated := fmt.Sprintf("tf-test-deprecated-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create (Deprecated)
			{
				Config: testAccNewRelicAlertChannelDeprecatedConfigByType(rNameDeprecated, "pagerduty", `{
					service_key = "abc123"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameDeprecated),
					resource.TestCheckResourceAttr(resourceName, "type", "pagerduty"),
					resource.TestCheckResourceAttr(resourceName, "configuration.service_key", "abc123"),
				),
				// The deprecated configuration attribute requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
			},
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
				// The config block requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
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
				// The config block requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
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

func TestAccNewRelicAlertChannel_OpsGenie(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	rNameDeprecated := fmt.Sprintf("tf-test-deprecated-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create (Deprecated)
			{
				Config: testAccNewRelicAlertChannelDeprecatedConfigByType(rNameDeprecated, "opsgenie", `{
					api_key = "abc123"
					teams = "example-team"
					tags = "tag1"
					recipients = "example@somedomain.com"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameDeprecated),
					resource.TestCheckResourceAttr(resourceName, "type", "opsgenie"),
					resource.TestCheckResourceAttr(resourceName, "configuration.api_key", "abc123"),
					resource.TestCheckResourceAttr(resourceName, "configuration.teams", "example-team"),
					resource.TestCheckResourceAttr(resourceName, "configuration.tags", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.recipients", "example@somedomain.com"),
				),
				// The deprecated configuration attribute requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
			},
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
				// The config block requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
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

func TestAccNewRelicAlertChannel_VictorOps(t *testing.T) {
	resourceName := "newrelic_alert_channel.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-test-%s", rand)
	rNameDeprecated := fmt.Sprintf("tf-test-deprecated-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			//Test: Create (Deprecated)
			{
				Config: testAccNewRelicAlertChannelDeprecatedConfigByType(rNameDeprecated, "victorops", `{
					key = "abc123"
					route_key = "/example-route"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameDeprecated),
					resource.TestCheckResourceAttr(resourceName, "type", "victorops"),
					resource.TestCheckResourceAttr(resourceName, "configuration.key", "abc123"),
					resource.TestCheckResourceAttr(resourceName, "configuration.route_key", "/example-route"),
				),
				// The deprecated configuration attribute requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
			},
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
					resource.TestCheckResourceAttr(resourceName, "config.0.key", "abc123"),
					resource.TestCheckResourceAttr(resourceName, "config.0.route_key", "/example-route"),
				),
				// The config block requires the resource being destroyed and recreated on every `apply`.
				ExpectNonEmptyPlan: true,
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

func testAccNewRelicAlertChannelDeprecatedConfig(rName string) string {
	return fmt.Sprintf(`
	resource "newrelic_alert_channel" "foo" {
		name = "%s"
		type = "email"

		configuration = {
			recipients = "terraform-acctest+foo@hashicorp.com"
			include_json_attachment = "1"
		}
	}
`, rName)
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

func testAccNewRelicAlertChannelDeprecatedConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_alert_channel" "foo" {
		name = "%s"
		type = "email"

		configuration = {
			recipients = "terraform-acctest+bar@hashicorp.com"
			include_json_attachment = "0"
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

func testAccNewRelicAlertChannelDeprecatedConfigByType(name string, channelType string, configuration string) string {
	return fmt.Sprintf(`
		resource "newrelic_alert_channel" "foo" {
			name = "%s"
			type = "%s"

			configuration = %s
		}
	`, name, channelType, configuration)
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
