//go:build integration || unit
// +build integration unit

// Test helpers
//

package newrelic

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/agentapplications"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apm"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logconfigurations"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

var (
	testAccExpectedAlertChannelName            string
	testAccExpectedApplicationName             string
	testAccExpectedSingleQuotedApplicationName string
	testAccExpectedAlertPolicyName             string
	testAccAPIKey                              string
	testAccProviders                           map[string]*schema.Provider
	testAccProvider                            *schema.Provider
	testAccountID                              int
	testSubAccountID                           int
	testAccountName                            string
	testAccAPMEntityCreated                    = false
	testAccAPMSingleQuotedEntityCreated        = false
	testAccCleanupComplete                     = false
	testCloudAccCleanupComplete                = map[string]bool{
		"aws":   false,
		"azure": false,
		"gcp":   false,
	}
	testAccBrowserApplicationCleanupComplete    = false
	testAccSyntheticTestEntitiesCleanupComplete = false
)

func init() {
	testAccExpectedAlertChannelName = fmt.Sprintf("%s tf-test@example.com", acctest.RandString(5))
	testAccExpectedApplicationName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccExpectedSingleQuotedApplicationName = fmt.Sprintf("tf_test_quote_%s%s%s", acctest.RandString(5), "'", acctest.RandString(5))
	testAccExpectedAlertPolicyName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"newrelic": testAccProvider,
	}
	testAccAPIKey = os.Getenv("NEW_RELIC_API_KEY")
	if v := os.Getenv("NEW_RELIC_API_KEY"); v == "" {
		testAccAPIKey = "foo"
	}

	if v, _ := strconv.Atoi(os.Getenv("NEW_RELIC_ACCOUNT_ID")); v != 0 {
		testAccountID = v
	}

	// Used for cross-account scenarios if needed, such as dashboard widgets.
	if v, _ := strconv.Atoi(os.Getenv("NEW_RELIC_SUBACCOUNT_ID")); v != 0 {
		testSubAccountID = v
	}
	if v := os.Getenv("NEW_RELIC_ACCOUNT_NAME"); v != "" {
		testAccountName = v
	} else {
		testAccountName = "New Relic Terraform Provider Acceptance Testing"
	}
}

func testAccNewRelicProviderConfig(region string, baseURL string, resourceName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
	alias = "integration-test-provider"
	region = "%[1]s"

	%[2]s
}

resource "newrelic_alert_policy" "foo" {
	provider = newrelic.integration-test-provider
  name = "tf-test-%[3]s"
}
`, region, baseURL, resourceName)
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckEnvVars(t)

	// Clean up old data partitions
	//testAccLogDataPartitionsCleanup(t)

	// Cleaning up the Parsing rules
	//testAccLogParsingRulesCleanup(t)

	if !testAccSyntheticTestEntitiesCleanupComplete {
		testAccSyntheticTestEntitiesCleanup(t)
	}

	// Create a test application for use in newrelic_alert_condition and other tests
	if !testAccAPMEntityCreated {
		testAccCreateEntity(t, testAccExpectedApplicationName)
		testAccAPMEntityCreated = true
	}
}

func testAccSingleQuotedPreCheck(t *testing.T) {
	testAccPreCheckEnvVars(t)

	// Clean up old data partitions
	//testAccLogDataPartitionsCleanup(t)

	// Cleaning up the Parsing rules
	//testAccLogParsingRulesCleanup(t)

	// Create a test application for use in tests in the "entity" data source
	// comprising an apostrophe "'".
	if !testAccAPMSingleQuotedEntityCreated {
		testAccCreateEntity(t, testAccExpectedSingleQuotedApplicationName)
		testAccAPMSingleQuotedEntityCreated = true
	}
}

func testAccCreateEntity(t *testing.T, name string) {
	// Clean up old instances of the applications
	testAccApplicationsCleanup(t)

	// Create the application, with the given 'name' in the argument.
	testAccCreateApplication(t, name)

	// We need to give the entity search engine time to index the app so
	// we try to get the entity, and retry if it fails for a certain amount
	// of time
	client := entities.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})
	params := entities.EntitySearchQueryBuilder{
		Name:   escapeSingleQuote(name),
		Type:   "APPLICATION",
		Domain: "APM",
	}

	retryErr := resource.RetryContext(context.Background(), 30*time.Second, func() *resource.RetryError {
		entityResults, err := client.GetEntitySearchWithContext(
			context.Background(),
			entities.EntitySearchOptions{},
			"",
			params,
			[]entities.EntitySearchSortCriteria{},
			[]entities.SortCriterionWithDirection{},
		)
		if err != nil {
			return resource.RetryableError(err)
		}

		if entityResults.Count != 1 {
			return resource.RetryableError(fmt.Errorf("Entity not found, or found more than one"))
		}

		return nil
	})

	if retryErr != nil {
		t.Fatalf("Unable to find application entity: %s", retryErr)
	}

	// We have to give time for the async nature of the entity creation to complete
	time.Sleep(1 * time.Second)
}

func testAccPreCheckEnvVars(t *testing.T) {
	if v := os.Getenv("NEW_RELIC_API_KEY"); v == "" {
		t.Skipf("[WARN] NEW_RELIC_API_KEY has not been set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_LICENSE_KEY"); v == "" {
		t.Skipf("NEW_RELIC_LICENSE_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("NEW_RELIC_ACCOUNT_ID"); v == "" {
		t.Skipf("NEW_RELIC_ACCOUNT_ID must be set for acceptance tests")
	}
}

func testAccCreateApplication(t *testing.T, name string) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigAppName(name),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
	)

	if err != nil {
		t.Fatalf("Error setting up New Relic application: %s", err)
	}

	if err := app.WaitForConnection(30 * time.Second); err != nil {
		t.Fatalf("Unable to setup New Relic application connection: %s", err)
	}

	app.RecordCustomEvent("terraform test", nil)

	app.Shutdown(30 * time.Second)
}

func testAccApplicationsCleanup(t *testing.T) {
	// Only run cleanup once per test run
	if testAccCleanupComplete {
		return
	}

	client := apm.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})

	params := apm.ListApplicationsParams{
		Name: "tf_test",
	}

	applications, err := client.ListApplications(&params)

	if err != nil {
		t.Logf("error fetching applications: %s", err)
	}

	deletedAppCount := 0

	for _, app := range applications {
		if !app.Reporting {
			// Applications that have reported in the past 12 hours will not be deleted
			// because of limitation on the API
			_, err = client.DeleteApplication(app.ID)

			if err == nil {
				deletedAppCount++
				t.Logf("deleted application %s (%d/%d)", app.Name, deletedAppCount, len(applications))
			}
		}
	}

	t.Logf("testacc cleanup of %d applications complete", deletedAppCount)

	testAccCleanupComplete = true
}

// Facilitates using a standardized name when creating test resources.
// The name will always be prefixed with "tf-test-". This ensures when
// we attempt to delete any dangling extraneous resources, we only delete
// resources with names that start with "tf-test-". This helps avoid
// deleting any resources that might be cross-account, such as workloads.
func generateNameForIntegrationTestResource() string {
	return fmt.Sprintf("tf-test-%s", acctest.RandString(15))
}

// testAccCloudLinkedAccountsCleanup handles cleaning up/deleting cloud accounts linked to the
// specified New Relic account, of the specified provider (aws, azure, gcp).
// Deleting linked accounts also deletes associated integrations
func testAccCloudLinkedAccountsCleanup(t *testing.T, provider string) {
	if testCloudAccCleanupComplete[provider] {
		return
	}
	client := cloud.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})
	t.Logf("***** unlinking '%s' cloud accounts linked to the subaccount ******", provider)

	cloudLinkedAccounts, err := client.GetLinkedAccounts(provider)
	if err != nil {
		t.Logf("error fetching '%s' cloud accounts linked to the subaccount: %s", provider, err)
	}

	if cloudLinkedAccounts == nil {
		t.Logf("no '%s' cloud accounts linked to the subaccount found to be deleted.", provider)
	} else {
		cloudLinkedAccountsToBeDisabled := make([]cloud.CloudUnlinkAccountsInput, 0)
		for _, cloudLinkedAccount := range *cloudLinkedAccounts {
			if cloudLinkedAccount.NrAccountId == testSubAccountID {
				cloudLinkedAccountsToBeDisabled = append(cloudLinkedAccountsToBeDisabled, cloud.CloudUnlinkAccountsInput{LinkedAccountId: cloudLinkedAccount.ID})
				t.Logf("identified '%s' cloud account '%d' linked to the subaccount '%d'", provider, cloudLinkedAccount.ID, testSubAccountID)
			}
		}
		_, err = client.CloudUnlinkAccount(testSubAccountID, cloudLinkedAccountsToBeDisabled)
		if err == nil {
			t.Logf("deleted %d '%s' cloud accounts linked to the subaccount: %v", len(cloudLinkedAccountsToBeDisabled), provider, cloudLinkedAccountsToBeDisabled)
		}
	}
	t.Logf("testacc cleanup of '%s' cloud accounts linked to the subaccount complete", provider)
	testCloudAccCleanupComplete[provider] = true
}

func testAccSyntheticTestEntitiesCleanup(t *testing.T) {
	testAccPreCheckEnvVars(t)

	if testAccSyntheticTestEntitiesCleanupComplete {
		return
	}
	testSyntheticEntitiesCleanedUpCount := 0

	client := entities.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})

	syntheticsClient := synthetics.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})

	syntheticEntityNameMatchingSuffixes := []string{"tf-test", "client-go-test"}
	for _, nameSuffix := range syntheticEntityNameMatchingSuffixes {
		t.Logf("***** fetching synthetic entities created by Terraform, comprising '%s' in their name ******", nameSuffix)
		testSynthEntities, err := client.GetEntitySearch(
			entities.EntitySearchOptions{},
			"",
			entities.EntitySearchQueryBuilder{
				Domain: "SYNTH",
				Name:   nameSuffix,
			},
			[]entities.EntitySearchSortCriteria{},
			[]entities.SortCriterionWithDirection{},
		)

		if err != nil {
			t.Logf("error fetching synthetic entities linked to the account: %s", err)
		}

		if testSynthEntities == nil {
			t.Logf("***** no synthetic entities linked to the account with the queried parameters found, to be deleted *****")
		} else {
			for _, a := range testSynthEntities.Results.Entities {
				entityType := a.GetType()
				entityDomain := a.GetDomain()
				entityName := a.GetName()
				entityGUID := a.GetGUID()
				if entityType == "SECURE_CRED" && entityDomain == "SYNTH" {
					_, err := syntheticsClient.SyntheticsDeleteSecureCredential(testAccountID, entityName)
					if err == nil {
						t.Logf("#%d: deleted %s '%s' with GUID %s", testSyntheticEntitiesCleanedUpCount+1, entityType, entityName, entityGUID)
						testSyntheticEntitiesCleanedUpCount += 1
					} else {
						t.Logf("failed to delete %s '%s' with GUID %s, error: %s", entityType, entityName, entityGUID, err.Error())
					}
				} else if entityType == "PRIVATE_LOCATION" && entityDomain == "SYNTH" {
					_, err := syntheticsClient.SyntheticsDeletePrivateLocation(synthetics.EntityGUID(entityGUID))
					if err == nil {
						t.Logf("#%d: deleted %s '%s' with GUID %s", testSyntheticEntitiesCleanedUpCount+1, entityType, entityName, entityGUID)
						testSyntheticEntitiesCleanedUpCount += 1
					} else {
						t.Logf("failed to delete %s '%s' with GUID %s, error: %s", entityType, entityName, entityGUID, err.Error())
					}
				} else if entityType == "MONITOR" && entityDomain == "SYNTH" {
					_, err := syntheticsClient.SyntheticsDeleteMonitor(synthetics.EntityGUID(entityGUID))
					if err == nil {
						t.Logf("#%d: deleted %s '%s' with GUID %s", testSyntheticEntitiesCleanedUpCount+1, entityType, entityName, entityGUID)
						testSyntheticEntitiesCleanedUpCount += 1
					} else {
						t.Logf("failed to delete %s '%s' with GUID %s, error: %s", entityType, entityName, entityGUID, err.Error())
					}
				}
			}
		}
		t.Logf("testacc cleanup of '%d' synthetic entities created by terraform acceptance tests complete", testSyntheticEntitiesCleanedUpCount)
	}
	testAccSyntheticTestEntitiesCleanupComplete = true
}

func testAccBrowserApplicationsCleanup(t *testing.T) {
	testAccPreCheckEnvVars(t)

	if testAccBrowserApplicationCleanupComplete {
		return
	}
	testBrowserApplicationEntitiesCleanedUpCount := 0

	client := entities.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})

	browserApplicationClient := agentapplications.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})

	t.Logf("***** fetching test browser application entities created by Terraform ******")

	testBrowserApplicationEntities, err := client.GetEntitySearch(
		entities.EntitySearchOptions{},
		"",
		entities.EntitySearchQueryBuilder{
			Domain: "BROWSER",
			Type:   "APPLICATION",
		},
		[]entities.EntitySearchSortCriteria{},
		[]entities.SortCriterionWithDirection{},
	)

	if err != nil {
		t.Logf("error fetching browser application entities linked to the account: %s", err)
	}

	if testBrowserApplicationEntities == nil {
		t.Logf("***** no browser application entities linked to the subaccount found to be deleted *****")
	} else {
		for _, a := range testBrowserApplicationEntities.Results.Entities {
			// Check if the entity is a Browser Application
			if a.GetType() == "APPLICATION" && a.GetDomain() == "BROWSER" {
				// Check if the name of the BrowserApplicationEntity starts with "tf-test-" or "nr-test-"
				// If yes, these may be deleted via a NerdGraph mutation as these were created by a relevant Terraform acceptance test
				if strings.Contains(a.GetName(), "tf-test-") || strings.Contains(a.GetName(), "nr-test-") {
					t.Logf("***** deleting browser application entity (%d) '%s' *****", testBrowserApplicationEntitiesCleanedUpCount+1, a.GetName())
					_, err := browserApplicationClient.AgentApplicationDelete(common.EntityGUID(a.GetGUID()))
					if err != nil {
						t.Logf("error deleting the browser application entity with GUID %s : %s", a.GetGUID(), err)
					}
					testBrowserApplicationEntitiesCleanedUpCount += 1
				}
			}
		}
		t.Logf("testacc cleanup of '%d' browser application entities created by terraform acceptance tests complete", testBrowserApplicationEntitiesCleanedUpCount)
	}
	testAccBrowserApplicationCleanupComplete = true
}

// Deleting the data partitions as they start with "Log_Test_"
// Only run if the data partitions limit exceeded
func testAccLogDataPartitionsCleanup(t *testing.T) {
	// Only run cleanup once per test run
	if testAccCleanupComplete {
		return
	}
	client := logconfigurations.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})
	t.Logf("***** Deleting data partitions ******")
	dataPartitions, err := client.GetDataPartitionRules(testAccountID)
	if err != nil {
		t.Logf("error fetching data partitions: %s", err)
	}
	deletedDataPartitionCount := 0

	for _, v := range *dataPartitions {
		str := string(v.TargetDataPartition)
		if (strings.Contains(str, "Log_Test_") || strings.Contains(str, "Log_testName_")) && v.Deleted != true {
			_, err = client.LogConfigurationsDeleteDataPartitionRule(testAccountID, v.ID)

			if err == nil {
				deletedDataPartitionCount++
				t.Logf("deleted data partition %s (%d/%d)", v.ID, deletedDataPartitionCount, len(*dataPartitions))
			}
		}
	}

	t.Logf("testacc cleanup of %d DataPartition complete", deletedDataPartitionCount)

	testAccCleanupComplete = true

}

// delete the parsing rules
// Only run if the limit is exceeded
func testAccLogParsingRulesCleanup(t *testing.T) {
	// Only run cleanup once per test run
	if testAccCleanupComplete {
		return
	}
	client := logconfigurations.New(config.Config{
		PersonalAPIKey: testAccAPIKey,
	})
	t.Logf("***** Deleting parsing rules ******")
	rules, err := client.GetParsingRules(testAccountID)
	if err != nil {
		t.Logf("error fetching data parsing rules: %s", err)
	}
	deletedCount := 0

	for _, v := range *rules {
		str := string(v.Description)
		if (strings.Contains(str, "testDescription_") || strings.Contains(str, "tf_test_")) && v.Deleted != true {
			_, err = client.LogConfigurationsDeleteParsingRule(testAccountID, v.ID)

			if err == nil {
				deletedCount++
				t.Logf("deleted parsing rules %s (%d/%d)", v.ID, deletedCount, len(*rules))
			}
		}
	}

	t.Logf("testacc cleanup of %d DataPartition complete", deletedCount)

	testAccCleanupComplete = true

}
