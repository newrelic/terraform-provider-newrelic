package newrelic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/contextkeys"
)

// Selects the proper accountID for usage within a resource. An account ID provided
// within a `resource` block will override a `provider` block account ID. This ensures
// resources can be scoped to specific accounts. Bear in mind those accounts must be
// accessible with the provided Personal API Key (APIKS).
func selectAccountID(providerConfig *ProviderConfig, d *schema.ResourceData) int {
	resourceAccountIDAttr := d.Get("account_id")

	if resourceAccountIDAttr != nil {
		resourceAccountID := resourceAccountIDAttr.(int)

		if resourceAccountID != 0 {
			return resourceAccountID
		}
	}

	return providerConfig.AccountID
}

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

func updateContextWithAccountID(ctx context.Context, accountID int) context.Context {
	if accountID > 0 {
		log.Printf("[INFO] Adding Account ID to X-Account-ID context %v", accountID)
		value := strconv.Itoa(accountID)

		updatedCtx := contextkeys.SetAccountID(ctx, value)
		return updatedCtx
	}

	return ctx
}

func mergeSchemas(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	schema := map[string]*schema.Schema{}
	for _, s := range schemas {
		for k, v := range s {
			schema[k] = v
		}
	}
	return schema
}

// This method helps identify single quotes in the argument 'name', to
// prefix the '\' escape character before single quotes, in order to allow
// NRQL to parse the query without errors caused by the single quote '.
func escapeSingleQuote(name string) string {
	unescapedSingleQuoteRegex := regexp.MustCompile(`'`)
	quoteFormattedName := unescapedSingleQuoteRegex.ReplaceAllString(name, "\\'")
	if strings.Compare(quoteFormattedName, name) != 0 {
		log.Printf("Changing the name ( %s ---> %s ) since a single quote has been identified.", name, quoteFormattedName)
		name = quoteFormattedName
	}

	return name
}

// This method helps identify escape characters preceding single quotes (')
// and eliminates the escape characters to match it with the NRQL expression
// parsed by the Terraform Provider.
func revertEscapedSingleQuote(name string) string {
	escapedSingleQuoteRegex := regexp.MustCompile(`\\'`)
	quoteFormattedName := escapedSingleQuoteRegex.ReplaceAllString(name, "'")
	if strings.Compare(quoteFormattedName, name) != 0 {
		log.Printf("Reverting the name ( %s ---> %s ) since a single quote has been identified.", quoteFormattedName, name)
		name = quoteFormattedName
	}

	return name
}

// This methods is a wrapper for Resource Data getter function
func fetchAttributeValueFromResourceConfig(d *schema.ResourceData, key string) (interface{}, bool) {
	return d.GetOk(key)
}

// Builds a condition entity guid of the format "[accountID]|AIOPS|CONDITION|[conditionID]"
func getConditionEntityGUID(conditionID int, accountID int) string {
	rawGUID := fmt.Sprintf("%d|AIOPS|CONDITION|%d", accountID, conditionID)
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(rawGUID))
}

func getBoolPointer(value bool) *bool {
	return &value
}

func getFloatPointer(value float64) *float64 {
	return &value
}

func convertInterfaceToStringSlice(v interface{}) []string {
	var result []string
	for _, item := range v.([]interface{}) {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}
