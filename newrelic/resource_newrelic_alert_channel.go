package newrelic

import (
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

var alertChannelTypes = map[string][]string{
	"email": {
		"include_json_attachment",
		"recipients",
	},
	"opsgenie": {
		"api_key",
		"recipients",
		"tags",
		"teams",
	},
	"pagerduty": {
		"service_key",
	},
	"slack": {
		"channel",
		"url",
	},
	"user": {
		"user_id",
	},
	"victorops": {
		"key",
		"route_key",
	},
	"webhook": {
		"auth_password",
		"auth_username",
		"base_url",
		"headers",
		"payload_type",
		"payload",
	},
}

func resourceNewRelicAlertChannel() *schema.Resource {
	validAlertChannelTypes := make([]string, 0, len(alertChannelTypes))
	for k := range alertChannelTypes {
		validAlertChannelTypes = append(validAlertChannelTypes, k)
	}

	return &schema.Resource{
		Create: resourceNewRelicAlertChannelCreate,
		Read:   resourceNewRelicAlertChannelRead,
		// Update: Not currently supported in API
		Delete: resourceNewRelicAlertChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validAlertChannelTypes, false),
			},
			"configuration": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				//TODO: ValidateFunc: (use list of keys from map above)
				Sensitive:     true,
				Deprecated:    "use `config` block instead",
				ConflictsWith: []string{"config"},
			},
			"config": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"configuration"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"auth_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"auth_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringInSlice([]string{"BASIC"}, false),
							ForceNew:     true,
						},
						"auth_username": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"base_url": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"channel": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"headers": {
							Type:          schema.TypeMap,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Optional:      true,
							Sensitive:     true,
							ForceNew:      true,
							ConflictsWith: []string{"config.0.headers_string"},
						},
						"headers_string": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							ForceNew:      true,
							ConflictsWith: []string{"config.0.headers"},
							// Suppress the diff shown if the differences are solely due to whitespace
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return StripWhitespace(old) == StripWhitespace(new)
							},
						},
						"key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"include_json_attachment": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"payload": {
							Type:          schema.TypeMap,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Sensitive:     true,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"config.0.payload_string"},
						},
						"payload_string": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							ForceNew:      true,
							ConflictsWith: []string{"config.0.payload"},
							// Suppress the diff shown if the differences are solely due to whitespace
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return StripWhitespace(old) == StripWhitespace(new)
							},
						},
						"payload_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"application/json", "application/x-www-form-urlencoded"}, false),
						},
						"recipients": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"region": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"US", "EU"}, false),
							ForceNew:     true,
						},
						"route_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"service_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"tags": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"teams": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"url": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicAlertChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	channel, err := expandAlertChannel(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating New Relic alert channel %s", channel.Name)

	channel, err = client.Alerts.CreateChannel(*channel)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(channel.ID))

	return nil
}

func resourceNewRelicAlertChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading New Relic alert channel %v", id)

	channel, err := client.Alerts.GetChannel(int(id))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenAlertChannel(channel, d)
}

func resourceNewRelicAlertChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting New Relic alert channel %v", id)

	if _, err := client.Alerts.DeleteChannel(int(id)); err != nil {
		return err
	}

	return nil
}

// StripWhitespace removes whitespace from a string.
func StripWhitespace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
