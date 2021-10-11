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
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
)

var (
	defaultTags = []string{
		"account",
		"accountId",
		"language",
		"trustedAccountId",
		"guid",
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

	guid := common.EntityGUID(d.Get("guid").(string))
	tags := expandEntityTags(d.Get("tag").(*schema.Set).List())

	_, err := client.Entities.TaggingAddTagsToEntityWithContext(ctx, guid, tags)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(guid))

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		t, err := client.Entities.GetTagsForEntity(guid)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error retrieving entity tags for guid %s: %s", d.Id(), err))
		}

		currentTags := convertTagTypes(t)

		for _, t := range tags {
			var tag *entities.TaggingTagInput
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

	t, err := client.Entities.GetTagsForEntity(common.EntityGUID(d.Id()))

	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	tags := convertTagTypes(t)

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

	_, err := client.Entities.TaggingReplaceTagsOnEntityWithContext(ctx, common.EntityGUID(d.Id()), tags)
	if err != nil {
		return diag.FromErr(err)
	}

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		t, err := client.Entities.GetTagsForEntity(common.EntityGUID(d.Id()))
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error retrieving entity tags for guid %s: %s", d.Id(), err))
		}

		currentTags := convertTagTypes(t)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error retrieving entity tags for guid %s: %s", d.Id(), err))
		}

		for _, t := range tags {
			var tag *entities.TaggingTagInput
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

// This is needed until the client implements a GetTags method with the same
// tag type as the rest of the methods.
func convertTagTypes(tags []*entities.EntityTag) []*entities.TaggingTagInput {
	var t []*entities.TaggingTagInput
	for _, tag := range tags {
		t = append(t, &entities.TaggingTagInput{
			Key:    tag.Key,
			Values: tag.Values,
		})
	}

	return t
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

	_, err := client.Entities.TaggingDeleteTagFromEntityWithContext(ctx, common.EntityGUID(d.Id()), tagKeys)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandEntityTags(tags []interface{}) []entities.TaggingTagInput {
	out := make([]entities.TaggingTagInput, len(tags))

	for i, rawCfg := range tags {
		cfg := rawCfg.(map[string]interface{})
		expanded := entities.TaggingTagInput{
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

func flattenEntityTags(d *schema.ResourceData, tags []*entities.TaggingTagInput) error {
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

func getTagKeys(tags []entities.TaggingTagInput) []string {
	tagKeys := []string{}

	for _, t := range tags {
		tagKeys = append(tagKeys, t.Key)
	}
	return tagKeys
}

func tagValuesExist(t *entities.TaggingTagInput, values []string) bool {
	for _, v := range values {
		if !stringInSlice(t.Values, v) {
			return false
		}
	}

	return true
}

func getTag(tags []*entities.TaggingTagInput, key string) *entities.TaggingTagInput {
	for _, t := range tags {
		if t.Key == key {
			return t
		}
	}

	return nil
}
