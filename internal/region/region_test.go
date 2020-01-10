// +build unit

package region

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	pairs := map[string]Region{
		"us":      US,
		"Us":      US,
		"uS":      US,
		"US":      US,
		"eu":      EU,
		"Eu":      EU,
		"eU":      EU,
		"EU":      EU,
		"staging": Staging,
		"Staging": Staging,
	}

	for k, v := range pairs {
		result := Parse(k)
		assert.Equal(t, result, v)
	}

	// Default is US
	result := Parse("")
	assert.Equal(t, result, US)
}

func TestString(t *testing.T) {
	t.Parallel()

	pairs := map[Region]string{
		US:      "US",
		EU:      "EU",
		Staging: "Staging",
	}

	for k, v := range pairs {
		result := k.String()
		assert.Equal(t, result, v)
	}

	// Verify that an uninitialized Region (should be 0) isn't known
	var unk Region
	result := unk.String()
	assert.Equal(t, result, "(Unknown)")
}
