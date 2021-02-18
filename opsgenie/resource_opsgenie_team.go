package opsgenie

import (
	"context"
	"errors"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/escalation"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"log"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func resourceOpsGenieTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieTeamCreate,
		Read:   handleNonExistentResource(resourceOpsGenieTeamRead),
		Update: resourceOpsGenieTeamUpdate,
		Delete: resourceOpsGenieTeamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieTeamName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ignore_members": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_default_resources": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"member": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"role": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "user",
						},
					},
				},
			},
		},
	}
}

func resourceOpsGenieTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	createRequest := &team.CreateTeamRequest{
		Name:        name,
		Description: description,
	}

	if len(d.Get("member").([]interface{})) > 0 && !d.Get("ignore_members").(bool) {
		createRequest.Members = expandOpsGenieTeamMembers(d)
	}

	log.Printf("[INFO] Creating OpsGenie team %q", name)

	_, err = client.Create(context.Background(), createRequest)
	if err != nil {
		return err
	}

	getRequest := &team.GetTeamRequest{
		IdentifierType:  team.Name,
		IdentifierValue: name,
	}

	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.SetId(getResponse.Id)

	shouldDeleteDefaultResources := d.Get("delete_default_resources").(bool)

	if shouldDeleteDefaultResources {
		err = findAndUpdateDefaultRoutingRule(name, meta.(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}

		err := findAndDeleteDefaultEscalation(name, meta.(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}

		err = findAndDeleteDefaultSchedule(name, meta.(*OpsgenieClient).client.Config)
		if err != nil {
			return err
		}
	}
	return resourceOpsGenieTeamRead(d, meta)
}

func resourceOpsGenieTeamRead(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	getRequest := &team.GetTeamRequest{
		IdentifierType:  team.Id,
		IdentifierValue: d.Id(),
	}

	log.Printf("[INFO] Retrieving state of OpsGenie team '%s'", d.Get("name"))

	getResponse, err := client.Get(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.Set("name", getResponse.Name)
	d.Set("description", getResponse.Description)

	if !d.Get("ignore_members").(bool) {
		d.Set("member", flattenOpsGenieTeamMembers(getResponse.Members))
	}

	return nil
}

func resourceOpsGenieTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	updateRequest := &team.UpdateTeamRequest{
		Id:          d.Id(),
		Name:        name,
		Description: description,
	}

	if len(d.Get("member").([]interface{})) > 0 && !d.Get("ignore_members").(bool) {
		updateRequest.Members = expandOpsGenieTeamMembers(d)
	}

	log.Printf("[INFO] Updating OpsGenie team '%s'", name)

	_, err = client.Update(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieTeamDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie team '%s'", d.Get("name").(string))
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &team.DeleteTeamRequest{
		IdentifierType:  team.Id,
		IdentifierValue: d.Id(),
	}

	_, err = client.Delete(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func flattenOpsGenieTeamMembers(input []team.Member) []map[string]interface{} {
	members := make([]map[string]interface{}, 0, len(input))
	for _, inputMember := range input {
		outputMember := make(map[string]interface{})
		outputMember["id"] = inputMember.User.ID
		outputMember["role"] = inputMember.Role
		members = append(members, outputMember)
	}

	return members
}

func expandOpsGenieTeamMembers(d *schema.ResourceData) []team.Member {
	input := d.Get("member").([]interface{})
	members := make([]team.Member, 0, len(input))

	if input == nil {
		return members
	}

	for _, v := range input {
		config := v.(map[string]interface{})

		userId := config["id"].(string)
		role := config["role"].(string)

		member := team.Member{
			User: team.User{
				ID: userId,
			},
			Role: role,
		}

		members = append(members, member)
	}

	return members
}

func validateOpsGenieTeamName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z 0-9_.-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alpha numeric characters, dots and underscores are allowed in %q: %q", k, value))
	}

	if len(value) >= 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
}

func findAndDeleteDefaultSchedule(teamName string, config *client.Config) error {
	scheduleClient, err := schedule.NewClient(config)
	if err != nil {
		return err
	}
	expand := true
	res, err := scheduleClient.List(context.Background(), &schedule.ListRequest{
		Expand: &expand,
	})
	if err != nil {
		return err
	}
	for _, sched := range res.Schedule {
		ownerTeam := sched.OwnerTeam
		if ownerTeam != nil {
			if ownerTeam.Name == teamName {
				_, err = scheduleClient.Delete(context.Background(), &schedule.DeleteRequest{
					IdentifierType:  schedule.Id,
					IdentifierValue: sched.Id,
				})
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	return errors.New("Could not find any schedule name for this team")
}

func findAndDeleteDefaultEscalation(teamName string, config *client.Config) error {
	escalationClient, err := escalation.NewClient(config)
	if err != nil {
		return err
	}
	res, err := escalationClient.List(context.Background())
	if err != nil {
		return err
	}
	for _, escal := range res.Escalations {
		ownerTeam := escal.OwnerTeam
		if ownerTeam != nil {
			if ownerTeam.Name == teamName {
				_, err = escalationClient.Delete(context.Background(), &escalation.DeleteRequest{
					IdentifierType: escalation.Id,
					Identifier:     escal.Id,
				})
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	return errors.New("Could not find any escalation for this team")
}

func findAndUpdateDefaultRoutingRule(teamName string, config *client.Config) error {
	teamClient, err := team.NewClient(config)
	if err != nil {
		return err
	}
	rules, err := teamClient.ListRoutingRules(context.Background(), &team.ListRoutingRulesRequest{
		TeamIdentifierType:  team.Name,
		TeamIdentifierValue: teamName,
	})
	if err != nil {
		return err
	}

	for _, rule := range rules.RoutingRules {
		_, err := teamClient.UpdateRoutingRule(context.Background(), &team.UpdateRoutingRuleRequest{
			TeamIdentifierType:  team.Name,
			TeamIdentifierValue: teamName,
			RoutingRuleId:       rule.Id,
			Notify: &team.Notify{
				Type: team.None,
			},
		})

		if err != nil {
			return err
		}
	}
	return nil
}
