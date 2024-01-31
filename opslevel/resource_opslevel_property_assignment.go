package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2024"
)

func resourcePropertyAssignment() *schema.Resource {
	return &schema.Resource{
		Description: "Manages properties assigned to entities (like Services)",
		Create:      wrap(resourcePropertyAssignmentCreate),
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
				ForceNew:    true,
			},
		},
	}
}

func resourcePropertyAssignmentCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.PropertyInput{
		Owner:      *opslevel.NewIdentifier(d.Get("owner").(string)),
		Definition: *opslevel.NewIdentifier(d.Get("definition").(string)),
		Value:      opslevel.JsonString(d.Get("value").(string)),
	}

	resource, err := client.PropertyAssign(input)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", resource.Owner.Id(), resource.Definition.Id))

	return resourcePropertyAssignmentRead(d, client)
}

func resourcePropertyAssignmentRead(d *schema.ResourceData, client *opslevel.Client) error {
	// an invalid id can be passed in by using 'terraform import', validate the it before attaching the id to the resource
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("[%s] invalid property assignment id, should be in format 'ownerId:definitionId' (only a single colon between both ids, no spaces or special characters)", d.Id())
	}
	ownerId := parts[0]
	definitionId := parts[1]
	if !opslevel.IsID(ownerId) {
		return fmt.Errorf("[%s] invalid ownerId", ownerId)
	}
	if !opslevel.IsID(definitionId) {
		return fmt.Errorf("[%s] invalid definitionId", definitionId)
	}

	resource, err := client.GetProperty(ownerId, definitionId)
	if err != nil {
		return err
	}
	// if resource was fetched correctly, attach the id to the resource
	d.SetId(fmt.Sprintf("%s:%s", ownerId, definitionId))

	if err := d.Set("definition", d.Get("definition")); err != nil {
		return err
	}
	if err := d.Set("owner", d.Get("owner")); err != nil {
		return err
	}
	if err := d.Set("value", string(*resource.Value)); err != nil {
		return err
	}

	return nil
}

func resourcePropertyAssignmentDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := strings.Split(d.Id(), ":")
	if len(id) != 2 {
		return fmt.Errorf("[%s] invalid property assignment id, should be in format 'ownerId:definitionId' (only a single colon between both ids, no spaces or special characters)", d.Id())
	}
	ownerId := id[0]
	definitionId := id[1]
	err := client.PropertyUnassign(ownerId, definitionId)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
