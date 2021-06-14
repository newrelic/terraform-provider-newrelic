package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

func dataSourceNewRelicEntity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicEntityRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the entity in New Relic One.  The first entity matching this name for the given search parameters will be returned.",
			},
			"ignore_case": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Ignore case when searching the entity name.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The entity's type. Valid values are APPLICATION, DASHBOARD, HOST, MONITOR, and WORKLOAD.",
				ValidateFunc: validation.StringInSlice([]string{"APPLICATION", "DASHBOARD", "HOST", "MONITOR", "WORKLOAD"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The entity's domain. Valid values are APM, BROWSER, INFRA, MOBILE, SYNTH, and VIZ. If not specified, all domains are searched.",
				ValidateFunc: validation.StringInSlice([]string{"APM", "BROWSER", "INFRA", "MOBILE", "SYNTH", "VIZ"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
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
			"application_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The domain-specific ID of the entity (only returned for APM and Browser applications).",
			},
			"serving_apm_application_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The browser-specific ID of the backing APM entity. (only returned for Browser applications)",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique entity identifier.",
			},
		},
	}
}

func dataSourceNewRelicEntityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic entities")

	name := d.Get("name").(string)
	ignoreCase := d.Get("ignore_case").(bool)
	entityType := entities.EntitySearchQueryBuilderType(strings.ToUpper(d.Get("type").(string)))
	tags := expandEntityTag(d.Get("tag").([]interface{}))
	domain := entities.EntitySearchQueryBuilderDomain(strings.ToUpper(d.Get("domain").(string)))

	params := entities.EntitySearchQueryBuilder{
		Name:   name,
		Type:   entityType,
		Tags:   tags,
		Domain: domain,
	}

	entityResults, err := client.Entities.GetEntitySearchWithContext(ctx, entities.EntitySearchOptions{}, "", params, []entities.EntitySearchSortCriteria{})
	if err != nil {
		return diag.FromErr(err)
	}

	var entity *entities.EntityOutlineInterface
	for _, e := range entityResults.Results.Entities {
		// Conditional on case sensitive match
		if e.GetName() == name || (ignoreCase && strings.EqualFold(e.GetName(), name)) {
			entity = &e
			break
		}
	}

	if entity == nil {
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic One entity for the given search parameters (ignore_case: %t)", name, ignoreCase))
	}

	return diag.FromErr(flattenEntityData(entity, d))
}

func flattenEntityData(entity *entities.EntityOutlineInterface, d *schema.ResourceData) error {
	var err error

	d.SetId(string((*entity).GetGUID()))

	if err = d.Set("name", (*entity).GetName()); err != nil {
		return err
	}

	if err = d.Set("guid", (*entity).GetGUID()); err != nil {
		return err
	}

	if err = d.Set("type", (*entity).GetType()); err != nil {
		return err
	}

	if err = d.Set("domain", (*entity).GetDomain()); err != nil {
		return err
	}

	if err = d.Set("account_id", (*entity).GetAccountID()); err != nil {
		return err
	}

	// store extra values per Entity Type, have to repeat code here due to
	// go handling of type switching
	switch e := (*entity).(type) {
	case *entities.ApmApplicationEntityOutline:
		if err = d.Set("application_id", e.ApplicationID); err != nil {
			return err
		}
	case *entities.MobileApplicationEntityOutline:
		if err = d.Set("application_id", e.ApplicationID); err != nil {
			return err
		}
	case *entities.BrowserApplicationEntityOutline:
		if err = d.Set("application_id", e.ApplicationID); err != nil {
			return err
		}

		if e.ServingApmApplicationID > 0 {
			err = d.Set("serving_apm_application_id", e.ServingApmApplicationID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func expandEntityTag(cfg []interface{}) []entities.EntitySearchQueryBuilderTag {
	var tags []entities.EntitySearchQueryBuilderTag

	if len(cfg) == 0 {
		return tags
	}

	tags = make([]entities.EntitySearchQueryBuilderTag, 0, len(cfg))

	for _, t := range cfg {
		cfgTag := t.(map[string]interface{})

		tag := entities.EntitySearchQueryBuilderTag{}

		if k, ok := cfgTag["key"]; ok {
			tag.Key = k.(string)
			if v, ok := cfgTag["value"]; ok {
				tag.Value = v.(string)

				tags = append(tags, tag)
			}
		}
	}

	return tags
}
