package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

var alertChannelTypes = map[string][]string{
	"campfire": {
		"room",
		"subdomain",
		"token",
	},
	"email": {
		"include_json_attachment",
		"recipients",
	},
	"hipchat": {
		"auth_token",
		"base_url",
		"room_id",
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
		"auth_type",
		"auth_username",
		"base_url",
		"payload_type",
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
				Required: true,
				ForceNew: true,
				//TODO: ValidateFunc: (use list of keys from map above)
				Sensitive: true,
			},
			"headers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Sensitive: true,
			},
			"payload": {
				Type:         schema.TypeMap,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: configBlockLengthGreaterThan(0),
			},
		},
	}
}

// This function verifies that the number of keys within a configuration
// map is greater than the provided parameter.
func configBlockLengthGreaterThan(minLength int) schema.SchemaValidateFunc {
	return func(i interface{}, st string) (s []string, errorSlice []error) {
		length := len(i.(map[string]interface{}))
		if length > minLength {
			return
		}
		errorSlice = append(errorSlice, fmt.Errorf("expected %s not to be empty", st))
		return
	}
}

func buildAlertChannelStruct(d *schema.ResourceData) *newrelic.AlertChannel {
	channel := newrelic.AlertChannel{
		Name:          d.Get("name").(string),
		Type:          d.Get("type").(string),
		Configuration: d.Get("configuration").(map[string]interface{}),
	}

	if headerMap, ok := d.GetOk("headers"); ok {
		channel.Configuration["headers"] = headerMap.(map[string]interface{})
	}

	if payload, ok := d.GetOk("payload"); ok {
		channel.Configuration["payload"] = payload.(map[string]interface{})
	}

	return &channel
}

func resourceNewRelicAlertChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	channel := buildAlertChannelStruct(d)

	log.Printf("[INFO] Creating New Relic alert channel %s", channel.Name)

	channel, err := client.CreateAlertChannel(*channel)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(channel.ID))

	return nil
}

func resourceNewRelicAlertChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading New Relic alert channel %v", id)

	channel, err := client.GetAlertChannel(int(id))
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", channel.Name)
	d.Set("type", channel.Type)

	// extract headers from Configuration before we try and set it in the resource
	if headers, ok := channel.Configuration["headers"]; ok {
		d.Set("headers", headers)
		delete(channel.Configuration, "headers")
	}

	// extract payload from Configuration before we try and set it in the resource
	if payload, ok := channel.Configuration["payload"]; ok {
		d.Set("payload", payload)
		delete(channel.Configuration, "payload")
	}

	if err := d.Set("configuration", channel.Configuration); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Alert Channel Configuration: %#v", err)
	}

	return nil
}

func resourceNewRelicAlertChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting New Relic alert channel %v", id)

	if err := client.DeleteAlertChannel(int(id)); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
