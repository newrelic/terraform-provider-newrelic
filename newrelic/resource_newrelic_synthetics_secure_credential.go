package newrelic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicSyntheticsSecureCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicSyntheticsSecureCredentialCreate,
		Read:   resourceNewRelicSyntheticsSecureCredentialRead,
		Update: resourceNewRelicSyntheticsSecureCredentialUpdate,
		Delete: resourceNewRelicSyntheticsSecureCredentialDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// Case fold this attribute when diffing
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceNewRelicSyntheticsSecureCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	sc := expandSyntheticsSecureCredential(d)

	log.Printf("[INFO] Creating New Relic Synthetics secure credential %s", sc.Key)

	sc, err := client.Synthetics.AddSecureCredential(sc.Key, sc.Value, sc.Description)
	if err != nil {
		return err
	}

	d.SetId(sc.Key)
	return resourceNewRelicSyntheticsSecureCredentialRead(d, meta)
}

func resourceNewRelicSyntheticsSecureCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics secure credential %s", d.Id())

	sc, err := client.Synthetics.GetSecureCredential(d.Id())
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenSyntheticsSecureCredential(sc, d)
}

func resourceNewRelicSyntheticsSecureCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Updating New Relic Synthetics secure credential %s", d.Id())

	sc := expandSyntheticsSecureCredential(d)

	_, err := client.Synthetics.UpdateSecureCredential(sc.Key, sc.Value, sc.Description)
	if err != nil {
		return err
	}

	return resourceNewRelicSyntheticsSecureCredentialRead(d, meta)
}

func resourceNewRelicSyntheticsSecureCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic Synthetics secure credential %s", d.Id())

	if err := client.Synthetics.DeleteSecureCredential(d.Id()); err != nil {
		return err
	}

	return nil
}
