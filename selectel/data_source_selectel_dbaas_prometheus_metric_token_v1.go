package selectel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/selectel/dbaas-go"
)

func dataSourceDBaaSPrometheusMetricTokenV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDBaaSPrometheusMetricTokenV1Read,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					ru1Region,
					ru2Region,
					ru3Region,
					ru7Region,
					ru8Region,
					ru9Region,
				}, false),
			},
			"prometheus_metrics_tokens": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDBaaSPrometheusMetricTokenV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dbaasClient, diagErr := getDBaaSClient(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	tokens, err := dbaasClient.PrometheusMetricTokens(ctx)
	if err != nil {
		return diag.FromErr(errGettingObjects(objectPrometheusMetricToken, err))
	}

	tokenIDs := []string{}
	for _, token := range tokens {
		tokenIDs = append(tokenIDs, token.ID)
	}

	tokensFlatten := flattenDBaaSPrometheusMetricTokenTypes(tokens)
	if err := d.Set("prometheus_metrics_tokens", tokensFlatten); err != nil {
		return diag.FromErr(err)
	}
	checksum, err := stringListChecksum(tokenIDs)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(checksum)

	return nil
}

func flattenDBaaSPrometheusMetricTokenTypes(tokens []dbaas.PrometheusMetricToken) []interface{} {
	tokensList := make([]interface{}, len(tokens))
	for i, token := range tokens {
		tokensMap := make(map[string]interface{})
		tokensMap["id"] = token.ID
		tokensMap["created_at"] = token.CreatedAt
		tokensMap["updated_at"] = token.UpdatedAt
		tokensMap["project_id"] = token.ProjectID
		tokensMap["name"] = token.Name
		tokensMap["value"] = token.Value

		tokensList[i] = tokensMap
	}

	return tokensList
}
