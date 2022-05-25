package opslevel

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a team",
		Create:      wrap(resourceTeamCreate),
		Read:        wrap(resourceTeamRead),
		Update:      wrap(resourceTeamUpdate),
		Delete:      wrap(resourceTeamDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alias": {
				Type:        schema.TypeString,
				Description: "The human-friendly, unique identifier for the team.",
				Computed:    true,
				Deprecated:  "field 'alias' on team is no longer supported please use the 'aliases' field which is a list",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The team's display name.",
				ForceNew:    false,
				Required:    true,
			},
			"manager_email": {
				Type:        schema.TypeString,
				Description: "The email of the user who manages the team.",
				ForceNew:    false,
				Optional:    true,
			},
			"responsibilities": {
				Type:        schema.TypeString,
				Description: "A description of what the team is responsible for.",
				ForceNew:    false,
				Optional:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "A list of human-friendly, unique identifiers for the team. Must be ordered alphabetically",
				ForceNew:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"group": {
				Type:        schema.TypeString,
				Description: "The group this team belongs to. Only accepts group's Alias",
				ForceNew:    false,
				Optional:    true,
			},
			"members": {
				Type:     schema.TypeSet,
				Description: "List of user emails that belong to the team.",
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: false,
				Optional: true,
			},
		},
	}
}

func reconcileTeamAliases(d *schema.ResourceData, team *opslevel.Team, client *opslevel.Client) error {
	expectedAliases := getStringArray(d, "aliases")
	existingAliases := team.Aliases
	for _, existingAlias := range existingAliases {
		if existingAlias == team.Alias {
			continue
		}
		if stringInArray(existingAlias, expectedAliases) {
			continue
		}
		// Delete
		err := client.DeleteTeamAlias(existingAlias)
		if err != nil {
			return err
		}
	}
	for _, expectedAlias := range expectedAliases {
		if stringInArray(expectedAlias, existingAliases) {
			continue
		}
		// Add
		_, err := client.CreateAliases(team.Id, []string{expectedAlias})
		if err != nil {
			return err
		}
	}
	return nil
}

func reconcileTeamMembership(d *schema.ResourceData, team *opslevel.Team, client *opslevel.Client) error {
	expectedMembers := expandStringArray(d.Get("members").(*schema.Set).List())
	existingMembers := collectMemberEmailsFromTeam(team)
	membersToRemove := make([]string, 0)
	membersToAdd := make([]string, 0)

	for _, existingMember := range existingMembers {
		if stringInArray(existingMember, expectedMembers) {
			continue
		}

		membersToRemove = append(membersToRemove, existingMember)
	}

	for _, expectedMember := range expectedMembers {
		if expectedMember == team.Manager.Email {
			continue
		}

		if stringInArray(expectedMember, existingMembers) {
			continue
		}
		membersToAdd = append(membersToAdd, expectedMember)
	}

	if len(membersToAdd) != 0 {
		_, err := client.AddMembers(&team.TeamId, membersToAdd)
		if err != nil {
			return err
		}
	}

	if len(membersToRemove) != 0 {
		_, err := client.RemoveMembers(&team.TeamId, membersToRemove)
		if err != nil {
			return err
		}
	}
	return nil
}

func collectMemberEmailsFromTeam(team *opslevel.Team) []string {
	memberEmails := make([]string, 0, len(team.Members.Nodes))

	for _, user := range team.Members.Nodes {
		// Team managers are members by default and should not be included as part of the members attribute
		if user.Email == team.Manager.Email {
			continue
		}

		memberEmails = append(memberEmails, user.Email)
	}
	return memberEmails
}

func resourceTeamCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.TeamCreateInput{
		Name:             d.Get("name").(string),
		ManagerEmail:     d.Get("manager_email").(string),
		Responsibilities: d.Get("responsibilities").(string),
	}
	if group, ok := d.GetOk("group"); ok {
		input.Group = opslevel.NewIdentifier(group.(string))
	}
	resource, err := client.CreateTeam(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	aliasesErr := reconcileTeamAliases(d, resource, client)
	if aliasesErr != nil {
		return aliasesErr
	}

	membersErr := reconcileTeamMembership(d, resource, client)
	if membersErr != nil {
		return membersErr
	}

	return resourceTeamRead(d, client)
}

func resourceTeamRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetTeam(id)
	if err != nil {
		return err
	}

	if err := d.Set("alias", resource.Alias); err != nil {
		return err
	}
	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("manager_email", resource.Manager.Email); err != nil {
		return err
	}
	if err := d.Set("responsibilities", resource.Responsibilities); err != nil {
		return err
	}
	if _, ok := d.GetOk("group"); ok {
		if err := d.Set("group", resource.Group.Alias); err != nil {
			return err
		}
	}
	if _, ok := d.GetOk("aliases"); ok {
		aliases := []string{}
		for _, alias := range resource.Aliases {
			if alias == resource.Alias {
				// If user specifies the auto-generated alias in terraform config, don't skip it
				if stringInArray(alias, getStringArray(d, "aliases")) != true {
					continue
				}
			}
			aliases = append(aliases, alias)
		}
		sort.Strings(aliases)
		if err := d.Set("aliases", aliases); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("members"); ok {
		if err := d.Set("members", collectMemberEmailsFromTeam(resource)); err != nil {
			return err
		}
	}

	return nil
}

func resourceTeamUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.TeamUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("manager_email") {
		input.ManagerEmail = d.Get("manager_email").(string)
	}
	if d.HasChange("responsibilities") {
		input.Responsibilities = d.Get("responsibilities").(string)
	}
	if d.HasChange("group") {
		if group, ok := d.GetOk("group"); ok {
			input.Group = opslevel.NewIdentifier(group.(string))
		} else {
			input.Group = nil
		}
	}

	resource, err := client.UpdateTeam(input)
	if err != nil {
		return err
	}

	if d.HasChange("aliases") {
		tagsErr := reconcileTeamAliases(d, resource, client)
		if tagsErr != nil {
			return tagsErr
		}
	}

	if d.HasChange("members") {
		membershipErr := reconcileTeamMembership(d, resource, client)
		if membershipErr != nil {
			return membershipErr
		}
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceTeamRead(d, client)
}

func resourceTeamDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteTeam(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
