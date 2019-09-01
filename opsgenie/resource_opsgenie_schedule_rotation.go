package opsgenie

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOpsgenieScheduleRotation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsgenieScheduleRotationCreate,
		Read:   resourceOpsgenieScheduleRotationRead,
		Update: resourceOpsgenieScheduleRotationUpdate,
		Delete: resourceOpsgenieScheduleRotationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"schedule_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_date": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDate,
			},
			"end_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDate,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsgenieScheduleRotationType,
			},
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"participant": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateScheduleRotationParticipantType,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
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

func resourceOpsgenieScheduleRotationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	scheduleIdentiferValue := d.Get("schedule_id").(string)

	name := d.Get("name").(string)
	start_date := d.Get("start_date").(string)
	end_date := d.Get("end_date").(string)
	rotationType := d.Get("type").(string)
	length := d.Get("length").(int)
	timeRestriction := d.Get("time_restriction").([]interface{})
	participants := d.Get("participant").([]interface{})
	layoutStr := "2006-01-02T15:04:05Z"
	startDate, err := time.Parse(layoutStr, start_date)
	if err != nil {
		return fmt.Errorf("Cannot parse date-time")
	}

	createRequest := &schedule.CreateRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: scheduleIdentiferValue,
		Rotation: &og.Rotation{
			StartDate:    &startDate,
			Length:       uint32(length),
			Type:         og.RotationType(rotationType),
			Participants: expandOpsgenieScheduleParticipants(participants),
		},
	}

	if name != "" {
		createRequest.Rotation.Name = name
	}
	if end_date != "" {
		endDate, err := time.Parse(layoutStr, end_date)
		if err != nil {
			log.Fatal(err)
		}
		createRequest.Rotation.EndDate = &endDate
	}
	if length != 0 {
		createRequest.Rotation.Length = uint32(length)
	}
	if len(timeRestriction) > 0 {
		createRequest.Rotation.TimeRestriction = expandOpsgenieScheduleTimeRestrictions(timeRestriction)
	}

	log.Printf("[INFO] Creating OpsGenie rotation '%s'", name)

	result, err := client.CreateRotation(context.Background(), createRequest)
	if err != nil {
		return err
	}

	d.SetId(result.Id)

	return resourceOpsgenieScheduleRotationRead(d, meta)
}

func resourceOpsgenieScheduleRotationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}

	scheduleIdentiferValue := d.Get("schedule_id").(string)

	getRequest := &schedule.GetRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: scheduleIdentiferValue,
		RotationId:              d.Id(),
	}
	getResponse, err := client.GetRotation(context.Background(), getRequest)
	if err != nil {
		return err
	}

	d.Set("id", getResponse.Rotation.Id)
	d.Set("participant", getResponse.Rotation.Participants)
	d.Set("type", getResponse.Rotation.Type)
	d.Set("start_date", getResponse.Rotation.StartDate)

	return nil
}

func resourceOpsgenieScheduleRotationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	scheduleIdentiferValue := d.Get("schedule_id").(string)

	name := d.Get("name").(string)
	start_date := d.Get("start_date").(string)
	end_date := d.Get("end_date").(string)
	rotationType := d.Get("type").(string)
	length := d.Get("length").(int)
	timeRestriction := d.Get("time_restriction").([]interface{})
	participants := d.Get("participant").([]interface{})
	layoutStr := "2006-01-02T15:04:05Z"
	startDate, err := time.Parse(layoutStr, start_date)
	if err != nil {
		log.Fatal(err)
	}

	updateRequest := &schedule.UpdateRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: scheduleIdentiferValue,
		RotationId:              d.Id(),
		Rotation: &og.Rotation{
			StartDate:    &startDate,
			Length:       uint32(length),
			Type:         og.RotationType(rotationType),
			Participants: expandOpsgenieScheduleParticipants(participants),
		},
	}

	if name != "" {
		updateRequest.Rotation.Name = name
	}
	if end_date != "" {
		endDate, err := time.Parse(layoutStr, end_date)
		if err != nil {
			log.Fatal(err)
		}
		updateRequest.Rotation.EndDate = &endDate
	}
	if len(timeRestriction) > 0 {
		updateRequest.Rotation.TimeRestriction = expandOpsgenieScheduleTimeRestrictions(timeRestriction)
	}
	log.Printf("[INFO] Updating OpsGenie schedule rotation '%s'", name)

	_, err = client.UpdateRotation(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsgenieScheduleRotationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie schedule rotation '%s'", d.Get("name").(string))
	client, err := schedule.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	scheduleIdentiferValue := d.Get("schedule_id").(string)

	deleteRequest := &schedule.DeleteRotationRequest{
		ScheduleIdentifierType:  schedule.Id,
		ScheduleIdentifierValue: scheduleIdentiferValue,
		RotationId:              d.Id(),
	}

	_, err = client.DeleteRotation(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func expandOpsgenieScheduleParticipants(input []interface{}) []og.Participant {
	participants := make([]og.Participant, 0, len(input))

	if input == nil {
		return participants
	}

	for _, v := range input {
		config := v.(map[string]interface{})

		participantType := config["type"].(string)
		Id := config["id"].(string)

		participant := og.Participant{
			Type: og.ParticipantType(participantType),
			Id:   Id,
		}

		participants = append(participants, participant)
	}

	return participants
}
func expandOpsgenieScheduleTimeRestrictions(d []interface{}) *og.TimeRestriction {

	timeRestriction := og.TimeRestriction{}

	for _, v := range d {
		config := v.(map[string]interface{})

		timeRestrictionType := config["type"].(string)
		timeRestriction.Type = og.RestrictionType(timeRestrictionType)

		if len(config["restrictions"].([]interface{})) > 0 {
			timeRestriction.RestrictionList = expandOpsgenieScheduleRestrictions(config["restrictions"].([]interface{}))
		} else {
			timeRestriction.Restriction = expandOpsgenieScheduleRestriction(config["restriction"].([]interface{}))
		}
	}

	return &timeRestriction
}

func expandOpsgenieScheduleRestrictions(input []interface{}) []og.Restriction {
	restrictionList := make([]og.Restriction, 0, len(input))

	if input == nil {
		return restrictionList
	}

	for _, v := range input {
		config := v.(map[string]interface{})

		restriction := og.Restriction{
			StartDay:  og.Day(config["start_day"].(string)),
			StartHour: uint32(config["start_hour"].(int)),
			StartMin:  uint32(config["start_min"].(int)),
			EndHour:   uint32(config["end_hour"].(int)),
			EndDay:    og.Day(config["end_day"].(string)),
			EndMin:    uint32(config["end_min"].(int)),
		}

		restrictionList = append(restrictionList, restriction)
	}

	return restrictionList
}

func expandOpsgenieScheduleRestriction(input []interface{}) og.Restriction {

	restriction := og.Restriction{}
	for _, v := range input {
		config := v.(map[string]interface{})

		restriction = og.Restriction{
			StartHour: uint32(config["start_hour"].(int)),
			StartMin:  uint32(config["start_min"].(int)),
			EndHour:   uint32(config["end_hour"].(int)),
			EndMin:    uint32(config["end_min"].(int)),
		}

	}

	return restriction
}

func validateOpsgenieScheduleRotationType(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"daily":  true,
		"weekly": true,
		"hourly": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Opsgenie Schedule Rotation Type  can only be 'Daily' ,'Weekly' or 'Hourly'"))
	}
	return
}

func validateScheduleRotationParticipantType(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"user":       true,
		"team":       true,
		"escalation": true,
		"none":       true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("it can only be one of these 'user', 'schedule', 'team', 'escalation'"))
	}
	return
}

func validateDay(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"monday":    true,
		"tuesday":   true,
		"wednesday": true,
		"thursday":  true,
		"friday":    true,
		"saturday":  true,
		"sunday":    true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("it can only be day of week (monday,tuesday...)"))
	}
	return
}

func validateHourParams(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if value < 0 || value > 24 {
		errors = append(errors, fmt.Errorf("hour must between 0-24"))
	}
	return
}
func validateMinParams(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if value < 0 || value > 59 {
		errors = append(errors, fmt.Errorf("minute must in between of 0-59"))

	}
	return
}
