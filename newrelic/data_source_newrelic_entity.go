package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func dataSourceNewRelicEntity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicEntityRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the entity in New Relic One. The first entity matching this name for the given search parameters will be returned.",
			},
			"ignore_case": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Ignore case when searching the entity name.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The entity's type. Valid values are APPLICATION, DASHBOARD, HOST, MONITOR, SERVICE and WORKLOAD.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The entity's domain. Valid values are APM, BROWSER, INFRA, MOBILE, SYNTH, and EXT. If not specified, all domains are searched.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // Case fold this attribute when diffing
				},
			},
			"tag": {
				Type:        schema.TypeList,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The New Relic account ID; if specified, constrains the data source to return an entity belonging to the account with this ID, of all matching entities retrieved.",
				ValidateFunc: validation.IntAtLeast(1),
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
	accountID := meta.(*ProviderConfig).AccountID
	if acc, ok := d.GetOk("account_id"); ok {
		accountID = acc.(int)
	}

	log.Printf("[INFO] Reading New Relic entities")

	name := d.Get("name").(string)
	name = escapeSingleQuote(name)
	ignoreCase := d.Get("ignore_case").(bool)
	entityType := strings.ToUpper(d.Get("type").(string))
	domain := strings.ToUpper(d.Get("domain").(string))
	tags := d.Get("tag").([]interface{})

	query := buildEntitySearchQuery(name, domain, entityType, tags)

	entityResults, err := client.Entities.GetEntitySearchByQueryWithContext(ctx,
		entities.EntitySearchOptions{
			CaseSensitiveTagMatching: ignoreCase,
		},
		query,
		[]entities.EntitySearchSortCriteria{},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	if entityResults == nil {
		return diag.FromErr(fmt.Errorf("GetEntitySearchByQuery response was nil"))
	}

	var entity *entities.EntityOutlineInterface
	for _, e := range entityResults.Results.Entities {
		// Conditional on case-sensitive match

		str := e.GetName()
		str = strings.TrimSpace(str)

		name = revertEscapedSingleQuote(name)
		if strings.Compare(str, name) == 0 || (ignoreCase && strings.EqualFold(str, name)) {
			if e.GetAccountID() != accountID {
				continue
			} else {
				entity = &e
				break
			}

		}
	}

	if entity == nil {
		return diag.FromErr(fmt.Errorf("no entities found for the provided search parameters, please ensure your schema attributes are valid"))
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

func buildEntitySearchQuery(name string, domain string, entityType string, tags []interface{}) string {
	var query string

	if name != "" {
		query = fmt.Sprintf("name = '%s'", name)
	}

	if domain != "" {
		query = fmt.Sprintf("%s AND domain = '%s'", query, domain)
	}

	if entityType != "" {
		query = fmt.Sprintf("%s AND type = '%s'", query, entityType)
	}

	if len(tags) > 0 {
		query = fmt.Sprintf("%s AND %s", query, buildTagsQueryFragment(tags))
	}

	return query
}

func buildTagsQueryFragment(tags []interface{}) string {
	var query string

	for i, t := range tags {
		tag := t.(map[string]interface{})

		var q string
		if i > 0 {
			q = fmt.Sprintf(" AND tags.`%s` = '%s'", tag["key"], tag["value"].(string))
		} else {
			q = fmt.Sprintf("tags.`%s` = '%s'", tag["key"], tag["value"].(string))
		}

		query = fmt.Sprintf("%s%s", query, q)
	}

	return query
}
