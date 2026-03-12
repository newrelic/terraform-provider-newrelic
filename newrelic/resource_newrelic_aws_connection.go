package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pipelinecontrol"
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
			"role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ARN of the IAM role to assume for this connection.",
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

	input := pipelinecontrol.EntityManagementAwsConnectionEntityCreateInput{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
		Credential: pipelinecontrol.EntityManagementAwsCredentialsCreateInput{
			AssumeRole: pipelinecontrol.EntityManagementAwsAssumeRoleConfigCreateInput{
				RoleArn: d.Get("role_arn").(string),
			},
		},
		Scope: pipelinecontrol.EntityManagementScopedReferenceInput{
			Type: pipelinecontrol.EntityManagementEntityScope(scopeType),
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

	resp, err := client.Pipelinecontrol.EntityManagementCreateAwsConnectionWithContext(ctx, input)
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

	resp, err := client.Pipelinecontrol.GetEntityWithContext(ctx, entityID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("entity with ID %s was nil", entityID))
	}

	switch entityType := (*resp).(type) {
	case *pipelinecontrol.EntityManagementAwsConnectionEntity:
		entity := (*resp).(*pipelinecontrol.EntityManagementAwsConnectionEntity)

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
		if err := d.Set("role_arn", entity.Credential.AssumeRole.RoleArn); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("settings", flattenAwsConnectionSettings(entity.Settings)); err != nil {
			return diag.FromErr(err)
		}
	default:
		return diag.Errorf("unexpected entity type %T for ID %s", entityType, d.Id())
	}
	return nil
}

func resourceNewRelicAwsConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Update not yet supported for AwsConnectionEntity in the client
	return resourceNewRelicAwsConnectionRead(ctx, d, meta)
}

func resourceNewRelicAwsConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic AWS Connection entity with ID %s", d.Id())

	entityID := d.Id()

	result, err := client.Pipelinecontrol.EntityManagementDelete(entityID)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.ID == "" {
		return diag.FromErr(fmt.Errorf("error in deleting entity with ID %s", entityID))
	}
	return nil
}

func expandAwsConnectionSettings(settings []interface{}) []pipelinecontrol.EntityManagementConnectionSettingsCreateInput {
	result := make([]pipelinecontrol.EntityManagementConnectionSettingsCreateInput, len(settings))
	for i, s := range settings {
		setting := s.(map[string]interface{})
		result[i] = pipelinecontrol.EntityManagementConnectionSettingsCreateInput{
			Key:   setting["key"].(string),
			Value: setting["value"].(string),
		}
	}
	return result
}

func flattenAwsConnectionSettings(settings []pipelinecontrol.EntityManagementConnectionSettings) []map[string]interface{} {
	result := make([]map[string]interface{}, len(settings))
	for i, s := range settings {
		result[i] = map[string]interface{}{
			"key":   s.Key,
			"value": s.Value,
		}
	}
	return result
}
