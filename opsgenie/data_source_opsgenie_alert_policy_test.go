package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceOpsGenieAlertPolicy_Basic(t *testing.T) {
	randomTeamName := acctest.RandString(6)
	randomAlertPolicyName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieAlertPolicyConfig(randomTeamName, randomAlertPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOpsGenieAlertPolicy("opsgenie_alert_policy.test", "data.opsgenie_alert_policy.existingAlertPolicy"),
				),
			},
		},
	})
}

func testAccDataSourceOpsGenieAlertPolicy(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a AlertPolicy ID from OpsGenie")
		}

		testAtts := []string{"name", "team_id"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the AlertPolicy %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieAlertPolicyConfig(randomTeamName, randomAlertPolicyName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_alert_policy" "test" {
  name               = "geniepolicy-%s"
  policy_description = "Perfect Alert policy for the team."
  message = "This is a test message"
  filter {}
  time_restriction {
    type = "weekday-and-time-of-day"
    restrictions {
      end_day    = "monday"
      end_hour   = 7
      end_min    = 0
      start_day  = "sunday"
      start_hour = 21
      start_min  = 0
    }
    restrictions {
      end_day    = "tuesday"
      end_hour   = 7
      end_min    = 0
      start_day  = "monday"
      start_hour = 22
      start_min  = 0
    }
  }
}
data "opsgenie_alert_policy" "existingAlertPolicy" {
  name = "${opsgenie_alert_policy.test.name}"
}
`, randomTeamName, randomAlertPolicyName)
}
