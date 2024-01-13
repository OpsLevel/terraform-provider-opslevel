package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceTeamTag() *schema.Resource {
	return &schema.Resource{
		Description:        "Manages a team tag",
		DeprecationMessage: "This resource is deprecated. Please use `opslevel_tag` instead.",
		Create:             wrap(resourceTeamTagCreate),
		Read:               wrap(resourceTeamTagRead),
		Update:             wrap(resourceTeamTagUpdate),
		Delete:             wrap(resourceTeamTagDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"team": {
				Type:        schema.TypeString,
				Description: "The id of the team that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"team_alias": {
				Type:        schema.TypeString,
				Description: "The alias of the team that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"key": {
				Type:         schema.TypeString,
				Description:  "The tag's key.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringMatch(TagKeyRegex, TagKeyErrorMsg),
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The tag's value.",
				ForceNew:    false,
				Required:    true,
			},
		},
	}
}

func resourceTeamTagCreate(d *schema.ResourceData, client *opslevel.Client) error {
	team, err := findTeam("team_alias", "team", d, client)
	if err != nil {
		return err
	}

	input := opslevel.TagCreateInput{
		Id: opslevel.RefOf(team.Id),

		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
	}
	resource, err := client.CreateTag(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	if err := d.Set("key", resource.Key); err != nil {
		return err
	}
	if err := d.Set("value", resource.Value); err != nil {
		return err
	}

	return nil
}

func resourceTeamTagRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	// Handle Import by spliting the ID into the 2 parts
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		d.Set("team", parts[0])
		id = parts[1]
		d.SetId(id)
	}

	team, err := findTeam("team_alias", "team", d, client)
	if err != nil {
		return err
	}

	var resource *opslevel.Tag
	for _, t := range team.Tags.Nodes {
		if string(t.Id) == id {
			resource = &t
			break
		}
	}
	if resource == nil {
		return fmt.Errorf("unable to find tag with id '%s' on team '%s'", id, team.Aliases[0])
	}

	if err := d.Set("key", resource.Key); err != nil {
		return err
	}
	if err := d.Set("value", resource.Value); err != nil {
		return err
	}

	return nil
}

func resourceTeamTagUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	input := opslevel.TagUpdateInput{
		Id: opslevel.ID(id),
	}

	if d.HasChange("key") {
		input.Key = opslevel.RefOf(d.Get("key").(string))
	}
	if d.HasChange("value") {
		input.Value = opslevel.RefOf(d.Get("value").(string))
	}

	resource, err := client.UpdateTag(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())

	if err := d.Set("key", resource.Key); err != nil {
		return err
	}
	if err := d.Set("value", resource.Value); err != nil {
		return err
	}
	return nil
}

func resourceTeamTagDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteTag(opslevel.ID(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
