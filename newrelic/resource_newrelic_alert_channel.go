package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "(Required) The name of the channel.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validAlertChannelTypes, false),
				Description:  fmt.Sprintf("(Required) The type of channel. One of: (%s).", strings.Join(validAlertChannelTypes, ", ")),
			},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The configuration block for the alert channel.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							ForceNew:    true,
							Description: "The API key for integrating with OpsGenie.",
						},
						"auth_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							ForceNew:    true,
							Description: "Specifies an authentication password for use with a channel. Supported by the webhook channel type.",
						},
						"auth_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringInSlice([]string{"BASIC"}, false),
							ForceNew:     true,
							Description:  "Specifies an authentication method for use with a channel. Supported by the webhook channel type. Only HTTP basic authentication is currently supported via the value BASIC.",
						},
						"auth_username": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Specifies an authentication username for use with a channel. Supported by the webhook channel type.",
						},
						"base_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							ForceNew:    true,
							Description: "The base URL of the webhook destination.",
						},
						"channel": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The Slack channel to send notifications to.",
						},
						"headers": {
							Type:          schema.TypeMap,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Optional:      true,
							Sensitive:     true,
							ForceNew:      true,
							ConflictsWith: []string{"config.0.headers_string"},
							Description:   "A map of key/value pairs that represents extra HTTP headers to be sent along with the webhook payload.",
						},
						"headers_string": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							ForceNew:      true,
							ConflictsWith: []string{"config.0.headers"},
							Description:   "Use instead of headers if the desired payload is more complex than a list of key/value pairs (e.g. a set of headers that makes use of nested objects). The value provided should be a valid JSON string with escaped double quotes. Conflicts with headers.",
							// Suppress the diff shown if the differences are solely due to whitespace
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return stripWhitespace(old) == stripWhitespace(new)
							},
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							ForceNew:    true,
							Description: "The key for integrating with VictorOps.",
						},
						"include_json_attachment": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "0 or 1. Flag for whether or not to attach a JSON document containing information about the associated alert to the email that is sent to recipients.",
						},
						"payload": {
							Type:          schema.TypeMap,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Sensitive:     true,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"config.0.payload_string"},
							Description:   "A map of key/value pairs that represents the webhook payload. Must provide payload_type if setting this argument.",
						},
						"payload_string": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							ForceNew:      true,
							ConflictsWith: []string{"config.0.payload"},
							Description:   "Use instead of payload if the desired payload is more complex than a list of key/value pairs (e.g. a payload that makes use of nested objects). The value provided should be a valid JSON string with escaped double quotes. Conflicts with payload.",
							// Suppress the diff shown if the differences are solely due to whitespace
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return stripWhitespace(old) == stripWhitespace(new)
							},
						},
						"payload_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"application/json", "application/x-www-form-urlencoded"}, false),
							Description:  "Can either be application/json or application/x-www-form-urlencoded. The payload_type argument is required if payload is set.",
						},
						"recipients": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "A set of recipients for targeting notifications. Multiple values are comma separated.",
						},
						"region": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"US", "EU"}, false),
							ForceNew:     true,
							Description:  "The data center region to store your data. Valid values are US and EU. Default is US.",
						},
						"route_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							ForceNew:    true,
							Description: "The route key for integrating with VictorOps.",
						},
						"service_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							ForceNew:    true,
							Description: "Specifies the service key for integrating with Pagerduty.",
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "A set of tags for targeting notifications. Multiple values are comma separated.",
						},
						"teams": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "A set of teams for targeting notifications. Multiple values are comma separated.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							ForceNew:    true,
							Description: "Your organization's Slack URL.",
						},
						"user_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The user ID for use with the user channel type.",
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

	return resourceNewRelicAlertChannelRead(d, meta)
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
