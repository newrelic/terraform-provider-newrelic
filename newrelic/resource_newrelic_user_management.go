package newrelic

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/usermanagement"
)

func resourceNewRelicUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicUserCreate,
		ReadContext:   resourceNewRelicUserRead,
		UpdateContext: resourceNewRelicUserUpdate,
		DeleteContext: resourceNewRelicUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the user.",
				Required:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The email ID of the user.",
				Required:    true,
			},
			"authentication_domain_id": {
				Type:        schema.TypeString,
				Description: "The ID of the authentication domain the user will belong to.",
				Required:    true,
			},
			"user_type": {
				Type:        schema.TypeString,
				Description: "The type of the user to be created.",
				Optional:    true,
				Default:     string(usermanagement.UserManagementRequestedTierNameTypes.BASIC_USER_TIER),
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(usermanagement.UserManagementRequestedTierNameTypes.BASIC_USER_TIER),
						string(usermanagement.UserManagementRequestedTierNameTypes.CORE_USER_TIER),
						string(usermanagement.UserManagementRequestedTierNameTypes.FULL_USER_TIER),
					},
					true),
			},
		},
	}
}

var userTypes = map[string]usermanagement.UserManagementRequestedTierName{
	"BASIC_USER_TIER": usermanagement.UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	"CORE_USER_TIER":  usermanagement.UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	"FULL_USER_TIER":  usermanagement.UserManagementRequestedTierNameTypes.FULL_USER_TIER,
}

var userTypesReadCompatible = map[string]usermanagement.UserManagementRequestedTierName{
	"Basic":         usermanagement.UserManagementRequestedTierNameTypes.BASIC_USER_TIER,
	"Core":          usermanagement.UserManagementRequestedTierNameTypes.CORE_USER_TIER,
	"Full platform": usermanagement.UserManagementRequestedTierNameTypes.FULL_USER_TIER,
}

// Create a user
func resourceNewRelicUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	createInput := usermanagement.UserManagementCreateUser{
		AuthenticationDomainId: d.Get("authentication_domain_id").(string),
		Email:                  d.Get("email").(string),
		Name:                   d.Get("name").(string),
		UserType:               userTypes[d.Get("user_type").(string)],
	}

	created, err := client.UserManagement.UserManagementCreateUserWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: user create result wasn't returned or user was not created.")
	}

	userID := created.CreatedUser.ID

	d.SetId(userID)

	return resourceNewRelicUserRead(ctx, d, meta)
}

// Read a created user
func resourceNewRelicUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	userID := d.Id()
	authenticationDomainID := d.Get("authentication_domain_id").(string)
	user, err := getUserID(ctx, client, authenticationDomainID, userID)

	if err != nil && user == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", user.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("user_type", userTypesReadCompatible[user.Type.DisplayName]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Iterate through users associated with the authentication domain
func getUserID(ctx context.Context, client *newrelic.NewRelic, authenticationDomainID string, userID string) (user *usermanagement.UserManagementUser, err error) {
	id := []string{authenticationDomainID}
	resp, err := client.UserManagement.GetAuthenticationDomainsWithContext(ctx, id)
	if err != nil {
		return nil, err
	}
	for i, userList := range *resp {
		if i == 0 {
			for _, user := range userList.Users.Users {
				if user.ID == userID {
					return &user, nil
				}
			}
		}
	}
	return nil, errors.New("error: user is not found")
}

// Update a created user
func resourceNewRelicUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	createInput := usermanagement.UserManagementUpdateUser{
		Email:    d.Get("email").(string),
		Name:     d.Get("name").(string),
		ID:       d.Id(),
		UserType: userTypes[d.Get("user_type").(string)],
	}

	created, err := client.UserManagement.UserManagementUpdateUserWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: user create result wasn't returned or user was not created.")
	}

	userID := created.User.ID

	d.SetId(userID)

	return resourceNewRelicUserRead(ctx, d, meta)
}

// Delete a created user
func resourceNewRelicUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic user user id %s", d.Id())

	//accountID := selectAccountID(meta.(*ProviderConfig), d)
	deleteConfig := usermanagement.UserManagementDeleteUser{
		ID: d.Id(),
	}

	_, err := client.UserManagement.UserManagementDeleteUserWithContext(ctx, deleteConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
