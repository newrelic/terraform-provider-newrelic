package utils

import (
	"strconv"
	"strings"
)

// IntArrayToString converts an array of integers
// to a comma-separated list string -
// e.g [1, 2, 3] will be converted to "1,2,3".
func IntArrayToString(integers []int) string {
	sArray := make([]string, len(integers))

	for i, n := range integers {
		sArray[i] = strconv.Itoa(n)
	}

	return strings.Join(sArray, ",")
}
