// +build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestNew(t *testing.T) {
	New(config.Config{})
}
