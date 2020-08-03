// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNewRelicAlertPolicy_Basic(t *testing.T) {
	resourceName := "newrelic_alert_policy.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertPolicyConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyExists(resourceName),
				),
			},
			// Test: Import
			{
				ImportState:       true,
				ImportStateIdFunc: testAccImportStateIDFunc(resourceName, strconv.Itoa(testAccountID)),
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}

func TestAccNewRelicAlertPolicy_NoDiffOnReapply(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertPolicyConfig(rName),
			},
			{
				Config:             testAccNewRelicAlertPolicyConfig(rName),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccNewRelicAlertPolicy_ResourceNotFound(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertPolicyConfig(rName),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(rName),
				Config:    testAccNewRelicAlertPolicyConfig(rName),
			},
		},
	})
}

func TestAccNewRelicAlertPolicy_WithChannels(t *testing.T) {
	resourceName := "newrelic_alert_policy.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccAlertPolicyConfigWithChannels(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyExists(resourceName),
				),
			},
		},
	})
}
