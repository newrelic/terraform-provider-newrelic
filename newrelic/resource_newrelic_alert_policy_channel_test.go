package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicAlertPolicyChannel_Basic(t *testing.T) {
	resourceName := "newrelic_alert_policy_channel.foo"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertPolicyChannelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyChannelExists("newrelic_alert_policy_channel.foo"),
				),
			},
			// Test: No diff on re-apply
			{
				Config:             testAccNewRelicAlertPolicyChannelConfig(rName),
				ExpectNonEmptyPlan: false,
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertPolicyChannelConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyChannelExists("newrelic_alert_policy_channel.foo"),
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

func TestAccNewRelicAlertPolicyChannel_AlertPolicyNotFound(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertPolicyConfig(rName),
			},
			{
				PreConfig: testAccDeleteAlertPolicy(rName),
				Config:    testAccNewRelicAlertPolicyChannelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyChannelExists("newrelic_alert_policy_channel.foo"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertPolicyChannel_AlertChannelNotFound(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertPolicyConfig(rName),
			},
			{
				PreConfig: testAccDeleteAlertChannel(rName),
				Config:    testAccNewRelicAlertPolicyChannelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyChannelExists("newrelic_alert_policy_channel.foo"),
				),
			},
		},
	})
}

func testAccCheckNewRelicAlertPolicyChannelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_policy_channel" {
			continue
		}

		ids, err := parseIDs(r.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		channelID := ids[1]

		exists, err := policyChannelExists(client, policyID, channelID)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("resource still exists")
		}
	}
	return nil
}

func testAccCheckNewRelicAlertPolicyChannelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).Client

		ids, err := parseIDs(rs.Primary.ID, 2)
		if err != nil {
			return err
		}

		policyID := ids[0]
		channelID := ids[1]

		exists, err := policyChannelExists(client, policyID, channelID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("resource not found: %v", rs.Primary.ID)
		}

		return nil
	}
}

func testAccNewRelicAlertPolicyChannelConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "%[1]s"
}

resource "newrelic_alert_channel" "foo" {
  name = "%[1]s"
	type = "email"
	
	configuration = {
		recipients = "terraform-acctest+foo@hashicorp.com"
		include_json_attachment = "1"
	}
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = "${newrelic_alert_policy.foo.id}"
  channel_id = "${newrelic_alert_channel.foo.id}"
}
`, name)
}

func testAccNewRelicAlertPolicyChannelConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-updated-%[1]s"
}

resource "newrelic_alert_channel" "foo" {
  name = "tf-test-updated-%[1]s"
	type = "email"
	
	configuration = {
		recipients = "terraform-acctest+bar@hashicorp.com"
		include_json_attachment = "0"
	}
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = "${newrelic_alert_policy.foo.id}"
  channel_id = "${newrelic_alert_channel.foo.id}"
}
`, rName)
}
