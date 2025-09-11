package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apiaccess"
)

var (
	keyTypeIngest        = string(apiaccess.APIAccessKeyTypeTypes.INGEST)
	keyTypeIngestBrowser = string(apiaccess.APIAccessIngestKeyTypeTypes.BROWSER)
	keyTypeIngestLicense = string(apiaccess.APIAccessIngestKeyTypeTypes.LICENSE)
	keyTypeUser          = string(apiaccess.APIAccessKeyTypeTypes.USER)
)

func resourceNewRelicAPIAccessKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ResourceNewRelicAPIAccessKeyAttributeLabels.AccountID:  APIAccessKeysSchemaAccountID,
		ResourceNewRelicAPIAccessKeyAttributeLabels.KeyType:    APIAccessKeysSchemaKeyType,
		ResourceNewRelicAPIAccessKeyAttributeLabels.IngestType: APIAccessKeysSchemaIngestType,
		ResourceNewRelicAPIAccessKeyAttributeLabels.UserID:     APIAccessKeysSchemaUserID,
		ResourceNewRelicAPIAccessKeyAttributeLabels.Name:       APIAccessKeysSchemaName,
		ResourceNewRelicAPIAccessKeyAttributeLabels.Notes:      APIAccessKeysSchemaNotes,
		ResourceNewRelicAPIAccessKeyAttributeLabels.Key:        APIAccessKeysSchemaKey,
	}
}

// type ResourceNewRelicAPIAccessKeyAttributeLabel string

// commented the above, as the usage of the following types as ResourceNewRelicAPIAccessKeyAttributeLabel
// and not directly as a string isn't allowing them to be recognized as strings in the schema

var ResourceNewRelicAPIAccessKeyAttributeLabels = struct {
	AccountID  string
	KeyType    string
	UserID     string
	IngestType string
	Name       string
	Notes      string
	Key        string
}{
	AccountID:  "account_id",
	KeyType:    "key_type",
	UserID:     "user_id",
	IngestType: "ingest_type",
	Name:       "name",
	Notes:      "notes",
	Key:        "key",
}

func validateAPIAccessKeyAttributes(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	var errorsList []error

	accountIDInConfiguration := getAPIAccessKeyAttributeValueFromConfiguration(
		d,
		ResourceNewRelicAPIAccessKeyAttributeLabels.AccountID,
	)

	keyTypeInConfiguration := getAPIAccessKeyAttributeValueFromConfiguration(
		d,
		ResourceNewRelicAPIAccessKeyAttributeLabels.KeyType,
	)

	ingestTypeInConfiguration := getAPIAccessKeyAttributeValueFromConfiguration(
		d,
		ResourceNewRelicAPIAccessKeyAttributeLabels.IngestType,
	)

	userIDInConfiguration := getAPIAccessKeyAttributeValueFromConfiguration(
		d,
		ResourceNewRelicAPIAccessKeyAttributeLabels.UserID,
	)

	// this is not needed since the attribute is required in the schema
	// however, this is being added as a double check, if the attribute were to support some sort of defaulting in the future
	// so the author would need to disable this check too, which would be a conscious decision
	if accountIDInConfiguration == 0 {
		errorsList = append(errorsList, fmt.Errorf("the `account_id` attribute is required, please specify a non-null account ID"))
	}

	if keyTypeInConfiguration != keyTypeIngest && keyTypeInConfiguration != keyTypeUser {
		errorsList = append(errorsList, fmt.Errorf("the `key_type` attribute must be set to either %s or %s", keyTypeIngest, keyTypeUser))
	}

	if keyTypeInConfiguration == keyTypeIngest {
		// extra function to ensure this works after the first apply too, and no false positives are returned by d.GetChange(), and the config is truly pointed to
		if !d.GetRawConfig().GetAttr(ResourceNewRelicAPIAccessKeyAttributeLabels.UserID).IsNull() {
			errorsList = append(errorsList, fmt.Errorf("the `user_id` attribute cannot be set when `key_type` is set to INGEST, please remove the `user_id` attribute from the configuration, or change `key_type` to USER"))
		}
		if ingestTypeInConfiguration == "" {
			errorsList = append(errorsList, fmt.Errorf("the `ingest_type` attribute is required when `key_type` is set to INGEST, please specify %s or %s as the the value of `ingest_type`", keyTypeIngestLicense, keyTypeIngestBrowser))
		} else if ingestTypeInConfiguration != keyTypeIngestLicense && ingestTypeInConfiguration != keyTypeIngestBrowser {
			errorsList = append(errorsList, fmt.Errorf("the `ingest_type` attribute must be set to either %s or %s when `key_type` is set to INGEST", keyTypeIngestLicense, keyTypeIngestBrowser))
		}
	}

	if keyTypeInConfiguration == keyTypeUser {
		// extra function to ensure this works after the first apply too, and no false positives are returned by d.GetChange(), and the config is truly pointed to
		if !d.GetRawConfig().GetAttr(ResourceNewRelicAPIAccessKeyAttributeLabels.IngestType).IsNull() {
			errorsList = append(errorsList, fmt.Errorf("the `ingest_type` attribute cannot be set when `key_type` is set to USER, please remove the `ingest_type` attribute from the configuration, or change `key_type` to INGEST"))
		}
		if userIDInConfiguration == 0 {
			errorsList = append(errorsList, fmt.Errorf("the `user_id` attribute is required when `key_type` is set to USER, please specify a valid user ID"))
		}
	}

	if len(errorsList) == 0 {
		return nil
	}

	errorsString := "the following validation errors have been identified with the configuration of the synthetic monitor: \n"

	for index, val := range errorsList {
		errorsString += fmt.Sprintf("(%d): %s\n", index+1, val)
	}

	return errors.New(errorsString)
}

func getAPIAccessKeyAccountID(d *schema.ResourceData) int {
	accountID := 0

	if v, ok := d.GetOk(ResourceNewRelicAPIAccessKeyAttributeLabels.AccountID); ok {
		accountID = v.(int)
	}

	log.Printf("[DEBUG] new api access account_id: %d", accountID)

	// 0 cannot be returned (though an initial value) in the event of `account_id` not being set in the configuration,
	// since the function validatedAPIAccessKeyAttributes prevents `account_id` from being null or 0
	return accountID
}

func getAPIAccessKeyType(d *schema.ResourceData) string {
	keyType := ""
	if v, ok := d.GetOk(ResourceNewRelicAPIAccessKeyAttributeLabels.KeyType); ok {
		keyType = v.(string)
	}

	log.Printf("[DEBUG] api access key_type: %s", keyType)

	return keyType
}

func getAPIAccessIngestType(d *schema.ResourceData) string {
	ingestType := ""
	if v, ok := d.GetOk(ResourceNewRelicAPIAccessKeyAttributeLabels.IngestType); ok {
		ingestType = v.(string)
	}

	log.Printf("[DEBUG] api access ingest_type: %s", ingestType)

	return ingestType
}

func getAPIAccessUserID(d *schema.ResourceData) int {
	userID := 0
	if v, ok := d.GetOk(ResourceNewRelicAPIAccessKeyAttributeLabels.UserID); ok {
		userID = v.(int)
	}

	log.Printf("[DEBUG] api access user_id: %d", userID)

	// 0 cannot be returned (though an initial value) in the event of `user_id` not being set in the configuration,
	// since the function validatedAPIAccessKeyAttributes prevents `user_id` from being null or 0 when key_type = USER
	return userID
}

func getAPIAccessKeyNotes(d *schema.ResourceData) string {
	notes := ""
	if v, ok := d.GetOk(ResourceNewRelicAPIAccessKeyAttributeLabels.Notes); ok {
		notes = v.(string)
	}

	log.Printf("[DEBUG] api access key notes: %s", notes)

	return notes
}

func getAPIAccessKeyName(d *schema.ResourceData) string {
	name := ""
	if v, ok := d.GetOk(ResourceNewRelicAPIAccessKeyAttributeLabels.Name); ok {
		name = v.(string)
	}

	log.Printf("[DEBUG] api access key name: %s", name)

	return name
}

func getAPIAccessKeyAttributeValueFromConfiguration(
	d *schema.ResourceDiff,
	attributeLabel string,
) interface{} {
	_, value := d.GetChange(attributeLabel)
	return value
}

var APIAccessKeysSchemaAccountID = &schema.Schema{
	Type: schema.TypeInt,
	// this is still 'Required', as applied in the custom validation written in validateAPIAccessKeyAttributes
	// however, setting this to Optional allows any future change to allow defaulting behaviour simpler
	Optional: true,
	ForceNew: true,
}

var APIAccessKeysSchemaKeyType = &schema.Schema{
	Type:     schema.TypeString,
	Required: true,
	ForceNew: true,
	// validation logic shifted to validateAPIAccessKeyAttributes
}

var APIAccessKeysSchemaIngestType = &schema.Schema{
	Type:          schema.TypeString,
	Optional:      true,
	ForceNew:      true,
	Computed:      true,
	ConflictsWith: []string{ResourceNewRelicAPIAccessKeyAttributeLabels.UserID},
	// validation logic shifted to validateAPIAccessKeyAttributes
}

var APIAccessKeysSchemaUserID = &schema.Schema{
	Type:          schema.TypeInt,
	Optional:      true,
	ForceNew:      true,
	Computed:      true,
	ConflictsWith: []string{ResourceNewRelicAPIAccessKeyAttributeLabels.IngestType},
	// validation logic shifted to validateAPIAccessKeyAttributes
}

var APIAccessKeysSchemaName = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	Computed: true,
}

var APIAccessKeysSchemaNotes = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	Computed: true,
}

var APIAccessKeysSchemaKey = &schema.Schema{
	Type:     schema.TypeString,
	Computed: true,
}
