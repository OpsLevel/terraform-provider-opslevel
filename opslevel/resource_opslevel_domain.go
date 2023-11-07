package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a domain",
		Create:      wrap(resourceDomainCreate),
		Read:        wrap(resourceDomainRead),
		Update:      wrap(resourceDomainUpdate),
		Delete:      wrap(resourceDomainDelete),
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
				Description: "The aliases of the domain.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name for the domain.",
				ForceNew:    false,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description for the domain.",
				ForceNew:    false,
				Optional:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The team id of the owner for the domain.",
				ForceNew:    false,
				Optional:    true,
			},
			"note": {
				Type:        schema.TypeString,
				Description: "Additional information about the domain.",
				ForceNew:    false,
				Optional:    true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, client *opslevel.Client) error {
	resource, err := client.CreateDomain(opslevel.DomainInput{
		Name:        GetString(d, "name"),
		Description: GetString(d, "description"),
		Owner:       opslevel.NewID(d.Get("owner").(string)),
		Note:        GetString(d, "note"),
	})
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))
	return resourceDomainRead(d, client)
}

func resourceDomainRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetDomain(id)
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
	if err := d.Set("owner", resource.Owner.Id()); err != nil {
		return err
	}
	if err := d.Set("note", resource.Note); err != nil {
		return err
	}

	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	input := opslevel.DomainInput{}

	if d.HasChange("name") {
		input.Name = GetString(d, "name")
	}
	if d.HasChange("description") {
		input.Description = GetString(d, "description")
	}
	if d.HasChange("owner") {
		input.Owner = opslevel.NewID(d.Get("owner").(string))
	}
	if d.HasChange("note") {
		input.Note = GetString(d, "note")
	}

	_, err := client.UpdateDomain(id, input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceDomainRead(d, client)
}

func resourceDomainDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteDomain(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
