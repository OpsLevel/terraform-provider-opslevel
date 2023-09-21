package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a tag.",
		Create:      wrap(resourceTagCreate),
		Read:        wrap(resourceTagRead),
		Update:      wrap(resourceTagUpdate),
		Delete:      wrap(resourceTagDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identifier": {
				Type:        schema.TypeString,
				Description: "The id or human-friendly, unique identifier for the tag.",
				ForceNew:    true,
				Optional:    true,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Description: "The id or human-friendly, unique identifier of the resource this tag belongs to.",
				ForceNew:    true,
				Optional:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The resource type that the tag applies to.",
				ForceNew:    false,
				Optional:    true,
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
		},
	}
}

func resourceTagCreate(d *schema.ResourceData, client *opslevel.Client) error {
	resourceId := d.Get("resource_id").(string)
	resourceType, err := getValidatedTaggableResource(d.Get("type").(string))
	if err != nil {
		return err
	}

	tagCreateInput := opslevel.TagCreateInput{
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
		Type:  resourceType,
	}
	if opslevel.IsID(resourceId) {
		tagCreateInput.Id = opslevel.ID(resourceId)
	} else {
		tagCreateInput.Alias = resourceId
	}

	newTag, _ := client.CreateTag(tagCreateInput)
	if err != nil {
		return err
	}

	d.Set("type", resourceType)
	d.SetId(string(newTag.Id))

	return resourceTagRead(d, client)
}

func resourceTagRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	resource, err := getResource(d.Get("type").(string), id, client)
	if err != nil {
		return err
	}

	if err := d.Set("identifier", id); err != nil {
		return err
	}

	if err := d.Set("resource_id", resource.ResourceId()); err != nil {
		return err
	}

	if err := d.Set("type", resource.ResourceType()); err != nil {
		return err
	}

	tag := resource.GetTag(*opslevel.NewID(id), client)
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
		Id: opslevel.ID(id),
	}

	if d.HasChange("key") {
		input.Key = d.Get("key").(string)
	}

	if d.HasChange("value") {
		input.Value = d.Get("value").(string)
	}

	_, err := client.UpdateTag(input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceTagRead(d, client)
}

func resourceTagDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteTag(opslevel.ID(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func getValidatedTaggableResource(providedResourceType string) (opslevel.TaggableResource, error) {
	lowerCaseResourceType := strings.ToLower(providedResourceType)
	for _, taggableResource := range opslevel.AllTaggableResource {
		if lowerCaseResourceType == strings.ToLower(taggableResource) {
			return opslevel.TaggableResource(taggableResource), nil
		}
	}
	return opslevel.TaggableResource(""), fmt.Errorf("Unknown resource type '%s'", providedResourceType)
}

func getResource(validResourceType string, identifier string, client *opslevel.Client) (opslevel.TaggableResourceInterface, error) {
	var taggableResource opslevel.TaggableResourceInterface

	switch validResourceType {
	case string(opslevel.TaggableResourceService):
		if opslevel.IsID(identifier) {
			client.GetService(opslevel.ID(identifier))
		} else {
			client.GetServiceWithAlias(identifier)
		}
	case string(opslevel.TaggableResourceRepository):
		if opslevel.IsID(identifier) {
			client.GetRepository(opslevel.ID(identifier))
		} else {
			client.GetRepositoryWithAlias(identifier)
		}
	case string(opslevel.TaggableResourceTeam):
		if opslevel.IsID(identifier) {
			client.GetTeam(opslevel.ID(identifier))
		} else {
			client.GetTeamWithAlias(identifier)
		}
	case string(opslevel.TaggableResourceUser):
		client.GetUser(identifier)
	case string(opslevel.TaggableResourceDomain):
		client.GetDomain(identifier)
	case string(opslevel.TaggableResourceSystem):
		client.GetSystem(identifier)
	case string(opslevel.TaggableResourceInfrastructureresource):
		client.GetInfrastructure(identifier)
	}
	return taggableResource, nil
}
