package newrelic

import (
	"context"
	"encoding/base64"
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
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "",
				Elem:        objectiveSchema(),
			},
			"sli_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"sli_guid": {
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
				ForceNew:     true,
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
			"select": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				MinItems:    0,
				MaxItems:    1,
				Elem:        eventsQuerySelectSchema(),
			},
		},
	}
}

func eventsQuerySelectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attribute": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"function": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "",
				ValidateFunc: validation.StringInSlice([]string{"COUNT", "SUM"}, false),
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
	entityGUID := d.Get("guid").(string)
	createInput := expandServiceLevelCreateInput(d)

	if createInput.Events.GoodEvents == nil && createInput.Events.BadEvents == nil {
		return diag.Errorf("err: Defining a new SLI requires a good or bad events query.")
	}
	if createInput.Events.GoodEvents != nil && createInput.Events.BadEvents != nil {
		return diag.Errorf("err: Only a good or bad events query can be defined for an SLI.")
	}

	log.Printf("[INFO] Creating New Relic One Service Level %s", createInput.Name)

	created, err := client.ServiceLevel.ServiceLevelCreateWithContext(ctx, common.EntityGUID(entityGUID), createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	identifier := serviceLevelIdentifier{
		AccountID:  createInput.Events.AccountID,
		ID:         created.ID,
		EntityGUID: entityGUID,
	}

	sliGUID := getSliGUID(&identifier)

	d.SetId(identifier.String())
	_ = d.Set("sli_id", created.ID)
	_ = d.Set("sli_guid", sliGUID)

	return diag.FromErr(flattenServiceLevelIndicator(*created, &identifier, d, sliGUID))
}

func resourceNewRelicServiceLevelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	identifier, err := parseIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sliGUID := getSliGUID(identifier)
	indicators, err := client.ServiceLevel.GetIndicatorsWithContext(ctx, common.EntityGUID(sliGUID))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	for _, indicator := range *indicators {
		if indicator.ID == identifier.ID {
			return diag.FromErr(flattenServiceLevelIndicator(indicator, identifier, d, sliGUID))
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

	_, err = client.ServiceLevel.ServiceLevelUpdateWithContext(ctx, common.EntityGUID(getSliGUID(identifier)), updateInput)
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

	if _, err := client.ServiceLevel.ServiceLevelDeleteWithContext(ctx, common.EntityGUID(getSliGUID(identifier))); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

type serviceLevelIdentifier struct {
	AccountID  int
	ID         string
	EntityGUID string
}

func (identifier *serviceLevelIdentifier) String() string {
	return fmt.Sprintf("%d:%s:%s", identifier.AccountID, identifier.ID, identifier.EntityGUID)
}

func parseIdentifier(ids string) (*serviceLevelIdentifier, error) {
	split := strings.Split(ids, ":")

	accountID, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return nil, err
	}

	return &serviceLevelIdentifier{
		AccountID:  int(accountID),
		ID:         split[1],
		EntityGUID: split[2],
	}, nil
}

func getSliGUID(identifier *serviceLevelIdentifier) string {
	rawGUID := fmt.Sprintf("%d|EXT|SERVICE_LEVEL|%s", identifier.AccountID, identifier.ID)
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(rawGUID))
}
