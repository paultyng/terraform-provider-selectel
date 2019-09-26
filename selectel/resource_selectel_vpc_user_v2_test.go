package selectel

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/users"
)

func TestAccVPCV2UserBasic(t *testing.T) {
	var user users.User
	userName := acctest.RandomWithPrefix("tf-acc")
	userNameUpdated := acctest.RandomWithPrefix("tf-acc")
	userPassword := acctest.RandString(8)
	userPasswordUpdated := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccSelectelPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCV2UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCV2UserBasic(userName, userPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCV2UserExists("selectel_vpc_user_v2.user_tf_acc_test_1", &user),
					resource.TestCheckResourceAttr("selectel_vpc_user_v2.user_tf_acc_test_1", "name", userName),
					resource.TestCheckResourceAttr("selectel_vpc_user_v2.user_tf_acc_test_1", "password", userPassword),
					resource.TestCheckResourceAttr("selectel_vpc_user_v2.user_tf_acc_test_1", "enabled", "true"),
				),
			},
			{
				Config: testAccVPCV2UserBasic(userNameUpdated, userPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "name", userNameUpdated),
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "password", userPassword),
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "enabled", "true"),
				),
			},
			{
				Config: testAccVPCV2UserBasic(userNameUpdated, userPasswordUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "name", userNameUpdated),
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "password", userPasswordUpdated),
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "enabled", "true"),
				),
			},
			{
				Config: testAccVPCV2UserDisabled(userNameUpdated, userPasswordUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "name", userNameUpdated),
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "password", userPasswordUpdated),
					resource.TestCheckResourceAttr(
						"selectel_vpc_user_v2.user_tf_acc_test_1", "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckVPCV2UserDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	resellV2Client := config.resellV2Client()
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "selectel_vpc_user_v2" {
			continue
		}

		_, _, err := users.Get(ctx, resellV2Client, rs.Primary.ID)

		if err == nil {
			return errors.New("user still exists")
		}
	}

	return nil
}

func testAccCheckVPCV2UserExists(n string, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		resellV2Client := config.resellV2Client()
		ctx := context.Background()

		foundUser, _, err := users.Get(ctx, resellV2Client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundUser.ID != rs.Primary.ID {
			return errors.New("user not found")
		}

		*user = *foundUser

		return nil
	}
}

func testAccVPCV2UserBasic(userName, userPassword string) string {
	return fmt.Sprintf(`
resource "selectel_vpc_user_v2" "user_tf_acc_test_1" {
  name        = "%s"
  password    = "%s"
}`, userName, userPassword)
}

func testAccVPCV2UserDisabled(userName, userPassword string) string {
	return fmt.Sprintf(`
resource "selectel_vpc_user_v2" "user_tf_acc_test_1" {
  name        = "%s"
  password    = "%s"
  enabled     = false
}`, userName, userPassword)
}
