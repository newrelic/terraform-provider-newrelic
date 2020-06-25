package newrelic

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

func TestAccNewRelicInsightsEvent_Basic(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	tNow := time.Now().Unix() * 1000
	eType := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicInsightsEventConfig(eType, tNow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInsightsEventExists(
						"newrelic_insights_event.foo",
						[]string{
							fmt.Sprintf(
								"SELECT * FROM tf_test_%s "+
									"WHERE event_test = 'checking floats' "+
									"AND a_float = 101.1 "+
									"AND timestamp = %d",
								eType, tNow),
							fmt.Sprintf(
								"SELECT * FROM tf_test_%s "+
									"WHERE event_test = 'checking ints' "+
									"AND an_int = 42 "+
									"AND timestamp = %d",
								eType, tNow),
							fmt.Sprintf(
								"SELECT * FROM tf_test_%s "+
									"WHERE event_test = 'checking strings' "+
									"AND a_string = 'a string' "+
									"AND another_string = 'another string' "+
									"AND timestamp = %d",
								eType, tNow),
						},
					),
					resource.TestCheckResourceAttr("newrelic_insights_event.foo", "event.#", "3"),
				),
			},
		},
	})
}

func testAccCheckNewRelicInsightsEventExists(n string, nrqls []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no event ID is set")
		}

		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		client := providerConfig.NewClient
		accountID := providerConfig.AccountID

		// Due to the asynchronous operation, we need to
		// wait for the events to propagate.
		time.Sleep(15 * time.Second)

		var errs []string
		for _, nrql := range nrqls {
			resp, err := client.Nrdb.Query(accountID, nrdb.Nrql(nrql))
			if err != nil {
				errs = append(errs, err.Error())
			}

			if len(resp.Results) != 1 {
				errs = append(errs, fmt.Sprintf("[Error] Insights event not found (likely due to async issue) - query: %v", nrql))
			}
		}

		if len(errs) > 0 {
			fmt.Printf("test case TestAccNewRelicInsightsEvent_Basic failed (likely due to asynchronous)")

			return fmt.Errorf("%v", strings.Join(errs, "\n"))
		}

		return nil
	}
}

func testAccCheckNewRelicInsightsEventConfig(eType string, t int64) string {
	return fmt.Sprintf(`
resource "newrelic_insights_event" "foo" {
  event {
    type = "tf_test_%[1]s"

    timestamp = %[2]d

    attribute {
      key   = "event_test"
      value = "checking strings"
    }
    attribute {
      key   = "a_string"
      value = "a string"
    }
    attribute {
      key   = "another_string"
      value = "another string"
      type  = "string"
    }
  }

  event {
    type = "tf_test_%[1]s"

    timestamp = %[2]d

    attribute {
      key   = "event_test"
      value = "checking floats"
    }
    attribute {
      key   = "a_float"
      value = 101.1
      type  = "float"
    }
  }

  event {
    type = "tf_test_%[1]s"

    timestamp = %[2]d

    attribute {
      key   = "event_test"
      value = "checking ints"
    }
    attribute {
      key   = "an_int"
      value = 42
      type  = "int"
    }
  }
}`, eType, t)
}
