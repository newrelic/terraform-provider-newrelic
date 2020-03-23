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
					testAccCheckNewRelicAlertPolicyChannelExists(resourceName),
				),
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
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccNewRelicAlertPolicyImportStateCheckFunc(resourceName),
			},
		},
	})
}

func TestAccNewRelicAlertPolicyChannel_MutipleChannels(t *testing.T) {
	resourceName := "newrelic_alert_policy_channel.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertPolicyChannelsConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyChannelExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertPolicyChannelsConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyChannelExists(resourceName),
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
				PreConfig: testAccDeleteNewRelicAlertPolicy(rName),
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
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_policy_channel" {
			continue
		}

		ids, err := parseHashedIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		policyID := ids[0]
		channelIDs := ids[1:]

		exists, err := policyChannelsExist(client, policyID, channelIDs)
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

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		ids, err := parseHashedIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		policyID := ids[0]
		channelIDs := ids[1:]

		exists, err := policyChannelsExist(client, policyID, channelIDs)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("resource not found: %v", rs.Primary.ID)
		}

		return nil
	}
}

func testAccNewRelicAlertPolicyImportStateCheckFunc(resourceName string) resource.ImportStateCheckFunc {
	return func(state []*terraform.InstanceState) error {
		expectedChannelsCount := "1"
		channelsCount := state[0].Attributes["channel_ids.#"]

		if channelsCount != expectedChannelsCount {
			return fmt.Errorf(
				"%s: Attribute '%s' expected %#v, got %#v",
				resourceName,
				"channel_ids.#",
				expectedChannelsCount,
				channelsCount,
			)
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

	config {
		recipients = "terraform-acctest+foo@hashicorp.com"
		include_json_attachment = "1"
	}
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = newrelic_alert_policy.foo.policy_id
  channel_id = newrelic_alert_channel.foo.id
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

	config {
		recipients = "terraform-acctest+bar@hashicorp.com"
		include_json_attachment = "0"
	}
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = newrelic_alert_policy.foo.policy_id
  channel_id = newrelic_alert_channel.foo.id
}
`, rName)
}

func testAccNewRelicAlertPolicyChannelsConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_alert_channel" "foo" {
  name = "tf-test-%[1]s"
	type = "email"
	config {
		recipients = "terraform-acctest+foo@hashicorp.com"
		include_json_attachment = "1"
	}
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = newrelic_alert_policy.foo.policy_id
  channel_ids = [
		newrelic_alert_channel.foo.id
	]
}
`, name)
}

func testAccNewRelicAlertPolicyChannelsConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%[1]s"
}

resource "newrelic_alert_channel" "foo" {
  name = "tf-test-%[1]s"
	type = "email"
	config {
		recipients = "terraform-acctest+foo@hashicorp.com"
		include_json_attachment = "1"
	}
}

resource "newrelic_alert_channel" "bar" {
  name = "tf-test-2-%[1]s"
	type = "email"
	config {
		recipients = "terraform-acctest+bar@hashicorp.com"
	}
}

resource "newrelic_alert_policy_channel" "foo" {
  policy_id  = newrelic_alert_policy.foo.policy_id
  channel_ids = [
		newrelic_alert_channel.foo.id,
		newrelic_alert_channel.bar.id
	]
}
`, name)
}
