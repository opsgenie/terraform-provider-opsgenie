package opsgenie

import (
	"context"
	"log"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	if !regexp.MustCompile(`^[a-zA-Z 0-9_-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alpha numeric characters and underscores are allowed in %q: %q", k, value))
	}

	if len(value) >= 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
}
