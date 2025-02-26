package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	tftest "github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func testHelpersReturnTestSyntheticMonitorGUIDsForDowntime(operation string) []string {
	switch operation {
	case "create":
		return []string{
			"MzgwNjUyNnxTWU5USHxNT05JVE9SfDQ3ZWI5YmYzLWRiOTEtNDljYy04MzM2LTBhZWJhNTE5MzhiOQ",
		}
	case "update":
		return []string{
			"MzgwNjUyNnxTWU5USHxNT05JVE9SfDkwODEwNTRhLWRhYTAtNGI0Mi05YmIwLTY3M2M1MDI2ZWYyOA",
			"MzgwNjUyNnxTWU5USHxNT05JVE9SfGEyNWRmOTIwLTcxYjUtNDlmYy1iZTgzLTBhOGE0NjdiYWNhMg",
			"MzgwNjUyNnxTWU5USHxNT05JVE9SfDAzNjQ0ZDNlLTg0YzMtNDQyMC1hYjM4LTc0ZjBjODI4NTk3ZA",
		}
	}
	return []string{}
}

func testHelpersFetchAPIKeyAndAccountID() (NewRelicAccountId int, NewRelicApiKey string, err error) {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	accountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")

	if apiKey == "" {
		return 0, "", fmt.Errorf("an empty/null NEW_RELIC_API_KEY has been found in the environment, please use a non-empty NEW_RELIC_API_KEY")
	}

	if accountID == "" {
		return 0, "", fmt.Errorf("an empty/null NEW_RELIC_ACCOUNT_ID has been found in the environment, please use a non-empty NEW_RELIC_ACCOUNT_ID")
	}

	accountIDAsInt, convErr := strconv.Atoi(accountID)
	if convErr != nil {
		return 0, "", fmt.Errorf("unable to convert NEW_RELIC_ACCOUNT_ID to an integer: %v", convErr)
	}

	return accountIDAsInt, apiKey, nil
}

func testHelpersWaitUntilAndConfirmEmptyPlan(t *testing.T, terraformOptions *tftest.Options) {
	// Implement retries for plan checks
	retryCount := 3
	retryDelay := 5 * time.Second

	for i := 0; i < retryCount; i++ {
		planOutput := tftest.Plan(t, terraformOptions)
		if !strings.Contains(planOutput, "to add") && !strings.Contains(planOutput, "to update") && !strings.Contains(planOutput, "to delete") {
			break
		}
		fmt.Printf("Plan was not empty, retrying after delay (%d/%d)\n", i+1, retryCount)
		time.Sleep(retryDelay)
	}

	// Assert that the plan is empty after retries
	finalPlanOutput := tftest.Plan(t, terraformOptions)
	assert.NotContains(t, finalPlanOutput, "to add")
	assert.NotContains(t, finalPlanOutput, "to update")
	assert.NotContains(t, finalPlanOutput, "to delete")
}

func testHelpersSetTerraformOptions(t *testing.T, directoryPath string, variables map[string]interface{}) (terraformOptions *tftest.Options) {
	return tftest.WithDefaultRetryableErrors(t, &tftest.Options{
		TerraformDir: directoryPath,
		Vars:         variables,
	})
}

func testHelpersSetIntegrationTestTerraformConfigDirectory(directoryLabel string) string {
	return fmt.Sprintf("./integration_test_configuration/%s", directoryLabel)
}
