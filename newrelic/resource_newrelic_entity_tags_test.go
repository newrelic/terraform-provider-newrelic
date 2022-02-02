//go:build integration
// +build integration

package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"testing"
	"time"
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
					testAccCheckNewRelicEntityUnmutableExists(resourceName, []string{"account", "guid", "language"}),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicEntityTagsConfigUpdated(testAccExpectedApplicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityTagsExist(resourceName, []string{"test_key_2"}),
					testAccCheckNewRelicEntityUnmutableExists(resourceName, []string{"account", "guid", "language"}),
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

		_, err := client.Entities.ListTags(common.EntityGUID(r.Primary.ID))

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

		retryErr := resource.RetryContext(context.Background(), 5*time.Second, func() *resource.RetryError {
			t, err := client.Entities.GetTagsForEntityMutable(common.EntityGUID(rs.Primary.ID))
			if err != nil {
				return resource.RetryableError(err)
			}

			tags := convertTagTypes(t)

			for _, keyToCheck := range keysToCheck {
				if tag := getTag(tags, keyToCheck); tag == nil {
					return resource.RetryableError(fmt.Errorf("entity tag %s not found for GUID %s", keyToCheck, rs.Primary.ID))
				}
			}

			return nil
		})

		if retryErr != nil {
			return retryErr
		}

		return nil
	}
}

func testAccCheckNewRelicEntityUnmutableExists(n string, keysToCheck []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no entity GUID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		retryErr := resource.RetryContext(context.Background(), 5*time.Second, func() *resource.RetryError {
			t, err := client.Entities.GetTagsForEntityMutable(common.EntityGUID(rs.Primary.ID))
			if err != nil {
				return resource.RetryableError(err)
			}

			tags := convertTagTypes(t)

			for _, keyToCheck := range keysToCheck {
				if tag := getTag(tags, keyToCheck); tag != nil {
					return resource.RetryableError(fmt.Errorf("unmutable entity tag %s found for GUID %s", keyToCheck, rs.Primary.ID))
				}
			}

			return nil
		})

		if retryErr != nil {
			return retryErr
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
