package opslevel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func resourceServiceRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service repository",
		Create:      wrap(resourceServiceRepositoryCreate),
		Read:        wrap(resourceServiceRepositoryRead),
		Update:      wrap(resourceServiceRepositoryUpdate),
		Delete:      wrap(resourceServiceRepositoryDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service": {
				Type:        schema.TypeString,
				Description: "The id of the service that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"service_alias": {
				Type:        schema.TypeString,
				Description: "The alias of the service that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "The id of the service that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"repository_alias": {
				Type:        schema.TypeString,
				Description: "The alias of the service that this will be added to.",
				ForceNew:    true,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name displayed in the UI for the service repository.",
				ForceNew:    false,
				Optional:    true,
			},
			"base_directory": {
				Type:        schema.TypeString,
				Description: "The directory in the repository containing opslevel.yml.",
				ForceNew:    false,
				Optional:    true,
			},
		},
	}
}

func resourceServiceRepositoryCreate(d *schema.ResourceData, client *opslevel.Client) error {
	service, err := findService("service_alias", "service", d, client)
	if err != nil {
		return err
	}
	repository, err := findRepository("repository_alias", "repository", d, client)
	if err != nil {
		return err
	}

	input := opslevel.ServiceRepositoryCreateInput{
		Service:    opslevel.IdentifierInput{Id: service.Id},
		Repository: opslevel.IdentifierInput{Id: repository.Id},

		DisplayName:   d.Get("name").(string),
		BaseDirectory: d.Get("base_directory").(string),
	}
	resource, err := client.CreateServiceRepository(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("base_directory", resource.BaseDirectory); err != nil {
		return err
	}

	return nil
}

func resourceServiceRepositoryRead(d *schema.ResourceData, client *opslevel.Client) error {
	service, err := findService("service_alias", "service", d, client)
	if err != nil {
		return err
	}

	id := d.Id()
	var resource *opslevel.ServiceRepository
	for _, edge := range service.Repositories.Edges {
		for _, repository := range edge.ServiceRepositories {
			if repository.Id == id {
				resource = &repository
				break
			}
		}
		if resource != nil {
			break
		}
	}
	if resource == nil {
		return fmt.Errorf("Unable to find service repository with id '%s' on service '%s'", id, service.Aliases[0])
	}

	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("base_directory", resource.BaseDirectory); err != nil {
		return err
	}

	return nil
}

func resourceServiceRepositoryUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.ServiceRepositoryUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("name") {
		input.DisplayName = d.Get("name").(string)
	}
	if d.HasChange("base_directory") {
		input.BaseDirectory = d.Get("base_directory").(string)
	}

	resource, err := client.UpdateServiceRepository(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())

	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("base_directory", resource.BaseDirectory); err != nil {
		return err
	}
	return nil
}

func resourceServiceRepositoryDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteServiceRepository(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
