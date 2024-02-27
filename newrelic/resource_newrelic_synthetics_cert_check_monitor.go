package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func resourceNewRelicSyntheticsCertCheckMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsCertCheckMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsCertCheckMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsCertCheckMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsCertCheckMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "ID of the newrelic account",
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of the cert check monitor",
				Required:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "",
				Required:    true,
			},
			"certificate_expiration": {
				Type:        schema.TypeInt,
				Description: "",
				Required:    true,
			},
			"locations_public": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"locations_public", "locations_private"},
				Description:  "The locations in which this monitor should be run.",
			},
			"locations_private": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"locations_public", "locations_private"},
				Description:  "The locations in which this monitor should be run.",
			},
			"status": {
				Type:         schema.TypeString,
				Description:  "The monitor status (ENABLED or DISABLED).",
				Required:     true,
				ValidateFunc: validateSyntheticMonitorStatus,
			},
			"tag": {
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Description: "The tags that will be associated with the monitor",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the tag key",
						},
						"values": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Values associated with the tag key",
						},
					},
				},
			},
			"period": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.",
				ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorPeriods(), false),
			},
			"period_in_minutes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The interval in minutes at which this monitor should run.",
			},
		},
	}
}

func resourceNewRelicSyntheticsCertCheckMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	var diags diag.Diagnostics

	monitorInput := buildSyntheticsCertCheckMonitorCreateInput(d)
	resp, err := client.Synthetics.SyntheticsCreateCertCheckMonitorWithContext(ctx, accountID, monitorInput)
	if err != nil {
		diag.FromErr(err)
	}

	if len(resp.Errors) > 0 {
		for _, err := range resp.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
			})
		}
	}

	d.SetId(string(resp.Monitor.GUID))
	_ = d.Set("account_id", accountID)
	_ = d.Set("certificate_expiration", resp.Monitor.NumberDaysToFailBeforeCertExpires)
	_ = d.Set("locations_public", resp.Monitor.Locations.Public)
	_ = d.Set("locations_private", resp.Monitor.Locations.Private)
	_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[resp.Monitor.Period])

	err = setSyntheticsMonitorAttributes(d, map[string]string{
		"domain": resp.Monitor.Domain,
		"name":   resp.Monitor.Name,
		"period": string(resp.Monitor.Period),
		"status": string(resp.Monitor.Status),
	})

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
	}

	return diags
}

func resourceNewRelicSyntheticsCertCheckMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// This should probably be in go-client so we can use *errors.NotFound
	if *resp == nil {
		d.SetId("")
		return nil
	}

	switch e := (*resp).(type) {
	case *entities.SyntheticMonitorEntity:
		entity := (*resp).(*entities.SyntheticMonitorEntity)

		d.SetId(string(e.GUID))
		_ = d.Set("account_id", accountID)
		_ = d.Set("locations_public", getPublicLocationsFromEntityTags(entity.GetTags()))
		_ = d.Set("period_in_minutes", int(entity.GetPeriod()))

		err = setSyntheticsMonitorAttributes(d, map[string]string{
			"name":   e.Name,
			"period": string(syntheticsMonitorPeriodValueMap[int(entity.GetPeriod())]),
			"status": string(entity.MonitorSummary.Status),
		})
	}

	return diag.FromErr(err)
}

func resourceNewRelicSyntheticsCertCheckMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	guid := synthetics.EntityGUID(d.Id())

	var diags diag.Diagnostics

	monitorInput := buildSyntheticsCertCheckMonitorUpdateInput(d)
	resp, err := client.Synthetics.SyntheticsUpdateCertCheckMonitorWithContext(ctx, guid, monitorInput)
	if err != nil {
		diag.FromErr(err)
	}
	if len(resp.Errors) > 0 {
		for _, err := range resp.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
			})
		}
	}

	_ = d.Set("certificate_expiration", resp.Monitor.NumberDaysToFailBeforeCertExpires)
	_ = d.Set("locations_public", resp.Monitor.Locations.Public)
	_ = d.Set("locations_private", resp.Monitor.Locations.Private)
	_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[resp.Monitor.Period])

	err = setSyntheticsMonitorAttributes(d, map[string]string{
		"domain": resp.Monitor.Domain,
		"name":   resp.Monitor.Name,
		"period": string(resp.Monitor.Period),
		"status": string(resp.Monitor.Status),
	})

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
	}

	return diags
}

func resourceNewRelicSyntheticsCertCheckMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	guid := synthetics.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	_, err := client.Synthetics.SyntheticsDeleteMonitorWithContext(ctx, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildSyntheticsCertCheckMonitorCreateInput(d *schema.ResourceData) (result synthetics.SyntheticsCreateCertCheckMonitorInput) {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateCertCheckMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
	}

	if v, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("domain"); ok {
		input.Domain = v.(string)
	}

	if v, ok := d.GetOk("certificate_expiration"); ok {
		input.NumberDaysToFailBeforeCertExpires = v.(int)
	}

	return input
}

func buildSyntheticsCertCheckMonitorUpdateInput(d *schema.ResourceData) (result synthetics.SyntheticsUpdateCertCheckMonitorInput) {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateCertCheckMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
	}

	if v, ok := d.GetOk("locations_public"); ok {
		input.Locations.Public = expandStringSlice(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("locations_private"); ok {
		input.Locations.Private = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("domain"); ok {
		input.Domain = v.(string)
	}

	if v, ok := d.GetOk("certificate_expiration"); ok {
		input.NumberDaysToFailBeforeCertExpires = v.(int)
	}
	return input
}
