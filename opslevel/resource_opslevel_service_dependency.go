// package opslevel

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceServiceDependency() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a service dependency",
// 		Create:      wrap(resourceServiceDependencyCreate),
// 		Read:        wrap(resourceServiceDependencyRead),
// 		Delete:      wrap(resourceServiceDependencyDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"service": {
// 				Type:        schema.TypeString,
// 				Description: "The ID or alias of the service with the dependency.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 			"depends_upon": {
// 				Type:        schema.TypeString,
// 				Description: "The ID or alias of the service that is depended upon.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 			"note": {
// 				Type:        schema.TypeString,
// 				Description: "Notes for service dependency.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 		},
// 	}
// }

// func resourceServiceDependencyCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	serviceIdentifier := d.Get("service").(string)
// 	dependsOn := d.Get("depends_upon").(string)

// 	input := opslevel.ServiceDependencyCreateInput{
// 		DependencyKey: opslevel.ServiceDependencyKey{
// 			SourceIdentifier:      opslevel.NewIdentifier(serviceIdentifier),
// 			DestinationIdentifier: opslevel.NewIdentifier(dependsOn),
// 		},
// 		Notes: opslevel.RefOf(d.Get("note").(string)),
// 	}
// 	resource, err := client.CreateServiceDependency(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceServiceDependencyRead(d, client)
// }

// func lookupService(identifier string, client *opslevel.Client) (*opslevel.Service, error) {
// 	if opslevel.IsID(identifier) {
// 		return client.GetService(*opslevel.NewID(identifier))
// 	} else {
// 		return client.GetServiceWithAlias(identifier)
// 	}
// }

// func resourceServiceDependencyRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	service, err := lookupService(d.Get("service").(string), client)
// 	if err != nil {
// 		return err
// 	}

// 	dependencies, err := service.GetDependencies(client, nil)
// 	if err != nil {
// 		return err
// 	}

// 	var resource *opslevel.ServiceDependenciesEdge
// 	for _, edge := range dependencies.Edges {
// 		if string(edge.Id) == id {
// 			resource = &edge
// 			break
// 		}
// 	}
// 	if resource == nil {
// 		d.SetId("")
// 		return nil
// 	}

// 	if err := d.Set("note", resource.Notes); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceServiceDependencyDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteServiceDependency(*opslevel.NewID(id))
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
