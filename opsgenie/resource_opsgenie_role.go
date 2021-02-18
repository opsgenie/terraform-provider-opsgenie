package opsgenie

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/custom_user_role"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var validCustomRolesRights = []string{
	"who-is-on-call-show-all",
	"notification-rules-edit",
	"quiet-hours-edit",
	"alerts-access-all",
	"reports-access",
	"logs-page-access",
	"maintenance-edit",
	"contacts-edit",
	"profile-edit",
	"login-email-edit",
	"profile-custom-fields-edit",
	"configurations-read-only",
	"configurations-edit",
	"configurations-delete",
	"billing-manage",
	"alert-action",
	"alert-create",
	"alert-add-attachment",
	"alert-delete-attachment",
	"alert-add-note",
	"alert-acknowledge",
	"alert-unacknowledge",
	"alert-snooze",
	"alert-escalate",
	"alert-close",
	"alert-delete",
	"alert-take-ownership",
	"alert-assign-ownership",
	"alert-add-recipient",
	"alert-add-team",
	"alert-edit-tags",
	"alert-edit-details",
	"alert-custom-action",
	"alert-update-priority",
	"alert-acknowledge-all",
	"alert-close-all",
	"incident-create",
	"incident-add-stakeholder",
	"incident-add-responder",
	"incident-resolve",
	"incident-reopen",
	"mass-notification-create",
	"service-access",
}

func resourceOpsGenieCustomUserRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieCustomUserRoleCreate,
		Read:   handleNonExistentResource(resourceOpsGenieCustomUserRoleRead),
		Update: resourceOpsGenieCustomUserRoleUpdate,
		Delete: resourceOpsGenieCustomUserRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"extended_role": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"user", "observer", "stakeholder",
				}, false),
			},
			"granted_rights": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validCustomRolesRights, false),
				},
				Set: schema.HashString,
			},
			"disallowed_rights": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validCustomRolesRights, false),
				},
				Set: schema.HashString,
			},
		},
	}
}

func flattenSet(input *schema.Set) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}

	for _, v := range input.List() {
		output = append(output, v.(string))
	}
	return output
}

func resourceOpsGenieCustomUserRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := custom_user_role.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	UserRoleName := d.Get("role_name").(string)
	ExtendedUserRole := d.Get("extended_role").(string)
	GrantedRights := flattenSet(d.Get("granted_rights").(*schema.Set))
	DisallowedRights := flattenSet(d.Get("disallowed_rights").(*schema.Set))

	log.Printf("[INFO] Creating OpsGenie custom user role '%s'", UserRoleName)
	result, err := client.Create(context.Background(), &custom_user_role.CreateRequest{
		Name:             UserRoleName,
		ExtendedRole:     custom_user_role.ExtendedRole(ExtendedUserRole),
		GrantedRights:    GrantedRights,
		DisallowedRights: DisallowedRights,
	})

	if err != nil {
		return err
	}

	d.SetId(result.Id)
	return resourceOpsGenieCustomUserRoleRead(d, meta)
}

func resourceOpsGenieCustomUserRoleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := custom_user_role.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	UserRoleName := d.Get("role_name").(string)

	log.Printf("[INFO] Reading OpsGenie custom role '%s'", UserRoleName)

	usrRole, err := client.Get(context.Background(), &custom_user_role.GetRequest{
		Identifier:     UserRoleName,
		IdentifierType: custom_user_role.Name,
	})
	if err != nil {
		return err
	}

	d.Set("role_name", usrRole.Name)
	d.Set("extended_role", usrRole.ExtendedRole)
	d.Set("granted_rights", usrRole.GrantedRights)
	d.Set("disallowed_rights", usrRole.DisallowedRights)

	return nil
}

func resourceOpsGenieCustomUserRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := custom_user_role.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	UserRoleName := d.Get("role_name").(string)
	ExtendedUserRole := d.Get("extended_role").(string)
	GrantedRights := flattenSet(d.Get("granted_rights").(*schema.Set))
	DisallowedRights := flattenSet(d.Get("disallowed_rights").(*schema.Set))

	log.Printf("[INFO] Updating OpsGenie custom user role '%s'", UserRoleName)

	_, err = client.Update(context.Background(), &custom_user_role.UpdateRequest{
		Identifier:       d.Id(),
		IdentifierType:   custom_user_role.Id,
		Name:             UserRoleName,
		ExtendedRole:     custom_user_role.ExtendedRole(ExtendedUserRole),
		GrantedRights:    GrantedRights,
		DisallowedRights: DisallowedRights,
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieCustomUserRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := custom_user_role.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting OpsGenie custom user role '%s'", d.Get("role_name").(string))

	_, err = client.Delete(context.Background(), &custom_user_role.DeleteRequest{
		Identifier:     d.Id(),
		IdentifierType: custom_user_role.Id,
	})

	if err != nil {
		return err
	}

	return nil
}
