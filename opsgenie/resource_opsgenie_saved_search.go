package opsgenie

import (
	"context"
	"log"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	ogClient "github.com/opsgenie/opsgenie-go-sdk-v2/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOpsgenieSavedSearch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpsgenieSavedSearchCreate,
		ReadContext:   resourceOpsgenieSavedSearchRead,
		Update:        resourceOpsGenieSavedSearchUpdate,
		Delete:        resourceOpsgenieSavedSearchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"query": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"owner": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 15000),
			},
			"teams": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceOpsgenieSavedSearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := alert.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	query := d.Get("query").(string)

	createRequest := &alert.CreateSavedSearchRequest{
		Description: description,
		Name:        name,
		Query:       query,
	}
	if d.Get("owner").(*schema.Set).Len() > 0 {
		createRequest.Owner = expandOpsGenieSavedSearchOwner(d.Get("owner").(*schema.Set))
	}

	if len(d.Get("teams").([]interface{})) > 0 {
		createRequest.Teams = expandOpsGenieSavedSearchTeams(d.Get("teams"))
	}
	log.Printf("[INFO] Creating OpsGenie savedSearch")

	result, err := client.CreateSavedSearch(context.Background(), createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Id)

	return resourceOpsgenieSavedSearchRead(ctx, d, meta)
}

func resourceOpsgenieSavedSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := alert.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	savedSearchRes := &alert.GetSavedSearchResult{}
	savedSearchRes, err = client.GetSavedSearch(context.Background(), &alert.GetSavedSearchRequest{
		IdentifierValue: d.Id(),
	})
	if err != nil {
		x := err.(*ogClient.ApiError)
		if x.StatusCode == 404 {
			log.Printf("[WARN] Removing SavedSearch because it's gone %s", d.Get("name"))
			d.SetId("")
			return nil
		}
	}

	d.Set("name", savedSearchRes.Name)
	d.Set("description", savedSearchRes.Description)
	d.Set("query", savedSearchRes.Query)
	d.Set("teams", flattenOpsGenieSeavedSearchTeams(savedSearchRes.Teams))

	return nil
}

func resourceOpsGenieSavedSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := alert.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	query := d.Get("query").(string)

	updateRequest := &alert.UpdateSavedSearchRequest{
		IdentifierValue: d.Id(),
		Description:     description,
		NewName:         name,
		Query:           query,
	}
	if d.Get("owner").(*schema.Set).Len() > 0 {
		updateRequest.Owner = expandOpsGenieSavedSearchOwner(d.Get("owner").(*schema.Set))
	}

	if len(d.Get("teams").([]interface{})) > 0 {
		updateRequest.Teams = expandOpsGenieSavedSearchTeams(d.Get("teams"))
	}
	log.Printf("[INFO] Creating OpsGenie savedSearch")

	_, err = client.UpdateSavedSearch(context.Background(), updateRequest)
	if err != nil {
		return err
	}

	return nil
}

func resourceOpsgenieSavedSearchDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie Saved Search ")
	client, err := alert.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return err
	}
	deleteRequest := &alert.DeleteSavedSearchRequest{
		IdentifierValue: d.Id(),
	}

	_, err = client.DeleteSavedSearch(context.Background(), deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func flattenOpsGenieSeavedSearchTeams(input []alert.Team) []map[string]interface{} {
	teams := make([]map[string]interface{}, 0)
	for _, v := range input {
		team := make(map[string]interface{})
		team["name"] = v.Name
		team["id"] = v.ID
		teams = append(teams, team)
	}

	return teams
}

func expandOpsGenieSavedSearchTeams(input interface{}) []alert.Team {
	teams := make([]alert.Team, 0)

	if input == nil {
		return teams
	}

	for _, v := range input.([]interface{}) {
		teamMap := v.(map[string]interface{})
		teamID := teamMap["id"].(string)
		name := teamMap["name"].(string)

		team := alert.Team{
			ID:   teamID,
			Name: name,
		}

		teams = append(teams, team)
	}

	return teams
}

func expandOpsGenieSavedSearchOwner(input *schema.Set) alert.User {
	var output alert.User
	if input == nil {
		return output
	}
	for _, v := range input.List() {
		config := v.(map[string]interface{})
		output.ID = config["id"].(string)
		output.Username = config["username"].(string)
	}

	return output
}
