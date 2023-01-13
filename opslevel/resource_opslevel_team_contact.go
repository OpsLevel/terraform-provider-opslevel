package opslevel

import (
	"github.com/hasura/go-graphql-client"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
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
				Description: "The id or alias of the team the contact belongs to.",
				ForceNew:    true,
				Required:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The method of contact [email, slack, slack_handle, web].",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllContactType(), false),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name shown in the UI for the contact.",
				ForceNew:    false,
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The contact value. Examples: support@company.com for type email, https://opslevel.com for type web, #devs for type slack",
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
	d.SetId(string(resource.Id))

	return resourceTeamContactRead(d, client)
}

func resourceTeamContactRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	// Handle Import by spliting the ID into the 2 parts
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		d.Set("team", parts[0])
		id = parts[1]
		d.SetId(id)
	}

	identifier := d.Get("team").(string)
	var err error
	var team *opslevel.Team
	if opslevel.IsID(identifier) {
		team, err = client.GetTeam(*opslevel.NewID(identifier))
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
		if string(t.Id) == id {
			resource = &t
			break
		}
	}
	if resource == nil {
		d.SetId("")
		return nil
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
	id := d.Id()
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

	_, err := client.UpdateContact(graphql.ID(id), input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())

	return resourceTeamContactRead(d, client)
}

func resourceTeamContactDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.RemoveContact(graphql.ID(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
