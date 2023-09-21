package opslevel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a tag. Only uses 'tag create' API due to reasons",
		Create:      wrap(resourceTagCreate),
		Read:        wrap(resourceTagRead),
		Update:      wrap(resourceTagUpdate),
		Delete:      wrap(resourceTagDelete),
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Description: "The id or human-friendly, unique identifier of the resource this tag belongs to.",
				ForceNew:    true,
				Optional:    true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Description:  "The resource type that the tag applies to.",
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllTaggableResource, false),
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The key of the tag.",
				ForceNew:    false,
				Optional:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The value of the tag.",
				ForceNew:    false,
				Optional:    true,
			},
			"returned_resource": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTagCreate(d *schema.ResourceData, client *opslevel.Client) error {
	resourceType := opslevel.TaggableResource(d.Get("resource_type").(string))

	tagCreateInput := opslevel.TagCreateInput{
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
		Type:  resourceType,
	}

	resourceId := d.Get("resource_id").(string)
	if opslevel.IsID(resourceId) {
		tagCreateInput.Id = opslevel.ID(resourceId)
	} else {
		tagCreateInput.Alias = resourceId
	}
	newTag, err := client.CreateTag(tagCreateInput)
	if err != nil {
		return err
	}

	resource, err := getResource(resourceType, resourceId, client)
	if err != nil {
		return err
	}

	d.SetId(string(newTag.Id))
	d.Set("returned_resource", resource.ResourceId())
	d.Set("resource_type", resourceType)

	return resourceTagRead(d, client)
}

func resourceTagRead(d *schema.ResourceData, client *opslevel.Client) error {
	resourceId := d.Get("returned_resource").(string)
	resourceType := opslevel.TaggableResource(d.Get("resource_type").(string))
	resource, err := getResource(resourceType, resourceId, client)
	if err != nil {
		return err
	}

	id := d.Id()
	tag := resource.GetTag(*opslevel.NewID(id), client)
	if tag == nil {
		return fmt.Errorf(
			"Tag '%s' for resource with id: '%s' of type '%s' not found",
			id,
			resource.ResourceId(),
			resource.ResourceType(),
		)
	}

	if err := d.Set("resource_id", d.Get("resource_id").(string)); err != nil {
		return err
	}

	if err := d.Set("resource_type", string(resource.ResourceType())); err != nil {
		return err
	}

	if err := d.Set("key", tag.Key); err != nil {
		return err
	}

	if err := d.Set("value", tag.Value); err != nil {
		return err
	}

	return nil
}

func resourceTagUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resourceId := d.Get("resource_id").(string)
	resourceType := opslevel.TaggableResource(d.Get("resource_type").(string))
	resource, err := getResource(resourceType, resourceId, client)
	if err != nil {
		return err
	}

	input := opslevel.TagUpdateInput{
		Id:    opslevel.ID(id),
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
	}

	if _, err := client.UpdateTag(input); err != nil {
		return err
	}

	d.Set("returned_resource", resource.ResourceId())
	d.Set("last_updated", timeLastUpdated())
	return resourceTagRead(d, client)
}

func resourceTagDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	if err := client.DeleteTag(opslevel.ID(id)); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func getResource(validResourceType opslevel.TaggableResource, identifier string, client *opslevel.Client) (opslevel.TaggableResourceInterface, error) {
	var err error
	var taggableResource opslevel.TaggableResourceInterface

	switch validResourceType {
	case opslevel.TaggableResourceService:
		if opslevel.IsID(identifier) {
			taggableResource, err = client.GetService(opslevel.ID(identifier))
		} else {
			taggableResource, err = client.GetServiceWithAlias(identifier)
		}
	case opslevel.TaggableResourceRepository:
		if opslevel.IsID(identifier) {
			taggableResource, err = client.GetRepository(opslevel.ID(identifier))
		} else {
			taggableResource, err = client.GetRepositoryWithAlias(identifier)
		}
	case opslevel.TaggableResourceTeam:
		if opslevel.IsID(identifier) {
			taggableResource, err = client.GetTeam(opslevel.ID(identifier))
		} else {
			taggableResource, err = client.GetTeamWithAlias(identifier)
		}
	case opslevel.TaggableResourceDomain:
		taggableResource, err = client.GetDomain(identifier)
	case opslevel.TaggableResourceSystem:
		taggableResource, err = client.GetSystem(identifier)
	}

	if err != nil {
		return nil, err
	}
	return taggableResource, nil
}
