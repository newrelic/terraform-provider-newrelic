package newrelic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/accountmanagement"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

const NewRelicAccountManagementSchemaName string = "name"
const NewRelicAccountManagementSchemaRegion string = "region"
const NewRelicAccountManagementSchemaStatus string = "status"

func resourceNewRelicAccountManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAccountCreate,
		ReadContext:   resourceNewRelicAccountRead,
		UpdateContext: resourceNewRelicAccountUpdate,
		DeleteContext: resourceNewRelicAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			NewRelicAccountManagementSchemaName: {
				Type:        schema.TypeString,
				Description: "Name of the account to be created",
				Required:    true,
			},
			NewRelicAccountManagementSchemaRegion: {
				Type:         schema.TypeString,
				Description:  "A description of what this parsing rule represents.",
				ValidateFunc: validation.StringInSlice([]string{"us01", "eu01"}, false),
				Required:     true,
			},
			NewRelicAccountManagementSchemaStatus: {
				Type:        schema.TypeString,
				Description: "Status of the account - active or canceled",
				//ValidateFunc: validation.StringInSlice([]string{
				//	string(customeradministration.OrganizationAccountStatusTypes.ACTIVE),
				//	string(customeradministration.OrganizationAccountStatusTypes.CANCELED),
				//}, false),
				Computed: true,
			},
		},
	}
}

func resourceNewRelicAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	createAccountInput := accountmanagement.AccountManagementCreateInput{
		Name:       d.Get(NewRelicAccountManagementSchemaName).(string),
		RegionCode: d.Get(NewRelicAccountManagementSchemaRegion).(string),
	}
	created, err := client.AccountManagement.AccountManagementCreateAccount(createAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("account creation failed, please check input details")
	}
	accountID := created.ManagedAccount.ID

	d.SetId(strconv.Itoa(accountID))

	// After successfully creating an account, the resource sleeps for 10 seconds
	// to allow the backend to update and populate the newly created account. This delay
	// ensures the account is indexed by the `customeradministration` NerdGraph endpoint.
	time.Sleep(time.Second * 10)

	return resourceNewRelicAccountRead(ctx, d, meta)
}

func resourceNewRelicAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	var diags diag.Diagnostics
	organization, getOrgError := client.Organization.GetOrganization()

	if getOrgError != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to fetch organization information upon trying to read details of the created account: %v", getOrgError),
		})
		return diags
	}

	organizationID := organization.ID

	accountID, accountIDConversionError := strconv.Atoi(d.Id())
	if accountIDConversionError != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to convert account ID string to integer upon trying to read details of the created account: %v", accountIDConversionError),
		})
		return diags
	}

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		getAccountsInOrganizationResponse, accountStatus, getAccountsInOrganizationError := fetchAccountsInOrganization(
			meta,
			organizationID,
			accountID,
		)

		if getAccountsInOrganizationError != nil {
			return resource.NonRetryableError(getAccountsInOrganizationError)
		}

		accountsInOrganizationResponse := getAccountsInOrganizationResponse.Items

		if len(accountsInOrganizationResponse) != 1 {
			return resource.RetryableError(fmt.Errorf("failed to read account details, retrying"))
		}

		accountInOrganizationFetched := accountsInOrganizationResponse[0].ID
		if accountInOrganizationFetched != accountID {
			return resource.RetryableError(fmt.Errorf("failed to read details of account %d, obtained details of account %d instead - retrying", accountID, accountInOrganizationFetched))
		}

		account := accountsInOrganizationResponse[0]

		_ = d.Set(NewRelicAccountManagementSchemaName, account.Name)
		_ = d.Set(NewRelicAccountManagementSchemaRegion, account.RegionCode)
		_ = d.Set(NewRelicAccountManagementSchemaStatus, accountStatus)

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}
	return nil
}

func resourceNewRelicAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	updateAccountInput := accountmanagement.AccountManagementUpdateInput{
		Name: d.Get(NewRelicAccountManagementSchemaName).(string),
		ID:   accountID,
	}
	updated, err := client.AccountManagement.AccountManagementUpdateAccount(updateAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		return diag.Errorf("account update failed, please check input details")
	}

	// After successfully updating an account, the resource sleeps for 10 seconds
	// to allow the backend to update and populate the newly created account. This delay
	// ensures the account is indexed by the `customeradministration` NerdGraph endpoint.
	time.Sleep(time.Second * 10)
	return resourceNewRelicAccountRead(ctx, d, meta)
}

func resourceNewRelicAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID, accountIDConversionError := strconv.Atoi(d.Id())
	if accountIDConversionError != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to convert account ID string to integer upon trying to read details of the created account: %v", accountIDConversionError),
		})
		return diags
	}

	cancelAccountResponse, cancelAccountError := client.AccountManagement.AccountManagementCancelAccount(accountID)
	if cancelAccountError != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to cancel account %d: %v", accountID, cancelAccountError),
		})
		return diags
	}

	if cancelAccountResponse == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("account cancellation response was nil for account %d", accountID),
		})
		return diags
	}

	if cancelAccountResponse.ID == accountID && cancelAccountResponse.IsCanceled {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf(`Please note that the 'terraform destroy' operation performed on this resource has resulted in the 'cancellation' of account %d, meaning it is no longer active. 
For more details, please refer to https://docs.newrelic.com/docs/apis/nerdgraph/examples/manage-accounts-nerdgraph/#cancel-an-account.`, accountID),
		})
		return diags
	}

	return nil
}

func fetchAccountsInOrganization(
	meta interface{},
	organizationID string,
	accountID int,
) (*customeradministration.OrganizationAccountCollection, string, error) {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	matchCurrentAccountWithActiveAccountsResponse, matchActiveAccountsError := client.CustomerAdministration.GetAccountsMinimized(
		"",
		customeradministration.OrganizationAccountFilterInput{
			OrganizationId: customeradministration.OrganizationAccountOrganizationIdFilterInput{
				Eq: organizationID,
			},
			ID: customeradministration.OrganizationAccountIdFilterInput{
				Eq: accountID,
			},
			Status: customeradministration.OrganizationAccountStatusFilterInput{
				Eq: customeradministration.OrganizationAccountStatusTypes.ACTIVE,
			},
		},
		[]customeradministration.OrganizationAccountSortInput{},
	)

	if matchActiveAccountsError != nil {
		return nil, "", matchActiveAccountsError
	}

	if len(matchCurrentAccountWithActiveAccountsResponse.Items) == 0 {
		matchCurrentAccountWithCanceledAccountsResponse, matchCanceledAccountsError := client.CustomerAdministration.GetAccountsMinimized(
			"",
			customeradministration.OrganizationAccountFilterInput{
				OrganizationId: customeradministration.OrganizationAccountOrganizationIdFilterInput{
					Eq: organizationID,
				},
				ID: customeradministration.OrganizationAccountIdFilterInput{
					Eq: accountID,
				},
				Status: customeradministration.OrganizationAccountStatusFilterInput{
					Eq: customeradministration.OrganizationAccountStatusTypes.CANCELED,
				},
			},
			[]customeradministration.OrganizationAccountSortInput{},
		)

		if matchCanceledAccountsError != nil {
			return nil, "", matchCanceledAccountsError
		}

		if len(matchCurrentAccountWithCanceledAccountsResponse.Items) == 0 {
			return nil, "", nil
		}

		return matchCurrentAccountWithCanceledAccountsResponse, string(customeradministration.OrganizationAccountStatusTypes.CANCELED), nil
	}

	return matchCurrentAccountWithActiveAccountsResponse, string(customeradministration.OrganizationAccountStatusTypes.ACTIVE), nil
}
