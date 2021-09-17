package selectel

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/selectel/dbaas-go"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/projects"
)

func TestAccDBaaSDataSourcePrometheusMetricTokenV1Basic(t *testing.T) {
	var (
		dbaasTokens []dbaas.PrometheusMetricToken
		project     projects.Project
	)

	projectName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBaaSDataSourcePrometheusMetricTokenV1Basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCV2ProjectExists("selectel_vpc_project_v2.project_tf_acc_test_1", &project),
					testAccDBaaSPrometheusMetricTokensV1Exists("data.selectel_dbaas_prometheus_metric_token_v1.prometheus_metric_token_tf_acc_test_1", &dbaasTokens),
				),
			},
		},
	})
}

func testAccDBaaSPrometheusMetricTokensV1Exists(n string, dbaasTokens *[]dbaas.PrometheusMetricToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		ctx := context.Background()

		dbaasClient, err := baseTestAccCheckDBaaSV1EntityExists(ctx, rs, testAccProvider)
		if err != nil {
			return err
		}

		tokens, err := dbaasClient.PrometheusMetricTokens(ctx)
		if err != nil {
			return err
		}

		*dbaasTokens = tokens

		return nil
	}
}

func testAccDBaaSDataSourcePrometheusMetricTokenV1Basic(projectName string) string {
	return fmt.Sprintf(`
resource "selectel_vpc_project_v2" "project_tf_acc_test_1" {
  name        = "%s"
  auto_quotas = true
}

data "selectel_dbaas_prometheus_metric_token_v1" "prometheus_metric_token_tf_acc_test_1" {
  project_id = "${selectel_vpc_project_v2.project_tf_acc_test_1.id}"
  region     = "ru-3"
}
`, projectName)
}
