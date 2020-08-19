package newrelic

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func parseIDs(serializedID string, count int) ([]int, error) {
	rawIDs := strings.SplitN(serializedID, ":", count)
	if len(rawIDs) != count {
		return []int{}, fmt.Errorf("unable to parse ID %v", serializedID)
	}

	ids := make([]int, count)

	for i, rawID := range rawIDs {
		id, err := strconv.ParseInt(rawID, 10, 32)
		if err != nil {
			return ids, err
		}

		ids[i] = int(id)
	}

	return ids, nil
}

func parseCompositeID(id string) (p1 string, p2 string, err error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		p1 = parts[0]
		p2 = parts[1]
	} else {
		err = fmt.Errorf("error: Import composite ID requires two parts separated by colon, eg x:y")
	}
	return
}

// Converts a hash of IDs into an array.
// Examples: "12345:54432:66564" -> []int{12345,54432,66564}
func parseHashedIDs(serializedID string) ([]int, error) {
	rawIDs := strings.Split(serializedID, ":")
	ids := make([]int, len(rawIDs))

	for i, rawID := range rawIDs {
		id, err := strconv.ParseInt(rawID, 10, 32)
		if err != nil {
			return ids, err
		}

		ids[i] = int(id)
	}

	return ids, nil
}

func serializeIDs(ids []int) string {
	idStrings := make([]string, len(ids))

	for i, id := range ids {
		idStrings[i] = strconv.Itoa(id)
	}

	return strings.Join(idStrings, ":")
}

// Helper for converting data to pretty JSON
// nolint:deadcode,unused
func toJSON(data interface{}) string {
	c, _ := json.MarshalIndent(data, "", "  ")

	return string(c)
}

func stripWhitespace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// Mutates original slice
func sortIntegerSlice(integers []int) {
	sort.Slice(integers, func(i, j int) bool {
		return integers[i] < integers[j]
	})
}

func stringInSlice(slice []string, str string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}

	return false
}
