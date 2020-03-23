package newrelic

import (
	"encoding/json"
	"fmt"
	"os"
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

// envAccountID implements the DefaultFunc to allow a resource to retrieve a number from the environment.
func envAccountID() (interface{}, error) {
	if v := os.Getenv("NEWRELIC_ACCOUNT_ID"); v != "" {

		accountID, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}

		return accountID, nil
	}

	return 0, nil
}
