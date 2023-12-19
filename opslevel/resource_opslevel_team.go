package opslevel

import (
	"errors"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
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
			"name": {
				Type:        schema.TypeString,
				Description: "The team's display name.",
				ForceNew:    false,
				Required:    true,
			},
			"responsibilities": {
				Type:        schema.TypeString,
				Description: "A description of what the team is responsible for.",
				ForceNew:    false,
				Optional:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "A list of human-friendly, unique identifiers for the team.",
				ForceNew:    false,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"group": {
				Type:        schema.TypeString,
				Description: "The group this team belongs to. Only accepts group's Alias",
				Deprecated:  "field 'group' on team is no longer supported please use the 'parent' field.",
				ForceNew:    false,
				Optional:    true,
			},
			"member": {
				Type:        schema.TypeList,
				Description: "List of members in the team with email address and role.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Description: "The email address or ID of the user to add to a team.",
							Required:    true,
						},
						"role": {
							Type:        schema.TypeString,
							Description: "The type of relationship this membership implies.",
							Required:    true,
						},
					},
				},
			},
			"parent": {
				Type:        schema.TypeString,
				Description: "The parent team. Only accepts team's Alias",
				ForceNew:    false,
				Optional:    true,
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

func collectMembersFromTeam(team *opslevel.Team) []opslevel.TeamMembershipUserInput {
	members := []opslevel.TeamMembershipUserInput{}

	for _, user := range team.Memberships.Nodes {
		member := opslevel.TeamMembershipUserInput{
			User: opslevel.NewUserIdentifier(user.User.Email),
			Role: string(user.Role),
		}
		members = append(members, member)
	}
	return members
}

func memberInArray(member opslevel.TeamMembershipUserInput, array []opslevel.TeamMembershipUserInput) bool {
	for _, m := range array {
		if m.User.Email == member.User.Email && m.Role == member.Role {
			return true
		}
	}
	return false
}

func reconcileTeamMembership(d *schema.ResourceData, team *opslevel.Team, client *opslevel.Client) error {
	expectedMembers := []opslevel.TeamMembershipUserInput{}
	existingMembers := collectMembersFromTeam(team)

	if members, ok := d.GetOk("member"); ok {
		membersInput := members.([]interface{})

		for _, m := range membersInput {
			memberInput := m.(map[string]interface{})
			member := opslevel.TeamMembershipUserInput{
				User: opslevel.NewUserIdentifier(memberInput["email"].(string)),
				Role: memberInput["role"].(string),
			}
			expectedMembers = append(expectedMembers, member)
		}
	}

	membersToRemove := []opslevel.TeamMembershipUserInput{}
	membersToAdd := []opslevel.TeamMembershipUserInput{}

	for _, existingMember := range existingMembers {
		if memberInArray(existingMember, expectedMembers) {
			continue
		}

		membersToRemove = append(membersToRemove, existingMember)
	}

	for _, expectedMember := range expectedMembers {
		if memberInArray(expectedMember, existingMembers) {
			continue
		}
		membersToAdd = append(membersToAdd, expectedMember)
	}

	// warning: must remove memberships before adding them.
	// this prevents a bug where the role of a user changes
	// but the user isn't added back and disappears.
	if len(membersToRemove) != 0 {
		_, err := client.RemoveMemberships(&team.TeamId, membersToRemove...)
		if err != nil {
			return err
		}
	}

	if len(membersToAdd) != 0 {
		_, err := client.AddMemberships(&team.TeamId, membersToAdd...)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTeamCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.TeamCreateInput{
		Name:             d.Get("name").(string),
		Responsibilities: d.Get("responsibilities").(string),
	}
	if _, ok := d.GetOk("group"); ok {
		return errors.New("groups are deprecated - create and update are disabled.")
	}
	if parentTeam, ok := d.GetOk("parent"); ok {
		input.ParentTeam = opslevel.NewIdentifier(parentTeam.(string))
	}

	resource, err := client.CreateTeam(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

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

	resource, err := client.GetTeam(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("responsibilities", resource.Responsibilities); err != nil {
		return err
	}
	if err := d.Set("group", resource.Group.Alias); err != nil {
		return err
	}
	if err := d.Set("parent", resource.ParentTeam.Alias); err != nil {
		return err
	}

	aliases := []string{}
	for _, alias := range resource.Aliases {
		if alias == resource.Alias {
			// If user specifies the auto-generated alias in terraform config, don't skip it
			if !stringInArray(alias, getStringArray(d, "aliases")) {
				continue
			}
		}
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	if err := d.Set("aliases", aliases); err != nil {
		return err
	}

	// only read members if it was set before in the configuration
	// some customers may not have any member {} blocks defined
	// in their config, and they cannot use terraform to manage
	// teams because of it without either adding the members into
	// the config or unassigning all the members (unwanted)
	if members, ok := d.GetOk("member"); members != nil || ok {
		members := collectMembersFromTeam(resource)
		memberOutput := []map[string]interface{}{}
		for _, m := range members {
			mOutput := make(map[string]interface{})
			mOutput["email"] = m.User.Email
			mOutput["role"] = m.Role
			memberOutput = append(memberOutput, mOutput)
		}
		if err := d.Set("member", memberOutput); err != nil {
			return err
		}
	}

	return nil
}

func resourceTeamUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	input := opslevel.TeamUpdateInput{
		Id: opslevel.ID(id),
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("responsibilities") {
		input.Responsibilities = d.Get("responsibilities").(string)
	}
	if d.HasChange("group") {
		return errors.New("groups are deprecated - create and update are disabled.")
	}
	if d.HasChange("parent") {
		if parentTeam, ok := d.GetOk("parent"); ok {
			input.ParentTeam = opslevel.NewIdentifier(parentTeam.(string))
		} else {
			input.ParentTeam = nil
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

	if d.HasChange("member") {
		membersErr := reconcileTeamMembership(d, resource, client)
		if membersErr != nil {
			return membersErr
		}
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceTeamRead(d, client)
}

func resourceTeamDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteTeam(opslevel.ID(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
