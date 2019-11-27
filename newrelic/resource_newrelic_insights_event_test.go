package newrelic

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicInsightsEvent_Basic(t *testing.T) {
	tNow := time.Now().Unix() * 1000
	eType := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

		client := testAccProvider.Meta().(*ProviderConfig).InsightsQueryClient

		for _, nrql := range nrqls {
			resp, err := client.QueryEvents(nrql)
			if err != nil {
				return err
			}
			if len(resp.Results) != 1 {
				return errors.New("did not find Insights event")
			}
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
