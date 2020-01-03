package utils

import (
	"strconv"
	"strings"
)

// IntArrayToString converts an array of integers
// to a comma-separated list string -
// e.g [1, 2, 3] will be converted to "1,2,3".
func IntArrayToString(integers []int) string {
	sArray := []string{}

	for _, n := range integers {
		sArray = append(sArray, strconv.Itoa(n))
	}

	return strings.Join(sArray, ",")
}
