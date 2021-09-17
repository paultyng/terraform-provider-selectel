package selectel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/selectel/dbaas-go"
)

func resourceDBaaSPrometheusMetricTokenV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDBaaSPrometheusMetricTokenV1Create,
		ReadContext:   resourceDBaaSPrometheusMetricTokenV1Read,
		UpdateContext: resourceDBaaSPrometheusMetricTokenV1Update,
		DeleteContext: resourceDBaaSPrometheusMetricTokenV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDBaaSPrometheusMetricTokenV1ImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					ru1Region,
					ru2Region,
					ru3Region,
					ru7Region,
					ru8Region,
					ru9Region,
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDBaaSPrometheusMetricTokenV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbaasClient, diagErr := getDBaaSClient(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	prometheusMetricsTokenCreateOpts := dbaas.PrometheusMetricTokenCreateOpts{
		Name: d.Get("name").(string),
	}

	log.Print(msgCreate(objectPrometheusMetricToken, prometheusMetricsTokenCreateOpts))
	token, err := dbaasClient.CreatePrometheusMetricToken(ctx, prometheusMetricsTokenCreateOpts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectPrometheusMetricToken, err))
	}

	d.SetId(token.ID)

	return resourceDBaaSPrometheusMetricTokenV1Read(ctx, d, meta)
}

func resourceDBaaSPrometheusMetricTokenV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbaasClient, diagErr := getDBaaSClient(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgGet(objectPrometheusMetricToken, d.Id()))
	token, err := dbaasClient.PrometheusMetricToken(ctx, d.Id())
	if err != nil {
		return diag.FromErr(errGettingObject(objectPrometheusMetricToken, d.Id(), err))
	}
	d.Set("name", token.Name)
	d.Set("value", token.Value)

	return nil
}

func resourceDBaaSPrometheusMetricTokenV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbaasClient, diagErr := getDBaaSClient(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)

		updateOpts := dbaas.PrometheusMetricTokenUpdateOpts{
			Name: name,
		}

		log.Print(msgUpdate(objectPrometheusMetricToken, d.Id(), updateOpts))
		_, err := dbaasClient.UpdatePrometheusMetricToken(ctx, d.Id(), updateOpts)
		if err != nil {
			return diag.FromErr(errUpdatingObject(objectPrometheusMetricToken, d.Id(), err))
		}
	}

	return resourceDBaaSPrometheusMetricTokenV1Read(ctx, d, meta)
}

func resourceDBaaSPrometheusMetricTokenV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbaasClient, diagErr := getDBaaSClient(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	log.Print(msgDelete(objectPrometheusMetricToken, d.Id()))
	err := dbaasClient.DeletePrometheusMetricToken(ctx, d.Id())
	if err != nil {
		return diag.FromErr(errDeletingObject(objectPrometheusMetricToken, d.Id(), err))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{strconv.Itoa(http.StatusOK)},
		Target:     []string{strconv.Itoa(http.StatusNotFound)},
		Refresh:    dbaasPrometheusMetricTokenV1DeleteStateRefreshFunc(ctx, dbaasClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] waiting for token %s to become deleted", d.Id())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for the token %s to become deleted: %s", d.Id(), err))
	}

	return nil
}

func resourceDBaaSPrometheusMetricTokenV1ImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if config.ProjectID == "" {
		return nil, errors.New("SEL_PROJECT_ID must be set for the resource import")
	}
	if config.Region == "" {
		return nil, errors.New("SEL_REGION must be set for the resource import")
	}

	d.Set("project_id", config.ProjectID)
	d.Set("region", config.Region)

	return []*schema.ResourceData{d}, nil
}

func dbaasPrometheusMetricTokenV1DeleteStateRefreshFunc(ctx context.Context, client *dbaas.API, prometheusMetricsTokenID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		d, err := client.PrometheusMetricToken(ctx, prometheusMetricsTokenID)
		if err != nil {
			var dbaasError *dbaas.DBaaSAPIError
			if errors.As(err, &dbaasError) {
				return d, strconv.Itoa(dbaasError.StatusCode()), nil
			}
			return nil, "", err
		}

		return d, strconv.Itoa(200), err
	}
}
