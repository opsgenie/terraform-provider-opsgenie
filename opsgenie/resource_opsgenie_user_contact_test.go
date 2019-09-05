package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/contact"
)

func init() {
	resource.AddTestSweepers("opsgenie_user_contact", &resource.Sweeper{
		Name: "opsgenie_user_contact",
		F:    testSweepUserContact,
	})

}

func testSweepUserContact(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := contact.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	userClient, err := user.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	respUser, err := userClient.List(context.Background(), &user.ListRequest{})
	if err != nil {
		return err
	}
	for _, u := range respUser.Users {
		if strings.HasPrefix(u.Username, "genietest-") {
			resp, err := client.List(context.Background(), &contact.ListRequest{
				UserIdentifier: u.Id,
			})
			if err != nil {
				return err
			}
			for _, c := range resp.Contact {
				deleteRequest := contact.DeleteRequest{
					ContactIdentifier: c.Id,
				}

				if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func TestAccOpsGenieUserContact_basic(t *testing.T) {
	config := testAccOpsGenieUserContact_basic()

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieUserContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieUserContactExists("opsgenie_user_contact.contact"),
				),
			},
		},
	})
}

func testCheckOpsGenieUserContactDestroy(s *terraform.State) error {
	client, err := contact.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_user_contact" {
			continue
		}
		req := contact.GetRequest{
			UserIdentifier:    "genietest+contact@opsgenie.com",
			ContactIdentifier: rs.Primary.Attributes["id"],
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("User contact still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieUserContactExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client, err := contact.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := contact.GetRequest{
			ContactIdentifier: rs.Primary.Attributes["id"],
			UserIdentifier:    "genietest+contact@opsgenie.com",
		}

		_, err = client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: User contact does not exist")
		} else {
			log.Printf("User contact found")
		}

		return nil
	}
}

func testAccOpsGenieUserContact_basic() string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest+contact@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_user_contact" "contact" {
  username = "${opsgenie_user.test.username}"
  to="90-123"
  method="sms"
  enabled="true"
}

`)
}
