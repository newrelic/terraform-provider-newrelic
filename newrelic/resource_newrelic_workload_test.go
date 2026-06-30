//go:build integration || WORKLOADS

package newrelic

import (
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func TestAccNewRelicWorkload_Basic(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	// TODO: Need to move this to Terraform sweeper so this runs
	//       after all tests have completed.
	// defer cleanupDanglingWorkloadResources()

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
			// {
			// 	ResourceName:            resourceName,
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{"status_config_automatic"},
			// },
		},
	})
}

func TestAccNewRelicWorkload_EntitiesOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

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
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkloadConfigWrongEntitySearchQueriesOnly(rName, ""),
				ExpectError: regexp.MustCompile("Invalid expression"),
			},
			{
				Config:      testAccNewRelicWorkloadConfigWrongEntitySearchQueriesOnly(rName, "\"\""),
				ExpectError: regexp.MustCompile("expected \"entity_search_query.0.query\" to not be an empty string"),
			},
			{
				Config:      testAccNewRelicWorkloadConfigWrongEntitySearchQueriesOnly(rName, "\"     \""),
				ExpectError: regexp.MustCompile("expected \"entity_search_query.0.query\" to not be an empty string or whitespace"),
			},
			{
				Config: testAccNewRelicWorkloadConfigEntitySearchQueriesOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_EntityMultiSearchQueriesOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadConfigEntityMultiSearchQueriesOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_EntityScopeAccountsOnly(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicWorkloadConfigScopeAccountsOnly(rName),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNewRelicWorkload_BasicOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadBasicConfigOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicWorkloadConfigUpdatedBasicConfigOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
			// Test: Import
			//{
			//	ResourceName:      resourceName,
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	//ImportStateVerifyIgnore: []string{"entity_search_query", "composite_entity_search_query", "description"},
			//},
		},
	})
}

func TestAccNewRelicWorkload_StaticOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadConfigStaticOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_AutomaticOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadAutomaticConfigOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_AutomaticEnabledOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadAutomaticConfig_EnabledOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_AutomaticRemainingEntitiesOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadAutomaticConfig_RemainingEntitiesOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_AutomaticRuleOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadAutomaticConfig_RuleOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_AutomaticRulesOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadAutomaticConfig_RulesOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_AutomaticRuleRollupOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadAutomaticConfig_RuleRollupOnly(rName),
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

		found, err := client.Workloads.GetCollection(ids.AccountID, (ids.GUID))
		if err != nil {
			return err
		}

		if found.GUID != (ids.GUID) {
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

		_, err = client.Workloads.GetCollection(ids.AccountID, (ids.GUID))
		if err == nil {
			return fmt.Errorf("workload still exists")
		}
	}
	return nil
}

func cleanupDanglingWorkloadResources() error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	query := "domain = 'NR1' AND type = 'WORKLOAD' AND (name LIKE '%tf-test-%' OR name LIKE '%tf_test_%')"

	fmt.Printf("\n[INFO] cleaning up any dangling integration test resources... \n")
	time.Sleep(1 * time.Second)

	for {
		matches, err := client.Entities.GetEntitySearchByQuery(
			entities.EntitySearchOptions{},
			query,
			[]entities.EntitySearchSortCriteria{},
		)

		if err != nil {
			return fmt.Errorf("error cleaning up dangling synthetics resources: %s", err)
		}

		if matches != nil {
			resources := matches.Results.Entities
			for _, r := range resources {
				_, err := client.Workloads.WorkloadDelete(common.EntityGUID(string(r.GetGUID())))
				if err != nil {
					log.Printf("[ERROR] error deleting dangling resource: %s", err)
				}
			}

			fmt.Printf("\n[INFO] deleted %d dangling resources", len(resources))
		}

		if matches.Results.NextCursor == "" {
			break
		}
	}

	return nil
}

func testAccNewRelicWorkloadConfig(name string) string {
	return fmt.Sprintf(`

data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	entity_guids = [data.newrelic_entity.app.guid]

	entity_search_query {
		query = "name like '%[3]s'"
	}

	scope_account_ids =  [%[1]d]

  	description = "something"

	status_config_automatic {
		enabled = true
		remaining_entities_rule{
			remaining_entities_rule_rollup {
			  strategy = "BEST_STATUS_WINS"
			  threshold_type = "FIXED"
			  threshold_value = 100
			  group_by = "ENTITY_TYPE"
			}
		}
		rule {
		 entity_guids = [data.newrelic_entity.app.guid]
		 nrql_query{
		   query = "name like 'ok'"
		 }
		rollup {
			strategy = "BEST_STATUS_WINS"
			threshold_type = "FIXED"
			threshold_value = 100
			}
		}
	}

	status_config_static {
	description = "test"
	enabled = true
	status = "OPERATIONAL"
	summary = "egetgykwesgksegkerh"
	}
}

`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadConfigUpdated(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s-updated"
	account_id = %[1]d

	entity_guids = [data.newrelic_entity.app.guid]

	entity_search_query {
		query = "name like '%[3]s'"
	}

	scope_account_ids =  [%[1]d]

  	description = "something"

	status_config_automatic {
		enabled = true
		remaining_entities_rule{
			remaining_entities_rule_rollup {
			  strategy = "WORST_STATUS_WINS"
			  threshold_type = "FIXED"
			  threshold_value = 100
			  group_by = "ENTITY_TYPE"
			}
		}
		rule{
		 entity_guids = [data.newrelic_entity.app.guid]
		 nrql_query{
		   query = "name like 'ok-updated'"
		 }
		rollup{
			strategy = "BEST_STATUS_WINS"
			threshold_type = "FIXED"
			threshold_value = 100
			}
		}
	}

	status_config_static {
	description = "test"
	enabled = true
	status = "OPERATIONAL"
	summary = "summary - updated"
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadConfigUpdatedBasicConfigOnly(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s-updated"
	account_id = %[1]d

	entity_guids = [data.newrelic_entity.app.guid]

	entity_search_query {
		query = "name like '%[3]s'"
	}

	scope_account_ids =  [%[1]d]

	description="something-updated"
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadConfigEntitiesOnly(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	entity_guids = [data.newrelic_entity.app.guid]
}
`, testAccountID, name, testAccExpectedApplicationName)
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

func testAccNewRelicWorkloadConfigWrongEntitySearchQueriesOnly(name string, esq string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	entity_search_query {
		query = %[3]s
	}
}
`, testAccountID, name, esq)
}

func testAccNewRelicWorkloadConfigEntityMultiSearchQueriesOnly(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	entity_search_query {
		query = "tags.namespace like '%%App%%' "
	}

	entity_search_query {
		query = "type = 'DASHBOARD' and name like '%%App%%' "
	}
}
`, testAccountID, name)
}

func testAccNewRelicWorkloadConfigScopeAccountsOnly(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	scope_account_ids =  [%[1]d]
}
`, testAccountID, name)
}

func testAccNewRelicWorkloadBasicConfigOnly(name string) string {
	return fmt.Sprintf(`

data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	entity_guids = [data.newrelic_entity.app.guid]

	entity_search_query {
		query = "name like 'App'"
	}

	scope_account_ids =  [%[1]d]

 	description = "something"
}

`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadConfigStaticOnly(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d
	entity_guids = [data.newrelic_entity.app.guid]


	status_config_static {
		description = "test"
		enabled = true
		status = "OPERATIONAL"
		summary = "egetgykwesgksegkerh"
	  }
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadAutomaticConfigOnly(name string) string {
	return fmt.Sprintf(`

data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d
	entity_guids = [data.newrelic_entity.app.guid]

 	description = "something"
	status_config_automatic {
		enabled = true
		remaining_entities_rule{
			remaining_entities_rule_rollup {
			  strategy = "BEST_STATUS_WINS"
			  threshold_type = "FIXED"
			  threshold_value = 100
			  group_by = "ENTITY_TYPE"
			}
		}
		rule{
		 entity_guids = [data.newrelic_entity.app.guid]
		 nrql_query{
		   query = "name like 'ok'"
		 }
		rollup{
			strategy = "BEST_STATUS_WINS"
			threshold_type = "FIXED"
			threshold_value = 100
			}
		}
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadAutomaticConfig_EnabledOnly(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d
	entity_guids = [data.newrelic_entity.app.guid]

	status_config_automatic {
		enabled = true
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadAutomaticConfig_RemainingEntitiesOnly(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d
	entity_guids = [data.newrelic_entity.app.guid]

 	description = "something"
	status_config_automatic {
		enabled = true
		remaining_entities_rule{
			remaining_entities_rule_rollup {
			  strategy = "BEST_STATUS_WINS"
			  threshold_type = "FIXED"
			  threshold_value = 100
			  group_by = "ENTITY_TYPE"
			}
		}
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadAutomaticConfig_RuleOnly(name string) string {
	return fmt.Sprintf(`

data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d
	entity_guids = [data.newrelic_entity.app.guid]

 	description = "something"
	status_config_automatic {
		enabled = true
		rule{
		 entity_guids = [data.newrelic_entity.app.guid]
		 nrql_query{
		   query = "name like 'ok'"
		 }
			rollup{
				strategy = "BEST_STATUS_WINS"
				threshold_type = "FIXED"
				threshold_value = 100
			}
		}
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadAutomaticConfig_RulesOnly(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d
	entity_guids = [data.newrelic_entity.app.guid]


	status_config_automatic {
		enabled = true
		rule{
		 entity_guids = [data.newrelic_entity.app.guid]
		 nrql_query{
		   query = "name like 'ok'"
		 }
			rollup{
				strategy = "BEST_STATUS_WINS"
				threshold_type = "FIXED"
				threshold_value = 100
			}
		}
		rule{
		 entity_guids = [data.newrelic_entity.app.guid]
		 nrql_query{
		   query = "name like 'ok'"
		 }
			rollup{
				strategy = "BEST_STATUS_WINS"
				threshold_type = "FIXED"
				threshold_value = 100
			}
		}
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadAutomaticConfig_RuleRollupOnly(name string) string {
	return fmt.Sprintf(`

data "newrelic_entity" "app" {
	name = "%[3]s"
	domain = "APM"
	type = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d
	entity_guids = [data.newrelic_entity.app.guid]

	status_config_automatic {
		enabled = true
		rule {
    		nrql_query {
     			query = "name like 'ok'"
    		}
			rollup {
				strategy = "BEST_STATUS_WINS"
				threshold_type = "FIXED"
				threshold_value = 100
			}
		}
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func TestAccNewRelicWorkload_CrossAccountEntityGUIDs(t *testing.T) {
	t.Skipf("Skipping this resource until the bug associated with cross-account entity GUIDs is resolved.")

	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkload_CrossAccountEntityGUIDs(rName, testAccountID, testSubAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "entity_guids.#", "1"),
				),
			},
		},
	})
}

func testAccNewRelicWorkload_CrossAccountEntityGUIDs(name string, testAccountID int, testSubAccountID int) string {
	return fmt.Sprintf(`

	data "newrelic_entity" "app" {
	  name       = "Dummy App Two Max"
	  domain     = "APM"
	  type       = "APPLICATION"
	  account_id = %[3]d
	}
	
	resource "newrelic_workload" "foo" {
	  name              = "%[1]s"
	  account_id        = %[2]d
	  entity_guids      = [data.newrelic_entity.app.guid]
	  scope_account_ids = [%[2]d]
	}
`, name, testAccountID, testSubAccountID)
}

// TestAccNewRelicWorkload_EntitySearchQueryBackticksNoDrift verifies that an
// un-backticked tag query (e.g. "tags.env = 'prod'") in entity_search_query
// does not cause perpetual drift on subsequent plans. The provider's expand
// rewrites tags.* -> `tags.*` before sending to NerdGraph, so the API stores
// the backticked form. Without the custom Set hash on entity_search_query,
// the SDK sees config (un-backticked) and Read's response (backticked) as
// different elements, leaving state mismatched with config and producing
// perpetual drift. With the Set hash that normalizes before hashing, both
// forms map to the same identity, the SDK keeps config's form in state
// (Optional-only set rule), and plan stays empty. The framework's automatic
// post-apply refresh-then-plan check is what enforces this — no explicit
// assertion is needed.
func TestAccNewRelicWorkload_EntitySearchQueryBackticksNoDrift(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadConfigEntitySearchQueryUnBackticked(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicWorkloadConfigEntitySearchQueryUnBackticked(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
  name       	= "%[1]s"
  account_id 	= %[2]d

  entity_search_query {
    query = "tags.environment = 'prod'"
  }

  scope_account_ids = [%[2]d]
}
`, name, testAccountID)
}

func TestAccNewRelicWorkload_IntelligentWorkload(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			// Test: Create intelligent workload with dynamic_flows and status_config_alert_policy
			{
				Config: testAccNewRelicWorkloadIntelligentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "dynamic_flows.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status_config_alert_policy.#", "1"),
				),
			},
			// Test: Update intelligent workload (update name, disable alert_policy)
			{
				Config: testAccNewRelicWorkloadIntelligentConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "dynamic_flows.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status_config_alert_policy.#", "1"),
				),
			},
		},
	})
}

func testAccNewRelicWorkloadIntelligentConfig(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name   = "%[3]s"
	domain = "APM"
	type   = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name       = "%[2]s"
	account_id = %[1]d

	dynamic_flows {
		entity_guid      = data.newrelic_entity.app.guid
		transaction_name = "WebTransaction/Action/index"
	}

	scope_account_ids = [%[1]d]

	status_config_alert_policy {
		enabled = true
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}

func testAccNewRelicWorkloadIntelligentConfigUpdated(name string) string {
	return fmt.Sprintf(`
data "newrelic_entity" "app" {
	name   = "%[3]s"
	domain = "APM"
	type   = "APPLICATION"
}

resource "newrelic_workload" "foo" {
	name       = "%[2]s-updated"
	account_id = %[1]d

	dynamic_flows {
		entity_guid      = data.newrelic_entity.app.guid
		transaction_name = "WebTransaction/Action/index"
	}

	scope_account_ids = [%[1]d]

	status_config_alert_policy {
		enabled = false
	}
}
`, testAccountID, name, testAccExpectedApplicationName)
}
