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
	resource.AddTestSweepers("opsgenie_alert_policy", &resource.Sweeper{
		Name: "opsgenie_alert_policy",
		F:    testSweepAlertPolicy,
	})

}

func testSweepAlertPolicy(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	fmt.Println("Starting testSweepAlertPolicy")
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	resp, err := client.List(context.Background(), &team.ListTeamRequest{})
	if err != nil {
		return err
	}

	fmt.Println("Checking all the teams to get our TeamID")
	for _, u := range resp.Teams {
		fmt.Printf("checking team: %s", u.Name)
		if strings.HasPrefix(u.Name, "genieteam") {
			log.Printf("Destroying alert policy for team %s", u.Name)
			client2, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
			if err != nil {
				return err
			}

			resp2, err := client2.ListAlertPolicies(context.Background(), &policy.ListAlertPoliciesRequest{
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

				fmt.Printf("Deleting policy %s", k.Name)
				if _, err := client2.DeletePolicy(context.Background(), &deleteRequest); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func TestAccOpsGenieAlertPolicy_basic(t *testing.T) {
	alertPolicyName := acctest.RandString(6)
	config := testAccOpsGenieAlertPolicy_basic(alertPolicyName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieAlertPolicyExists("opsgenie_alert_policy.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieAlertPolicy_complete(t *testing.T) {
	randomTeam := acctest.RandString(6)
	randomAlertPolicyName := acctest.RandString(6)

	config := testAccOpsGenieAlertPolicy_complete(randomTeam, randomAlertPolicyName)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckOpsGenieAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieAlertPolicyExists("opsgenie_alert_policy.test2"),
				),
			},
		},
	})
}

func testCheckOpsGenieAlertPolicyDestroy(s *terraform.State) error {
	client, err := policy.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_alert_policy" {
			continue
		}
		req := policy.GetAlertPolicyRequest{
			Id: rs.Primary.Attributes["id"],
		}
		_, err := client.GetAlertPolicy(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Alert policy still exists : %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieAlertPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		fmt.Printf("Running testCheckOpsGenieAlertPolicyExists")
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]

		fmt.Printf("Got ID for policy: %s", name)
		client, err := policy.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := policy.GetAlertPolicyRequest{
			TeamId: rs.Primary.Attributes["team_id"],
			Id:     rs.Primary.Attributes["id"],
		}

		_, err = client.GetAlertPolicy(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Bad: Alert policy %q does not exist", id)
		}
		return nil
	}
}

func testAccOpsGenieAlertPolicy_basic(alertPolicyName string) string {
	return fmt.Sprintf(`
resource "opsgenie_alert_policy" "test" {
  name               = "genie-alert-policy-%s"
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
`, alertPolicyName)
}

func testAccOpsGenieAlertPolicy_complete(randomTeam, randomAlertPolicyName string) string {

	return fmt.Sprintf(`
	resource "opsgenie_team" "test" {
	  name        = "genieteam-%s"
	  description = "This team deals with all the things"
	}
	resource "opsgenie_alert_policy" "test2" {
	  name               = "genie-alert-policy-%s"
	  policy_description = "Perfect Alert policy for the team."
	  message = "This is a test message"
	  team_id = opsgenie_team.test.id
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
	  continue_policy = true
	  alias = "alias"
	  entity = "test"
	  source = "new source"
	  ignore_original_actions = false
	  ignore_original_responders = false
	  ignore_original_tags = false
	  priority = "P3"
	  responders {
		type = "team"
		id = "${opsgenie_team.test.id}"
	  }
	  tags = ["test"]
	  actions = ["test_action"]
	}
	`, randomTeam, randomAlertPolicyName)

}
