package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/custom_user_role"
	"log"
	"regexp"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("opsgenie_role", &resource.Sweeper{
		Name: "opsgenie_role",
		F:    testSweepUserRole,
	})
}

func testSweepUserRole(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := custom_user_role.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background(), &custom_user_role.ListRequest{})
	if err != nil {
		return err
	}

	for _, u := range resp.CustomUserRoles {
		if strings.HasPrefix(u.Name, "genietest-") {
			log.Printf("Destroying user %s", u.Name)

			deleteRequest := custom_user_role.DeleteRequest{
				Identifier: u.Id,
			}

			if _, err := client.Delete(context.Background(), &deleteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieUserRole_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUserRole_basic(rs)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieUserRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieUserRoleExists("opsgenie_custom_role.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieUserRole_complete(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUserRole_complete(rs)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieUserRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieUserRoleExists("opsgenie_custom_role.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieUserRoleDestroy(s *terraform.State) error {
	client, err := custom_user_role.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_custom_role" {
			continue
		}
		req := custom_user_role.GetRequest{
			Identifier:     rs.Primary.Attributes["role_name"],
			IdentifierType: custom_user_role.Name,
		}
		_, err := client.Get(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("User role still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieUserRoleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]
		userRoleName := rs.Primary.Attributes["role_name"]

		client, err := custom_user_role.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := custom_user_role.GetRequest{
			Identifier:     userRoleName,
			IdentifierType: custom_user_role.Name,
		}

		result, err := client.Get(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: userrole %q (userRoleName: %q) does not exist", id, userRoleName)
		} else {
			log.Printf("User role found :%s ", result.Name)
		}

		return nil
	}
}

func TestAccOpsGenieUserRole_extendedRoleValidationError(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUserRole_ExtendedRoleValidationError(rs)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Error: expected extended_role to be one of \[user observer stakeholder\], got invalid-role`)),
			},
		},
	})
}

func TestAccOpsGenieUserRole_grantedRightsValidationError(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUserRole_grantedRightsValidationError(rs)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(fmt.Sprintf(`config is invalid: expected granted_rights.0 to be one of \[who-is-on-call-show-all notification-rules-edit quiet-hours-edit alerts-access-all reports-access logs-page-access maintenance-edit contacts-edit profile-edit login-email-edit profile-custom-fields-edit configurations-read-only configurations-read-only configurations-edit configurations-delete billing-manage alert-action alert-create alert-add-attachment alert-delete-attachment alert-add-note alert-acknowledge alert-unacknowledge alert-snooze alert-escalate alert-close alert-delete alert-take-ownership alert-assign-ownership alert-add-recipient alert-add-team alert-edit-tags alert-edit-details alert-custom-action alert-update-priority alert-acknowledge-all alert-close-all incident-create incident-add-stakeholder incident-add-responder incident-resolve incident-reopen mass-notification-create service-access\], got invalid-right`)),
			},
		},
	})
}

func TestAccOpsGenieUserRole_disallowedRightsValidationError(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUserRole_disallowedRightsValidationError(rs)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(fmt.Sprintf(`config is invalid: expected disallowed_rights.0 to be one of \[who-is-on-call-show-all notification-rules-edit quiet-hours-edit alerts-access-all reports-access logs-page-access maintenance-edit contacts-edit profile-edit login-email-edit profile-custom-fields-edit configurations-read-only configurations-read-only configurations-edit configurations-delete billing-manage alert-action alert-create alert-add-attachment alert-delete-attachment alert-add-note alert-acknowledge alert-unacknowledge alert-snooze alert-escalate alert-close alert-delete alert-take-ownership alert-assign-ownership alert-add-recipient alert-add-team alert-edit-tags alert-edit-details alert-custom-action alert-update-priority alert-acknowledge-all alert-close-all incident-create incident-add-stakeholder incident-add-responder incident-resolve incident-reopen mass-notification-create service-access\], got invalid-right`)),
			},
		},
	})
}

func testAccOpsGenieUserRole_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_custom_role" "test" {
  role_name  = "opsgenie-%s"
  extended_role = "user"
}
`, rString)
}

func testAccOpsGenieUserRole_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_custom_role" "test" {
  role_name  = "genietest-%s"
  extended_role = "user"
  granted_rights = ["alert-delete"]
  disallowed_rights = ["profile-edit", "contacts-edit"]
}
`, rString)
}

func testAccOpsGenieUserRole_ExtendedRoleValidationError(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_custom_role" "test" {
  role_name  = "genietest-%s"
  extended_role = "invalid-role"
}
`, rString)
}

func testAccOpsGenieUserRole_grantedRightsValidationError(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_custom_role" "test" {
  role_name  = "genietest-%s"
  extended_role = "user"
  granted_rights = ["invalid-right"]
  disallowed_rights = ["profile-edit", "contacts-edit"]
}
`, rString)
}

func testAccOpsGenieUserRole_disallowedRightsValidationError(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_custom_role" "test" {
  role_name  = "genietest-%s"
  extended_role = "user"
  disallowed_rights = ["invalid-right"]
}
`, rString)
}
