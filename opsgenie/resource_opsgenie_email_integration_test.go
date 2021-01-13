package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/integration"
)

func init() {
	resource.AddTestSweepers("opsgenie_email_integration", &resource.Sweeper{
		Name: "opsgenie_email_integration",
		F:    testSweepEmailIntegration,
	})

}

func testSweepEmailIntegration(region string) error {
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

func TestAccOpsGenieEmailIntegration_basic(t *testing.T) {
	randomName := acctest.RandString(6)
	randomMail := acctest.RandString(6)

	config := testAccOpsGenieEmailIntegration_basic(randomName, randomMail)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieEmailIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieEmailIntegrationExists("opsgenie_email_integration.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieEmailIntegration_complete(t *testing.T) {
	randomName := acctest.RandString(6)
	randomTeam := acctest.RandString(6)
	randomTeam2 := acctest.RandString(6)
	randomSchedule := acctest.RandString(6)
	randomIntegration := acctest.RandString(6)
	randomEscalation := acctest.RandString(6)
	config := testAccOpsGenieEmailIntegration_complete(randomName, randomTeam, randomTeam2, randomSchedule, randomEscalation, randomIntegration)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieEmailIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieEmailIntegrationExists("opsgenie_email_integration.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieEmailIntegrationDestroy(s *terraform.State) error {
	client, err := integration.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_email_integration" {
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
				return errors.New(fmt.Sprintf("Email Integration still exists : %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieEmailIntegrationExists(name string) resource.TestCheckFunc {
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
			return errors.New(fmt.Sprintf("Bad: EmailIntegration with id %q does not exist", Id))
		}
		return nil
	}
}

func testAccOpsGenieEmailIntegration_basic(randomMail, randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_email_integration" "test" {
  name = "genieintegration-%s"
  email_username="fahri-%s"
}
`, randomName, randomMail)
}

func testAccOpsGenieEmailIntegration_complete(randomUser, randomTeam, randomTeam2, randomSchedule, randomEscalation, randomIntegration string) string {
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
resource "opsgenie_email_integration" "test" {
  name = "genieintegration-%s"
  responders {
    type ="escalation"
    id = "${opsgenie_escalation.test.id}"
  }
  responders {
    type ="schedule"
    id = "${opsgenie_schedule.test.id}"
  }
  responders {
    type ="team"
    id = "${opsgenie_team.test2.id}"
  }
  responders {
    type ="user"
    id = "${opsgenie_user.test.id}"
  }

  email_username="%stest"
  ignore_responders_from_payload = true
  suppress_notifications = true
}
`, randomUser, randomTeam, randomTeam2, randomSchedule, randomEscalation, randomIntegration, randomUser)
}
