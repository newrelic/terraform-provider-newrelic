package newrelic

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
)

func expandNotificationDestination(d *schema.ResourceData) (*notifications.AiNotificationsDestinationInput, error) {
	destination := notifications.AiNotificationsDestinationInput{
		Name: d.Get("name").(string),
		Type: notifications.AiNotificationsDestinationType(d.Get("type").(string)),
	}

	var auth, authOk = d.GetOk("auth")

	if !authOk {
		return nil, errors.New("notification destination requires an auth attribute")
	}

	if authOk {
		a, err := expandNotificationDestinationAuth(auth)
		if err != nil {
			return nil, err
		}

		destination.Auth = *a
	}

	err := validateDestinationAuth(destination.Auth)
	if err != nil {
		return nil, err
	}

	properties, propertiesOk := d.GetOk("properties")
	isPagerDutyType := validatePagerDutyDestinationType(destination.Type)

	if !propertiesOk && !isPagerDutyType {
		return nil, errors.New("notification destination requires a properties attribute")
	}

	if propertiesOk {
		var destinationProperty map[string]interface{}

		x := properties.([]interface{})

		for _, property := range x {
			destinationProperty = property.(map[string]interface{})
			if val, err := expandNotificationDestinationProperty(destinationProperty); err == nil {
				destination.Properties = append(destination.Properties, *val)
			}
		}
	} else if isPagerDutyType {
		destination.Properties = []notifications.AiNotificationsPropertyInput{
			{
				Key:   "two_way_integration",
				Value: "false",
			},
		}
	}

	return &destination, nil
}

func expandNotificationDestinationAuth(authList interface{}) (*notifications.AiNotificationsCredentialsInput, error) {
	auth := notifications.AiNotificationsCredentialsInput{}
	authConfig := authList.(map[string]interface{})

	if typeAuth, ok := authConfig["type"]; ok {
		auth.Type = notifications.AiNotificationsAuthType(typeAuth.(string))
	}

	if auth.Type == notifications.AiNotificationsAuthTypeTypes.TOKEN {
		if prefix, ok := authConfig["prefix"]; ok {
			auth.Token.Prefix = prefix.(string)
		}

		if token, ok := authConfig["token"]; ok {
			auth.Token.Token = notifications.SecureValue(token.(string))
		}
	}

	if auth.Type == notifications.AiNotificationsAuthTypeTypes.BASIC {
		if user, ok := authConfig["user"]; ok {
			auth.Basic.User = user.(string)
		}

		if password, ok := authConfig["password"]; ok {
			auth.Basic.Password = notifications.SecureValue(password.(string))
		}
	}

	return &auth, nil
}

func expandNotificationDestinationProperty(cfg map[string]interface{}) (*notifications.AiNotificationsPropertyInput, error) {
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

func flattenNotificationDestination(destination *notifications.AiNotificationsDestination, d *schema.ResourceData) error {
	if destination == nil {
		return nil
	}

	var err error

	if err = d.Set("name", destination.Name); err != nil {
		return err
	}

	if err = d.Set("type", destination.Type); err != nil {
		return err
	}

	if err := d.Set("auth", flattenNotificationDestinationAuth(&destination.Auth)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting notification auth: %#v", err)
	}

	properties, propertiesErr := flattenNotificationDestinationProperties(&destination.Properties)
	if propertiesErr != nil {
		return propertiesErr
	}

	if err := d.Set("properties", properties); err != nil {
		return err
	}

	return nil
}

func flattenNotificationDestinationProperties(p *[]notifications.AiNotificationsProperty) ([]map[string]interface{}, error) {
	if p == nil {
		return nil, nil
	}

	var properties []map[string]interface{}

	for _, property := range *p {
		if val, err := flattenNotificationDestinationProperty(&property); err == nil {
			properties = append(properties, val)
		}
	}

	return properties, nil
}

func flattenNotificationDestinationProperty(p *notifications.AiNotificationsProperty) (map[string]interface{}, error) {
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

func flattenNotificationDestinationAuth(a *ai.AiNotificationsAuth) interface{} {

	authConfig := map[string]interface{}{
		"authType": a.AuthType,
	}

	if authConfig["authType"] == notifications.AiNotificationsAuthTypeTypes.BASIC {
		authConfig["user"] = a.User
	}

	if authConfig["authType"] == notifications.AiNotificationsAuthTypeTypes.TOKEN {
		authConfig["prefix"] = a.Prefix
	}

	if authConfig["authType"] == notifications.AiNotificationsAuthTypeTypes.OAUTH2 {
		authConfig["access_token_url"] = a.AccessTokenURL
	}

	return authConfig
}

func validateDestinationAuth(auth notifications.AiNotificationsCredentialsInput) error {
	if auth.Type == "" {
		return errors.New("auth type is required")
	}

	if auth.Type != notifications.AiNotificationsAuthTypeTypes.TOKEN && auth.Type != notifications.AiNotificationsAuthTypeTypes.BASIC {
		return errors.New("auth type must be token or basic")
	}

	if auth.Type == notifications.AiNotificationsAuthTypeTypes.TOKEN && (auth.Token.Token == "" || auth.Token.Prefix == "") {
		return errors.New("token and prefix are required when using token auth type")
	}

	if auth.Type == notifications.AiNotificationsAuthTypeTypes.BASIC && (auth.Basic.User == "" || auth.Basic.Password == "") {
		return errors.New("user and password are required when using basic auth type")
	}

	return nil
}

func validatePagerDutyDestinationType(destinationType notifications.AiNotificationsDestinationType) bool {
	if destinationType == notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_ACCOUNT_INTEGRATION || destinationType == notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_SERVICE_INTEGRATION {
		return true
	}

	return false
}
