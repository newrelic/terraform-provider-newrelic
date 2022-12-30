package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNewRelicInstall() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicInstallRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the account in New Relic.",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API key to use, will default to the one provided to Terraform.",
			},
			"os": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Target operating system (linux, windows, darwin)",
				ValidateFunc: validation.StringInSlice([]string{"linux", "windows", "darwin"}, false),
			},
			"recipe": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Define a list of recipes you want to run. By default the CLI will automatically detected the right recipes to run for your system.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"download_cli": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `By default the command will download the latest version of the CLI. You can disable this behaviour to use the one already on the system.`,
			},
			"assume_yes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Assume yes on all prompts by the installation.`,
			},
			"tag": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A tag applied to the agents installed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The tag value.",
						},
					},
				},
			},
			"command": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auto generated command to use in other parts of your Terraform environment, or externally.",
			},
		},
	}
}

func dataSourceNewRelicInstallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	os := d.Get("os").(string)
	command := ""

	if os == "linux" || os == "darwin" {
		command = buildForLinux(d, meta)
	} else if os == "windows" {
		command = buildForWindow(d, meta)
	} else {
		return diag.Errorf("Unknown OS")
	}

	id := uuid.New()
	d.SetId(id.String())

	if err := d.Set("command", command); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getAPIKey(providerConfig *ProviderConfig, d *schema.ResourceData) string {
	resourceAPIKey := d.Get("api_key").(string)

	if resourceAPIKey != "" {
		return resourceAPIKey
	}

	return providerConfig.PersonalAPIKey
}

func buildForLinux(d *schema.ResourceData, meta interface{}) string {
	providerConfig := meta.(*ProviderConfig)

	command := ""
	if attr, ok := d.GetOk("download_cli"); ok && attr.(bool) {
		command += "curl -Ls https://download.newrelic.com/install/newrelic-cli/scripts/install.sh | bash && "
	}

	command += fmt.Sprintf(
		"sudo NEW_RELIC_API_KEY=%s NEW_RELIC_ACCOUNT_ID=%d /usr/local/bin/newrelic install %s %s",
		getAPIKey(providerConfig, d),
		selectAccountID(providerConfig, d),
		buildAdditionalParams(d, meta),
		buildRecipe(d, meta),
	)

	return command
}

func buildForWindow(d *schema.ResourceData, meta interface{}) string {
	providerConfig := meta.(*ProviderConfig)

	command := ""
	if attr, ok := d.GetOk("download_cli"); ok && attr.(bool) {
		command += "[Net.ServicePointManager]::SecurityProtocol = 'tls12, tls'; (New-Object System.Net.WebClient).DownloadFile(\"https://download.newrelic.com/install/newrelic-cli/scripts/install.ps1\", \"$env:TEMP\\install.ps1\"); & PowerShell.exe -ExecutionPolicy Bypass -File $env:TEMP\\install.ps1; "
	}

	command += fmt.Sprintf(
		"$env:NEW_RELIC_API_KEY='%s'; $env:NEW_RELIC_ACCOUNT_ID='%d'; & 'C:\\Program Files\\New Relic\\New Relic CLI\\newrelic.exe' install %s %s",
		getAPIKey(providerConfig, d),
		selectAccountID(providerConfig, d),
		buildAdditionalParams(d, meta),
		buildRecipe(d, meta),
	)

	return command
}

func buildRecipe(d *schema.ResourceData, meta interface{}) string {
	if recipes, ok := d.GetOk("recipe"); ok {
		recipeNames := recipes.([]interface{})
		recipeArray := make([]string, len(recipes.([]interface{})))
		for i := range recipeNames {
			recipeArray[i] = recipeNames[i].(string)
		}

		return "-n " + strings.Join(recipeArray, ",")
	}

	return ""
}

func buildAdditionalParams(d *schema.ResourceData, meta interface{}) string {
	additionalParams := ""

	// Assume yes
	if attr, ok := d.GetOk("assume_yes"); ok && attr.(bool) {
		additionalParams += "-y "
	}

	// Tags
	tagsRaw := d.Get("tag").([]interface{})
	tags := make([]string, 0, len(tagsRaw))

	for _, t := range tagsRaw {
		tag := t.(map[string]interface{})
		if k, ok := tag["key"]; ok {
			if v, ok := tag["value"]; ok {
				tags = append(tags, fmt.Sprintf("%s:%s", k.(string), v.(string)))
			}
		}
	}

	if len(tags) > 0 {
		additionalParams += fmt.Sprintf("--tag %s", strings.Join(tags, ",")) + " "
	}

	return additionalParams
}
