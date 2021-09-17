package selectel

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/selectel/dbaas-go"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/projects"
)

func TestAccDBaaSPrometheusMetricTokenV1Basic(t *testing.T) {
	var (
		dbaasToken dbaas.PrometheusMetricToken
		project    projects.Project
	)

	projectName := acctest.RandomWithPrefix("tf-acc")
	tokenName := acctest.RandomWithPrefix("tf-acc-token")
	updatedTokenName := acctest.RandomWithPrefix("tf-acc-token-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBaaSPrometheusMetricTokenV1Basic(projectName, tokenName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCV2ProjectExists("selectel_vpc_project_v2.project_tf_acc_test_1", &project),
					testAccDBaaSPrometheusMetricTokenV1Exists("selectel_dbaas_prometheus_metric_token_v1.prometheus_metric_token_tf_acc_test_1", &dbaasToken),
					resource.TestCheckResourceAttr("selectel_dbaas_prometheus_metric_token_v1.prometheus_metric_token_tf_acc_test_1", "name", tokenName),
				),
			},
			{
				Config: testAccDBaaSPrometheusMetricTokenV1Update(projectName, updatedTokenName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCV2ProjectExists("selectel_vpc_project_v2.project_tf_acc_test_1", &project),
					testAccDBaaSPrometheusMetricTokenV1Exists("selectel_dbaas_prometheus_metric_token_v1.prometheus_metric_token_tf_acc_test_1", &dbaasToken),
					resource.TestCheckResourceAttr("selectel_dbaas_prometheus_metric_token_v1.prometheus_metric_token_tf_acc_test_1", "name", updatedTokenName),
				),
			},
		},
	})
}

func testAccDBaaSPrometheusMetricTokenV1Exists(n string, dbaasToken *dbaas.PrometheusMetricToken) resource.TestCheckFunc {
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

		token, err := dbaasClient.PrometheusMetricToken(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if token.ID != rs.Primary.ID {
			return errors.New("token not found")
		}

		*dbaasToken = token

		return nil
	}
}

func testAccDBaaSPrometheusMetricTokenV1Basic(projectName, name string) string {
	return fmt.Sprintf(`
resource "selectel_vpc_project_v2" "project_tf_acc_test_1" {
  name        = "%s"
  auto_quotas = true
}

resource "selectel_dbaas_prometheus_metric_token_v1" "prometheus_metric_token_tf_acc_test_1" {
  project_id = "${selectel_vpc_project_v2.project_tf_acc_test_1.id}"
  region     = "ru-3"
  name       = "%s"
}
`, projectName, name)
}

func testAccDBaaSPrometheusMetricTokenV1Update(projectName, name string) string {
	return fmt.Sprintf(`
resource "selectel_vpc_project_v2" "project_tf_acc_test_1" {
  name        = "%s"
  auto_quotas = true
}

resource "selectel_dbaas_prometheus_metric_token_v1" "prometheus_metric_token_tf_acc_test_1" {
  project_id = "${selectel_vpc_project_v2.project_tf_acc_test_1.id}"
  region     = "ru-3"
  name       = "%s"
}
`, projectName, name)
}
