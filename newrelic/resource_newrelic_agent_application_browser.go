package newrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/agentapplications"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func resourceNewRelicBrowserApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicBrowserApplicationCreate,
		ReadContext:   resourceNewRelicBrowserApplicationRead,
		UpdateContext: resourceNewRelicBrowserApplicationUpdate,
		DeleteContext: resourceNewRelicBrowserApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the application to monitor.",
			},
			"cookies_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Configure cookies. The default is enabled: true.",
			},
			"distributed_tracing_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Configure distributed tracing in browser apps. The default is enabled: true.",
			},
			"loader_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     agentapplications.AgentApplicationBrowserLoaderTypes.SPA,
				Description: `Determines which browser loader is configured. The default is "SPA".`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The GUID of the browser application.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the browser application.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The account ID.",
			},
			"js_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JavaScript configuration of the browser application encoded into a string.",
			},
		},
	}
}

func resourceNewRelicBrowserApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)
	appName := d.Get("name").(string)
	settingsInput := agentapplications.AgentApplicationBrowserSettingsInput{
		CookiesEnabled:            d.Get("cookies_enabled").(bool),
		DistributedTracingEnabled: d.Get("distributed_tracing_enabled").(bool),
		LoaderType:                agentapplications.AgentApplicationBrowserLoader(strings.ToUpper(d.Get("loader_type").(string))),
	}

	resp, err := client.AgentApplications.AgentApplicationCreateBrowserWithContext(ctx, accountID, appName, settingsInput)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("potential error creating browser application resource, response was nil"))
	}

	d.SetId(string(resp.GUID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("account_id", accountID)
	_ = d.Set("cookies_enabled", resp.Settings.CookiesEnabled)
	_ = d.Set("distributed_tracing_enabled", resp.Settings.DistributedTracingEnabled)
	_ = d.Set("loader_type", string(resp.Settings.LoaderType))
	_ = d.Set("guid", string(resp.GUID))

	return resourceNewRelicBrowserApplicationRead(ctx, d, meta)
}

func resourceNewRelicBrowserApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Id()

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(guid))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("entity with GUID %s was nil", guid))
	}

	switch (*resp).(type) {
	case *entities.BrowserApplicationEntity:
		entity := (*resp).(*entities.BrowserApplicationEntity)

		d.SetId(string(entity.GUID))
		_ = d.Set("name", entity.Name)
		_ = d.Set("cookies_enabled", entity.BrowserSettings.BrowserMonitoring.Privacy.CookiesEnabled)
		_ = d.Set("distributed_tracing_enabled", entity.BrowserSettings.BrowserMonitoring.DistributedTracing.Enabled)
		_ = d.Set("loader_type", string(entity.BrowserSettings.BrowserMonitoring.Loader))
		_ = d.Set("guid", string(entity.GUID))
		_ = d.Set("account_id", entity.AccountID)
		_ = d.Set("application_id", strconv.Itoa(entity.ApplicationID))

		// The following block of code encodes the JavaScript configuration of the browser application into a JSON,
		// that is then exported as a string, which can be accessed from the Terraform Configuration using jsondecode().
		jsonOutput, err := json.Marshal(entity.BrowserProperties.JsConfig)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("js_config", string(jsonOutput)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceNewRelicBrowserApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	settingsInput := agentapplications.AgentApplicationSettingsUpdateInput{
		BrowserMonitoring: &agentapplications.AgentApplicationSettingsBrowserMonitoringInput{
			DistributedTracing: &agentapplications.AgentApplicationSettingsBrowserDistributedTracingInput{
				Enabled: d.Get("distributed_tracing_enabled").(bool),
			},
			Privacy: &agentapplications.AgentApplicationSettingsBrowserPrivacyInput{
				CookiesEnabled: d.Get("cookies_enabled").(bool),
			},
		},
	}

	guid := d.Id()

	resp, err := client.AgentApplications.AgentApplicationSettingsUpdateWithContext(ctx, common.EntityGUID(guid), settingsInput)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("entity with GUID %s was nil", guid))
	}

	d.SetId(string(resp.GUID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("cookies_enabled", resp.BrowserSettings.BrowserMonitoring.Privacy.CookiesEnabled)
	_ = d.Set("distributed_tracing_enabled", resp.BrowserSettings.BrowserMonitoring.DistributedTracing.Enabled)
	_ = d.Set("loader_type", string(resp.BrowserSettings.BrowserMonitoring.Loader))
	_ = d.Set("guid", string(resp.GUID))

	return resourceNewRelicBrowserApplicationRead(ctx, d, meta)
}

func resourceNewRelicBrowserApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	guid := d.Id()

	log.Printf("[INFO] Deleting New Relic browser application %s", guid)

	_, err := client.AgentApplications.AgentApplicationDeleteWithContext(ctx, common.EntityGUID(guid))
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}
