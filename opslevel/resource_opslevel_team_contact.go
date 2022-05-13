package opslevel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go"
)

func resourceTeamContact() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a team contact",
		Create:      wrap(resourceTeamContactCreate),
		Read:        wrap(resourceTeamContactRead),
		Update:      wrap(resourceTeamContactUpdate),
		Delete:      wrap(resourceTeamContactDelete),
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
				Description: "The id or alias of the team to add the contact to.",
				ForceNew:    true,
				Required:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The type of contact",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllContactType(), false),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the contact.",
				ForceNew:    false,
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "A value of the contact.",
				ForceNew:    false,
				Required:    true,
			},
		},
	}
}

func resourceTeamContactCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.ContactInput{
		Type:        opslevel.ContactType(d.Get("type").(string)),
		DisplayName: d.Get("name").(string),
		Address:     d.Get("value").(string),
	}
	resource, err := client.AddContact(d.Get("team").(string), input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceTeamContactRead(d, client)
}

func resourceTeamContactRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	identifier := d.Get("team").(string)
	var err error
	var team *opslevel.Team
	if opslevel.IsID(identifier) {
		team, err = client.GetTeam(opslevel.NewID(identifier))
		if err != nil {
			return err
		}
	} else {
		team, err = client.GetTeamWithAlias(identifier)
		if err != nil {
			return err
		}
	}

	var resource *opslevel.Contact
	for _, t := range team.Contacts {
		if t.Id == id {
			resource = &t
			break
		}
	}
	if resource == nil {
		return fmt.Errorf("unable to find contact with id '%s' on team '%s'", id, team.Aliases[0])
	}

	if err := d.Set("type", resource.Type); err != nil {
		return err
	}
	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("value", resource.Address); err != nil {
		return err
	}

	return nil
}

func resourceTeamContactUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.ContactInput{}

	if d.HasChange("type") {
		input.Type = opslevel.ContactType(d.Get("type").(string))
	}
	if d.HasChange("name") {
		input.DisplayName = d.Get("name").(string)
	}
	if d.HasChange("value") {
		input.Address = d.Get("value").(string)
	}

	_, err := client.UpdateContact(d.Id(), input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())

	return resourceTeamContactRead(d, client)
}

func resourceTeamContactDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.RemoveContact(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
