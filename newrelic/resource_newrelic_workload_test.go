package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testEntityGUID      = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
	testApplicationName = "Dummy App"
)

func TestAccNewRelicWorkload_Basic(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicWorkloadConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicWorkloadConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// "entity_search_query" is returned as nil after a delete.
				ImportStateVerifyIgnore: []string{"entity_search_query", "composite_entity_search_query"},
			},
		},
	})
}

func TestAccNewRelicWorkload_EntitiesOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadConfigEntitiesOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_EntitySearchQueriesOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadConfigEntitySearchQueriesOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_EntityScopeAccountsOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadConfigScopeAccountsOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckNewRelicWorkloadExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no workload ID is set")
		}

		ids, err := parseWorkloadIDs(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		found, err := client.Workloads.GetWorkload(ids.AccountID, ids.GUID)
		if err != nil {
			return err
		}

		if found.GUID != ids.GUID {
			return fmt.Errorf("workload not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckNewRelicWorkloadDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_workload" {
			continue
		}

		ids, err := parseWorkloadIDs(r.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.Workloads.GetWorkload(ids.AccountID, ids.GUID)
		if err == nil {
			return fmt.Errorf("workload still exists")
		}
	}
	return nil
}

func testAccNewRelicWorkloadConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	entity_guids = ["%[3]s"]

	entity_search_query {
		query = "name like '%[4]s'"
	}

	scope_account_ids =  [1, %[1]d]
}
`, testAccountID, name, testEntityGUID, testApplicationName)
}

func testAccNewRelicWorkloadConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s-updated"
	account_id = %[1]d

	entity_guids = ["%[3]s"]

	entity_search_query {
		query = "name like '%[4]s'"
	}

	scope_account_ids =  [1, %[1]d]
}
`, testAccountID, name, testEntityGUID, testApplicationName)
}

func testAccNewRelicWorkloadConfigEntitiesOnly(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	entity_guids = ["%[3]s"]
}
`, testAccountID, name, testEntityGUID)
}

func testAccNewRelicWorkloadConfigEntitySearchQueriesOnly(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	entity_search_query {
		query = "name like 'App'"
	}
}
`, testAccountID, name)
}

func testAccNewRelicWorkloadConfigScopeAccountsOnly(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	scope_account_ids =  [1, %[1]d]
}
`, testAccountID, name)
}
