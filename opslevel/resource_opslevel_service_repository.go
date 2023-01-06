package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
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
	d.SetId(string(resource.Id))

	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("base_directory", resource.BaseDirectory); err != nil {
		return err
	}

	return nil
}

func resourceServiceRepositoryRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	// Handle Import by spliting the ID into the 2 parts
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		d.Set("service", parts[0])
		id = parts[1]
		d.SetId(id)
	}

	service, err := findService("service_alias", "service", d, client)
	if err != nil {
		return err
	}

	var resource *opslevel.ServiceRepository
	for _, edge := range service.Repositories.Edges {
		for _, repository := range edge.ServiceRepositories {
			if string(repository.Id) == id {
				resource = &repository
				break
			}
		}
		if resource != nil {
			break
		}
	}
	if resource == nil {
		return fmt.Errorf("unable to find service repository with id '%s' on service '%s'", id, service.Aliases[0])
	}

	if err := d.Set("name", resource.DisplayName); err != nil {
		return err
	}
	if err := d.Set("repository", resource.Repository.Id); err != nil {
		return err
	}
	if err := d.Set("base_directory", resource.BaseDirectory); err != nil {
		return err
	}

	return nil
}

func resourceServiceRepositoryUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.ServiceRepositoryUpdateInput{
		Id: *opslevel.NewID(d.Id()),
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
	err := client.DeleteServiceRepository(*opslevel.NewID(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
