package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
	"github.com/shurcooL/graphql"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a group",
		Create:      wrap(resourceGroupCreate),
		Read:        wrap(resourceGroupRead),
		Update:      wrap(resourceGroupUpdate),
		Delete:      wrap(resourceGroupDelete),
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
				Description: "The human-friendly, unique identifier for the group.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the group.",
				ForceNew:    false,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the group.",
				ForceNew:    false,
				Optional:    true,
			},
			"parent": {
				Type:        schema.TypeString,
				Description: "The parent of the group.",
				ForceNew:    false,
				Optional:    true,
			},
			"members": {
				Type:        schema.TypeList,
				Description: "The users who are members of the group.",
				ForceNew:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"teams": {
				Type:        schema.TypeList,
				Description: "The teams where this group is the direct parent.",
				ForceNew:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, client *opslevel.Client) error {
	members := []opslevel.MemberInput{}
	for _, member := range getStringArray(d, "members") {
		members = append(members, opslevel.MemberInput{Email: member})
	}

	teams := []opslevel.IdentifierInput{}
	for _, team := range getStringArray(d, "teams") {
		if opslevel.IsID(team) {
			teams = append(teams, opslevel.IdentifierInput{Id: team})
		} else {
			teams = append(teams, opslevel.IdentifierInput{Alias: graphql.String(team)})
		}
	}

	input := opslevel.GroupInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Parent:      opslevel.NewIdentifier(d.Get("parent").(string)),
		Members:     &members,
		Teams:       &teams,
	}
	resource, err := client.CreateGroup(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceGroupRead(d, client)
}

func resourceGroupRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetGroup(id)
	if err != nil {
		return err
	}

	groupMembers, err := resource.Members(client)
	if err != nil {
		return err
	}
	members := flattenMembersArray(groupMembers)

	descendantTeams, err := resource.DescendantTeams(client)
	if err != nil {
		return err
	}
	teams := flattenTeamsArray(descendantTeams)

	if err := d.Set("alias", resource.Alias); err != nil {
		return err
	}
	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("parent", resource.Parent.Alias); err != nil {
		return err
	}
	if err := d.Set("members", members); err != nil {
		return err
	}
	if err := d.Set("teams", teams); err != nil {
		return err
	}

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.GroupInput{}

	members := []opslevel.MemberInput{}
	for _, member := range getStringArray(d, "members") {
		members = append(members, opslevel.MemberInput{Email: member})
	}

	teams := []opslevel.IdentifierInput{}
	for _, team := range getStringArray(d, "teams") {
		if opslevel.IsID(team) {
			teams = append(teams, opslevel.IdentifierInput{Id: team})
		} else {
			teams = append(teams, opslevel.IdentifierInput{Alias: graphql.String(team)})
		}
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		input.Description = d.Get("description").(string)
	}
	if d.HasChange("parent") {
		input.Parent = opslevel.NewIdentifier(d.Get("parent").(string))
	}
	if d.HasChange("members") {
		input.Members = &members
	}
	if d.HasChange("teams") {
		input.Teams = &teams
	}

	_, err := client.UpdateGroup(d.Id(), input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceGroupRead(d, client)
}

func resourceGroupDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteGroup(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
