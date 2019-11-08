package newrelic

import (
	"fmt"
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
				),
			},
		},
	})
}

func TestAccNewRelicAlertChannel(t *testing.T) {
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
