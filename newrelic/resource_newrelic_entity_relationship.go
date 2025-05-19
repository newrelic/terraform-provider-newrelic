package newrelic

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entityrelationship"
	"log"
	"reflect"
	"strings"
)

func resourceNewRelicEntityRelationship() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicEntityRelationshipCreateOrUpdate,
		ReadContext:   resourceNewRelicEntityRelationshipRead,
		UpdateContext: resourceNewRelicEntityRelationshipCreateOrUpdate,
		DeleteContext: resourceNewRelicEntityRelationshipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"source_entity_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the source entity to tag.",
			},
			"target_entity_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the target entity to tag.",
			},
			"relation_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(getEdgeTypeStrings(entityrelationship.EntityRelationshipEdgeTypeTypes), false),
				Description:  "The type of relationship to create between the source and target entities.",
			},
		},
	}
}

func resourceNewRelicEntityRelationshipCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	sourceEntityGuid := common.EntityGUID(d.Get("source_entity_guid").(string))
	targetEntityGuid := common.EntityGUID(d.Get("target_entity_guid").(string))
	relationType := d.Get("relation_type").(string)

	log.Printf("[INFO] Creating Entity Relationship between source entity with GUID %+v and target entity with guid %+v having relation type as %+v", sourceEntityGuid, targetEntityGuid, relationType)

	_, err := client.EntityRelationship.EntityRelationshipUserDefinedCreateOrReplace(sourceEntityGuid, targetEntityGuid, entityrelationship.EntityRelationshipEdgeType(relationType))

	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%s:%s", sourceEntityGuid, targetEntityGuid)

	d.SetId(id)
	if err = d.Set("source_entity_guid", sourceEntityGuid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("target_entity_guid", targetEntityGuid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("relation_type", relationType); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicEntityRelationshipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	sourceEntityGUID, targetEntityGUID, err := getEntityRelationshipGUIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading New Relic entity relationship for guids %s", d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(sourceEntityGUID))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("no New Relic application found with given guid %s", sourceEntityGUID))
	}

	var dig diag.Diagnostics
	switch (*resp).(type) {
	case *entities.ExternalEntity:
		relatedEntities := (*resp).(*entities.ExternalEntity).RelatedEntities
		for _, relationship := range relatedEntities.Results {
			if userDefinedEdge, ok := relationship.(*entities.EntityRelationshipUserDefinedEdge); ok {
				// Now 'userDefinedEdge' is a pointer to an EntityRelationshipUserDefinedEdge and you can access its fields.
				if userDefinedEdge.Target.GUID == common.EntityGUID(targetEntityGUID) && userDefinedEdge.Type == d.Get("relation_type") {
					//Set the relationship related fields in the resource data
					if err = d.Set("source_entity_guid", userDefinedEdge.Source.GUID); err != nil {
						return diag.FromErr(err)
					}
					if err = d.Set("target_entity_guid", userDefinedEdge.Target.GUID); err != nil {
						return diag.FromErr(err)
					}
					if err := d.Set("relation_type", userDefinedEdge.Type); err != nil {
						return diag.FromErr(err)
					}
					break
				}
			}
		}

	default:
		dig = diag.FromErr(fmt.Errorf("problem in retrieving application with GUID %s", sourceEntityGUID))
	}
	return dig
}

func resourceNewRelicEntityRelationshipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic entity relationship for entity guid %s", d.Id())

	sourceEntityGuid := common.EntityGUID(d.Get("source_entity_guid").(string))
	targetEntityGuid := common.EntityGUID(d.Get("target_entity_guid").(string))
	relationType := d.Get("relation_type").(string)

	_, err := client.EntityRelationship.EntityRelationshipUserDefinedDelete(sourceEntityGuid, targetEntityGuid, entityrelationship.EntityRelationshipEdgeType(relationType))

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getEdgeTypeStrings(edgeTypeStruct interface{}) []string {
	var edgeTypeStrings []string
	val := reflect.ValueOf(edgeTypeStruct)
	for i := 0; i < val.NumField(); i++ {
		edgeTypeStrings = append(edgeTypeStrings, val.Field(i).String())
	}
	return edgeTypeStrings
}

func getEntityRelationshipGUIDs(id string) (string, string, error) {
	strIDs := strings.Split(id, ":")

	if len(strIDs) != 2 {
		return "", "", errors.New("could not parse entity relationship GUIDs")
	}

	return strIDs[0], strIDs[1], nil
}
