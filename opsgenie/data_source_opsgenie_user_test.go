package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceOpsGenieUser_Basic(t *testing.T) {
	randomName := acctest.RandString(6)
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieUserConfig(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOpsGenieUser("opsgenie_user.test", "data.opsgenie_user.by_username"),
				),
			},
		},
	})
}

func testAccDataSourceOpsGenieUser(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a user ID from OpsGenie")
		}

		testAtts := []string{"username", "full_name", "role", "timezone", "locale"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the user %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieUserConfig(randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "acctest-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}

data "opsgenie_user" "by_username" {
  username = "${opsgenie_user.test.username}"
  depends_on = [opsgenie_user.test]
}
`, randomName)
}
