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
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
)

func init() {
	resource.AddTestSweepers("opsgenie_schedule", &resource.Sweeper{
		Name: "opsgenie_schedule",
		F:    testSweepSchedule,
	})

}

func testSweepSchedule(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	expand := false
	resp, err := client.List(context.Background(), &schedule.ListRequest{Expand: &expand})
	if err != nil {
		return err
	}

	for _, u := range resp.Schedule {
		if strings.HasPrefix(u.Name, "genieschedule-") {
			log.Printf("Destroying schedule %s", u.Name)

			deleteRequest := schedule.DeleteRequest{
				IdentifierType:  schedule.Id,
				IdentifierValue: u.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieSchedule_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieSchedule_basic(rs)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieScheduleRotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieScheduleExists("opsgenie_schedule.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieSchedule_complete(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieSchedule_complete(rs)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieScheduleExists("opsgenie_schedule.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieScheduleDestroy(s *terraform.State) error {
	client, err := schedule.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_schedule" {
			continue
		}

		req := schedule.GetRequest{
			IdentifierType:  schedule.Id,
			IdentifierValue: rs.Primary.Attributes["id"],
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Schedule still exists : %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieScheduleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]

		client, err := schedule.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := schedule.GetRequest{
			IdentifierType:  schedule.Id,
			IdentifierValue: id,
		}

		_, err = client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Schedule with id %q does not exist", id)
		}
		return nil
	}
}

func testAccOpsGenieSchedule_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
}
`, rString)
}

func testAccOpsGenieSchedule_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "geniescheduleteam"
  description = "This team deals with all the things"
}
resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
  owner_team_id = "${opsgenie_team.test.id}"
}

`, rString)
}
