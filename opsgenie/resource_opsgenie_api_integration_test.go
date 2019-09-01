package opsgenie

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/integration"
	"github.com/pkg/errors"
)

func init() {
	resource.AddTestSweepers("opsgenie_api_integration", &resource.Sweeper{
		Name: "opsgenie_api_integration",
		F:    testSweepApiIntegration,
	})

}

func testSweepApiIntegration(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := integration.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background())
	if err != nil {
		return err
	}

	for _, u := range resp.Integrations {
		if strings.HasPrefix(u.Name, "genieintegration-") {
			log.Printf("Destroying integration %s", u.Name)

			deleteRequest := integration.DeleteIntegrationRequest{
				Id: u.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieApiIntegration_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieApiIntegration_basic(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieApiIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieApiIntegrationExists("opsgenie_api_integration.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieApiIntegration_complete(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieApiIntegration_complete(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieApiIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieApiIntegrationExists("opsgenie_api_integration.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieApiIntegrationDestroy(s *terraform.State) error {
	client, err := integration.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_api_integration" {
			continue
		}
		Id := rs.Primary.Attributes["id"]

		req := integration.GetRequest{
			Id: Id,
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Api Integration still exists: %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieApiIntegrationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client, err := integration.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		Id := rs.Primary.Attributes["id"]

		req := integration.GetRequest{
			Id: Id,
		}

		_, err = client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: ApiIntegration with id %q does not exist", Id)
		}
		return nil
	}
}

func testAccOpsGenieApiIntegration_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_api_integration" "test" {
  name = "genieintegration-%s"
}
`, rString)
}

func testAccOpsGenieApiIntegration_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_team" "test" {
  name        = "genieteam"
  description = "This team deals with all the things"
}
resource "opsgenie_team" "test2" {
  name        = "genieteam2"
  description = "This team deals with all the things"
}
resource "opsgenie_schedule" "test" {
  name = "genieschedule"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
}

resource "opsgenie_escalation" "test" {
 name ="genieescalation"
 rules {
  condition =   "if-not-acked"
    notify_type  =   "default"
    recipient {
      type  = "user"
      id  = "${opsgenie_user.test.id}"
		}
    delay = 1
	}
}
resource "opsgenie_api_integration" "test" {
  name = "genieintegration-%s"
  responders {
    type ="user"
    id = "${opsgenie_user.test.id}"
  }
  responders {
    type ="schedule"
    id = "${opsgenie_schedule.test.id}"
  }
  responders {
    type ="escalation"
    id = "${opsgenie_escalation.test.id}"
  }
  responders {
    type ="team"
    id = "${opsgenie_team.test2.id}"
  }
  enabled = true
  allow_write_access = false
  ignore_responders_from_payload = true
  suppress_notifications = true
  owner_team_id = "${opsgenie_team.test.id}"
}
`, rString)
}
