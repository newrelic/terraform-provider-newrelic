//go:build integration || CLOUD
// +build integration CLOUD

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicCloudOciLinkAccount_Basic(t *testing.T) {
	testOciLinkAccountName := fmt.Sprintf("tf_cloud_link_account_test_oci_%s", acctest.RandString(5))
	resourceName := "newrelic_cloud_oci_link_account.foo"

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testOciTenantID := os.Getenv("INTEGRATION_TESTING_OCI_TENANT_ID")
	if testOciTenantID == "" {
		t.Skipf("INTEGRATION_TESTING_OCI_TENANT_ID must be set for this acceptance test")
	}

	testOciCompartmentOcid := os.Getenv("INTEGRATION_TESTING_OCI_COMPARTMENT_OCID")
	if testOciCompartmentOcid == "" {
		t.Skipf("INTEGRATION_TESTING_OCI_COMPARTMENT_OCID must be set for this acceptance test")
	}

	testOciClientId := os.Getenv("INTEGRATION_TESTING_OCI_CLIENT_ID")
	if testOciClientId == "" {
		t.Skipf("INTEGRATION_TESTING_OCI_CLIENT_ID must be set for this acceptance test")
	}

	testOciClientSecret := os.Getenv("INTEGRATION_TESTING_OCI_CLIENT_SECRET")
	if testOciClientSecret == "" {
		t.Skipf("INTEGRATION_TESTING_OCI_CLIENT_SECRET must be set for this acceptance test")
	}

	testOciDomainUrl := os.Getenv("INTEGRATION_TESTING_OCI_DOMAIN_URL")
	if testOciDomainUrl == "" {
		t.Skipf("INTEGRATION_TESTING_OCI_DOMAIN_URL must be set for this acceptance test")
	}

	testOciHomeRegion := os.Getenv("INTEGRATION_TESTING_OCI_HOME_REGION")
	if testOciHomeRegion == "" {
		t.Skipf("INTEGRATION_TESTING_OCI_HOME_REGION must be set for this acceptance test")
	}

	// Optional fields
	testOciRegion := os.Getenv("INTEGRATION_TESTING_OCI_REGION")
	testOciMetricStackOcid := os.Getenv("INTEGRATION_TESTING_OCI_METRIC_STACK_OCID")
	testOciIngestVaultOcid := os.Getenv("INTEGRATION_TESTING_OCI_INGEST_VAULT_OCID")
	testOciUserVaultOcid := os.Getenv("INTEGRATION_TESTING_OCI_USER_VAULT_OCID")
	testOciLoggingStackOcid := os.Getenv("INTEGRATION_TESTING_OCI_LOGGING_STACK_OCID")
	testOciInstrumentationType := "METRICS" // Default to metrics for testing

	OciLinkAccountTestConfig := map[string]string{
		"name":                 testOciLinkAccountName,
		"account_id":           strconv.Itoa(testSubAccountID),
		"tenant_id":            testOciTenantID,
		"compartment_ocid":     testOciCompartmentOcid,
		"oci_client_id":        testOciClientId,
		"oci_client_secret":    testOciClientSecret,
		"oci_domain_url":       testOciDomainUrl,
		"oci_home_region":      testOciHomeRegion,
		"oci_region":           testOciRegion,
		"metric_stack_ocid":    testOciMetricStackOcid,
		"ingest_vault_ocid":    testOciIngestVaultOcid,
		"user_vault_ocid":      testOciUserVaultOcid,
		"instrumentation_type": testOciInstrumentationType,
		"logging_stack_ocid":   testOciLoggingStackOcid,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "oci") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudOciLinkAccountDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicOciLinkAccountConfig(OciLinkAccountTestConfig, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudOciLinkAccountExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicOciLinkAccountConfig(OciLinkAccountTestConfig, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudOciLinkAccountExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oci_client_secret"},
			},
		},
	})
}

func testAccCheckNewRelicCloudOciLinkAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		resourceId, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if err != nil && linkedAccount == nil {
			return err
		}

		return nil
	}
}

func testAccCheckNewRelicCloudOciLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_oci_link_account" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked oci account still exists: #{err}")
		}
	}

	return nil
}

func testAccNewRelicOciLinkAccountConfig(OciLinkAccountTestConfig map[string]string, updated bool) string {
	if updated == true {
		OciLinkAccountTestConfig["name"] += "_updated"
	}

	config := fmt.Sprintf(`
	provider "newrelic" {
		account_id = "%s"
		alias      = "cloud-integration-provider"
	}

	resource "newrelic_cloud_oci_link_account" "foo" {
		provider               = newrelic.cloud-integration-provider
		tenant_id              = "%s"
		name                   = "%s"
		account_id             = "%s"
		compartment_ocid       = "%s"
		oci_client_id          = "%s"
		oci_client_secret      = "%s"
		oci_domain_url         = "%s"
		oci_home_region        = "%s"
		ingest_vault_ocid      = "%s"
		user_vault_ocid        = "%s"
		`,
		OciLinkAccountTestConfig["account_id"],
		OciLinkAccountTestConfig["tenant_id"],
		OciLinkAccountTestConfig["name"],
		OciLinkAccountTestConfig["account_id"],
		OciLinkAccountTestConfig["compartment_ocid"],
		OciLinkAccountTestConfig["oci_client_id"],
		OciLinkAccountTestConfig["oci_client_secret"],
		OciLinkAccountTestConfig["oci_domain_url"],
		OciLinkAccountTestConfig["oci_home_region"],
		OciLinkAccountTestConfig["ingest_vault_ocid"],
		OciLinkAccountTestConfig["user_vault_ocid"])

	// Add optional fields if they exist
	if OciLinkAccountTestConfig["oci_region"] != "" && updated == true {
		config += fmt.Sprintf(`
		oci_region             = "%s"`, OciLinkAccountTestConfig["oci_region"])
	}

	if OciLinkAccountTestConfig["metric_stack_ocid"] != "" && updated == true {
		config += fmt.Sprintf(`
		metric_stack_ocid      = "%s"`, OciLinkAccountTestConfig["metric_stack_ocid"])
	}

	if OciLinkAccountTestConfig["instrumentation_type"] != "" {
		config += fmt.Sprintf(`
		instrumentation_type   = "%s"`, OciLinkAccountTestConfig["instrumentation_type"])
	}

	if OciLinkAccountTestConfig["logging_stack_ocid"] != "" && updated == true {
		config += fmt.Sprintf(`
		logging_stack_ocid     = "%s"`, OciLinkAccountTestConfig["logging_stack_ocid"])
	}

	config += `
	}`

	return config
}
