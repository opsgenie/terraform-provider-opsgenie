package opsgenie

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"log"
	"strings"
)

func resourceOpsGenieTeamMembership() *schema.Resource { //TODO encode the e-mail addrs (e.g. because of +)? https://github.com/opsgenie/opsgenie-go-sdk-v2/issues/62
	return &schema.Resource{
		Create: resourceOpsGenieTeamMembershipCreate,
		Read:   handleNonExistentResource(resourceOpsGenieTeamMembershipRead),
		//Update: resourceOpsGenieTeamMembershipUpdate, // requires https://github.com/opsgenie/opsgenie-go-sdk-v2/issues/59
		Delete: resourceOpsGenieTeamMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "user",
				ForceNew: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOpsGenieTeamMembershipCreate(d *schema.ResourceData, meta interface{}) error {

	userID := d.Get("user_id").(string)
	role := d.Get("role").(string)
	teamID := d.Get("team_id").(string)

	log.Printf("[INFO] Adding user %q to team %q", teamID, userID)

	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	// add member to team
	_, err = client.AddMember(context.Background(), &team.AddTeamMemberRequest{
		TeamIdentifierType:  team.Id,
		TeamIdentifierValue: teamID,
		User: team.User{
			ID: userID,
		},
		Role: role, //TODO lowercase?
	})
	if err != nil {
		return err
	}

	d.SetId(buildTwoPartID(teamID, userID))

	return resourceOpsGenieTeamMembershipRead(d, meta)
}

func resourceOpsGenieTeamMembershipRead(d *schema.ResourceData, meta interface{}) error {

	teamID, userID, err := parseTwoPartID(d.Id(), "teamID", "userID")
	if err != nil {
		return err
	}

	getRequest := &team.GetTeamRequest{
		IdentifierType:  team.Id,
		IdentifierValue: teamID,
	}

	log.Printf("[INFO] Retrieving membership of user %q in team %q", userID, teamID)

	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		return err
	}

	role, err := getUserRole(userID, teamID, getResponse.Members)
	if err != nil {
		return err
	}

	d.Set("user_id", userID)
	d.Set("role", role)
	d.Set("team_id", teamID)

	return nil
}

func resourceOpsGenieTeamMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	userID := d.Get("user_id").(string)
	teamID := d.Get("team_id").(string)

	log.Printf("[INFO] Deleting membership of user %q in team %q", userID, teamID)

	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	_, err = client.RemoveMember(context.Background(), &team.RemoveTeamMemberRequest{
		TeamIdentifierType:    team.Id,
		TeamIdentifierValue:   teamID,
		MemberIdentifierType:  team.Id,
		MemberIdentifierValue: userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func getUserRole(userID string, teamID string, input []team.Member) (string, error) {
	role := ""

	for _, inputMember := range input {
		if inputMember.User.ID == userID {
			role = inputMember.Role
			break
		}
	}

	if len(role) == 0 {
		return "", fmt.Errorf("did not found user %q in team %q", userID, teamID)
	}

	return role, nil
}

// format the strings into an id `a:b`
func buildTwoPartID(a, b string) string {
	return fmt.Sprintf("%s:%s", a, b)
}

// return the pieces of id `left:right` as left, right
func parseTwoPartID(id, left, right string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected ID format %q, expected %s:%s", id, left, right)
	}

	return parts[0], parts[1], nil
}
