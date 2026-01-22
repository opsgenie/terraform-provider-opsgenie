package opsgenie

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceOpsGenieContact_Basic(t *testing.T) {
	randomUserName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpsGenieContactConfig(randomUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOpsGenieContact("opsgenie_user_contact.test", "data.opsgenie_contact.existing"),
				),
			},
		},
	})
}
func testAccDataSourceOpsGenieContact(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["contact_list"] == "" {
			return fmt.Errorf("Expected to get contact_list")
		}

		testAtts := []string{"id", "method", "to"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the user %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}

		return nil
	}
}

func testAccDataSourceOpsGenieContactConfig(randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
	username  = "acctest-%s@opsgenie.com"
	full_name = "Acceptance Test User"
	role      = "User"
}
resource "opsgenie_user_contact" "test" {
  username = "${opsgenie_user.test.username}"
  to       = "39-123"
  method   = "sms"
}

data "opsgenie_contact" "existing" {
  username = "${opsgenie_user.test.username}"
}
`, randomName)
}
