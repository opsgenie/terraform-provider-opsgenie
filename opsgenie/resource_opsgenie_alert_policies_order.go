package opsgenie

import (
	"context"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/policy"
)

func resourceOpsGeniePoliciesOrder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpsGeniePoliciesOrderSet,
		ReadContext:   resourceOpsGeniePoliciesOrderRead,
		UpdateContext: resourceOpsGeniePoliciesOrderSet,
		DeleteContext: resourceOpsGeniePoliciesOrderDelete,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceOpsGeniePoliciesOrderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// The map makes it easier to check if an id is set in the order list
	policy_ids := make(map[string]struct{})
	for _, id := range d.Get("policy_ids").([]interface{}) {
		policy_ids[id.(string)] = struct{}{}
	}

	current_policy_ids, err := readPoliciesOrder(d.Get("team_id").(string), meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Ignore current ids that aren't set in the order list
	// Avoids systematic diff when other policies aren't ordered through the ressource
	// (which should be avoided anyway)
	filtered_current_policy_ids := make([]string, len(policy_ids))
	idx := 0
	for _, id := range current_policy_ids {
		_, ok := policy_ids[id]
		if ok {
			filtered_current_policy_ids[idx] = id
			idx++
		}
	}

	d.Set("policy_ids", filtered_current_policy_ids)

	return nil
}

func resourceOpsGeniePoliciesOrderSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return diag.FromErr(err)
	}

	team_id := d.Get("team_id").(string)
	policy_ids := d.Get("policy_ids").([]interface{})

	for i, policy_id := range policy_ids {
		orderRequest := &policy.ChangeOrderRequest{
			Id:          policy_id.(string),
			Type:        policy.AlertPolicy,
			TargetIndex: i + 1,
		}
		if d.Get("team_id").(string) != "" {
			orderRequest.TeamId = d.Get("team_id").(string)
		}
		log.Printf("[INFO] Updating order for Alert Policy '%s' with index %d", policy_id, i)
		_, err = client.ChangeOrder(context.Background(), orderRequest)
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	d.SetId(team_id)

	return nil
}

func resourceOpsGeniePoliciesOrderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func readPoliciesOrder(team_id string, meta interface{}) ([]string, error) {
	log.Printf("[INFO] Reading OpsGenie Alert Policies for team %s", team_id)

	client, err := policy.NewClient(meta.(*OpsgenieClient).client.Config)
	if err != nil {
		return nil, err
	}

	listRequest := &policy.ListAlertPoliciesRequest{}
	if team_id != "" {
		listRequest.TeamId = team_id
	}

	policiesRes, err := client.ListAlertPolicies(context.Background(), listRequest)
	if err != nil {
		return nil, err
	}

	policies := policiesRes.Policies

	sort.SliceStable(policies, func(i, j int) bool {
		return policies[i].Order < policies[j].Order
	})

	policy_ids := make([]string, len(policies))

	for i, policy := range policies {
		policy_ids[i] = policy.Id
	}

	return policy_ids, nil
}
