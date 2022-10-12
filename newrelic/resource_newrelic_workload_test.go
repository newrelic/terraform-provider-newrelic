//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entity_search_query", "composite_entity_search_query", "description", "status_config_automatic", "status_config_static"},
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

func TestAccNewRelicWorkload_EntityMultiSearchQueriesOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := acctest.RandString(5)

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

func TestAccNewRelicWorkload_BasicOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := acctest.RandString(5)

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
			// Test: Update
			//{
			//	Config: testAccNewRelicWorkloadConfigUpdatedBasicConfigOnly(rName),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckNewRelicWorkloadExists(resourceName),
			//	),
			//},
			// Test: Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entity_search_query", "composite_entity_search_query", "description"},
			},
		},
	})
}

func TestAccNewRelicWorkload_StaticOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicWorkloadDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicWorkloadConfigSaticOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicWorkloadExists(resourceName),
				),
			},
		},
	})
}

func TestAccNewRelicWorkload_AutomaticOnly(t *testing.T) {
	resourceName := "newrelic_workload.foo"
	rName := acctest.RandString(5)

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
	rName := acctest.RandString(5)

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
	rName := acctest.RandString(5)

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
	rName := acctest.RandString(5)

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
	rName := acctest.RandString(5)

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
	rName := acctest.RandString(5)

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

		found, err := client.Workloads.GetWorkload(ids.AccountID, string(ids.GUID))
		if err != nil {
			return err
		}

		if found.GUID != string(ids.GUID) {
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

		_, err = client.Workloads.GetWorkload(ids.AccountID, string(ids.GUID))
		if err == nil {
			return fmt.Errorf("workload still exists")
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

func testAccNewRelicWorkloadConfigSaticOnly(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	status_config_static {
		description = "test"
		enabled = true
		status = "OPERATIONAL"
		summary = "egetgykwesgksegkerh"
	  }
}
`, testAccountID, name)
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
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	status_config_automatic {
		enabled = true
	}
}
`, testAccountID, name)
}

func testAccNewRelicWorkloadAutomaticConfig_RemainingEntitiesOnly(name string) string {
	return fmt.Sprintf(`
resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

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
`, testAccountID, name)
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

resource "newrelic_workload" "foo" {
	name = "%[2]s"
	account_id = %[1]d

	status_config_automatic {
		enabled = true
		rule{
			rollup{
				strategy = "BEST_STATUS_WINS"
				threshold_type = "FIXED"
				threshold_value = 100
			}
		}
	}
}
`, testAccountID, name)
}
