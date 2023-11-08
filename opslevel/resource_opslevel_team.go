package opslevel

import (
	"errors"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/rs/zerolog/log"
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
				Deprecated:  "field 'group' on team is no longer supported please use the 'parent' field.",
				ForceNew:    false,
				Optional:    true,
			},
			"member": {
				Type:        schema.TypeList,
				Description: "List of members in the team with email address and role. At least one member with role 'manager' must be present.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role": {
							Type:     schema.TypeString,
							Required: true,
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

func collectMembersFromTeam(team *opslevel.Team) []map[string]interface{} {
	members := []map[string]interface{}{}

	for _, user := range team.Members.Nodes {
		member := make(map[string]interface{})
		member["email"] = user.Email
		member["role"] = string(user.Role)
		members = append(members, member)
	}
	return members
}

func memberInArray(member map[string]interface{}, array []map[string]interface{}) bool {
	for _, m := range array {
		if m["email"] == member["email"] && m["role"] == member["role"] {
			return true
		}
	}
	return false
}

func reconcileTeamMembership(d *schema.ResourceData, team *opslevel.Team, client *opslevel.Client) error {
	log.Info().Msg("(hello world) getting started")

	var expectedMembers []interface{}

	if members, ok := d.GetOk("member"); ok {
		expectedMembers = members.([]interface{})

		for _, m := range expectedMembers {
			member := m.(map[string]interface{})
			log.Info().Msgf("(read member) email=%v role=%v", m["email"], m["role"])
		}
	}

	return nil
}

func validateMembershipState(d *schema.ResourceData) error {
	// TODO: what to do here?
	if membersSet, ok := d.GetOk("members"); ok {
		if managerEmail, ok := d.GetOk("manager_email"); ok {
			memberEmails := expandStringArray(membersSet.(*schema.Set).List())
			if !stringInArray(managerEmail.(string), memberEmails) {
				return errors.New("The 'manager_email' value is required as a member")
			}
		}
	}
	return nil
}

func resourceTeamCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.TeamCreateInput{
		Name:             d.Get("name").(string),
		ManagerEmail:     d.Get("manager_email").(string),
		Responsibilities: d.Get("responsibilities").(string),
	}
	if _, ok := d.GetOk("group"); ok {
		return errors.New("groups are deprecated - create and update are disabled.")
	}
	if parentTeam, ok := d.GetOk("parent"); ok {
		input.ParentTeam = opslevel.NewIdentifier(parentTeam.(string))
	}

	membershipValidationErr := validateMembershipState(d)
	if membershipValidationErr != nil {
		return membershipValidationErr
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
	if _, ok := d.GetOk("parent"); ok {
		if err := d.Set("parent", resource.ParentTeam.Alias); err != nil {
			return err
		}
	}
	if _, ok := d.GetOk("aliases"); ok {
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
	}

	return nil
}

func resourceTeamUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	input := opslevel.TeamUpdateInput{
		Id: opslevel.ID(id),
	}

	membershipValidationErr := validateMembershipState(d)
	if membershipValidationErr != nil {
		return membershipValidationErr
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

	if d.HasChange("members") {
		// TODO: update this...
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
