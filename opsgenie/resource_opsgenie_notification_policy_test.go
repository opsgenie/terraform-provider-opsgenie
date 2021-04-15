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
	"github.com/opsgenie/opsgenie-go-sdk-v2/policy"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func init() {
	resource.AddTestSweepers("opsgenie_notification_policy", &resource.Sweeper{
		Name: "opsgenie_notification_policy",
		F:    testSweepNotificationPolicy,
	})

}

func testSweepNotificationPolicy(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background(), &team.ListTeamRequest{})
	if err != nil {
		return err
	}

	for _, u := range resp.Teams {
		if strings.HasPrefix(u.Name, "genieteam") {
			log.Printf("Destroying notification policy for team %s", u.Name)
			client2, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
			if err != nil {
				return err
			}

			resp2, err := client2.ListNotificationPolicies(context.Background(), &policy.ListNotificationPoliciesRequest{
				TeamId: u.Id,
			})
			if err != nil {
				return err
			}
			for _, k := range resp2.Policies {

				deleteRequest := policy.DeletePolicyRequest{
					Type:   policy.PolicyType(k.Type),
					TeamId: u.Id,
					Id:     k.Id,
				}

				if _, err := client2.DeletePolicy(context.Background(), &deleteRequest); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func TestAccOpsGenieNotificationPolicy_basic(t *testing.T) {
	teamName := acctest.RandString(6)
	notificationPolicyName := acctest.RandString(6)

	config := testAccOpsGenieNotificationPolicy_basic(teamName, notificationPolicyName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieNotificationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieNotificationPolicyExists("opsgenie_notification_policy.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieDeDuplicationNotificationPolicy_basic(t *testing.T) {
	teamName := acctest.RandString(6)
	notificationPolicyName := acctest.RandString(6)

	config := testAccOpsGenieDeDuplicationActionNotificationPolicy_basic(teamName, notificationPolicyName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieNotificationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieNotificationPolicyExists("opsgenie_notification_policy.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieNotificationPolicyDestroy(s *terraform.State) error {
	client, err := policy.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_notification_policy" {
			continue
		}
		req := policy.GetNotificationPolicyRequest{
			TeamId: rs.Primary.Attributes["team_id"],
			Id:     rs.Primary.Attributes["id"],
		}
		_, err := client.GetNotificationPolicy(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Notification policy still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieNotificationPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]

		client, err := policy.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := policy.GetNotificationPolicyRequest{
			TeamId: rs.Primary.Attributes["team_id"],
			Id:     rs.Primary.Attributes["id"],
		}

		_, err = client.GetNotificationPolicy(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Notification policy %q does not exist", id)
		}
		return nil
	}
}

func testAccOpsGenieNotificationPolicy_basic(teamName, notificationPolicyName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}

resource "opsgenie_notification_policy" "test" {
  name               = "geniepolicy-%s"
  team_id            = opsgenie_team.test.id
  policy_description = "Perfect notification policy for the team."
  delay_action {
    delay_option = "next-time"
    until_minute = 30
    until_hour   = 7
  }
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
`, teamName, notificationPolicyName)
}

func testAccOpsGenieDeDuplicationActionNotificationPolicy_basic(teamName, notificationPolicyName string) string {
	return fmt.Sprintf(`
resource "opsgenie_team" "test" {
  name        = "genieteam-%s"
  description = "This team deals with all the things"
}

resource "opsgenie_notification_policy" "test" {
  name               = "geniepolicy-%s"
  team_id            = opsgenie_team.test.id
  policy_description = "Perfect notification policy for the team."
  de_duplication_action {
		count                      = 20
		de_duplication_action_type = "value-based"
    }
	filter {}
}
`, teamName, notificationPolicyName)
}
