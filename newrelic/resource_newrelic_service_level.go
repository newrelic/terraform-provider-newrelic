package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicServiceLevel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicServiceLevelCreate,
		ReadContext:   resourceNewRelicServiceLevelRead,
		UpdateContext: resourceNewRelicServiceLevelUpdate,
		DeleteContext: resourceNewRelicServiceLevelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"events": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				MinItems:    1,
				MaxItems:    1,
				Elem:        eventsSchema(),
			},
			"objective": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "",
				Elem:        objectiveSchema(),
			},
			"sli_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func eventsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"valid_events": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				MinItems:    1,
				MaxItems:    1,
				Elem:        eventsQuerySchema(),
			},
			"good_events": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				MinItems:    0,
				MaxItems:    1,
				Elem:        eventsQuerySchema(),
			},
			"bad_events": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				MinItems:    0,
				MaxItems:    1,
				Elem:        eventsQuerySchema(),
			},
		},
	}
}

func eventsQuerySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"from": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"where": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func objectiveSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"target": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "",
			},
			"time_window": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				MinItems:    1,
				MaxItems:    1,
				Elem:        timeWindowSchema(),
			},
		},
	}
}

func timeWindowSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rolling": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				MinItems:    1,
				MaxItems:    1,
				Elem:        rollingTimeWindowSchema(),
			},
		},
	}
}

func rollingTimeWindowSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "",
				ValidateFunc: intInSlice([]int{1, 7, 28}),
			},
			"unit": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "",
				ValidateFunc: validation.StringInSlice([]string{"DAY"}, false),
			},
		},
	}
}

func resourceNewRelicServiceLevelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	guid := d.Get("guid").(string)
	createInput := expandServiceLevelCreateInput(d)

	if createInput.Events.GoodEvents == nil && createInput.Events.BadEvents == nil {
		return diag.Errorf("err: Defining a new SLI requires a good or bad events query.")
	}
	if createInput.Events.GoodEvents != nil && createInput.Events.BadEvents != nil {
		return diag.Errorf("err: Only a good or bad events query can be defined for an SLI.")
	}

	log.Printf("[INFO] Creating New Relic One Service Level %s", createInput.Name)

	created, err := client.ServiceLevel.ServiceLevelCreateWithContext(ctx, common.EntityGUID(guid), createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	identifier := serviceLevelIdentifier{
		AccountID: createInput.Events.AccountID,
		ID:        created.ID,
		GUID:      guid,
	}

	d.SetId(identifier.String())
	_ = d.Set("sli_id", created.ID)

	return diag.FromErr(nil)
}

func resourceNewRelicServiceLevelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	identifier, err := parseIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	indicators, err := client.ServiceLevel.GetIndicatorsWithContext(ctx, common.EntityGUID(identifier.GUID))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return diag.Errorf("err: SLI with id=%s not found.", d.Id())
		}
		return diag.FromErr(err)
	}

	for _, indicator := range *indicators {
		if indicator.ID == identifier.ID {
			return diag.FromErr(flattenServiceLevelIndicator(indicator, identifier, d))
		}
	}

	return diag.Errorf("err: SLI with id=%s not found.", d.Id())
}

func resourceNewRelicServiceLevelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput := expandServiceLevelUpdateInput(d)

	log.Printf("[INFO] Updating New Relic One Service Level %s", d.Id())

	identifier, err := parseIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ServiceLevel.ServiceLevelUpdateWithContext(ctx, identifier.ID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicServiceLevelRead(ctx, d, meta)
}

func resourceNewRelicServiceLevelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One Service Level %s", d.Id())

	identifier, err := parseIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err := client.ServiceLevel.ServiceLevelDeleteWithContext(ctx, identifier.ID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type serviceLevelIdentifier struct {
	AccountID int
	ID        string
	GUID      string
}

func (identifier *serviceLevelIdentifier) String() string {
	return fmt.Sprintf("%d:%s:%s", identifier.AccountID, identifier.ID, identifier.GUID)
}

func parseIdentifier(ids string) (*serviceLevelIdentifier, error) {
	split := strings.Split(ids, ":")

	accountID, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return nil, err
	}

	return &serviceLevelIdentifier{
		AccountID: int(accountID),
		ID:        split[1],
		GUID:      split[2],
	}, nil
}
