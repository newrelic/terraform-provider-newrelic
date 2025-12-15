//go:build integration || unit || ALERTS
// +build integration unit ALERTS

//
// Test helpers
//

package newrelic

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
)

func testAccCheckNewRelicAlertCompoundConditionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no compound alert condition ID is set")
		}

		conditionID := rs.Primary.ID
		accountID := providerConfig.AccountID

		if rs.Primary.Attributes["account_id"] != "" {
			var err error
			accountID, err = strconv.Atoi(rs.Primary.Attributes["account_id"])
			if err != nil {
				return err
			}
		}

		filter := &alerts.AlertsCompoundConditionFilterInput{
			Id: &alerts.AlertsCompoundConditionIDFilter{
				Eq: &conditionID,
			},
		}

		found, err := client.Alerts.SearchCompoundConditions(accountID, filter, nil, nil)
		if err != nil {
			return err
		}

		if len(found) == 0 || found[0].ID != conditionID {
			return fmt.Errorf("compound alert condition not found: %v", conditionID)
		}

		return nil
	}
}

func testAccCheckNewRelicAlertCompoundConditionDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_compound_condition" {
			continue
		}

		conditionID := r.Primary.ID
		accountID := providerConfig.AccountID

		if r.Primary.Attributes["account_id"] != "" {
			var err error
			accountID, err = strconv.Atoi(r.Primary.Attributes["account_id"])
			if err != nil {
				return err
			}
		}

		filter := &alerts.AlertsCompoundConditionFilterInput{
			Id: &alerts.AlertsCompoundConditionIDFilter{
				Eq: &conditionID,
			},
		}

		found, err := client.Alerts.SearchCompoundConditions(accountID, filter, nil, nil)
		if err == nil && len(found) > 0 {
			return fmt.Errorf("compound alert condition still exists")
		}
	}

	return nil
}
