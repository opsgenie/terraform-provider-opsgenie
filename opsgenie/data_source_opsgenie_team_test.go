package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceOpsGenieTeam_Basic(t *testing.T) {
	randomName := acctest.RandString(6)
	randomTeamName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieTeamConfig(randomName, randomTeamName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOpsGenieTeam("opsgenie_team.test", "data.opsgenie_team.existingteam"),
				),
			},
		},
	})
}

func testAccDataSourceOpsGenieTeam(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a team ID from OpsGenie")
		}

		testAtts := []string{"name", "description", "member"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the team %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieTeamConfig(randomName, randomTeamName string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genietest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_user" "test2" {
  username  = "genietest-2%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"

  member {
    id = "${opsgenie_user.test.id}"
    role     = "admin"
  }
  member {
    id = "${opsgenie_user.test2.id}"
    role     = "admin"
  }
}
data "opsgenie_team" "existingteam" {
  name = "${opsgenie_team.test.name}"
  depends_on = [opsgenie_team.test]
}
`, randomName, randomName, randomTeamName)
}
