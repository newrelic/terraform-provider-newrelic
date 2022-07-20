package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

type SyntheticsMonitorType string

func resourceNewRelicSyntheticsBrokenLinksMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsBrokenLinksMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsBrokenLinksMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsBrokenLinksMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsBrokenLinksMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: mergeSchemas(
			syntheticsMonitorBrokenLinksMonitorSchema(),
			syntheticsMonitorCommonSchema(),
			syntheticsMonitorLocationsAsStringsSchema(),
		),
	}
}

// Returns the common schema attributes shared by all Synthetics monitor types.
func syntheticsMonitorCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Description: "ID of the newrelic account",
			Computed:    true,
			Optional:    true,
		},
		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique entity identifier of the monitor in New Relic.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The title of this monitor.",
		},
		"status": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The monitor status (i.e. ENABLED, MUTED, DISABLED).",
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
	}
}

// NOTE: This can likely be a shared schema partial for other synthetics monitor resources
func syntheticsMonitorLocationsAsStringsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"locations_private": {
			Type:         schema.TypeList,
			Elem:         &schema.Schema{Type: schema.TypeString},
			Description:  "List private location GUIDs for which the monitor will run.",
			Optional:     true,
			AtLeastOneOf: []string{"locations_public", "locations_private"},
		},
		"locations_public": {
			Type:         schema.TypeList,
			Elem:         &schema.Schema{Type: schema.TypeString},
			Description:  "Publicly available location names in which the monitor will run.",
			Optional:     true,
			AtLeastOneOf: []string{"locations_public", "locations_private"},
		},
	}
}

func syntheticsMonitorBrokenLinksMonitorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// "createdAt": {
		// 	Type:     schema.TypeInt,
		// 	Computed: true,
		// },
		// "modifiedAt": {
		// 	Type:     schema.TypeInt,
		// 	Computed: true,
		// },
		"uri": {
			Type:        schema.TypeString,
			Description: "The URI the monitor runs against.",
			Required:    true,
		},
	}
}

func resourceNewRelicSyntheticsBrokenLinksMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	monitorInput := buildSyntheticsBrokenLinksMonitorCreateInput(d)
	resp, err := client.Synthetics.SyntheticsCreateBrokenLinksMonitorWithContext(ctx, accountID, *monitorInput)
	if err != nil {
		return diag.FromErr(err)
	}

	errors := buildCreateSyntheticsMonitorResponseErrors(resp.Errors)
	if len(errors) > 0 {
		return errors
	}

	// Set attributes
	d.SetId(string(resp.Monitor.GUID))
	_ = d.Set("account_id", accountID)
	attrs := map[string]string{"guid": string(resp.Monitor.GUID)}
	if err = setSyntheticsMonitorAttributes(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicSyntheticsBrokenLinksMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	switch e := (*resp).(type) {
	case *entities.SyntheticMonitorEntity:
		entity := (*resp).(*entities.SyntheticMonitorEntity)

		d.SetId(string(e.GUID))
		_ = d.Set("account_id", accountID)

		err = setSyntheticsMonitorAttributes(d, map[string]string{
			"guid":   string(e.GUID),
			"name":   entity.Name,
			"period": string(syntheticsMonitorPeriodValueMap[int(entity.GetPeriod())]),
			"status": string(entity.MonitorSummary.Status),
			"uri":    entity.MonitoredURL,
		})
	}

	return diag.FromErr(err)
}

func resourceNewRelicSyntheticsBrokenLinksMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	guid := synthetics.EntityGUID(d.Id())

	monitorInput := buildSyntheticsBrokenLinksMonitorUpdateInput(d)
	resp, err := client.Synthetics.SyntheticsUpdateBrokenLinksMonitorWithContext(ctx, guid, *monitorInput)
	if err != nil {
		return diag.FromErr(err)
	}

	errors := buildUpdateSyntheticsMonitorResponseErrors(resp.Errors)
	if len(errors) > 0 {
		return errors
	}

	err = setSyntheticsMonitorAttributes(d, map[string]string{
		"name": resp.Monitor.Name,
		"guid": string(resp.Monitor.GUID),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// NOTE: We can make rename this to reusable function for all new monitor types,
//       but the legacy function already has a good generic name (`resourceNewRelicSyntheticsMonitorDelete()`)
func resourceNewRelicSyntheticsBrokenLinksMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
