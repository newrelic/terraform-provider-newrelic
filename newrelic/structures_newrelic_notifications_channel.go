package newrelic

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
)

func expandNotificationChannel(d *schema.ResourceData) (*notifications.AiNotificationsChannelInput, error) {
	channel := notifications.AiNotificationsChannelInput{
		Name:          d.Get("name").(string),
		DestinationId: d.Get("destination_id").(string),
		Type:          notifications.AiNotificationsChannelType(d.Get("type").(string)),
		Product:       notifications.AiNotificationsProduct(d.Get("product").(string)),
	}

	properties, propertiesOk := d.GetOk("properties")
	isEmailType := validateEmailChannelType(channel.Type)

	if !propertiesOk && !isEmailType {
		return nil, errors.New("notification channel requires a properties attribute")
	}

	if propertiesOk {
		var destinationProperty map[string]interface{}

		x := properties.([]interface{})

		for _, property := range x {
			destinationProperty = property.(map[string]interface{})
			if val, err := expandNotificationChannelProperty(destinationProperty); err == nil {
				channel.Properties = append(channel.Properties, *val)
			}
		}
	} else if isEmailType {
		channel.Properties = []notifications.AiNotificationsPropertyInput{{Key: "subject", Value: "{{ issueTitle }}"}} // Default subject
	}

	// Create terraform source property
	terraformProperty := notifications.AiNotificationsPropertyInput{
		Key:   "source",
		Value: "terraform",
		Label: "terraform-source-internal",
	}
	channel.Properties = append(channel.Properties, terraformProperty)

	return &channel, nil
}

func expandNotificationChannelUpdate(d *schema.ResourceData) (*notifications.AiNotificationsChannelUpdate, error) {
	channel := notifications.AiNotificationsChannelUpdate{
		Name:   d.Get("name").(string),
		Active: d.Get("active").(bool),
	}
	channelType := notifications.AiNotificationsChannelType(d.Get("type").(string))

	properties, propertiesOk := d.GetOk("properties")
	isEmailType := validateEmailChannelType(channelType)

	if !propertiesOk && !isEmailType {
		return nil, errors.New("notification channel requires a properties attribute")
	}

	if propertiesOk {
		var destinationProperty map[string]interface{}

		x := properties.([]interface{})

		for _, property := range x {
			destinationProperty = property.(map[string]interface{})
			if val, err := expandNotificationChannelProperty(destinationProperty); err == nil {
				channel.Properties = append(channel.Properties, *val)
			}
		}
	} else if isEmailType {
		channel.Properties = []notifications.AiNotificationsPropertyInput{{Key: "subject", Value: "{{ issueTitle }}"}} // Default subject
	}

	return &channel, nil
}

func expandNotificationChannelProperty(cfg map[string]interface{}) (*notifications.AiNotificationsPropertyInput, error) {
	property := notifications.AiNotificationsPropertyInput{}

	if propertyKey, ok := cfg["key"]; ok {
		property.Key = propertyKey.(string)
	}

	if propertyValue, ok := cfg["value"]; ok {
		property.Value = propertyValue.(string)
	}

	if propertyDisplayValue, ok := cfg["display_value"]; ok {
		property.DisplayValue = propertyDisplayValue.(string)
	}

	if propertyLabel, ok := cfg["label"]; ok {
		property.Label = propertyLabel.(string)
	}

	return &property, nil
}

func flattenNotificationChannel(channel *notifications.AiNotificationsChannel, d *schema.ResourceData) error {
	if channel == nil {
		return nil
	}

	var err error

	if err = d.Set("name", channel.Name); err != nil {
		return err
	}

	if err = d.Set("type", channel.Type); err != nil {
		return err
	}

	if err = d.Set("product", channel.Product); err != nil {
		return err
	}

	if err = d.Set("destination_id", channel.DestinationId); err != nil {
		return err
	}

	properties, propertiesErr := flattenNotificationChannelProperties(&channel.Properties)
	if propertiesErr != nil {
		return propertiesErr
	}

	if err := d.Set("properties", properties); err != nil {
		return err
	}

	return nil
}

func flattenNotificationChannelProperties(p *[]notifications.AiNotificationsProperty) ([]map[string]interface{}, error) {
	if p == nil {
		return nil, nil
	}

	var properties []map[string]interface{}

	for _, property := range *p {
		if val, err := flattenNotificationChannelProperty(&property); err == nil {
			properties = append(properties, val)
		}
	}

	return properties, nil
}

func flattenNotificationChannelProperty(p *notifications.AiNotificationsProperty) (map[string]interface{}, error) {
	if p == nil {
		return nil, nil
	}

	propertyResult := make(map[string]interface{})

	propertyResult["key"] = p.Key
	propertyResult["value"] = p.Value

	if p.DisplayValue != "" {
		propertyResult["display_value"] = p.DisplayValue
	}

	if p.Label != "" {
		propertyResult["label"] = p.Label
	}

	return propertyResult, nil
}

func validateEmailChannelType(channelType notifications.AiNotificationsChannelType) bool {
	return channelType == notifications.AiNotificationsChannelTypeTypes.EMAIL
}
