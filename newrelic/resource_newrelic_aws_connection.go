package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/federatedlogs"
)

func resourceNewRelicAwsConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAwsConnectionCreate,
		ReadContext:   resourceNewRelicAwsConnectionRead,
		UpdateContext: resourceNewRelicAwsConnectionUpdate,
		DeleteContext: resourceNewRelicAwsConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The account ID where the AWS connection will be created. Used when scope_type is ACCOUNT.",
			},
			"scope_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The scope ID (account ID or organization ID) for the AWS connection.",
			},
			"scope_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "ACCOUNT",
				Description: "The scope type for the AWS connection. Valid values are ACCOUNT and ORGANIZATION.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the AWS connection.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the AWS connection.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag to indicate if the connection is enabled. True by default.",
			},
			"external_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional field representing an identifier managed by the consumer.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default region for this connection.",
			},
			"credential": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Credentials for accessing the AWS account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assume_role": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "AssumeRole configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_arn": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ARN of the IAM role New Relic should assume.",
									},
									"external_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "External ID supplied by New Relic during AssumeRole.",
									},
								},
							},
						},
					},
				},
			},
			"settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional list of connection settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key or name of the setting.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the setting.",
						},
					},
				},
			},
			"tag": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Tags applied to the AWS Connection entity.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag key.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The tag values.",
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicAwsConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	// Determine scope
	scopeType := d.Get("scope_type").(string)
	var scopeID string

	if v, ok := d.GetOk("scope_id"); ok {
		scopeID = v.(string)
	} else {
		accountID := selectAccountID(providerConfig, d)
		scopeID = strconv.Itoa(accountID)
		_ = d.Set("account_id", accountID)
	}

	input := federatedlogs.EntityManagementAwsConnectionEntityCreateInput{
		Name:       d.Get("name").(string),
		Enabled:    getBoolPointer(d.Get("enabled").(bool)),
		Credential: expandAwsConnectionCredential(d.Get("credential").([]interface{})),
		Scope: federatedlogs.EntityManagementScopedReferenceInput{
			Type: federatedlogs.EntityManagementEntityScope(scopeType),
			ID:   scopeID,
		},
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("external_id"); ok {
		input.ExternalId = v.(string)
	}
	if v, ok := d.GetOk("region"); ok {
		input.Region = v.(string)
	}
	if v, ok := d.GetOk("settings"); ok {
		input.Settings = expandAwsConnectionSettings(v.([]interface{}))
	}
	if v, ok := d.GetOk("tag"); ok {
		input.Tags = expandAwsConnectionTags(v.(*schema.Set).List())
	}

	resp, err := client.Federatedlogs.EntityManagementCreateAwsConnectionWithContext(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Entity.ID)

	return resourceNewRelicAwsConnectionRead(ctx, d, meta)
}

func resourceNewRelicAwsConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic AWS Connection for %s", d.Id())

	entityID := d.Id()

	resp, err := client.Federatedlogs.GetEntityWithContext(ctx, entityID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil || *resp == nil {
		log.Printf("[WARN] AWS Connection %s not found", entityID)
		d.SetId("")
		return nil
	}

	switch entityType := (*resp).(type) {
	case *federatedlogs.EntityManagementAwsConnectionEntity:
		entity := (*resp).(*federatedlogs.EntityManagementAwsConnectionEntity)

		accountIDStr := entity.Scope.ID
		accountIDInt, accountIDIntErr := strconv.Atoi(accountIDStr)
		if accountIDIntErr != nil {
			log.Printf("[ERROR] Failed to convert account ID to integer: %v", accountIDIntErr)
			accountIDInt = selectAccountID(providerConfig, d)
			log.Printf("[INFO] Assigning the value of account_id from the state to prevent a panic: %d", accountIDInt)
		}

		if err := d.Set("account_id", accountIDInt); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("scope_id", entity.Scope.ID); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("scope_type", string(entity.Scope.Type)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("name", entity.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("description", entity.Description); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("enabled", entity.Enabled); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("external_id", entity.ExternalId); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("region", entity.Region); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("credential", flattenAwsConnectionCredential(entity.Credential)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("settings", flattenAwsConnectionSettings(entity.Settings)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("tag", flattenAwsConnectionTags(entity.Tags)); err != nil {
			return diag.FromErr(err)
		}
	default:
		return diag.Errorf("unexpected entity type %T for ID %s", entityType, d.Id())
	}
	return nil
}

func resourceNewRelicAwsConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	entityID := d.Id()
	log.Printf("[INFO] Updating New Relic AWS Connection entity with ID %s", entityID)

	getResp, err := client.Federatedlogs.GetEntityWithContext(ctx, entityID)
	if err != nil {
		return diag.FromErr(err)
	}
	if getResp == nil {
		d.SetId("")
		return nil
	}
	awsEntity, ok := (*getResp).(*federatedlogs.EntityManagementAwsConnectionEntity)
	if !ok {
		return diag.Errorf("unexpected entity type %T for ID %s", *getResp, entityID)
	}

	// Build the update input from changed fields only. Each field has
	// json:",omitempty" on the wire, so unset fields are not sent.
	input := federatedlogs.EntityManagementAwsConnectionEntityUpdateInput{}
	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		input.Description = d.Get("description").(string)
	}
	if d.HasChange("enabled") {
		input.Enabled = getBoolPointer(d.Get("enabled").(bool))
	}
	if d.HasChange("external_id") {
		input.ExternalId = d.Get("external_id").(string)
	}
	if d.HasChange("region") {
		input.Region = d.Get("region").(string)
	}
	if d.HasChange("credential") {
		input.Credential = expandAwsConnectionCredentialUpdate(d.Get("credential").([]interface{}))
	}
	if d.HasChange("settings") {
		input.Settings = expandAwsConnectionSettingsUpdate(d.Get("settings").([]interface{}))
	}
	if d.HasChange("tag") {
		input.Tags = expandAwsConnectionTags(d.Get("tag").(*schema.Set).List())
	}

	if _, err := client.Federatedlogs.EntityManagementUpdateAwsConnectionWithContext(ctx, input, entityID, awsEntity.Metadata.Version); err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicAwsConnectionRead(ctx, d, meta)
}

func resourceNewRelicAwsConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic AWS Connection entity with ID %s", d.Id())

	entityID := d.Id()

	// EntityManagementDelete requires the entity's current metadata.version
	// for optimistic concurrency control. Fetch it first.
	getResp, err := client.Federatedlogs.GetEntityWithContext(ctx, entityID)
	if err != nil {
		return diag.FromErr(err)
	}
	if getResp == nil {
		d.SetId("")
		return nil
	}
	awsEntity, ok := (*getResp).(*federatedlogs.EntityManagementAwsConnectionEntity)
	if !ok {
		return diag.Errorf("unexpected entity type %T for ID %s", *getResp, entityID)
	}

	result, err := client.Federatedlogs.EntityManagementDeleteWithContext(ctx, entityID, awsEntity.Metadata.Version)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.ID == "" {
		return diag.FromErr(fmt.Errorf("error in deleting entity with ID %s", entityID))
	}
	return nil
}

func expandAwsConnectionSettings(settings []interface{}) []federatedlogs.EntityManagementConnectionSettingsCreateInput {
	result := make([]federatedlogs.EntityManagementConnectionSettingsCreateInput, len(settings))
	for i, s := range settings {
		setting := s.(map[string]interface{})
		result[i] = federatedlogs.EntityManagementConnectionSettingsCreateInput{
			Key:   setting["key"].(string),
			Value: setting["value"].(string),
		}
	}
	return result
}

// expandAwsConnectionSettingsUpdate is the update-path counterpart to
// expandAwsConnectionSettings. The Create / Update input types have identical
// fields (Key, Value) but distinct Go types; client-go keeps them separate
// because the GraphQL schema declares them separately.
func expandAwsConnectionSettingsUpdate(settings []interface{}) []federatedlogs.EntityManagementConnectionSettingsUpdateInput {
	result := make([]federatedlogs.EntityManagementConnectionSettingsUpdateInput, len(settings))
	for i, s := range settings {
		setting := s.(map[string]interface{})
		result[i] = federatedlogs.EntityManagementConnectionSettingsUpdateInput{
			Key:   setting["key"].(string),
			Value: setting["value"].(string),
		}
	}
	return result
}

func flattenAwsConnectionSettings(settings []federatedlogs.EntityManagementConnectionSettings) []map[string]interface{} {
	result := make([]map[string]interface{}, len(settings))
	for i, s := range settings {
		result[i] = map[string]interface{}{
			"key":   s.Key,
			"value": s.Value,
		}
	}
	return result
}

func expandAwsConnectionCredential(in []interface{}) federatedlogs.EntityManagementAwsCredentialsCreateInput {
	if len(in) == 0 {
		return federatedlogs.EntityManagementAwsCredentialsCreateInput{}
	}
	cred := in[0].(map[string]interface{})
	assumeRoleList := cred["assume_role"].([]interface{})
	if len(assumeRoleList) == 0 {
		return federatedlogs.EntityManagementAwsCredentialsCreateInput{}
	}
	ar := assumeRoleList[0].(map[string]interface{})
	return federatedlogs.EntityManagementAwsCredentialsCreateInput{
		AssumeRole: federatedlogs.EntityManagementAwsAssumeRoleConfigCreateInput{
			RoleArn:    ar["role_arn"].(string),
			ExternalId: federatedlogs.EntityManagementDynamicString(ar["external_id"].(string)),
		},
	}
}

func expandAwsConnectionCredentialUpdate(in []interface{}) *federatedlogs.EntityManagementAwsCredentialsUpdateInput {
	if len(in) == 0 {
		return nil
	}
	cred := in[0].(map[string]interface{})
	assumeRoleList := cred["assume_role"].([]interface{})
	if len(assumeRoleList) == 0 {
		return nil
	}
	ar := assumeRoleList[0].(map[string]interface{})
	return &federatedlogs.EntityManagementAwsCredentialsUpdateInput{
		AssumeRole: federatedlogs.EntityManagementAwsAssumeRoleConfigUpdateInput{
			RoleArn:    ar["role_arn"].(string),
			ExternalId: federatedlogs.EntityManagementDynamicString(ar["external_id"].(string)),
		},
	}
}

func flattenAwsConnectionCredential(c federatedlogs.EntityManagementAwsCredentials) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"assume_role": []map[string]interface{}{
				{
					"role_arn":    c.AssumeRole.RoleArn,
					"external_id": string(c.AssumeRole.ExternalId),
				},
			},
		},
	}
}

func expandAwsConnectionTags(in []interface{}) []federatedlogs.EntityManagementTagInput {
	result := make([]federatedlogs.EntityManagementTagInput, 0, len(in))
	for _, raw := range in {
		m := raw.(map[string]interface{})
		valuesRaw := m["values"].(*schema.Set).List()
		values := make([]string, 0, len(valuesRaw))
		for _, v := range valuesRaw {
			values = append(values, v.(string))
		}
		result = append(result, federatedlogs.EntityManagementTagInput{
			Key:    m["key"].(string),
			Values: values,
		})
	}
	return result
}

func flattenAwsConnectionTags(tags []federatedlogs.EntityManagementTag) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(tags))
	for _, t := range tags {
		result = append(result, map[string]interface{}{
			"key":    t.Key,
			"values": t.Values,
		})
	}
	return result
}
