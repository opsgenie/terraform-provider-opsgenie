package opsgenie

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opsgenie/opsgenie-go-sdk-v2/incident"
	"log"
)

func resourceOpsgenieIncidentTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieIncidentTemplateCreate,
		Read:   handleNonExistentResource(resourceOpsgenieIncidentTemplateRead),
		Update: resourceOpsgenieIncidentTemplateUpdate,
		Delete: resourceOpsgenieIncidentTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"details": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"priority": {
				Type:     schema.TypeString,
				Required: true,
			},
			"impacted_services": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"stakeholder_properties": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"message": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 15000),
						},
					},
				},
			},
		},
	}
}

func resourceOpsgenieIncidentTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := incident.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	createRequest := &incident.CreateIncidentTemplateRequest{
		Name:                  d.Get("name").(string),
		Message:               d.Get("message").(string),
		Description:           d.Get("description").(string),
		Tags:                  expandOpsgenieIncidentTemplateTags(d.Get("tags").(*schema.Set)),
		Details:               expandOpsgenieIncidentTemplateDetails(d.Get("details").(map[string]interface{})),
		Priority:              incident.Priority(d.Get("priority").(string)),
		ImpactedServices:      expandOpsgenieIncidentTemplateImpactedServices(d.Get("impacted_services").(*schema.Set)),
		StakeholderProperties: expandOpsgenieIncidentTemplateStakeholderProperties(d.Get("stakeholder_properties").([]interface{})),
	}
	result, err := client.CreateIncidentTemplate(context.Background(), createRequest)
	if err != nil {
		return err
	}
	d.SetId(result.IncidentTemplateId)
	return resourceOpsgenieIncidentTemplateRead(d, meta)
}

func resourceOpsgenieIncidentTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client, err := incident.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	result, err := client.GetIncidentTemplate(context.Background(), &incident.GetIncidentTemplateRequest{})
	if err != nil {
		return err
	}
	if result != nil {
		for _, value := range result.IncidentTemplates["incidentTemplates"] {
			if d.Id() == value.IncidentTemplateId {
				d.Set("name", value.Name)
				d.Set("message", value.Message)
				d.Set("tags", value.Tags)
				d.Set("description", value.Description)
				d.Set("details", value.Details)
				d.Set("priority", value.Priority)
				d.Set("stakeholder_properties", flattenIncidentStakeHolderProperties(value.StakeholderProperties))
				if value.ImpactedServices != nil {
					d.Set("impacted_services", schema.NewSet(schema.HashString, flattenIncidentImpactedServices(value.ImpactedServices)))
				}
				break
			}
		}
	} else {
		d.SetId("")
		log.Printf("[INFO] Incident template not found")
	}
	return nil
}

func resourceOpsgenieIncidentTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := incident.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	updateRequest := &incident.UpdateIncidentTemplateRequest{
		IncidentTemplateId:    d.Id(),
		Name:                  d.Get("name").(string),
		Message:               d.Get("message").(string),
		Description:           d.Get("description").(string),
		Tags:                  expandOpsgenieIncidentTemplateTags(d.Get("tags").(*schema.Set)),
		Details:               expandOpsgenieIncidentTemplateDetails(d.Get("details").(map[string]interface{})),
		Priority:              incident.Priority(d.Get("priority").(string)),
		ImpactedServices:      expandOpsgenieIncidentTemplateImpactedServices(d.Get("impacted_services").(*schema.Set)),
		StakeholderProperties: expandOpsgenieIncidentTemplateStakeholderProperties(d.Get("stakeholder_properties").([]interface{})),
	}
	_, err = client.UpdateIncidentTemplate(context.Background(), updateRequest)
	if err != nil {
		return err
	}
	return nil
}

func resourceOpsgenieIncidentTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := incident.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &incident.DeleteIncidentTemplateRequest{IncidentTemplateId: d.Id()}
	_, err = client.DeleteIncidentTemplate(context.Background(), deleteRequest)
	if err != nil {
		return err
	}
	return nil
}

func expandOpsgenieIncidentTemplateTags(input *schema.Set) []string {
	tags := make([]string, len(input.List()))
	if input != nil {
		for k, v := range input.List() {
			tags[k] = v.(string)
		}
	}
	return tags
}

func expandOpsgenieIncidentTemplateDetails(input map[string]interface{}) map[string]string {
	details := make(map[string]string)
	if input != nil {
		for k, v := range input {
			details[k] = v.(string)
		}
	}
	return details
}

func expandOpsgenieIncidentTemplateImpactedServices(input *schema.Set) []string {
	impactedService := make([]string, len(input.List()))
	if input != nil {
		for k, v := range input.List() {
			impactedService[k] = v.(string)
		}
	}
	return impactedService
}

func expandOpsgenieIncidentTemplateStakeholderProperties(input []interface{}) incident.StakeholderProperties {
	stakeholderProperties := incident.StakeholderProperties{}
	if input != nil {
		for _, v := range input {
			config := v.(map[string]interface{})
			enable := config["enable"].(bool)
			stakeholderProperties.Enable = &enable
			stakeholderProperties.Message = config["message"].(string)
			stakeholderProperties.Description = config["description"].(string)
		}
	}
	return stakeholderProperties
}

func flattenIncidentStakeHolderProperties(prop incident.StakeholderProperties) []map[string]interface{} {
	return []map[string]interface{}{{
		"enable":      *prop.Enable,
		"message":     prop.Message,
		"description": prop.Description,
	}}
}

func flattenIncidentImpactedServices(impactedServices []string) []interface{} {
	n := make([]interface{}, len(impactedServices))
	for i, service := range impactedServices {
		n[i] = service
	}
	return n
}
