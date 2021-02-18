package opsgenie

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieApiIntegrationDestroy,
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

func TestAccOpsGenieApiIntegration_limits(t *testing.T) {
	randomLongName := acctest.RandString(245)
	// include a backtick here as it's not possible to escape it in the multiline string
	randomName := "`" + acctest.RandString(6)
	config := testAccOpsGenieApiIntegration_limits(randomLongName, randomName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieApiIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieApiIntegrationExists("opsgenie_api_integration.test_length"),
					testCheckOpsGenieApiIntegrationExists("opsgenie_api_integration.test_format"),
				),
			},
		},
	})
}

func TestAccOpsGenieApiIntegration_complete(t *testing.T) {
	randomUsername := acctest.RandString(6)
	randomTeam := acctest.RandString(6)
	randomTeam2 := acctest.RandString(6)
	randomSchedule := acctest.RandString(6)
	randomEscalation := acctest.RandString(6)
	randomIntegration := acctest.RandString(6)
	randomIntegration2 := acctest.RandString(6)
	randomIntegration3 := acctest.RandString(6)

	config := testAccOpsGenieApiIntegration_complete(randomUsername, randomTeam, randomTeam2, randomSchedule, randomEscalation, randomIntegration, randomIntegration2, randomIntegration3)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieApiIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieApiIntegrationExists("opsgenie_api_integration.test"),
					testCheckOpsGenieApiIntegrationExists("opsgenie_api_integration.test3"),
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
  type = "API"
  name = "genieintegration-%s"
}
`, rString)
}

func testAccOpsGenieApiIntegration_limits(randomLongName, randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_api_integration" "test_length" {
  type = "API"
  name = "test-%s"
}

resource "opsgenie_api_integration" "test_format" {
  type = "API"
  name = "[] () {} 🚒 %s"
}
`, randomLongName, randomName)
}

func testAccOpsGenieApiIntegration_complete(randomUsername, randomTeam, randomTeam2, randomSchedule, randomEscalation, randomIntegration, randomIntegration2, randomIntegration3 string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_team" "test2" {
  name        = "genieteam2-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
}

resource "opsgenie_escalation" "test" {
 name ="genieescalation-%s"
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
  type = "API"

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
resource "opsgenie_api_integration" "test2" {
	name          = "genieintegration-prometheus-%s"
	type          = "Prometheus"
	owner_team_id = "${opsgenie_team.test.id}"
	enabled       = true
}
resource "opsgenie_api_integration" "test3" {
	name = "genieintegration-webhook-%s"
  	owner_team_id = "${opsgenie_team.test.id}"
  	type = "Webhook"
  	enabled                        = true
  	allow_write_access             = false
  	suppress_notifications         = false
  	webhook_url                    = "https://example.com/v1"
  	headers = {
		header = "value1"
	}
}
`, randomUsername, randomTeam, randomTeam2, randomSchedule, randomEscalation, randomIntegration, randomIntegration2, randomIntegration3)
}
