package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					UserIdentifier:    u.Id,
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
	randomName := acctest.RandString(6)
	config := testAccOpsGenieUserContact_basic(randomName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieUserContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieUserContactExists("opsgenie_user_contact.contact", randomName),
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
			UserIdentifier:    rs.Primary.Attributes["username"],
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

func testCheckOpsGenieUserContactExists(name, username string) resource.TestCheckFunc {
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
			UserIdentifier:    fmt.Sprintf("genietest+contact-%s@opsgenie.com", username),
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

func testAccOpsGenieUserContact_basic(randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest+contact-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_user_contact" "contact" {
  username = "${opsgenie_user.test.username}"
  to="90-123"
  method="sms"
  enabled="true"
}

`, randomName)
}
