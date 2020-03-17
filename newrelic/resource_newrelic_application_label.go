package newrelic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicApplicationLabel() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "The `newrelic_application_label` resource is deprecated. Use at your own risk.",
		Create:             resourceNewRelicApplicationLabelCreate,
		Update:             resourceNewRelicApplicationLabelUpdate,
		Read:               resourceNewRelicApplicationLabelRead,
		Delete:             resourceNewRelicApplicationLabelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A string representing the label key/category.",
				// Case fold this attribute when diffing
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A string that will be assigned to the label.",
				// Case fold this attribute when diffing
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"links": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The resources to which label should be assigned to. At least one of the inner attributes must be set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"applications": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							Description: "An array of application IDs.",
						},
						"servers": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							Description: "An array of server IDs.",
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicApplicationLabelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	label := expandApplicationLabel(d)

	log.Printf("[INFO] Creating New Relic Application label %s:%s", label.Category, label.Name)

	_, err := client.APM.CreateLabel(label)
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{label.Category, label.Name}, ":"))

	return nil
}

func resourceNewRelicApplicationLabelUpdate(d *schema.ResourceData, meta interface{}) error {
	label := expandApplicationLabel(d)

	log.Printf("[INFO] Updating New Relic Application label %s:%s", label.Category, label.Name)
	errDelete := resourceNewRelicApplicationLabelDelete(d, meta)
	if errDelete != nil {
		return errDelete
	}
	errCreate := resourceNewRelicApplicationLabelCreate(d, meta)
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func resourceNewRelicApplicationLabelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	key := d.Id()
	log.Printf("[INFO] Reading New Relic Application label %s", key)

	label, err := client.APM.GetLabel(key)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	if label == nil {
		d.SetId("")
		return nil
	}

	return flattenApplicationLabel(label, d)
}

func resourceNewRelicApplicationLabelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	key := d.Id()

	log.Printf("[INFO] Deleting New Relic label %s", key)

	_, err := client.APM.DeleteLabel(key)
	if err != nil {
		return err
	}

	return nil
}

func expandApplicationLabel(d *schema.ResourceData) apm.Label {
	label := apm.Label{
		Category: d.Get("category").(string),
		Name:     d.Get("name").(string),
		Links:    expandLinks(d.Get("links").([]interface{})[0].(map[string]interface{})),
	}

	return label
}

func flattenApplicationLabel(label *apm.Label, d *schema.ResourceData) error {
	d.Set("category", label.Category)
	d.Set("name", label.Name)
	d.Set("links", flattenLinks(&label.Links))
	return nil
}

func flattenLinks(links *apm.LabelLinks) interface{} {
	flattenedLinks := make(map[string]interface{})

	flattenedLinks["applications"] = links.Applications
	flattenedLinks["servers"] = links.Servers

	return []interface{}{flattenedLinks}
}

func expandLinks(d map[string]interface{}) apm.LabelLinks {
	appsDef := d["applications"].(*schema.Set).List()
	serversDef := d["servers"].(*schema.Set).List()
	apps := make([]int, len(appsDef))
	servers := make([]int, len(serversDef))

	for i, t := range appsDef {
		apps[i] = t.(int)
	}
	for i, t := range serversDef {
		servers[i] = t.(int)
	}
	links := apm.LabelLinks{
		Applications: apps,
		Servers:      servers,
	}
	return links
}
