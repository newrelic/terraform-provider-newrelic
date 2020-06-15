package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

func dataSourceNewRelicEntity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicEntityRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the entity in New Relic One.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The entity's type.",
				ValidateFunc: validation.StringInSlice([]string{"APPLICATION", "DASHBOARD", "HOST", "MONITOR", "WORKLOAD"}, false),
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The entity's domain.",
				ValidateFunc: validation.StringInSlice([]string{"APM", "BROWSER", "INFRA", "MOBILE", "SYNTH"}, false),
			},
			"tag": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "A tag applied to the entity.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag value.",
						},
					},
				},
			},
			"account_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The New Relic account ID associated with this entity.",
			},
			"domain_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The domain-specific ID of the entity (only returned for APM and Browser applications)",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique entity identifier.",
			},
		},
	}
}

func dataSourceNewRelicEntityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic entities")

	name := d.Get("name").(string)
	entityType := entities.EntityType(d.Get("type").(string))
	domain := entities.EntityDomainType(d.Get("domain").(string))
	tags := expandEntityTag(d.Get("tag").([]interface{}))

	params := entities.SearchEntitiesParams{
		Name:   name,
		Type:   entityType,
		Tags:   tags,
		Domain: domain,
	}

	entityResults, err := client.Entities.SearchEntities(params)
	if err != nil {
		return err
	}

	var entity *entities.Entity
	for _, e := range entityResults {
		if e.Name == name && e.Type == entities.Type(entityType) && e.Domain == domain {
			entity = e
			break
		}
	}

	if entity == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic One entity", name)
	}

	return flattenEntityData(entity, d)
}

func flattenEntityData(e *entities.Entity, d *schema.ResourceData) error {
	d.SetId(e.GUID)
	var err error

	err = d.Set("name", e.Name)
	if err != nil {
		return err
	}

	err = d.Set("guid", e.GUID)
	if err != nil {
		return err
	}

	err = d.Set("type", e.Type)
	if err != nil {
		return err
	}

	err = d.Set("domain", e.Domain)
	if err != nil {
		return err
	}

	err = d.Set("account_id", e.AccountID)
	if err != nil {
		return err
	}

	err = d.Set("domain_id", e.ApplicationID)
	if err != nil {
		return err
	}

	return nil
}

func expandEntityTag(cfg []interface{}) *entities.TagValue {
	if len(cfg) == 0 {
		return nil
	}

	cfgTag := cfg[0].(map[string]interface{})

	tag := &entities.TagValue{}

	if k, ok := cfgTag["key"]; ok {
		tag.Key = k.(string)
	}

	if v, ok := cfgTag["value"]; ok {
		tag.Value = v.(string)
	}

	return tag
}
