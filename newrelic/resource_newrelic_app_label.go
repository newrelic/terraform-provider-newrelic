package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"log"
	"strings"
)

func resourceNewRelicAppLabel() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicAppLabelCreate,
		Update: resourceNewRelicAppLabelUpdate,
		Read:   resourceNewRelicAppLabelRead,
		Delete: resourceNewRelicAppLabelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Required: true,
				// Case fold this attribute when diffing
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// Case fold this attribute when diffing
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"links": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"applications": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Optional: true,
						},
						"servers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicAppLabelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	label := expandAppLabel(d)

	log.Printf("[INFO] Creating New Relic App label %s:%s", label.Category, label.Name)

	_, err := client.APM.CreateLabel(label)
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{label.Category, label.Name}, ":"))

	return nil
}

func resourceNewRelicAppLabelUpdate(d *schema.ResourceData, meta interface{}) error {
	label := expandAppLabel(d)

	log.Printf("[INFO] Updating New Relic App label %s:%s", label.Category, label.Name)
	errDelete:= resourceNewRelicAppLabelDelete(d, meta)
	if errDelete != nil {
		return errDelete
	}
	errCreate:=resourceNewRelicAppLabelCreate(d,meta)
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func resourceNewRelicAppLabelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	key := d.Id()
	log.Printf("[INFO] Reading New Relic App label %s", key)

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

	return flattenAppLabel(label, d)
}

func resourceNewRelicAppLabelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	key := d.Id()

	log.Printf("[INFO] Deleting New Relic label %s", key)

	_, err := client.APM.DeleteLabel(key)
	if err != nil {
		return err
	}

	return nil
}

func expandAppLabel(d *schema.ResourceData) apm.Label {
	label := apm.Label{
		Category: d.Get("category").(string),
		Name:     d.Get("name").(string),
		Links:    expandLinks(d.Get("links").([]interface{})[0].(map[string]interface{})),
	}

	return label
}

func flattenAppLabel(label *apm.Label, d *schema.ResourceData) error {
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
	appsDef := d["applications"].([]interface{})
	serversDef := d["servers"].([]interface{})
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
