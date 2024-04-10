package newrelic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"
)

// migrateStateNewRelicNotificationDestinationV0toV1 currently facilitates migrating:
// remove is_user_authenticated argument
func migrateStateNewRelicNotificationDestinationV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	delete(rawState, "is_user_authenticated")

	return rawState, nil
}

func expandNotificationDestination(d *schema.ResourceData) (*notifications.AiNotificationsDestinationInput, error) {
	destination := notifications.AiNotificationsDestinationInput{
		Name: d.Get("name").(string),
		Type: notifications.AiNotificationsDestinationType(d.Get("type").(string)),
	}

	if attr, ok := d.GetOk("auth_basic"); ok {
		destination.Auth = expandNotificationDestinationAuthBasic(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("auth_token"); ok {
		destination.Auth = expandNotificationDestinationAuthToken(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("auth_custom_header"); ok {
		//customHeadersList := attr.(*schema.Set).List()
		destination.Auth = expandNotificationDestinationAuthCustomHeaders(attr.([]interface{}))
	}

	properties := d.Get("property")
	props := properties.(*schema.Set).List()
	destination.Properties = expandNotificationDestinationProperties(props)

	if attr, ok := d.GetOk("secure_url"); ok {
		destination.SecureURL = expandNotificationDestinationSecureUrlInput(attr.([]interface{}))
	}

	return &destination, nil
}

func expandNotificationDestinationSecureUrlInput(url []interface{}) *notifications.AiNotificationsSecureURLInput {
	secureUrlInput := notifications.AiNotificationsSecureURLInput{}

	for _, u := range url {
		uu := u.(map[string]interface{})
		secureUrlInput.Prefix = uu["prefix"].(string)
		secureUrlInput.SecureSuffix = notifications.SecureValue(uu["secure_suffix"].(string))
	}

	return &secureUrlInput
}

func expandNotificationDestinationSecureUrlUpdate(url []interface{}) *notifications.AiNotificationsSecureURLUpdate {
	secureUrlInput := notifications.AiNotificationsSecureURLUpdate{}

	for _, u := range url {
		uu := u.(map[string]interface{})
		secureUrlInput.Prefix = uu["prefix"].(string)
		secureUrlInput.SecureSuffix = notifications.SecureValue(uu["secure_suffix"].(string))
	}

	return &secureUrlInput
}

func expandNotificationDestinationAuthBasic(authRaw []interface{}) *notifications.AiNotificationsCredentialsInput {
	authInput := notifications.AiNotificationsCredentialsInput{}
	authInput.Type = notifications.AiNotificationsAuthTypeTypes.BASIC

	for _, a := range authRaw {
		aa := a.(map[string]interface{})
		authInput.Basic.User = aa["user"].(string)
		authInput.Basic.Password = notifications.SecureValue(aa["password"].(string))
	}

	return &authInput
}

func expandNotificationDestinationAuthToken(authRaw []interface{}) *notifications.AiNotificationsCredentialsInput {
	authInput := notifications.AiNotificationsCredentialsInput{}
	authInput.Type = notifications.AiNotificationsAuthTypeTypes.TOKEN

	for _, a := range authRaw {
		aa := a.(map[string]interface{})
		authInput.Token.Token = notifications.SecureValue(aa["token"].(string))
		authInput.Token.Prefix = aa["prefix"].(string)
	}

	return &authInput
}

func expandNotificationDestinationAuthCustomHeaders(authRaw []interface{}) *notifications.AiNotificationsCredentialsInput {
	authInput := notifications.AiNotificationsCredentialsInput{}
	authInput.Type = notifications.AiNotificationsAuthTypeTypes.CUSTOM_HEADERS

	customHeadersList := []notifications.AiNotificationsCustomHeaderInput{}

	for _, h := range authRaw {
		customHeadersList = append(customHeadersList, expandNotificationDestinationAuthCustomHeader(h.(map[string]interface{})))
	}

	customHeaders := notifications.AiNotificationsCustomHeadersAuthInput{CustomHeaders: customHeadersList}
	authInput.CustomHeaders = &customHeaders

	return &authInput
}

func expandNotificationDestinationAuthCustomHeader(authRaw map[string]interface{}) notifications.AiNotificationsCustomHeaderInput {
	customHeader := notifications.AiNotificationsCustomHeaderInput{}

	if key, ok := authRaw["key"]; ok {
		customHeader.Key = key.(string)
	}

	if value, ok := authRaw["value"]; ok {
		customHeader.Value = notifications.SecureValue(value.(string))
	}

	return customHeader
}

func expandNotificationDestinationUpdate(d *schema.ResourceData) (*notifications.AiNotificationsDestinationUpdate, error) {
	destination := notifications.AiNotificationsDestinationUpdate{
		Name:   d.Get("name").(string),
		Active: d.Get("active").(bool),
	}

	if attr, ok := d.GetOk("auth_basic"); ok {
		destination.Auth = expandNotificationDestinationAuthBasic(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("auth_token"); ok {
		destination.Auth = expandNotificationDestinationAuthToken(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("auth_custom_header"); ok {
		destination.Auth = expandNotificationDestinationAuthCustomHeaders(attr.([]interface{}))
	}

	secureUrl := d.Get("secure_url")
	destination.SecureURL = expandNotificationDestinationSecureUrlUpdate(secureUrl.([]interface{}))

	properties := d.Get("property")
	props := properties.(*schema.Set).List()
	destination.Properties = expandNotificationDestinationProperties(props)

	return &destination, nil
}

func expandNotificationDestinationProperties(properties []interface{}) []notifications.AiNotificationsPropertyInput {
	props := []notifications.AiNotificationsPropertyInput{}

	for _, p := range properties {
		props = append(props, expandNotificationDestinationProperty(p.(map[string]interface{})))
	}

	return props
}

func expandNotificationDestinationProperty(cfg map[string]interface{}) notifications.AiNotificationsPropertyInput {
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

	return property
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

	if err := d.Set("guid", destination.GUID); err != nil {
		return err
	}

	auth := flattenNotificationDestinationAuth(destination.Auth, d)

	var authAttr string
	switch destination.Auth.AuthType {
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.BASIC):
		authAttr = "auth_basic"
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.OAUTH2):
		authAttr = "auth_oauth2"
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.TOKEN):
		authAttr = "auth_token"
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.CUSTOM_HEADERS):
		authAttr = "auth_custom_header"
	}

	if authAttr != "" {
		if err := d.Set(authAttr, auth); err != nil {
			return fmt.Errorf("[DEBUG] Error setting notification auth: %#v", err)
		}
	}

	if err := d.Set("property", flattenNotificationDestinationProperties(destination.Properties)); err != nil {
		return err
	}

	if err := d.Set("active", destination.Active); err != nil {
		return err
	}

	if err := d.Set("account_id", destination.AccountID); err != nil {
		return err
	}

	if err := d.Set("status", destination.Status); err != nil {
		return err
	}

	if err := d.Set("last_sent", destination.LastSent); err != nil {
		return err
	}

	if err := d.Set("secure_url", flattenNotificationDestinationSecureUrl(&destination.SecureURL, d)); err != nil {
		return err
	}

	return nil
}

func flattenNotificationDestinationAuth(a ai.AiNotificationsAuth, d *schema.ResourceData) []map[string]interface{} {
	authConfig := make([]map[string]interface{}, 1)

	switch a.AuthType {
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.BASIC):
		authConfig[0] = map[string]interface{}{
			"user":     a.User,
			"password": d.Get("auth_basic.0.password"),
		}
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.TOKEN):
		authConfig[0] = map[string]interface{}{
			"prefix": a.Prefix,
			"token":  d.Get("auth_token.0.token"),
		}
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.CUSTOM_HEADERS):
		customHeaders := flattenNotificationDestinationCustomHeaders(a.CustomHeaders, d)
		authConfig = customHeaders
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.OAUTH2):
		// This auth type is not supported
	}

	return authConfig
}

func flattenNotificationDestinationCustomHeaders(a []ai.AiNotificationsCustomHeaders, d *schema.ResourceData) []map[string]interface{} {
	customHeaders := []map[string]interface{}{}

	for i, customHeader := range a {
		customHeaders = append(customHeaders, flattenNotificationDestinationCustomHeader(customHeader, d, i))
	}

	return customHeaders
}

func flattenNotificationDestinationCustomHeader(header ai.AiNotificationsCustomHeaders, d *schema.ResourceData, i int) map[string]interface{} {
	customHeaderResult := make(map[string]interface{})

	customHeaderResult["key"] = header.Key
	valuePath := fmt.Sprintf("auth_custom_header.%d.value", i)
	customHeaderResult["value"] = d.Get(valuePath)

	return customHeaderResult
}

func flattenNotificationDestinationProperties(p []notifications.AiNotificationsProperty) []map[string]interface{} {
	properties := []map[string]interface{}{}

	for _, property := range p {
		properties = append(properties, flattenNotificationDestinationProperty(property))
	}

	return properties
}

func flattenNotificationDestinationProperty(p notifications.AiNotificationsProperty) map[string]interface{} {
	propertyResult := make(map[string]interface{})

	propertyResult["key"] = p.Key
	propertyResult["value"] = p.Value
	propertyResult["display_value"] = p.DisplayValue
	propertyResult["label"] = p.Label

	return propertyResult
}

func flattenNotificationDestinationSecureUrl(url *notifications.AiNotificationsSecureURL, d *schema.ResourceData) []map[string]interface{} {
	if url == nil || url.Prefix == "" {
		return nil
	}

	secureUrlResult := make([]map[string]interface{}, 1)

	secureUrlResult[0] = map[string]interface{}{
		"prefix":        url.Prefix,
		"secure_suffix": d.Get("secure_url.0.secure_suffix"),
	}

	return secureUrlResult
}

func flattenNotificationDestinationDataSource(destination *notifications.AiNotificationsDestination, d *schema.ResourceData) error {
	if destination == nil {
		return nil
	}

	var err error

	d.SetId(destination.ID)

	if err = d.Set("name", destination.Name); err != nil {
		return err
	}

	if err = d.Set("type", destination.Type); err != nil {
		return err
	}

	if err := d.Set("property", flattenNotificationDestinationProperties(destination.Properties)); err != nil {
		return err
	}

	if err := d.Set("active", destination.Active); err != nil {
		return err
	}

	if err := d.Set("account_id", destination.AccountID); err != nil {
		return err
	}

	if err := d.Set("status", destination.Status); err != nil {
		return err
	}

	if err := d.Set("guid", destination.GUID); err != nil {
		return err
	}

	return nil
}
