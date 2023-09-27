package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/newrelic/newrelic-client-go/v2/pkg/accounts"
	"log"
	"strings"
)

//func dataSourceNewRelicAccount() *schema.Resource {
//	return &schema.Resource{
//		ReadContext: dataSourceNewRelicAccountRead,
//		Schema: map[string]*schema.Schema{
//			"scope": {
//				Type:         schema.TypeString,
//				Optional:     true,
//				Default:      string(accounts.RegionScopeTypes.IN_REGION),
//				Description:  `The scope of the account in New Relic.  Valid values are "global" and "in_region".  Defaults to "in_region".`,
//				ValidateFunc: validation.StringInSlice([]string{"global", "in_region"}, true),
//			},
//			"name": {
//				Type:        schema.TypeString,
//				Optional:    true,
//				Description: "The name of the account in New Relic.",
//			},
//			"account_id": {
//				Type:        schema.TypeInt,
//				Optional:    true,
//				Description: "The ID of the account in New Relic.",
//			},
//		},
//	}
//}
//
//func dataSourceNewRelicAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	providerConfig := meta.(*ProviderConfig)
//	client := providerConfig.NewClient
//
//	log.Printf("[INFO] Reading New Relic accounts")
//
//	scope := accounts.RegionScope(strings.ToUpper(d.Get("scope").(string)))
//
//	id, idOk := d.GetOk("account_id")
//	name, nameOk := d.GetOk("name")
//
//	params := accounts.ListAccountsParams{
//		Scope: &scope,
//	}
//
//	accts, err := client.Accounts.ListAccountsWithContext(ctx, params)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	var account *accounts.AccountOutline
//
//	if !idOk && !nameOk {
//		// Default to the provider's AccountID if no lookup attributes are provided.
//		id, idOk = selectAccountID(providerConfig, d), true
//	}
//
//	if idOk && nameOk {
//		return diag.FromErr(fmt.Errorf(`exactly one of "name" or "account_id" is required to locate a New Relic account`))
//	}
//
//	if nameOk {
//		for _, a := range accts {
//			if a.Name == name.(string) {
//				account = &a
//				break
//			}
//		}
//
//		if account == nil {
//			return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic accounts", name))
//		}
//	}
//
//	if idOk {
//		for _, a := range accts {
//			if a.ID == id.(int) {
//				account = &a
//				break
//			}
//		}
//
//		if account == nil {
//			return diag.FromErr(fmt.Errorf("the id '%d' does not match any New Relic accounts", id))
//		}
//	}
//
//	return diag.FromErr(flattenAccountData(account, d))
//}
//
//func flattenAccountData(a *accounts.AccountOutline, d *schema.ResourceData) error {
//	d.SetId(strconv.Itoa(a.ID))
//	var err error
//
//	err = d.Set("name", a.Name)
//	if err != nil {
//		return err
//	}
//
//	err = d.Set("account_id", a.ID)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

type dataSourceNewRelicAccount struct{}

type dataSourceNewRelicAccountModel struct {
	Scope     types.String `tfsdk:"scope"`
	Name      types.String `tfsdk:"name"`
	AccountID types.Int64  `tfsdk:"account_id"`
}

func (d *dataSourceNewRelicAccount) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "newrelic_account"
}

func (d *dataSourceNewRelicAccount) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"scope:": schema.StringAttribute{
				Optional:    true,
				Description: `The scope of the account in New Relic.  Valid values are "global" and "in_region".  Defaults to "in_region".`,
				Validators:  []validator.String{stringvalidator.OneOfCaseInsensitive("global", "in_region")},
				// How to set default value?
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the account in New Relic.",
			},
			"account_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the account in New Relic.",
			},
		},
	}
}

func (d *dataSourceNewRelicAccount) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	var data dataSourceNewRelicAccountModel

	log.Printf("[INFO] Reading New Relic accounts")

	scope := accounts.RegionScope(strings.ToUpper(data.Scope.ValueString()))

	params := accounts.ListAccountsParams{
		Scope: &scope,
	}

	accts, err := client.Accounts.ListAccountsWithContext(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("an error occurred", string(err.Error()))
	}

	var account *accounts.AccountOutline

	for _, a := range accts {
		if a.Name == data.Name.ValueString() {
			account = &a
			break
		}
	}

	if account == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("the name '%s' does not match any New Relic accounts", data.Name.ValueString()), "")
	}

	for _, a := range accts {
		if int64(a.ID) == data.AccountID.ValueInt64() {
			account = &a
			break
		}
	}

	if account == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("the id '%d' does not match any New Relic accounts", data.AccountID.ValueInt64()), "")
	}

	data.Name = types.StringValue(account.Name)
	data.AccountID = types.Int64Value(int64(account.ID))
}
