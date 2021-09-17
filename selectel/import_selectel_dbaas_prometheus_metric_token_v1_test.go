package selectel

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDBaaSPrometheusMetricTokenV1ImportBasic(t *testing.T) {
	resourceName := "selectel_dbaas_prometheus_metric_token_v1.prometheus_metric_token_tf_acc_test_1"
	projectName := acctest.RandomWithPrefix("tf-acc")
	tokenName := acctest.RandomWithPrefix("tf-acc-token")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBaaSPrometheusMetricTokenV1Basic(projectName, tokenName),
				Check:  testAccCheckSelectelImportEnv(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
