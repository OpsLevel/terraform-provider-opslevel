package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceSystem() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a system",
		Create:      wrap(resourceSystemCreate),
		Read:        wrap(resourceSystemRead),
		Update:      wrap(resourceSystemUpdate),
		Delete:      wrap(resourceSystemDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "The aliases of the system.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name for the system.",
				ForceNew:    false,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description for the system.",
				ForceNew:    false,
				Optional:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The id of the owner for the system.  Can be a team or group",
				ForceNew:    false,
				Optional:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The id of the parent domain this system is a child for.",
				ForceNew:    false,
				Optional:    true,
			},
			"note": {
				Type:        schema.TypeString,
				Description: "Additional information about the system.",
				ForceNew:    false,
				Optional:    true,
			},
		},
	}
}

func resourceSystemCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.SystemInput{
		Name:        GetString(d, "name"),
		Description: GetString(d, "description"),
		Note:        GetString(d, "note"),
	}
	if owner := d.Get("owner"); owner != "" {
		input.Owner = opslevel.NewID(owner.(string))
	}
	if domain := d.Get("domain"); domain != "" {
		input.Parent = opslevel.NewIdentifier(domain.(string))
	}

	resource, err := client.CreateSystem(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))
	return resourceSystemRead(d, client)
}

func resourceSystemRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetSystem(id)
	if err != nil {
		return err
	}

	if err := d.Set("aliases", resource.Aliases); err != nil {
		return err
	}
	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("note", resource.Note); err != nil {
		return err
	}

	// only read in changes to optional fields if they have been set before
	if owner, ok := d.GetOk("owner"); ok {
		var ownerValue string
		if opslevel.IsID(owner.(string)) {
			ownerValue = string(resource.Owner.Id())
		} else {
			ownerValue = string(resource.Owner.Alias())
		}

		if err := d.Set("owner", ownerValue); err != nil {
			return err
		}
	}
	if _, ok := d.GetOk("domain"); ok {
		if err := d.Set("domain", resource.Parent.Id); err != nil {
			return err
		}
	}

	return nil
}

func resourceSystemUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	input := opslevel.SystemInput{
		Name:        GetString(d, "name"),
		Description: GetString(d, "description"),
		Note:        GetString(d, "note"),
	}

	if d.HasChange("owner") {
		if owner := d.Get("owner"); owner != "" {
			input.Owner = opslevel.NewID(owner.(string))
		} else {
			// TODO: how can we unset an ID?
			input.Owner = nil
		}
	}
	if d.HasChange("domain") {
		if domain := d.Get("domain"); domain != "" {
			input.Parent = opslevel.NewIdentifier(domain.(string))
		} else {
			input.Parent = opslevel.EmptyIdentifier()
		}
	}

	_, err := client.UpdateSystem(id, input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceSystemRead(d, client)
}

func resourceSystemDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteSystem(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
