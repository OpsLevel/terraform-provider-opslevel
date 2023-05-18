package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
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
				Description: "The parent of the group. Accepts only Alias",
				ForceNew:    false,
				Optional:    true,
			},
			"members": {
				Type:        schema.TypeList,
				Description: "List of users' email who are members of the group.",
				ForceNew:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"teams": {
				Type:        schema.TypeList,
				Description: "The teams where this group is the direct parent. Accepts only Alias.",
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
		teams = append(teams, *opslevel.NewIdentifier(team))
	}

	input := opslevel.GroupInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Members:     &members,
		Teams:       &teams,
	}
	if parent, ok := d.GetOk("parent"); ok {
		input.Parent = opslevel.NewIdentifier(parent.(string))
	}
	resource, err := client.CreateGroup(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceGroupRead(d, client)
}

func resourceGroupRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetGroup(opslevel.ID(id))
	if err != nil {
		return err
	}

	groupMembers, err := resource.Members(client, nil)
	if err != nil {
		return err
	}
	members := flattenMembersArray(groupMembers)

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
	if _, ok := d.GetOk("teams"); ok {
		childTeams, err := resource.ChildTeams(client, nil)
		if err != nil {
			return err
		}
		teams := flattenTeamsArray(childTeams)
		if err := d.Set("teams", teams); err != nil {
			return err
		}
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
		teams = append(teams, *opslevel.NewIdentifier(team))
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		input.Description = d.Get("description").(string)
	}

	if parent, ok := d.GetOk("parent"); ok {
		input.Parent = opslevel.NewIdentifier(parent.(string))
	} else {
		input.Parent = nil
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
