package newrelic

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandAlertChannel(d *schema.ResourceData) (*alerts.Channel, error) {
	channel := alerts.Channel{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	config, configOk := d.GetOk("config")
	configuration, configurationOk := d.GetOk("configuration")

	if !configOk && !configurationOk {
		return nil, errors.New("alert channel requires a config or configuration attribute")
	}

	if configOk {
		c, err := expandAlertChannelConfiguration(config.([]interface{})[0].(map[string]interface{}))

		if err != nil {
			return nil, err
		}

		channel.Configuration = *c
	}

	if configurationOk {
		c, err := expandAlertChannelConfiguration(configuration.(map[string]interface{}))

		if err != nil {
			return nil, err
		}

		channel.Configuration = *c
	}

	err := validateChannelConfiguration(channel.Configuration)

	if err != nil {
		return nil, err
	}

	return &channel, nil
}

//nolint:gocyclo
func expandAlertChannelConfiguration(cfg map[string]interface{}) (*alerts.ChannelConfiguration, error) {
	config := alerts.ChannelConfiguration{}

	if apiKey, ok := cfg["api_key"]; ok {
		config.APIKey = apiKey.(string)
	}

	if authPassword, ok := cfg["auth_password"]; ok {
		config.AuthPassword = authPassword.(string)
	}

	if authUsername, ok := cfg["auth_username"]; ok {
		config.AuthUsername = authUsername.(string)
	}

	if baseURL, ok := cfg["base_url"]; ok {
		config.BaseURL = baseURL.(string)
	}

	if channel, ok := cfg["channel"]; ok {
		config.Channel = channel.(string)
	}

	if key, ok := cfg["key"]; ok {
		config.Key = key.(string)
	}

	if headers, ok := cfg["headers"]; ok {
		h := headers.(map[string]interface{})
		config.Headers = h
	}

	if headers, ok := cfg["headers_string"]; ok && headers != "" {
		s := []byte(headers.(string))
		var h map[string]interface{}
		err := json.Unmarshal(s, &h)

		if err != nil {
			return nil, err
		}

		config.Headers = h
	}

	if includeJSONAttachment, ok := cfg["include_json_attachment"]; ok {
		config.IncludeJSONAttachment = includeJSONAttachment.(string)
	}

	if payload, ok := cfg["payload"]; ok {
		p := payload.(map[string]interface{})
		config.Payload = p
	}

	if payload, ok := cfg["payload_string"]; ok && payload != "" {
		s := []byte(payload.(string))
		var p map[string]interface{}
		err := json.Unmarshal(s, &p)

		if err != nil {
			return nil, err
		}

		config.Payload = p
	}

	if payloadType, ok := cfg["payload_type"]; ok {
		config.PayloadType = payloadType.(string)
	}

	if recipients, ok := cfg["recipients"]; ok {
		config.Recipients = recipients.(string)
	}

	if region, ok := cfg["region"]; ok {
		config.Region = region.(string)
	}

	if routeKey, ok := cfg["route_key"]; ok {
		config.RouteKey = routeKey.(string)
	}

	if serviceKey, ok := cfg["service_key"]; ok {
		config.ServiceKey = serviceKey.(string)
	}

	if tags, ok := cfg["tags"]; ok {
		config.Tags = tags.(string)
	}

	if teams, ok := cfg["teams"]; ok {
		config.Teams = teams.(string)
	}

	if url, ok := cfg["url"]; ok {
		config.URL = url.(string)
	}

	if userID, ok := cfg["user_id"]; ok {
		config.UserID = userID.(string)
	}

	return &config, nil
}

func expandAlertChannelIDs(channelIDs []interface{}) []int {
	ids := make([]int, len(channelIDs))

	for i := range ids {
		ids[i] = channelIDs[i].(int)
	}

	return ids
}

func flattenAlertChannelDataSource(channel *alerts.Channel, d *schema.ResourceData) error {
	d.SetId(strconv.Itoa(channel.ID))
	d.Set("policy_ids", channel.Links.PolicyIDs)

	return flattenAlertChannel(channel, d)
}

func flattenAlertChannel(channel *alerts.Channel, d *schema.ResourceData) error {
	d.Set("name", channel.Name)
	d.Set("type", channel.Type)

	config, err := flattenAlertChannelConfiguration(&channel.Configuration, d)
	if err != nil {
		return err
	}

	configuration, err := flattenDeprecatedAlertChannelConfiguration(&channel.Configuration)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("configuration"); ok {
		if err := d.Set("configuration", configuration); err != nil {
			return err
		}
	} else {
		if err := d.Set("config", config); err != nil {
			return err
		}
	}

	return nil
}

func flattenAlertChannelConfiguration(c *alerts.ChannelConfiguration, d *schema.ResourceData) ([]interface{}, error) {
	if c == nil {
		return nil, nil
	}

	configResult := make(map[string]interface{})

	configResult["api_key"] = c.APIKey
	configResult["auth_password"] = c.AuthPassword
	configResult["auth_username"] = c.AuthUsername
	configResult["base_url"] = c.BaseURL
	configResult["channel"] = c.Channel
	configResult["key"] = c.Key
	configResult["include_json_attachment"] = c.IncludeJSONAttachment
	configResult["payload_type"] = c.PayloadType
	configResult["recipients"] = c.Recipients
	configResult["region"] = c.Region
	configResult["route_key"] = c.RouteKey
	configResult["service_key"] = c.ServiceKey
	configResult["tags"] = c.Tags
	configResult["teams"] = c.Teams
	configResult["url"] = c.URL
	configResult["user_id"] = c.UserID

	if _, ok := d.GetOk("config.0.headers"); ok {
		configResult["headers"] = c.Headers
	} else if _, ok := d.GetOk("config.0.headers_string"); ok {
		h, err := json.Marshal(c.Headers)

		if err != nil {
			return nil, err
		}

		configResult["headers_string"] = string(h)
	}

	if _, ok := d.GetOk("config.0.payload"); ok {
		configResult["payload"] = c.Payload
	} else if _, ok := d.GetOk("config.0.payload_string"); ok {
		h, err := json.Marshal(c.Payload)

		if err != nil {
			return nil, err
		}

		configResult["payload_string"] = string(h)
	}

	return []interface{}{configResult}, nil
}

func flattenDeprecatedAlertChannelConfiguration(c *alerts.ChannelConfiguration) (map[string]interface{}, error) {
	if c == nil {
		return nil, nil
	}

	configResult := make(map[string]interface{})

	configResult["api_key"] = c.APIKey
	configResult["auth_password"] = c.AuthPassword
	configResult["auth_username"] = c.AuthUsername
	configResult["base_url"] = c.BaseURL
	configResult["channel"] = c.Channel
	configResult["key"] = c.Key
	configResult["include_json_attachment"] = c.IncludeJSONAttachment
	configResult["payload_type"] = c.PayloadType
	configResult["recipients"] = c.Recipients
	configResult["region"] = c.Region
	configResult["route_key"] = c.RouteKey
	configResult["service_key"] = c.ServiceKey
	configResult["tags"] = c.Tags
	configResult["teams"] = c.Teams
	configResult["url"] = c.URL
	configResult["user_id"] = c.UserID

	headers, err := json.Marshal(c.Headers)
	if err != nil {
		return nil, err
	}

	configResult["headers"] = string(headers)

	payload, err := json.Marshal(c.Payload)
	if err != nil {
		return nil, err
	}

	configResult["payload"] = string(payload)

	return configResult, nil
}

func validateChannelConfiguration(config alerts.ChannelConfiguration) error {
	if len(config.Payload) != 0 && config.PayloadType == "" {
		return errors.New("payload_type is required when using payload")
	}

	return nil
}
