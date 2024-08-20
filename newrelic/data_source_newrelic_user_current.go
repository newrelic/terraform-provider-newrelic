package newrelic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNewRelicCurrentUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicCurrentUserRead,
		Schema: map[string]*schema.Schema{
			UserDataSourceUserNameAttrLabel:  dataSourceNewRelicCurrentUserSchemaConstructor(UserDataSourceUserNameAttrLabel),
			UserDataSourceUserEmailAttrLabel: dataSourceNewRelicCurrentUserSchemaConstructor(UserDataSourceUserEmailAttrLabel),
		},
	}
}

func dataSourceNewRelicCurrentUserSchemaConstructor(attributeLabel string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: fmt.Sprintf("The %s of the current user, i.e. the user owning the API key the Terraform Provider has been initialised with.", attributeLabel),
	}
}

func dataSourceNewRelicCurrentUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	resp, err := client.CustomerAdministration.GetUser()

	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("failed to fetch current user"))
	}

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set(UserDataSourceUserNameAttrLabel, resp.Name)
	_ = d.Set(UserDataSourceUserEmailAttrLabel, resp.Email)
	return nil
}
