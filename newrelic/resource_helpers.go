package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Generates a compound ID out of a slice of strings.
// This ID could contain metadata as the last string in the slice.
// e.g. 425235:2384930:someMetadata
func resourceGenerateCompoundID(idItems []string) ([]int, error) {
	nonMetadataIDCount := len(idItems) - 1
	ids := make([]int, nonMetadataIDCount)
	for i, id := range idItems[:nonMetadataIDCount] {
		intID, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}

		ids[i] = intID
	}

	return ids, nil
}

// Handles importing of resources that might utilize a compound ID,
// more specifically when metadata might be passed in the compoundID.
//
// The `defaultIDCount` argument represents the number of items that
// make up a compound ID. This count excludes any appended metadata.
//
// e.g. "425235:2384930" contains 2 items as the default
//
// The optional `attribute` argument provides an opportunity to
// set a schema attribute with a metadata value if metadata is provided
// in the compound ID as the LAST string in the compound ID.
func resourceImportStateWithMetadata(defaultIDCount int, attribute string) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		if defaultIDCount == 0 {
			return []*schema.ResourceData{d}, nil
		}

		idItems := strings.Split(d.Id(), ":")
		idItemsCount := len(idItems)
		nonMetadataIDCount := idItemsCount - 1

		if idItemsCount < defaultIDCount {
			return []*schema.ResourceData{}, fmt.Errorf("compound ID item count cannot be less than expected default ID item count: %v", idItems)
		}

		if idItemsCount == defaultIDCount {
			return []*schema.ResourceData{d}, nil
		}

		// The last item of a compound ID is the metadata
		metadataValue := idItems[nonMetadataIDCount]

		// Generate compound ID without the metadata suffix
		// e.g. 583922:5231245 (<policyID>:<conditionID>)
		ids, err := resourceGenerateCompoundID(idItems)
		if err != nil {
			return []*schema.ResourceData{}, nil
		}

		// If an attribute is supplied, attempt to set
		// with any provided metadata
		if attribute != "" {

			// In order to set the attribute correctly, we need to make sure that we
			// understand the attribute type, so that we don't set a string to an int
			// value.  First we fetch the Get() the attribute to determine what
			// interface type is returned, and then perform any necessary conversion
			// for that attribute before finally setting the attribute.

			var v int

			x := d.Get(attribute)

			switch t := x.(type) {
			case int:
				v, err = strconv.Atoi(metadataValue)
				if err != nil {
					log.Printf("[ERROR] unable to convert type %T: %s", t, err)
					return []*schema.ResourceData{}, err
				}
				err = d.Set(attribute, v)
			default:
				err = d.Set(attribute, metadataValue)
			}

			if err != nil {
				log.Printf("[ERROR] setting attribute %s: %s", attribute, err)
				return []*schema.ResourceData{}, err
			}
		}

		d.SetId(serializeIDs(ids))

		return []*schema.ResourceData{d}, nil
	}
}

// Selects the proper accountID for usage within a resource. An account ID provided
// within a `resource` block will override a `provider` block account ID. This ensures
// resources can be scoped to specific accounts. Bear in mind those accounts must be
// accessible with the provided Personal API Key (APIKS).
func selectAccountID(providerCondig *ProviderConfig, d *schema.ResourceData) int {
	resourceAccountIDAttr := d.Get("account_id")

	if resourceAccountIDAttr != nil {
		resourceAccountID := resourceAccountIDAttr.(int)

		if resourceAccountID != 0 {
			return resourceAccountID
		}
	}

	return providerCondig.AccountID
}
