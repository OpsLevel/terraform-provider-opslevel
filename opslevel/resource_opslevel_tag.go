package opslevel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a tag. Only uses API's 'tag create', not 'tag assign'.",
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
			"resource_identifier": {
				Type:        schema.TypeString,
				Description: "The id or human-friendly, unique identifier of the resource this tag belongs to.",
				ForceNew:    true,
				Required:    true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Description:  "The resource type that the tag applies to.",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllTaggableResource, false),
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The key of the tag.",
				ForceNew:    false,
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The value of the tag.",
				ForceNew:    false,
				Required:    true,
			},
			"returned_resource": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTagCreate(d *schema.ResourceData, client *opslevel.Client) error {
	resourceId := d.Get("resource_identifier").(string)
	resourceType := opslevel.TaggableResource(d.Get("resource_type").(string))
	resource, err := client.GetTaggableResource(resourceType, resourceId)
	if err != nil {
		return err
	}

	tagCreateInput := opslevel.TagCreateInput{
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
		Type:  opslevel.RefOf(resourceType),
	}

	if opslevel.IsID(resourceId) {
		tagCreateInput.Id = opslevel.NewID(resourceId)
	} else {
		tagCreateInput.Alias = opslevel.RefOf(resourceId)
	}
	newTag, err := client.CreateTag(tagCreateInput)
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
	resource, err := client.GetTaggableResource(resourceType, resourceId)
	if err != nil {
		return err
	}

	tags, err := resource.GetTags(client, nil)
	if err != nil {
		return fmt.Errorf("Unable to get tags from '%s' with id '%s'", resourceType, resourceId)
	}
	id := d.Id()
	tag, err := tags.GetTagById(*opslevel.NewID(id))
	if err != nil || tag == nil {
		return fmt.Errorf(
			"Tag '%s' for type %s with id '%s' not found. %s",
			id,
			resource.ResourceType(),
			resource.ResourceId(),
			err,
		)
	}

	if err := d.Set("resource_identifier", d.Get("resource_identifier").(string)); err != nil {
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

	input := opslevel.TagUpdateInput{
		Id:    opslevel.ID(id),
		Key:   opslevel.RefOf(d.Get("key").(string)),
		Value: opslevel.RefOf(d.Get("value").(string)),
	}
	if _, err := client.UpdateTag(input); err != nil {
		return err
	}

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
