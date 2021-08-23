package opslevel

import (
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
		},
	}
}

func resourceTeamCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.TeamCreateInput{
		Name:             d.Get("name").(string),
		ManagerEmail:     d.Get("manager_email").(string),
		Responsibilities: d.Get("responsibilities").(string),
	}
	resource, err := client.CreateTeam(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

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

	_, err := client.UpdateTeam(input)
	if err != nil {
		return err
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
