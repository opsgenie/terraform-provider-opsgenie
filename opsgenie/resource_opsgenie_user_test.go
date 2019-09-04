package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
)

func init() {
	resource.AddTestSweepers("opsgenie_user", &resource.Sweeper{
		Name: "opsgenie_user",
		F:    testSweepUser,
	})

}

func testSweepUser(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := user.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background(), &user.ListRequest{})
	if err != nil {
		return err
	}

	for _, u := range resp.Users {
		if strings.HasPrefix(u.Username, "genietest-") {
			log.Printf("Destroying user %s", u.Username)

			deleteRequest := user.DeleteRequest{
				Identifier: u.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieUser_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUser_basic(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieUserExists("opsgenie_user.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieUser_complete(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUser_complete(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieUserExists("opsgenie_user.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieUserDestroy(s *terraform.State) error {
	client, err := user.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_user" {
			continue
		}
		req := user.GetRequest{
			Identifier: rs.Primary.Attributes["id"],
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("User still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieUserExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]
		username := rs.Primary.Attributes["username"]

		client, err := user.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := user.GetRequest{
			Identifier: rs.Primary.Attributes["id"],
		}

		result, err := client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: User %q (username: %q) does not exist", id, username)
		} else {
			log.Printf("User found :%s ", result.Username)
		}

		return nil
	}
}

func testAccOpsGenieUser_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
`, rString)
}

func testAccOpsGenieUser_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
  locale    = "en_GB"
  timezone = "Europe/Rome"
}
`, rString)
}
