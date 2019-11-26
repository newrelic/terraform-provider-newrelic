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
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertChannelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "configuration.recipients", "terraform-acctest+foo@hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "configuration.include_json_attachment", "1"),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertChannelConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "configuration.recipients", "terraform-acctest+bar@hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "configuration.include_json_attachment", "0"),
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
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "slack", `{
					url = "https://example.slack.com"
					channel = "example-channel"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "slack"),
					resource.TestCheckResourceAttr(resourceName, "configuration.url", "https://example.slack.com"),
					resource.TestCheckResourceAttr(resourceName, "configuration.channel", "example-channel"),
				),
				// This notification channel resource requires being destroyed and recreated on every `apply`.
				// This is due to how the New Relic API has to handle this scenario, hence we need to set this to `true`.
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
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "pagerduty", `{
					service_key = "abc123"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "pagerduty"),
					resource.TestCheckResourceAttr(resourceName, "configuration.service_key", "abc123"),
				),
				// This notification channel resource requires being destroyed and recreated on every `apply`.
				// This is due to how the New Relic API has to handle this scenario, hence we need to set this to `true`.
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
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "opsgenie", `{
					api_key = "abc123"
					teams = "example-team"
					tags = "tag1"
					recipients = "example@somedomain.com"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "opsgenie"),
					resource.TestCheckResourceAttr(resourceName, "configuration.api_key", "abc123"),
					resource.TestCheckResourceAttr(resourceName, "configuration.teams", "example-team"),
					resource.TestCheckResourceAttr(resourceName, "configuration.tags", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.recipients", "example@somedomain.com"),
				),
				// This notification channel resource requires being destroyed and recreated on every `apply`.
				// This is due to how the New Relic API has to handle this scenario, hence we need to set this to `true`.
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
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertChannelConfigByType(rName, "victorops", `{
					key = "abc123"
					route_key = "/example-route"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "victorops"),
					resource.TestCheckResourceAttr(resourceName, "configuration.key", "abc123"),
					resource.TestCheckResourceAttr(resourceName, "configuration.route_key", "/example-route"),
				),
				// This notification channel resource requires being destroyed and recreated on every `apply`.
				// This is due to how the New Relic API has to handle this scenario, hence we need to set this to `true`.
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

func TestAccNewRelicAlertChannel_NoDiffOnReapply(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertChannelConfig(rName),
			},
			{
				Config:             testAccNewRelicAlertChannelConfig(rName),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccNewRelicAlertChannel_ResourceNotFound(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

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

func testAccNewRelicAlertChannelConfig(rName string) string {
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

func testAccNewRelicAlertChannelConfigUpdated(rName string) string {
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

func testAccNewRelicAlertChannelConfigByType(rName string, channelType string, configuration string) string {
	return fmt.Sprintf(`
		resource "newrelic_alert_channel" "foo" {
			name = "tf-test-%s"
			type = "%s"

			configuration = %s
		}
	`, rName, channelType, configuration)
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
			return fmt.Errorf("channel not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccDeleteAlertChannel(name string) func() {
	return func() {
		client := testAccProvider.Meta().(*ProviderConfig).Client
		alertChannels, _ := client.ListAlertChannels()

		for _, d := range alertChannels {
			if d.Name == name {
				_ = client.DeleteAlertChannel(d.ID)
				break
			}
		}
	}
}
