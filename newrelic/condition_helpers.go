package newrelic

import (
	"encoding/base64"
	"fmt"
)

// Builds a condition entity guid of the format "[accountID]|AIOPS|CONDITION|[conditionID]"
func getConditionEntityGUID(conditionID int, accountID int) string {
	rawGUID := fmt.Sprintf("%d|AIOPS|CONDITION|%d", accountID, conditionID)
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(rawGUID))
}
