package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/kubectl-opslevel/opslevel"
)

var opsLevelTeamSchema = map[string]*schema.Schema{
	"alias": {
		Type:        schema.TypeString,
		Description: "The human-friendly, unique identifier for the team",
		ForceNew:    false,
		Optional:    true,
	},
	"contacts": {
		Type:        schema.TypeList,
		Description: "The contacts for the team.",
		ForceNew:    false,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"address": {
					Type:        schema.TypeString,
					Description: "The contact address. Examples: support@company.com for type email, https://opslevel.com for type web.",
					ForceNew:    false,
					Required:    true,
				},
				"display_name": {
					Type:        schema.TypeString,
					Description: "The name shown in the UI for the contact.",
					ForceNew:    false,
					Optional:    true,
				},
				//"id": {
				//	Type:        schema.TypeString,
				//	Description: "The unique identifier for the contact.",
				//	ForceNew:    false,
				//	Required:    true,
				//},
				//"type": {
				//	Type:        schema.TypeString,
				//	Description: "The method of contact.",
				//	ForceNew:    false,
				//	Required:    true,
				//},
			},
		},
	},
	//"manager_id": {
	//	Type:        schema.TypeString,
	//	Description: "The ID of the user who manages the team.",
	//	ForceNew:    false,
	//	Optional:    true,
	//},
	//"member_id_list": {
	//	Type:        schema.TypeList,
	//	Description: "The users that are on the team.",
	//	ForceNew:    false,
	//	Optional:    true,
	//	Elem:        schema.TypeString,
	//},
	"name": {
		Type:        schema.TypeString,
		Description: "The display name of the service.",
		ForceNew:    false,
		Required:    true,
	},
	"responsibilities": {
		Type:        schema.TypeString,
		Description: "A description of what the team is responsible for.",
		ForceNew:    false,
		Optional:    true,
	},
}

func flattenTeam(team *opslevel.Team) map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = team.Id.(string)
	m["name"] = string(team.Name)
	m["alias"] = string(team.Alias)
	m["responsibilities"] = string(team.Responsibilities)
	m["contacts"] = flattenContacts(team.Contacts)

	return m
}

func flattenContacts(in []opslevel.Contact) []interface{} {
	out := make([]interface{}, len(in))

	for i, c := range in {
		m := make(map[string]interface{})
		m["display_name"] = string(c.DisplayName)
		m["address"] = string(c.Address)
		out[i] = m
	}

	return out
}
