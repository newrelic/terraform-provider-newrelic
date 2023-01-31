package newrelic

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/usermanagement"
)

func resourceNewRelicUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicUserGroupCreate,
		ReadContext:   resourceNewRelicUserGroupRead,
		UpdateContext: resourceNewRelicUserGroupUpdate,
		DeleteContext: resourceNewRelicUserGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			//"account_id": {
			//	Type:        schema.TypeInt,
			//	Description: "The New Relic account ID where you want to create the user groups",
			//	Computed:    true,
			//	Optional:    true,
			//},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the group",
				Required:    true,
			},
			"authentication_domain_id": {
				Type:        schema.TypeString,
				Description: "The id of the authentication domain the group will belong to.",
				Required:    true,
			},
		},
	}
}

// Create a [group](https://docs.newrelic.com/docs/accounts/accounts-billing/new-relic-one-user-management/user-management-concepts/#groups)
func resourceNewRelicUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	//accountID := selectAccountID(providerConfig, d)

	createInput := usermanagement.UserManagementCreateGroup{
		AuthenticationDomainId: d.Get("authentication_domain_id").(string),
		DisplayName:            d.Get("name").(string),
	}
	created, err := client.UserManagement.UserManagementCreateGroupWithContext(ctx, createInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: group create result wasn't returned or group was not created.")
	}

	groupID := created.Group.ID

	d.SetId(groupID)

	return resourceNewRelicUserGroupRead(ctx, d, meta)
}

func resourceNewRelicUserGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	groupId := d.Id()
	authenticationDomainId := d.Get("authentication_domain_id").(string)
	group, err := getUserGroupID(ctx, client, authenticationDomainId, groupId)

	if err != nil && group == nil {
		d.SetId("")
		return nil
	}

	//if err := d.Set("account_id", accountID); err != nil {
	//	return diag.FromErr(err)
	//}

	if err := d.Set("name", group.DisplayName); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getUserGroupID(ctx context.Context, client *newrelic.NewRelic, authId string, groupId string) (groups *usermanagement.UserManagementGroup, err error) {
	id := []string{authId}
	resp, err := client.UserManagement.GetAuthenticationDomainsWithContext(ctx, id)
	if err != nil {
		return nil, err
	}
	for i, groupList := range *resp {
		if i == 0 {
			for _, group := range groupList.Groups.Groups {
				if group.ID == groupId {
					return &group, nil
				}
			}
		}
	}
	return nil, errors.New("err: group is not found")
}

// Update the obfuscation expression
func resourceNewRelicUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput := usermanagement.UserManagementUpdateGroup{}

	if e, ok := d.GetOk("name"); ok {
		updateInput.DisplayName = e.(string)
	}
	updateInput.ID = d.Id()

	log.Printf("[INFO] Updating New Relic user group %s", d.Id())

	//accountID := selectAccountID(meta.(*ProviderConfig), d)

	_, err := client.UserManagement.UserManagementUpdateGroupWithContext(ctx, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicObfuscationExpressionRead(ctx, d, meta)
}

// Delete the UserGroup
func resourceNewRelicUserGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic user group id %s", d.Id())

	//accountID := selectAccountID(meta.(*ProviderConfig), d)
	deleteConfig := usermanagement.UserManagementDeleteGroup{
		ID: d.Id(),
	}

	_, err := client.UserManagement.UserManagementDeleteGroupWithContext(ctx, deleteConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
