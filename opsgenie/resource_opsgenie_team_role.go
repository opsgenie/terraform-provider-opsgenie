package opsgenie

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/team"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validTeamRolesRights = []string{
	"manage-members",
	"edit-team-roles",
	"delete-team-roles",
	"access-member-profiles",
	"edit-member-profiles",
	"edit-routing-rules",
	"delete-routing-rules",
	"edit-escalations",
	"delete-escalations",
	"edit-schedules",
	"delete-schedules",
	"edit-integrations",
	"delete-integrations",
	"edit-heartbeats",
	"delete-heartbeats",
	"access-reports",
	"edit-services",
	"delete-services",
	"edit-rooms",
	"delete-rooms",
	"send-service-status-update",
}

func resourceOpsGenieTeamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieTeamRoleCreate,
		Read:   handleNonExistentResource(resourceOpsGenieTeamRoleRead),
		Update: resourceOpsGenieTeamRoleUpdate,
		Delete: resourceOpsGenieTeamRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"rights": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"right": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(validTeamRolesRights, false),
						},
						"granted": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func mapRights(rawRights []any) []team.Right {
	rights := make([]team.Right, 0, len(rawRights))
	for _, r := range rawRights {
		m := r.(map[string]interface{})

		granted := m["granted"].(bool)

		rights = append(rights, team.Right{
			Right:   m["right"].(string),
			Granted: &granted,
		})
	}
	return rights
}

func resourceOpsGenieTeamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	name := d.Get("role_name").(string)
	createRequest := &team.CreateTeamRoleRequest{
		TeamIdentifierType:  team.Id,
		TeamIdentifierValue: d.Get("team_id").(string),
		Name:                name,
		Rights:              mapRights(d.Get("rights").(*schema.Set).List()),
	}

	log.Printf("[INFO] Creating OpsGenie team role '%s'", name)
	result, err := client.CreateRole(context.Background(), createRequest)

	if err != nil {
		return err
	}

	d.SetId(result.Id)
	return resourceOpsGenieTeamRoleRead(d, meta)
}

func resourceOpsGenieTeamRoleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	getRequest := &team.GetTeamRoleRequest{
		TeamID:   d.Get("team_id").(string),
		RoleID:   d.Id(),
		RoleName: d.Get("role_name").(string),
	}

	teamRole, err := client.GetRole(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.Set("role_name", teamRole.Name)
	d.Set("rights", teamRole.Rights)

	return nil
}

func resourceOpsGenieTeamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	updateRequest := &team.UpdateTeamRoleRequest{
		TeamID:   d.Get("team_id").(string),
		RoleID:   d.Id(),
		RoleName: d.Get("role_name").(string),
		Rights:   mapRights(d.Get("rights").(*schema.Set).List()),
	}

	_, err = client.UpdateRole(context.Background(), updateRequest)

	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieTeamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	deleteRequest := &team.DeleteTeamRoleRequest{
		TeamID:   d.Get("team_id").(string),
		RoleID:   d.Id(),
		RoleName: d.Get("role_name").(string),
	}

	_, err = client.DeleteRole(context.Background(), deleteRequest)

	if err != nil {
		return err
	}

	return nil
}
