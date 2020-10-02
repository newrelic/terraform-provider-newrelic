// +build integration

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					resource.TestCheckResourceAttr("newrelic_insights_event.foo", "event.#", "3"),
				),
			},
		},
	})
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
