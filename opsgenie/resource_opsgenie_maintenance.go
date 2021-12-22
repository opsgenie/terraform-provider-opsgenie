package opsgenie

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/maintenance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpsgenieMaintenance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieMaintenanceCreate,
		Read:   handleNonExistentResource(resourceOpsgenieMaintenanceRead),
		Update: resourceOpsgenieMaintenanceUpdate,
		Delete: resourceOpsgenieMaintenanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_date": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateDate,
						},
						"end_date": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateDate,
						},
					},
				},
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"state": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceOpsgenieMaintenanceCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := maintenance.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	description := d.Get("description").(string)

	createRequest := &maintenance.CreateRequest{
		Description: description,
		Time:        expandOpsgenieMaintenanceTime(d),
		Rules:       expandOpsgenieMaintenanceRules(d),
	}

	log.Printf("[INFO] Creating OpsGenie maintenance")

	result, err := client.Create(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	return resourceOpsgenieMaintenanceRead(d, meta)
}

func resourceOpsgenieMaintenanceRead(d *schema.ResourceData, meta interface{}) error {
	client, err := maintenance.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	listResponse, err := client.List(context.Background(), &maintenance.ListRequest{})
	if err != nil {
		return err
	}

	found := maintenance.GetResult{}
	for _, maintenanceResp := range listResponse.Maintenances {
		if maintenanceResp.Id == d.Id() {
			found.Time = maintenanceResp.Time
			found.Id = maintenanceResp.Id
			found.Description = maintenanceResp.Description
			found.Status = maintenanceResp.Status

			break
		}
	}

	if found.Id == "" {
		d.SetId("")
		log.Printf("[INFO] Maintenance not found. Removing from state")
		return nil
	}

	d.Set("time", flattenMaintenanceTime(found.Time))
	d.Set("description", found.Description)

	return nil
}

func resourceOpsgenieMaintenanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := maintenance.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	mnt, err := client.Get(context.Background(), &maintenance.GetRequest{
		Id: d.Id(),
	})
	if err != nil {
		log.Printf("[ERROR] Maintenance could not fetch")
		return err

	}
	maintenanceTime := expandOpsgenieMaintenanceTime(d)
	if mnt.Status == "active" {

		_, err := client.ChangeEndDate(context.Background(), &maintenance.ChangeEndDateRequest{
			Id:      d.Id(),
			EndDate: maintenanceTime.EndDate,
		})
		if err != nil {
			return err
		}

	} else if mnt.Status == "planned" {
		description := d.Get("description").(string)

		updateRequest := &maintenance.UpdateRequest{
			Id:          d.Id(),
			Description: description,
			Rules:       expandOpsgenieMaintenanceRules(d),
			Time:        maintenanceTime,
		}

		log.Printf("[INFO] Updating OpsGenie maintenance")

		_, err = client.Update(context.Background(), updateRequest)
		if err != nil {
			log.Printf("%s", err.Error())
			return err
		}
	} else {
		log.Printf("[ERROR] You cannot edit past maintenance")
		return errors.New("You cannot edit" + mnt.Status + " maintenances")

	}

	return nil
}

func resourceOpsgenieMaintenanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie escalation ")
	client, err := maintenance.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &maintenance.DeleteRequest{
		Id: d.Id(),
	}

	_, err = client.Delete(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func expandOpsgenieMaintenanceRules(d *schema.ResourceData) []maintenance.Rule {
	input := d.Get("rules").([]interface{})

	rules := make([]maintenance.Rule, 0, len(input))
	if input == nil {
		return rules
	}
	for _, v := range input {
		config := v.(map[string]interface{})

		state := config["state"].(string)
		entity := config["entity"].([]interface{})
		entityObj := expandOpsgenieMaintenanceEntity(entity)
		rule := maintenance.Rule{
			Entity: entityObj,
		}
		if entityObj.Type != maintenance.Integration {
			rule.State = maintenance.RuleState(state)
		} else {
			rule.State = maintenance.Disabled
		}

		rules = append(rules, rule)
	}

	return rules
}

func expandOpsgenieMaintenanceEntity(d []interface{}) maintenance.Entity {
	entity := maintenance.Entity{}

	for _, e := range d {
		ent := e.(map[string]interface{})
		entity.Id = ent["id"].(string)
		entity.Type = maintenance.RuleEntityType(ent["type"].(string))
	}
	return entity
}

func expandOpsgenieMaintenanceTime(d *schema.ResourceData) maintenance.Time {
	input := d.Get("time").([]interface{})

	maintenanceTime := maintenance.Time{}
	if input == nil {
		return maintenanceTime
	}
	for _, v := range input {
		config := v.(map[string]interface{})

		maintenanceType := config["type"].(string)
		start_Date := config["start_date"]
		end_Date := config["end_date"]

		maintenanceTime.Type = maintenance.TimeType(maintenanceType)
		if maintenanceTime.Type == maintenance.Schedule {
			if start_Date == nil || end_Date == nil {
				log.Fatal("Schedule type maintenance's must have start and end dates")
			}
		}

		layoutStr := "2006-01-02T15:04:05Z"
		startDate, err := time.Parse(layoutStr, start_Date.(string))
		if err != nil {
			log.Fatal(err)
		}
		endDate, err := time.Parse(layoutStr, end_Date.(string))
		if err != nil {
			log.Fatal(err)
		}
		maintenanceTime.StartDate = &startDate
		maintenanceTime.EndDate = &endDate
	}

	return maintenanceTime
}

func flattenMaintenanceTime(time maintenance.Time) []map[string]interface{} {
	timeLayout := "2006-01-02T15:04:05Z"
	return []map[string]interface{}{{
		"type":       time.Type,
		"start_date": time.StartDate.Format(timeLayout),
		"end_date":   time.EndDate.Format(timeLayout),
	}}
}
