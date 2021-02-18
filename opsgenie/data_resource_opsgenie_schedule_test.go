package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceOpsGenieSchedule_Basic(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomSchedule := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieScheduleConfig(randomTeam, randomSchedule),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOpsGenieSchedule("opsgenie_schedule.test", "data.opsgenie_schedule.existingschedule"),
				),
			},
		},
	})
}

func testAccDataSourceOpsGenieSchedule(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["name"] == "" {
			return fmt.Errorf("expected to get an schedule name from OpsGenie")
		}

		testAtts := []string{"name", "owner_team_id", "description", "enabled", "timezone"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("expected the schedule %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieScheduleConfig(randomTeam, randomSchedule string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_schedule" "test" {
  name = "genieschedule-%s"
  description = "schedule test"
  timezone = "Europe/Rome"
  owner_team_id = "${opsgenie_team.test.id}"
}
data "opsgenie_schedule" "existingschedule" {
  name = opsgenie_schedule.test.name
  depends_on = [opsgenie_schedule.test]
}
`, randomTeam, randomSchedule)
}
