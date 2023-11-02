package newrelic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/accountmanagement"
)

func resourceNewRelicMonitorDowntime() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicMonitorDowntimeCreate,
		ReadContext:   resourceNewRelicMonitorDowntimeRead,
		UpdateContext: resourceNewRelicMonitorDowntimeUpdate,
		DeleteContext: resourceNewRelicMonitorDowntimeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "A name to identify the Monitor Downtime to be created.",
				Required:    true,
			},
			"monitor_guids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of GUIDs of monitors, to which the created Monitor Downtime shall be applied.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The ID of the New Relic account in which the Monitor Downtime shall be created. Defaults to NEW_RELIC_ACCOUNT_ID if not specified.",
				Optional:    true,
			},
			"start_time": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A datetime stamp signifying the start of the Monitor Downtime.",
				ValidateFunc: validateNaiveDateTime,
			},
			"end_time": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A datetime stamp signifying the end of the Monitor Downtime.",
				ValidateFunc: validateNaiveDateTime,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The timezone that applies to the Monitor Downtime schedule.",
			},
			// used with daily, weekly and monthly monitor downtime
			"end_repeat": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "A specification of when the Monitor Downtime should end its repeat cycle, by number of occurrences or date.",
				// TODO: define validation to not use this with createOnce monitor downtime and keep this optional with other three types
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_date": {
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"on_date", "on_repeat"},
							Description:  "A date, on which the Monitor Downtime's repeat cycle is expected to end.",
							// TODO: define date validation here (possibly YYYY-MM-DD), didn't do it yet as the mutation is broken
						},
						"on_repeat": {
							Type:         schema.TypeInt,
							Optional:     true,
							ExactlyOneOf: []string{"on_date", "on_repeat"},
							Description:  "Number of repetitions after which the Monitor Downtime's repeat cycle is expected to end.",
						},
					},
				},
			},
			// used with weekly monitor downtime
			"maintenance_days": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of maintenance days to be included with the created weekly Monitor Downtime.",
				// TODO: define validation in such a way that this is used only with weekly monitor downtime
				// TODO: (reqd. only for weekly and not allowed for the rest)
				// TODO: in that function, include an "if" to check if it is "MONDAY", "TUESDAY", ... "SUNDAY"
				// TODO: !! Also check if this works as a list or as a single string; NG is not clear !!
			},
			// used with monthly monitor downtime
			"frequency": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "Configuration options for which days of the month a monitor downtime will occur",
				// TODO: define validation to use this only with monthly monitor downtime
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days_of_month": {
							Type:         schema.TypeList,
							Elem:         &schema.Schema{Type: schema.TypeInt},
							Optional:     true,
							ExactlyOneOf: []string{"days_of_month", "days_of_week"},
							Description:  "A numerical list of days of a month on which the Monitor Downtime is scheduled to run.",
							// TODO: define validation to have these values between 1 and 31
						},
						"days_of_week": {
							Type:         schema.TypeList,
							MinItems:     1,
							MaxItems:     1,
							Optional:     true,
							ExactlyOneOf: []string{"days_of_month", "days_of_week"},
							Description:  "A list of days of the week on which the Monitor Downtime is scheduled to run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ordinal_day_of_month": {
										Type:         schema.TypeString,
										Required:     true,
										ExactlyOneOf: []string{"on_date", "on_repeat"},
										Description:  "An occurrence of the day selected within the month.",
										// TODO: define this to belong to ["FIRST", "SECOND", "THIRD", "FOURTH", "LAST"]
									},
									"week_day": {
										Type:         schema.TypeInt,
										Required:     true,
										ExactlyOneOf: []string{"on_date", "on_repeat"},
										Description:  "The day of the week on which the Monitor Downtime would run.",
										// TODO: define this to belong to ["MONDAY", "TUESDAY", ... "SUNDAY"]
									},
								},
							},
						},
					},
				},
			},
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "An identifier of the type of Monitor Downtime to be created.",
				ValidateFunc: validation.StringInSlice([]string{"ONCE", "DAILY", "MONTHLY", "WEEKLY"}, false),
			},
		},
	}
}

// getCommonAttributes helps obtain the attributes required, and common to all monitor downtimes
func getCommonAttributes(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// do something here
	return diag.FromErr(errors.New("Some dysfunctional error"))
}

func resourceNewRelicMonitorDowntimeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: WRITE THE CREATE METHOD

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	var downtimeMode string
	if m, ok := d.GetOk("mode"); ok {
		downtimeMode = m.(string)
	}

	switch downtimeMode {
	case "ONCE":
		break
	case "DAILY":
		break
	case "WEEKLY":
		break
	case "MONTHLY":
		break
	default:
		return diag.FromErr(errors.New("invalid mode of operation: 'mode' can be 'ONCE', 'DAILY', 'WEEKLY' or 'MONTHLY'"))

	}

	createAccountInput := accountmanagement.AccountManagementCreateInput{
		Name:       d.Get("name").(string),
		RegionCode: d.Get("region").(string),
	}
	created, err := client.AccountManagement.AccountManagementCreateAccount(createAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: Account not created. Please check the input details")
	}
	accountID := created.ManagedAccount.ID

	d.SetId(strconv.Itoa(accountID))
	return resourceNewRelicAccountRead(ctx, d, meta)
}

func resourceNewRelicMonitorDowntimeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: WRITE THE READ METHOD

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		account, err := getCreatedAccountByID(client, d.Id())
		//		fmt.Println("read", account.ID, err)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if account == nil {
			return resource.RetryableError(fmt.Errorf("account not found"))
		}
		_ = d.Set("region", account.RegionCode)
		_ = d.Set("name", account.Name)

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}
	return nil
}

func resourceNewRelicMonitorDowntimeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: WRITE THE UPDATE METHOD

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	updateAccountInput := accountmanagement.AccountManagementUpdateInput{
		Name: d.Get("name").(string),
		ID:   accountID,
	}
	updated, err := client.AccountManagement.AccountManagementUpdateAccount(updateAccountInput)

	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		return diag.Errorf("err: Account not Updated. Please check the input details")
	}

	return resourceNewRelicAccountRead(ctx, d, meta)
}

func resourceNewRelicMonitorDowntimeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: WRITE THE DELETE METHOD

	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Account cannot be deleted via Terraform. https://docs.newrelic.com/docs/apis/nerdgraph/examples/manage-accounts-nerdgraph/#delete",
	})
	return diags
}
