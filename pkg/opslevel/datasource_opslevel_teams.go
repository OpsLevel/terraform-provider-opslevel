package opslevel

import (
	"fmt"
	"log"

	"github.com/kr/pretty"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceOpsLevelTeams() *schema.Resource {
	teamSchema := datasourceSchemaFromResourceSchema(opsLevelTeamSchema)

	teamSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The team ID",
		Computed:    true,
	}

	return &schema.Resource{
		Read: datasourceOpsLevelTeamsRead,
		Schema: map[string]*schema.Schema{
			//"filter": {
			//	Type:        schema.TypeList,
			//	Description: "The teams matching filter.",
			//	ForceNew:    true,
			//	Required:    true,
			//	MaxItems:    1,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"field": {
			//				Type:        schema.TypeString,
			//				Description: "The team field to filter on.",
			//				ForceNew:    true,
			//				Required:    true,
			//			},
			//			"value": {
			//				Type:        schema.TypeString,
			//				Description: "The field value to match.",
			//				ForceNew:    true,
			//				Optional:    true,
			//			},
			//		},
			//	},
			//},
			"teams": {
				Type:        schema.TypeList,
				Description: "The teams matching filter.",
				ForceNew:    false,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: teamSchema,
				},
			},
		},
	}
}

func datasourceOpsLevelTeamsRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)
	teams := []interface{}{}

	resp, err := p.client.ListTeams()
	if err != nil {
		return err
	}
	for _, team := range resp {
		teams = append(teams, flattenTeam(&team))
	}
	log.Printf("[DEBUG] teams len=%d cap=%d %s", len(teams), cap(teams), pretty.Sprint(teams))

	d.SetId(fmt.Sprintf("teams"))
	d.Set("teams", teams)

	return nil
}
