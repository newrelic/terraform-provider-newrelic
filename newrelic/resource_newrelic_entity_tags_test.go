// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

func TestAccNewRelicEntityTags_Basic(t *testing.T) {
	resourceName := "newrelic_entity_tags.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicEntityTagsDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicEntityTagsConfig(testAccExpectedApplicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityTagsExist(resourceName, []string{"test_key"}),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicEntityTagsConfigUpdated(testAccExpectedApplicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityTagsExist(resourceName, []string{"test_key_2"}),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func testAccCheckNewRelicEntityTagsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_entity_tags" {
			continue
		}

		_, err := client.Entities.ListTags(entities.EntityGUID(r.Primary.ID))

		if err != nil {
			return fmt.Errorf("entity tags still exist: %s", err)
		}

	}
	return nil
}

func testAccCheckNewRelicEntityTagsExist(n string, keysToCheck []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no entity GUID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		t, err := client.Entities.GetTagsForEntity(entities.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}

		tags := convertTagTypes(t)

		for _, keyToCheck := range keysToCheck {
			if tag := getTag(tags, keyToCheck); tag == nil {
				return fmt.Errorf("entity tag %s not found for GUID %s", keyToCheck, rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccNewRelicEntityTagsConfig(appName string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "foo" {
  name = "%s"
  type = "APPLICATION"
  domain = "APM"
}

resource "newrelic_entity_tags" "foo" {
  guid = data.newrelic_entity.foo.guid

  tag {
	key = "test_key"
	values = ["test_value"]
  }
}
`, appName)
}

func testAccNewRelicEntityTagsConfigUpdated(appName string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "foo" {
  name = "%s"
  type = "APPLICATION"
  domain = "APM"
}

resource "newrelic_entity_tags" "foo" {
  guid = data.newrelic_entity.foo.guid

  tag {
	key = "test_key_2"
	values = ["test_value_2"]
  }
}
`, appName)
}
