package selectel

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/projects"
)

const (
	dataSourceFeatureGates         = "selectel_mks_available_feature_gates_v1"
	dataSourceAdmissionControllers = "selectel_mks_available_admission_controllers_v1"
)

func TestAccMKSAvailableFeatureGatesV1Basic(t *testing.T) {
	var (
		project projects.Project
	)

	projectName := acctest.RandomWithPrefix("tf-acc")
	kubeVersion := "1.17.3"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKubeOptionsV1BasicConfig(projectName, dataSourceFeatureGates, kubeVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCV2ProjectExists("selectel_vpc_project_v2.project_tf_acc_test_1", &project),
					resource.TestCheckResourceAttr("data."+dataSourceFeatureGates+".dt", "available_feature_gates.0.kube_version_minor", kubeVersion),
					testFeatureGatesIsNotEmpty("data."+dataSourceFeatureGates+".dt"),
				),
			},
		},
	})
}

func TestAccMKSAvailableFeatureGatesV1NoFilter(t *testing.T) {
	var (
		project projects.Project
	)

	projectName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKubeOptionsV1ConfigNoFilter(projectName, dataSourceFeatureGates),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCV2ProjectExists("selectel_vpc_project_v2.project_tf_acc_test_1", &project),
					testFeatureGatesNoFilter("data."+dataSourceFeatureGates+".dt"),
				),
			},
		},
	})
}

func TestAccMKSAvailableAdmissionControllersV1Basic(t *testing.T) {
	var (
		project projects.Project
	)

	projectName := acctest.RandomWithPrefix("tf-acc")
	kubeVersion := "1.17.3"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKubeOptionsV1BasicConfig(projectName, dataSourceAdmissionControllers, kubeVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCV2ProjectExists("selectel_vpc_project_v2.project_tf_acc_test_1", &project),
					resource.TestCheckResourceAttr("data."+dataSourceAdmissionControllers+".dt", "available_admission_controllers.0.kube_version_minor", kubeVersion),
					testAdmissionControllersIsNotEmpty("data."+dataSourceAdmissionControllers+".dt"),
				),
			},
		},
	})
}

func TestAccMKSAvailableAdmissionControllersV1NoFilter(t *testing.T) {
	var (
		project projects.Project
	)

	projectName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPCV2ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKubeOptionsV1ConfigNoFilter(projectName, dataSourceAdmissionControllers),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCV2ProjectExists("selectel_vpc_project_v2.project_tf_acc_test_1", &project),
					testAdmissionControllersNoFilter("data."+dataSourceAdmissionControllers+".dt"),
				),
			},
		},
	})
}

func testFeatureGatesIsNotEmpty(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		availableFeatureGates, ok := rs.Primary.Attributes["available_feature_gates.0.names"]
		if !ok {
			return fmt.Errorf("attribute 'available_feature_gates' is not found")
		}
		if availableFeatureGates == "" {
			return fmt.Errorf("names is not set at 'available_feature_gates'")
		}

		return nil
	}
}
func testFeatureGatesNoFilter(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		availableFeatureGates, ok := rs.Primary.Attributes["available_feature_gates.0.names"]
		if !ok {
			return fmt.Errorf("attribute 'available_feature_gates' is not found")
		}
		if availableFeatureGates == "" {
			return fmt.Errorf("names is not set at 'available_feature_gates'")
		}

		fgLen, err := strconv.Atoi(rs.Primary.Attributes["available_feature_gates.#"])
		if err != nil {
			return fmt.Errorf("failed to get len of 'available_feature_gates': %w", err)
		}
		if fgLen <= 1 {
			return fmt.Errorf("received only one or less item in 'available_feature_gates'")
		}

		return nil
	}
}

func testAdmissionControllersIsNotEmpty(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		availableAdmissionControllers, ok := rs.Primary.Attributes["available_admission_controllers.0.names"]
		if !ok {
			return fmt.Errorf("attribute 'available_admission_controllers' is not found")
		}
		if availableAdmissionControllers == "" {
			return fmt.Errorf("names is not set at 'available_admission_controllers'")
		}

		return nil
	}
}
func testAdmissionControllersNoFilter(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		availableAdmissionControllers, ok := rs.Primary.Attributes["available_admission_controllers.0.names"]
		if !ok {
			return fmt.Errorf("attribute 'available_admission_controllers' is not found")
		}
		if availableAdmissionControllers == "" {
			return fmt.Errorf("names is not set at 'available_admission_controllers'")
		}

		fgLen, err := strconv.Atoi(rs.Primary.Attributes["available_admission_controllers.#"])
		if err != nil {
			return fmt.Errorf("failed to get len of 'available_admission_controllers': %w", err)
		}
		if fgLen <= 1 {
			return fmt.Errorf("received only one or less item in 'available_admission_controllers'")
		}

		return nil
	}
}

func testKubeOptionsV1BasicConfig(projectName, dataSource, kubeVersion string) string {
	return fmt.Sprintf(`
resource "selectel_vpc_project_v2" "project_tf_acc_test_1" {
  name        = "%s"
  auto_quotas = true
}

data "%s" "dt" {
  project_id = "${selectel_vpc_project_v2.project_tf_acc_test_1.id}"
  region     = "ru-3"
  filter {
    kube_version_minor = "%s"
  }
}
`, projectName, dataSource, kubeVersion)
}

func testKubeOptionsV1ConfigNoFilter(projectName, dataSource string) string {
	return fmt.Sprintf(`
resource "selectel_vpc_project_v2" "project_tf_acc_test_1" {
  name        = "%s"
  auto_quotas = true
}

data "%s" "dt" {
  project_id = "${selectel_vpc_project_v2.project_tf_acc_test_1.id}"
  region     = "ru-3"
}
`, projectName, dataSource)
}
