package opsgenie

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/notification"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
)

func init() {
	resource.AddTestSweepers("opsgenie_notification_rule", &resource.Sweeper{
		Name: "opsgenie_notification_rule",
		F:    testSweepNotificationRule,
	})

}

func testSweepNotificationRule(region string) error {
	meta, err := sharedConfigForRegion()
	if err != nil {
		return err
	}

	client, err := notification.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	userClient, err := user.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	respUser, err := userClient.List(context.Background(), &user.ListRequest{})
	if err != nil {
		return err
	}
	for _, u := range respUser.Users {
		if strings.HasPrefix(u.Username, "genietest-") {
			resp, err := client.ListRule(context.Background(), &notification.ListRuleRequest{
				UserIdentifier: u.Id,
			})
			if err != nil {
				return err
			}
			for _, r := range resp.SimpleNotificationRules {
				deleteRequest := &notification.DeleteRuleRequest{
					UserIdentifier: u.Id,
					RuleId:         r.Id,
				}
				if _, err := client.DeleteRule(context.Background(), deleteRequest); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func TestAccOpsGenieNotificationRule_basic(t *testing.T) {
	randomName := acctest.RandString(6)
	config := testAccOpsGenieNotificationRule_basic(randomName)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieNotificationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieNotificationRuleExists("opsgenie_notification_rule.test", randomName),
				),
			},
		},
	})
}

func testCheckOpsGenieNotificationRuleDestroy(s *terraform.State) error {
	client, err := notification.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_notification_rule" {
			continue
		}
		req := notification.GetRuleRequest{
			UserIdentifier: rs.Primary.Attributes["username"],
			RuleId:         rs.Primary.Attributes["id"],
		}
		_, err := client.GetRule(context.Background(), &req)
		if err != nil {
			x := err.(*ogClient.ApiError)
			if x.StatusCode != 404 {
				return errors.New(fmt.Sprintf("Notification rule still exists: %s", x.Error()))
			}
		}
	}

	return nil
}

func testCheckOpsGenieNotificationRuleExists(name, username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client, err := notification.NewClient(testAccProvider.Meta().(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
		req := notification.GetRuleRequest{
			RuleId:         rs.Primary.Attributes["id"],
			UserIdentifier: rs.Primary.Attributes["username"],
		}

		_, err = client.GetRule(context.Background(), &req)
		if err != nil {
			return fmt.Errorf("Notification rule does not exist (and it should).")
		} else {
			log.Printf("Notification rule found.")
		}

		return nil
	}
}

func testAccOpsGenieNotificationRule_basic(randomName string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "genieuser-%s@opsgenie.com"
  full_name = "Acceptance Test User"
  role      = "User"
}

resource "opsgenie_notification_rule" "test" {
  name = "genierule-%s"
  username = opsgenie_user.test.username
  action_type = "schedule-end"
  notification_time = ["just-before", "15-minutes-ago"]
  steps {
    contact {
      method = "email"
      to = "genieuser-%s@opsgenie.com"
    }
  }
  order = 0
  enabled = true
  repeat {
    loop_after = 2
  }
  timeRestriction {
    type = "time-of-day"
    restriction {
      startHour = 3
      startMin = 15
      endHour = 5
      endMin = 30
    }
  }
}

`, randomName, randomName, randomName)
}
