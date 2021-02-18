package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceOpsGenieHeartbeat_Basic(t *testing.T) {
	randomName := acctest.RandString(6)
	randomTeamName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieHeartbeatConfig(randomTeamName, randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOpsGenieHeartbeat("opsgenie_heartbeat.test", "data.opsgenie_heartbeat.existingheartbeat"),
				),
			},
		},
	})
}

func testAccDataSourceOpsGenieHeartbeat(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["name"] == "" {
			return fmt.Errorf("expected to get a heartbeat name from OpsGenie")
		}

		testAtts := []string{"name", "description", "interval_unit", "interval", "enabled", "alert_message", "alert_priority", "alert_tags", "owner_team_id"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("expected the heartbeat %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieHeartbeatConfig(randomTeamName, randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}
resource "opsgenie_heartbeat" "test" {
	name = "genieheartbeat-%s"
	description = "test opsgenie heartbeat terraform"
	interval_unit = "minutes"
	interval = 10
	enabled = false
	alert_message = "Test"
	alert_priority = "P3"
	alert_tags = ["test","fahri"]
	owner_team_id = "${opsgenie_team.test.id}"

}
data "opsgenie_heartbeat" "existingheartbeat" {
  name = "${opsgenie_heartbeat.test.name}"
  depends_on = [opsgenie_heartbeat.test]

}
`, randomTeamName, randomName)
}
