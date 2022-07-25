package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
	"log"
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
			"location_public": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"location_public", "location_private"},
				Description:  "The locations in which this monitor should be run.",
			},
			"location_private": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"location_public", "location_private"},
				Description:  "The locations in which this monitor should be run.",
			},
			"status": {
				Type:         schema.TypeString,
				Description:  "The monitor status (i.e. ENABLED, MUTED, DISABLED).",
				Required:     true,
				ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorStatuses(), false),
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

	setSyntheticsCertCheckMonitorCreateAttributes(d, resp)

	return nil
}

func setSyntheticsCertCheckMonitorCreateAttributes(d *schema.ResourceData, resp *synthetics.SyntheticsCertCheckMonitorCreateMutationResult) {
	_ = d.Set("name", resp.Monitor.Name)
	_ = d.Set("status", resp.Monitor.Status)
	_ = d.Set("domain", resp.Monitor.Domain)
	_ = d.Set("period", resp.Monitor.Period)
	_ = d.Set("certificate_expiration", resp.Monitor.NumberDaysToFailBeforeCertExpires)
	_ = d.Set("location_public", resp.Monitor.Locations.Public)
	_ = d.Set("location_private", resp.Monitor.Locations.Private)
}

func buildSyntheticsCertCheckMonitorCreateInput(d *schema.ResourceData) (result synthetics.SyntheticsCreateCertCheckMonitorInput) {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsCreateCertCheckMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
	}

	if v, ok := d.GetOk("location_public"); ok {
		input.Locations.Public = expandSyntheticsCertMonitorLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsCertMonitorLocations(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("domain"); ok {
		input.Domain = v.(string)
	}

	if v, ok := d.GetOk("certificate_expiration"); ok {
		input.NumberDaysToFailBeforeCertExpires = v.(int)
	}

	return input
}

//function to expand synthetics locations.
func expandSyntheticsCertMonitorLocations(locations []interface{}) []string {
	locationsOut := make([]string, len(locations))

	for i, v := range locations {
		locationsOut[i] = v.(string)
	}
	return locationsOut
}

func resourceNewRelicSyntheticsCertCheckMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("account_id", accountID)

	switch e := (*resp).(type) {
	case *entities.SyntheticMonitorEntity:
		entity := (*resp).(*entities.SyntheticMonitorEntity)

		d.SetId(string(e.GUID))
		_ = d.Set("account_id", accountID)
		//_ = d.Set("locations_public", getPublicLocationsFromEntityTags(entity.GetTags()))

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

	setSyntheticsCertCheckMonitorUpdateAttributes(d, resp)

	return nil
}

func setSyntheticsCertCheckMonitorUpdateAttributes(d *schema.ResourceData, resp *synthetics.SyntheticsCertCheckMonitorUpdateMutationResult) {
	_ = d.Set("name", resp.Monitor.Name)
	_ = d.Set("status", resp.Monitor.Status)
	_ = d.Set("domain", resp.Monitor.Domain)
	_ = d.Set("period", resp.Monitor.Period)
	_ = d.Set("certificate_expiration", resp.Monitor.NumberDaysToFailBeforeCertExpires)
	_ = d.Set("location_public", resp.Monitor.Locations.Public)
	_ = d.Set("location_private", resp.Monitor.Locations.Private)
}

func buildSyntheticsCertCheckMonitorUpdateInput(d *schema.ResourceData) (result synthetics.SyntheticsUpdateCertCheckMonitorInput) {
	inputBase := expandSyntheticsMonitorBase(d)

	input := synthetics.SyntheticsUpdateCertCheckMonitorInput{
		Name:   inputBase.Name,
		Period: inputBase.Period,
		Status: inputBase.Status,
		Tags:   inputBase.Tags,
	}

	if v, ok := d.GetOk("location_public"); ok {
		input.Locations.Public = expandSyntheticsCertMonitorLocations(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("location_private"); ok {
		input.Locations.Private = expandSyntheticsCertMonitorLocations(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("domain"); ok {
		input.Domain = v.(string)
	}

	if v, ok := d.GetOk("certificate_expiration"); ok {
		input.NumberDaysToFailBeforeCertExpires = v.(int)
	}
	return input
}

func resourceNewRelicSyntheticsCertCheckMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	guid := synthetics.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	_, err := client.Synthetics.SyntheticsDeleteMonitorWithContext(ctx, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
