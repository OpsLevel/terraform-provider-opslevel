package opslevel

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2024"
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
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"member": {
				Type:        schema.TypeList,
				Description: "List of members in the team with email address and role.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Description: "The email address of the team member.",
							Required:    true,
						},
						"role": {
							Type:        schema.TypeString,
							Description: "The role of the team member.",
							Required:    true,
						},
					},
				},
			},
			"parent": {
				Type:        schema.TypeString,
				Description: "The id or alias of the parent team.",
				ForceNew:    false,
				Optional:    true,
			},
		},
	}
}

func reconcileTeamAliases(d *schema.ResourceData, team *opslevel.Team, client *opslevel.Client) error {
	expectedAliases := getStringArray(d, "aliases")
	existingAliases := team.ManagedAliases
	for _, existingAlias := range existingAliases {
		if !slices.Contains(expectedAliases, existingAlias) {
			err := client.DeleteTeamAlias(existingAlias)
			if err != nil {
				return err
			}
		}
	}
	for _, expectedAlias := range expectedAliases {
		if !slices.Contains(existingAliases, expectedAlias) {
			_, err := client.CreateAliases(team.Id, []string{expectedAlias})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func collectMembersFromTeam(team *opslevel.Team) []opslevel.TeamMembershipUserInput {
	members := make([]opslevel.TeamMembershipUserInput, 0)

	for _, user := range team.Memberships.Nodes {
		newUserIdentifier := opslevel.NewUserIdentifier(user.User.Email)
		member := opslevel.TeamMembershipUserInput{
			User: newUserIdentifier,
			Role: opslevel.RefOf(user.Role),
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
	expectedMembers := make([]opslevel.TeamMembershipUserInput, 0)
	existingMembers := collectMembersFromTeam(team)

	if members, ok := d.GetOk("member"); ok {
		membersInput := members.([]interface{})

		for _, m := range membersInput {
			memberInput := m.(map[string]interface{})
			newUserIdentifier := opslevel.NewUserIdentifier(memberInput["email"].(string))
			member := opslevel.TeamMembershipUserInput{
				User: newUserIdentifier,
				Role: opslevel.RefOf(memberInput["role"].(string)),
			}
			expectedMembers = append(expectedMembers, member)
		}
	}

	membersToRemove := make([]opslevel.TeamMembershipUserInput, 0)
	membersToAdd := make([]opslevel.TeamMembershipUserInput, 0)

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
		Responsibilities: opslevel.RefOf(d.Get("responsibilities").(string)),
	}
	if parent := d.Get("parent"); parent != "" {
		input.ParentTeam = opslevel.NewIdentifier(parent.(string))
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

	// only read in changes to optional fields if they have been set before
	if parent, ok := d.GetOk("parent"); ok || parent != "" {
		var parentValue string
		if opslevel.IsID(parent.(string)) {
			parentValue = string(resource.ParentTeam.Id)
		} else {
			parentValue = string(resource.ParentTeam.Alias)
		}

		if err := d.Set("parent", parentValue); err != nil {
			return err
		}
	}

	slices.Sort(resource.ManagedAliases)
	if err := d.Set("aliases", resource.ManagedAliases); err != nil {
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
		Id: opslevel.NewID(id),
	}

	if d.HasChange("name") {
		input.Name = opslevel.RefOf(d.Get("name").(string))
	}
	if d.HasChange("responsibilities") {
		input.Responsibilities = opslevel.RefOf(d.Get("responsibilities").(string))
	}

	if d.HasChange("parent") {
		if parent := d.Get("parent"); parent != "" {
			input.ParentTeam = opslevel.NewIdentifier(parent.(string))
		} else {
			input.ParentTeam = opslevel.NewIdentifier()
		}
	}

	resource, err := client.UpdateTeam(input)
	if err != nil {
		return err
	}

	if d.HasChange("aliases") {
		err = reconcileTeamAliases(d, resource, client)
		if err != nil {
			return err
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
	err := client.DeleteTeam(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
