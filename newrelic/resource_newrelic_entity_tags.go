package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
)

var (
	defaultTags = []string{
		"account",
		"accountId",
		"language",
		"trustedAccountId",
	}
)

func resourceNewRelicEntityTags() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicEntityTagsCreate,
		ReadContext:   resourceNewRelicEntityTagsRead,
		UpdateContext: resourceNewRelicEntityTagsUpdate,
		DeleteContext: resourceNewRelicEntityTagsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The guid of the entity to tag.",
			},
			"tag": {
				Type:        schema.TypeSet,
				MinItems:    1,
				Required:    true,
				Description: "A set of key-value pairs to represent a tag. For example: Team:TeamName",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag key.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							MinItems:    1,
							Required:    true,
							Description: "The tag values.",
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
		},
	}
}

func resourceNewRelicEntityTagsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Create")
	}

	client := providerConfig.NewClient

	guid := entities.EntityGUID(d.Get("guid").(string))
	tags := expandEntityTags(d.Get("tag").(*schema.Set).List())

	if err := client.Entities.AddTags(guid, tags); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(guid))

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		currentTags, err := client.Entities.ListTags(guid)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error retrieving entity tags for guid %s: %s", d.Id(), err))
		}

		for _, t := range tags {
			var tag *entities.Tag
			if tag = getTag(currentTags, t.Key); tag == nil {
				return resource.RetryableError(fmt.Errorf("expected entity tag %s to have been updated but was not found", t.Key))
			}

			if ok := tagValuesExist(tag, t.Values); !ok {
				return resource.RetryableError(fmt.Errorf("expected entity tag values %s to have been updated for tag %s but were not found", t.Values, t.Key))
			}
		}

		diag := resourceNewRelicEntityTagsRead(ctx, d, meta)
		if diag.HasError() {
			return resource.RetryableError(errors.New("error reading tag values after creation"))
		}

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return nil
}

func resourceNewRelicEntityTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Read")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic entity tags for entity guid %s", d.Id())

	tags, err := client.Entities.ListTags(entities.EntityGUID(d.Id()))

	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenEntityTags(d, tags))
}

func resourceNewRelicEntityTagsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Update")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Updating New Relic entity tags for entity guid %s", d.Id())

	tags := expandEntityTags(d.Get("tag").(*schema.Set).List())

	if err := client.Entities.ReplaceTags(entities.EntityGUID(d.Id()), tags); err != nil {
		return diag.FromErr(err)
	}

	retryErr := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		currentTags, err := client.Entities.ListTags(entities.EntityGUID(d.Id()))

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error retrieving entity tags for guid %s: %s", d.Id(), err))
		}

		for _, t := range tags {
			var tag *entities.Tag
			if tag = getTag(currentTags, t.Key); tag == nil {
				return resource.RetryableError(fmt.Errorf("expected entity tag %s to have been updated but was not found", t.Key))
			}

			if ok := tagValuesExist(tag, t.Values); !ok {
				return resource.RetryableError(fmt.Errorf("expected entity tag values %s to have been created for tag %s but were not found", t.Values, t.Key))
			}
		}

		diag := resourceNewRelicEntityTagsRead(ctx, d, meta)
		if diag.HasError() {
			return resource.RetryableError(errors.New("error reading tag values after creation"))
		}

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return nil
}

func resourceNewRelicEntityTagsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Delete")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic entity tags from entity guid %s", d.Id())

	tags := expandEntityTags(d.Get("tag").(*schema.Set).List())
	tagKeys := getTagKeys(tags)

	if err := client.Entities.DeleteTags(entities.EntityGUID(d.Id()), tagKeys); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandEntityTags(tags []interface{}) []entities.Tag {
	out := make([]entities.Tag, len(tags))

	for i, rawCfg := range tags {
		cfg := rawCfg.(map[string]interface{})
		expanded := entities.Tag{
			Key:    cfg["key"].(string),
			Values: expandEntityTagValues(cfg["values"].(*schema.Set).List()),
		}

		out[i] = expanded
	}

	return out
}

func expandEntityTagValues(values []interface{}) []string {
	perms := make([]string, len(values))

	for i, v := range values {
		perms[i] = v.(string)
	}

	return perms
}

func flattenEntityTags(d *schema.ResourceData, tags []*entities.Tag) error {
	out := []map[string]interface{}{}
	for _, t := range tags {
		if stringInSlice(defaultTags, t.Key) {
			continue
		}

		m := make(map[string]interface{})
		m["key"] = t.Key
		m["values"] = t.Values

		out = append(out, m)
	}

	if err := d.Set("guid", d.Id()); err != nil {
		return err
	}

	if err := d.Set("tag", out); err != nil {
		return err
	}

	return nil
}

func getTagKeys(tags []entities.Tag) []string {
	tagKeys := []string{}

	for _, t := range tags {
		tagKeys = append(tagKeys, t.Key)
	}
	return tagKeys
}

func tagValuesExist(t *entities.Tag, values []string) bool {
	for _, v := range values {
		if !stringInSlice(t.Values, v) {
			return false
		}
	}

	return true
}

func getTag(tags []*entities.Tag, key string) *entities.Tag {
	for _, t := range tags {
		if t.Key == key {
			return t
		}
	}

	return nil
}
