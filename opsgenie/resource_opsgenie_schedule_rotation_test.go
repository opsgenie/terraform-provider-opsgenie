package opsgenie

import (
	"context"
	"errors"
	"fmt"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("opsgenie_schedule_rotation", &resource.Sweeper{
		Name: "opsgenie_schedule_rotation",
		F:    testSweepScheduleRotations,
	})

}

func testSweepScheduleRotations(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.ListRotations(context.Background(), &schedule.ListRotationsRequest{})
	if err != nil {
		return err
	}

	for _, u := range resp.Rotations {
		if strings.HasPrefix(u.Name, "genierotation-") {
			log.Printf("Destroying schedule rotation %s", u.Name)

			deleteRequest := schedule.DeleteRotationRequest{
				RotationId:              u.Id,
			}

			if _, err := client.DeleteRotation(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieScheduleRotation_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieScheduleRotation_basic(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieScheduleRotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieScheduleRotationExists("opsgenie_schedule_rotation.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieScheduleRotation_complete(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieScheduleRotation_complete(rs)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieScheduleRotationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieScheduleRotationExists("opsgenie_schedule_rotation.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieScheduleRotationDestroy(s *terraform.State) error {
	client, err := schedule.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_schedule_rotation" {
			continue
		}

		req := schedule.GetRotationRequest{
			ScheduleIdentifierType:  schedule.Id,
			ScheduleIdentifierValue: rs.Primary.Attributes["schedule_id"],
			RotationId:              rs.Primary.Attributes["id"],
		}
		_, err := client.GetRotation(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Schedule rotation still exists : %s", x.Error()))
			}
		}

	}

	return nil
}

func testCheckOpsGenieScheduleRotationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]
		scheduleId := rs.Primary.Attributes["schedule_id"]

		client, err := schedule.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := schedule.GetRotationRequest{
			ScheduleIdentifierType:  schedule.Id,
			ScheduleIdentifierValue: scheduleId,
			RotationId:              id,
		}

		_, err = client.GetRotation(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: ScheduleRotation with id %q (scheduleId: %q) does not exist", id, scheduleId)
		}
		return nil
	}
}

func testAccOpsGenieScheduleRotation_basic(rString string) string {
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

resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
  owner_team_id = "${opsgenie_team.test.id}"
}

resource "opsgenie_schedule_rotation" "test" { 
    schedule_id = "${opsgenie_schedule.test.id}"
    name = "test"
    start_date = "2019-06-18T17:45:00Z"
    end_date ="2019-06-20T17:45:00Z"
    type ="hourly"
    length = 6
    participant {
      type = "user"
      id = "${opsgenie_user.test.id}"
    }

    time_restriction {
      type ="time-of-day"
      restriction {
        start_hour = 1
        start_min = 1
        end_hour = 10
        end_min = 1
      }
}
}
`, rString)
}

func testAccOpsGenieScheduleRotation_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest232@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}

resource "opsgenie_team" "test" {
  name        = "genieteam222"
  description = "This team deals with all the things"
}

resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
  owner_team_id = "${opsgenie_team.test.id}"
}

resource "opsgenie_schedule_rotation" "test" { 
    schedule_id = "${opsgenie_schedule.test.id}"
    name = "t214est"
    start_date = "2019-06-18T17:45:00Z"
    end_date ="2019-06-20T17:45:00Z"
    type ="hourly"
    length = 6
    participant {
      type = "user"
      id = "${opsgenie_user.test.id}"
    }

    time_restriction {
      type ="time-of-day"
      restriction {
        start_hour = 1
        start_min = 1
        end_hour = 10
        end_min = 1
      }
}
}

resource "opsgenie_schedule_rotation" "test2" { 
    schedule_id = "${opsgenie_schedule.test.id}"
    name = "t12st"
    start_date = "2019-06-18T17:45:00Z"
    end_date ="2019-06-20T17:45:00Z"
    type ="hourly"
    length = 6
    participant {
      type = "user"
      id = "${opsgenie_user.test.id}"
    }

    time_restriction {
      type ="time-of-day"
      restriction {
        start_hour = 1
        start_min = 1
        end_hour = 10
        end_min = 1
      }
}
}
`, rString)
}
