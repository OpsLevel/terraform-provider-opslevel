package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func resourceTeamMembership() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a team's user membership",
		Create:      wrap(resourceTeamMembershipCreate),
		Read:        wrap(resourceTeamMembershipRead),
		Update:      wrap(resourceTeamMembershipUpdate),
		Delete:      wrap(resourceTeamMembershipDelete),
		Importer: &schema.ResourceImporter{
			State: resourceTeamMembershipImport,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"team_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				Description: "List of user emails that belong to the team.",
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTeamMembershipCreate(d *schema.ResourceData, client *opslevel.Client) error {
	teamAlias := d.Get("team_alias").(string)
	members := expandStringArray(d.Get("members").(*schema.Set).List())

	team, err := getTeamFromAlias(teamAlias, client)
	if err != nil {
		return err
	}

	if err = addMembers(team, members, client); err != nil {
		return err
	}

	d.SetId(resource.PrefixedUniqueId(""))

	return resourceTeamMembershipRead(d, client)
}

func resourceTeamMembershipRead(d *schema.ResourceData, client *opslevel.Client) error {
	teamAlias := d.Get("team_alias").(string)
	members, err := getCurrentMembersFromTeamAlias(teamAlias, client)
	if err != nil {
		return nil
	}

	d.Set("members", members)

	return nil
}

func resourceTeamMembershipUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	if d.HasChange("members") {
		teamAlias := d.Get("team_alias").(string)
		o, n := d.GetChange("members")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		removeSlice := expandStringArray(os.Difference(ns).List())
		addSlice := expandStringArray(ns.Difference(os).List())

		team, err := getTeamFromAlias(teamAlias, client)
		if err != nil {
			return err
		}

		if len(removeSlice) != 0 {
			if err = removeMembers(team, removeSlice, client); err != nil {
				return err
			}
		}

		if len(addSlice) != 0 {
			if err = addMembers(team, addSlice, client); err != nil {
				return err
			}
		}
	}
	d.Set("last_updated", timeLastUpdated())

	return resourceTeamMembershipRead(d, client)
}

func resourceTeamMembershipDelete(d *schema.ResourceData, client *opslevel.Client) error {
	teamAlias := d.Get("team_alias").(string)
	members := expandStringArray(d.Get("members").(*schema.Set).List())

	team, err := getTeamFromAlias(teamAlias, client)
	if err != nil {
		return err
	}

	if err = removeMembers(team, members, client); err != nil {
		return err
	}

	return nil
}

func resourceTeamMembershipImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*opslevel.Client)
	teamAlias := d.Id()
	currentMembers, err := getCurrentMembersFromTeamAlias(teamAlias, client)
	if err != nil {
		return nil, err
	}

	d.Set("team_alias", teamAlias)
	d.Set("members", currentMembers)

	d.SetId(resource.PrefixedUniqueId(""))
	return []*schema.ResourceData{d}, nil
}

func getTeamFromAlias(teamAlias string, client *opslevel.Client) (*opslevel.Team, error) {
	team, err := client.GetTeamWithAlias(teamAlias)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func getCurrentMembersFromTeamAlias(getCurrentMembersFromTeamAlias string, client *opslevel.Client) ([]string, error) {
	team, err := getTeamFromAlias(getCurrentMembersFromTeamAlias, client)
	if err != nil {
		return nil, err
	}

	members, err := collectMembersForTeam(team, client)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func collectMembersForTeam(team *opslevel.Team, client *opslevel.Client) ([]string, error) {
	memberEmails := make([]string, 0, len(team.Members.Nodes))

	for _, user := range team.Members.Nodes {
		memberEmails = append(memberEmails, user.Email)
	}
	return memberEmails, nil
}

func removeMembers(team *opslevel.Team, members []string, client *opslevel.Client) error {
	_, err := client.RemoveMembers(&team.TeamId, members)
	if err != nil {
		return err
	}
	return nil
}

func addMembers(team *opslevel.Team, members []string, client *opslevel.Client) error {
	_, err := client.AddMembers(&team.TeamId, members)
	if err != nil {
		return err
	}
	return nil
}
