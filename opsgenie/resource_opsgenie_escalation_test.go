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
	"github.com/opsgenie/opsgenie-go-sdk-v2/escalation"
)

func init() {
	resource.AddTestSweepers("opsgenie_escalation", &resource.Sweeper{
		Name: "opsgenie_escalation",
		F:    testSweepEscalation,
	})

}

func testSweepEscalation(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := escalation.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background())
	if err != nil {
		return err
	}

	for _, u := range resp.Escalations {
		if strings.HasPrefix(u.Description, "genieescalation-") {
			log.Printf("Destroying escalation %s", u.Description)

			deleteRequest := escalation.DeleteRequest{
				IdentifierType: escalation.Id,
				Identifier:     u.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieEscalation_basic(t *testing.T) {
	randomName := acctest.RandString(6)
	randomEscalation := acctest.RandString(6)

	config := testAccOpsGenieEscalation_basic(randomName, randomEscalation)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieEscalationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieEscalationExists("opsgenie_escalation.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieEscalation_complete(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomSchedule := acctest.RandString(6)
	randomEscalation := acctest.RandString(6)

	config := testAccOpsGenieEscalation_complete(randomTeam, randomSchedule, randomEscalation)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieEscalationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieEscalationExists("opsgenie_escalation.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieEscalationDestroy(s *terraform.State) error {
	client, err := escalation.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_escalation" {
			continue
		}
		Id := rs.Primary.Attributes["id"]

		req := escalation.GetRequest{
			IdentifierType: escalation.Id,
			Identifier:     Id,
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Escalation still exists : %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieEscalationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client, err := escalation.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		Id := rs.Primary.Attributes["id"]

		req := escalation.GetRequest{
			IdentifierType: escalation.Id,
			Identifier:     Id,
		}

		_, err = client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Escalation with id %q does not exist", Id)
		}
		return nil
	}
}

func testAccOpsGenieEscalation_basic(randomName, randomEscalation string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
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
`, randomName, randomEscalation)
}

func testAccOpsGenieEscalation_complete(randomTeam, randomSchedule, randomEscalation string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
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
 description="test"
 rules {
  condition =   "if-not-acked"
    notify_type  =   "default"
	recipient {
      type  = "schedule"
      id  = "${opsgenie_schedule.test.id}"
    }   
    delay = 1
 }
owner_team_id = "${opsgenie_team.test.id}"
repeat  {
  wait_interval = 10
  count = 20
  reset_recipient_states = true
  close_alert_after_all = false
  }
}
`, randomTeam, randomSchedule, randomEscalation)
}
