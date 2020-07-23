// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/stretchr/testify/require"
)

func TestAccNewRelicEntityData_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity"),
				),
			},
		},
	})
}

func TestExpandEntityTag(t *testing.T) {
	flattened := []interface{}{
		map[string]interface{}{
			"key":   "my-key",
			"value": "my-value",
		},
	}

	expected := entities.TagValue{
		Key:   "my-key",
		Value: "my-value",
	}

	expanded := expandEntityTag(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, &expected, expanded)
}

func testAccCheckNewRelicEntityDataExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["guid"] == "" {
			return fmt.Errorf("expected to get an entity GUID")
		}

		if a["application_id"] == "" {
			return fmt.Errorf("expected to get an application ID")
		}

		if a["name"] != testAccExpectedApplicationName {
			return fmt.Errorf("expected the entity name to be: %s, but got: %s", testAccExpectedApplicationName, a["name"])
		}

		return nil
	}
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig() string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	type = "application"
	domain = "apm"
	tag {
		key = "accountId"
		value = "%d"
	}
}
`, testAccExpectedApplicationName, testAccountID)
}
