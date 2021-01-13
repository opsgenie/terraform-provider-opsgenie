package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	expand := false
	scheduleResp, err := client.List(context.Background(), &schedule.ListRequest{Expand: &expand})
	if err != nil {
		return err
	}
	for _, u := range scheduleResp.Schedule {
		if strings.HasPrefix(u.Name, "genieschedule-") {

			resp, err := client.ListRotations(context.Background(), &schedule.ListRotationsRequest{
				ScheduleIdentifierType:  schedule.Name,
				ScheduleIdentifierValue: "acceptance-test",
			})
			if err != nil {
				return err
			}
			for _, u := range resp.Rotations {
				if strings.HasPrefix(u.Name, "genierotation-") {
					log.Printf("Destroying schedule rotation %s", u.Name)

					deleteRequest := schedule.DeleteRotationRequest{
						RotationId: u.Id,
					}

					if _, err := client.DeleteRotation(context.Background(), &deleteRequest); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func TestAccOpsGenieScheduleRotation_basic(t *testing.T) {
	randomUser := acctest.RandString(6)
	randomTeam := acctest.RandString(6)
	randomSchedule := acctest.RandString(6)
	randomRotation := acctest.RandString(6)
	config := testAccOpsGenieScheduleRotation_basic(randomUser, randomTeam, randomSchedule, randomRotation)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieScheduleRotationDestroy,
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
	randomUser := acctest.RandString(6)
	randomTeam := acctest.RandString(6)
	randomSchedule := acctest.RandString(6)
	randomRotation := acctest.RandString(6)
	randomRotation2 := acctest.RandString(6)

	config := testAccOpsGenieScheduleRotation_complete(randomUser, randomTeam, randomSchedule, randomRotation, randomRotation2)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieScheduleRotationDestroy,
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

func testAccOpsGenieScheduleRotation_basic(randomUser, randomTeam, randomSchedule, randomRotation string) string {
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

resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
  owner_team_id = "${opsgenie_team.test.id}"
}

resource "opsgenie_schedule_rotation" "test" { 
    schedule_id = "${opsgenie_schedule.test.id}"
    name = "test-%s"
    start_date = "2019-06-18T17:30:00Z"
    end_date = "2019-06-20T17:30:00Z"
    type = "hourly"
    length = 6
    participant {
      type = "user"
      id = "${opsgenie_user.test.id}"
    }

}
`, randomUser, randomTeam, randomSchedule, randomRotation)
}

func testAccOpsGenieScheduleRotation_complete(randomName, randomTeam, randomSchedule, randomRotation, randomRotation2 string) string {
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

resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  enabled = false
  owner_team_id = "${opsgenie_team.test.id}"
}

resource "opsgenie_escalation" "test" {
 name ="genieescalationsched-%s"
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

resource "opsgenie_schedule_rotation" "test" { 
    schedule_id = "${opsgenie_schedule.test.id}"
    name = "test-%s"
    start_date = "2019-06-18T17:30:00Z"
    end_date ="2019-06-20T17:30:00Z"
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
        start_min = 0
        end_hour = 10
        end_min = 0
	  }
	}
}

resource "opsgenie_schedule_rotation" "test2" { 
    schedule_id = "${opsgenie_schedule.test.id}"
    name = "schedule2-%s"
    start_date = "2019-06-18T17:30:00Z"
    end_date ="2019-06-20T17:00:00Z"
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
        start_min = 0
        end_hour = 10
        end_min = 0
      }
    }
}

resource "opsgenie_schedule_rotation" "none" { 
    schedule_id = "${opsgenie_schedule.test.id}"
    name = "schedulenone-%s"
    start_date = "2019-06-18T17:30:00Z"
    end_date ="2019-06-20T17:00:00Z"
    type ="hourly"
    length = 6
    participant {
      type = "none"
    }

    time_restriction {
      type ="time-of-day"
      restriction {
        start_hour = 1
        start_min = 0
        end_hour = 10
        end_min = 0
      }
    }
}

resource "opsgenie_schedule_rotation" "esc" { 
    schedule_id = "${opsgenie_schedule.test.id}"
    name = "schedule-escalation-%s"
    start_date = "2019-06-18T17:30:00Z"
    end_date ="2019-06-20T17:00:00Z"
    type ="hourly"
    length = 6
    participant {
      type = "escalation"
      id = "${opsgenie_escalation.test.id}"
    }

    time_restriction {
      type ="time-of-day"
      restriction {
        start_hour = 1
        start_min = 0
        end_hour = 10
        end_min = 0
      }
    }
}
`, randomName, randomTeam, randomSchedule, randomRotation, randomRotation2, randomRotation2, randomRotation2, randomRotation2)
}
