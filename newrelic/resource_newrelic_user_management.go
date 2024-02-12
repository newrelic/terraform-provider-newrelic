package newrelic

import (
	"context"
	"fmt"
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
				ForceNew:    true,
				// ForceNew has been added as the authentication_domain_id of a user cannot be updated post creation
				// This is because the `authenticationDomainId` field does not exist in the userManagementUpdateUser mutation
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

	createUserInput := usermanagement.UserManagementCreateUser{
		AuthenticationDomainId: d.Get("authentication_domain_id").(string),
		Email:                  d.Get("email").(string),
		Name:                   d.Get("name").(string),
		UserType:               userTypes[d.Get("user_type").(string)],
	}

	createUserResponse, err := client.UserManagement.UserManagementCreateUserWithContext(ctx, createUserInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if createUserResponse == nil {
		return diag.Errorf("err: user create result wasn't returned or user was not created.")
	}

	userID := createUserResponse.CreatedUser.ID
	d.SetId(userID)

	return nil
}

// Read a created user
func resourceNewRelicUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	userID := d.Id()
	user, authenticationDomainID, err := getUserByID(ctx, client, userID)

	if err != nil && user == nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("name", user.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("authentication_domain_id", authenticationDomainID); err != nil {
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

// iterate through users to spot the user with the specified user id
func getUserByID(ctx context.Context, client *newrelic.NewRelic, userID string) (user *usermanagement.UserManagementUser, authenticationDomainID string, err error) {
	userIDs := []string{userID}
	resp, err := client.UserManagement.UserManagementGetUsersWithContext(ctx, []string{}, userIDs, "", "")
	if err != nil {
		return nil, "", err
	}

	for _, authDomain := range resp.AuthenticationDomains {
		for _, u := range authDomain.Users.Users {
			if u.ID == userID {
				return &u, authDomain.ID, nil
			}
		}
	}

	return nil, "", fmt.Errorf("user with id %s not found", userID)
}

// Update a created user
func resourceNewRelicUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	updateUserInput := usermanagement.UserManagementUpdateUser{
		Email:    d.Get("email").(string),
		Name:     d.Get("name").(string),
		ID:       d.Id(),
		UserType: userTypes[d.Get("user_type").(string)],
	}

	updateUserResponse, err := client.UserManagement.UserManagementUpdateUserWithContext(ctx, updateUserInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if updateUserResponse == nil {
		return diag.Errorf("err: user update result wasn't returned or user was not created.")
	}

	userID := updateUserResponse.User.ID
	d.SetId(userID)

	return nil
}

// Delete a created user
func resourceNewRelicUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic user with user id %s\n", d.Id())

	deleteUserInput := usermanagement.UserManagementDeleteUser{
		ID: d.Id(),
	}

	_, err := client.UserManagement.UserManagementDeleteUserWithContext(ctx, deleteUserInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
