package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourcePropertyAssignment() *schema.Resource {
	return &schema.Resource{
		Description: "Manages properties assigned to entities (like Services)",
		Create:      wrap(resourcePropertyAssignmentCreate),
		Update:      wrap(resourcePropertyAssignmentUpdate),
		Read:        wrap(resourcePropertyAssignmentRead),
		Delete:      wrap(resourcePropertyAssignmentDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"definition": {
				Type:        schema.TypeString,
				Description: "The custom property definition's ID or alias.",
				Required:    true,
				ForceNew:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The ID or alias of the entity that the property has been assigned to.",
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The value of the custom property.",
				Required:    true,
			},
		},
	}
}

func resourcePropertyAssignmentCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.PropertyInput{
		Owner:      *opslevel.NewIdentifier(d.Get("owner").(string)),
		Definition: *opslevel.NewIdentifier(d.Get("definition").(string)),
		Value:      opslevel.JSONString(d.Get("value").(string)),
	}

	resource, err := client.PropertyAssign(input)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", resource.Owner.Id(), resource.Definition.Id))

	return resourcePropertyAssignmentRead(d, client)
}

func resourcePropertyAssignmentUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	// TODO: this is a hack
	// cannot update an existing property assignment, so instead
	// delete the current assignment and create a new one.
	id := strings.Split(d.Id(), ":")
	ownerId := id[0]
	definitionId := id[1]
	valueEncoded := opslevel.JSONString(d.Get("value").(string))
	input := opslevel.PropertyInput{
		Owner:      *opslevel.NewIdentifier(ownerId),
		Definition: *opslevel.NewIdentifier(definitionId),
		Value:      valueEncoded,
	}

	err := client.PropertyUnassign(ownerId, definitionId)
	if err != nil {
		return err
	}
	_, err = client.PropertyAssign(input)
	if err != nil {
		return err
	}

	return resourcePropertyAssignmentRead(d, client)
}

func resourcePropertyAssignmentRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := strings.Split(d.Id(), ":")
	ownerId := id[0]
	definitionId := id[1]

	resource, err := client.GetProperty(ownerId, definitionId)
	if err != nil {
		return err
	}

	if err := d.Set("definition", string(resource.Definition.Id)); err != nil {
		return err
	}
	if err := d.Set("owner", string(resource.Owner.Id())); err != nil {
		return err
	}
	if err := d.Set("value", string(*resource.Value)); err != nil {
		return err
	}

	return nil
}

func resourcePropertyAssignmentDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := strings.Split(d.Id(), ":")
	ownerId := id[0]
	definitionId := id[1]
	err := client.PropertyUnassign(ownerId, definitionId)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
