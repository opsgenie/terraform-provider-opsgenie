package opsgenie

import (
	"context"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func resourceOpsGenieTeamRoutingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieTeamRoutingRuleCreate,
		Read:   resourceOpsGenieTeamRoutingRuleRead,
		Update: resourceOpsGenieTeamRoutingRuleUpdate,
		Delete: resourceOpsGenieTeamRoutingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"criteria": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"conditions": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "user",
						},
					},
				},
			},
			"time_restriction": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"restrictions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_day": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateDay,
									},
									"end_day": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateDay,
									},
									"start_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateHourParams,
									},
									"start_min": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateMinParams,
									},
									"end_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateHourParams,
									},
									"end_min": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateMinParams,
									},
								},
							},
						},
						"restriction": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateHourParams,
									},
									"start_min": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateMinParams,
									},
									"end_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateHourParams,
									},
									"end_min": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateMinParams,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceOpsGenieTeamRoutingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	teamId := d.Get("team_id").(string)
	order := d.Get("order").(int)
	//timezone := d.Get("timezone").(string)
	timeRestriction := d.Get("time_restriction").([]interface{})
	criteria := d.Get("criteria").(map[string]interface{}) //TODO DÖN
	createRequest := &team.CreateRoutingRuleRequest{
		TeamIdentifierType:  team.Id,
		TeamIdentifierValue: teamId,
		Name:                name,
		Order:               &order,
		//	Timezone:            timezone,
		Criteria: &og.Criteria{
			CriteriaType: criteria["type"],
			Conditions:   nil,
		}
		//Notify:              expandOpsgenieNotify(),
	}

	_ = createRequest

	if len(timeRestriction) > 0 {
		createRequest.TimeRestriction = expandTimeRestrictions(timeRestriction)
	}

	log.Printf("[INFO] Creating OpsGenie team '%s'", name)

	_, err = client.CreateRoutingRule(context.Background(), createRequest)
	if err != nil {
		return err
	}

	getRequest := &team.GetRoutingRuleRequest{
		TeamIdentifierType:  0,
		TeamIdentifierValue: "",
		RoutingRuleId:       "",
	}

	getResponse, err := client.GetRoutingRule(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.SetId(getResponse.Id)

	return resourceOpsGenieTeamRead(d, meta)
}

func resourceOpsGenieTeamRoutingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	getRequest := &team.GetRoutingRuleRequest{
		TeamIdentifierType:  0,
		TeamIdentifierValue: "",
		RoutingRuleId:       "",
	}

	_, err = client.GetRoutingRule(context.Background(), getRequest)
	if err != nil {
		return err
	}

	/*d.Set("name", getResponse.Name)
	d.Set("description", getResponse.Description)
	d.Set("member", flattenOpsGenieTeamMembers(getResponse.Members))
	*/
	return nil
}

func resourceOpsGenieTeamRoutingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	//description := d.Get("description").(string)

	updateRequest := &team.UpdateRoutingRuleRequest{
		TeamIdentifierType:  0,
		TeamIdentifierValue: "",
		RoutingRuleId:       "",
		Name:                "",
		Timezone:            "",
		Criteria: &og.Filter{
			ConditionMatchType: "",
			Conditions:         nil,
		},
		TimeRestriction: &og.TimeRestriction{
			Type:            "",
			RestrictionList: nil,
			Restriction: og.Restriction{
				StartDay:  "",
				StartHour: nil,
				StartMin:  nil,
				EndHour:   nil,
				EndDay:    "",
				EndMin:    nil,
			},
		},
		Notify: &team.Notify{
			Type: "",
			Name: "",
			Id:   "",
		},
	}

	log.Printf("[INFO] Updating OpsGenie team '%s'", name)

	_, err = client.UpdateRoutingRule(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsGenieTeamRoutingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie team '%s'", d.Get("name").(string))
	client, err := team.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &team.DeleteRoutingRuleRequest{
		TeamIdentifierType:  0,
		TeamIdentifierValue: "",
		RoutingRuleId:       "",
	}

	_, err = client.DeleteRoutingRule(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func expandOpsgenieNotify(input []team.Notify) *team.Notify {

	return &team.Notify{
		Type: "",
		Name: "",
		Id:   "",
	}
}

func expandOpsgenieCriteria(input []interface{}) *og.Filter {
	//crit := input.(og.Criteria)
	//_ = crit
	return &og.Filter{}
}
