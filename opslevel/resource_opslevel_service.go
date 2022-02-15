package opslevel

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service",
		Create:      wrap(resourceServiceCreate),
		Read:        wrap(resourceServiceRead),
		Update:      wrap(resourceServiceUpdate),
		Delete:      wrap(resourceServiceDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the service.",
				ForceNew:    false,
				Required:    true,
			},
			"product": {
				Type:        schema.TypeString,
				Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
				ForceNew:    false,
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A brief description of the service.",
				ForceNew:    false,
				Optional:    true,
			},
			"language": {
				Type:        schema.TypeString,
				Description: "The primary programming language that the service is written in.",
				ForceNew:    false,
				Optional:    true,
			},
			"framework": {
				Type:        schema.TypeString,
				Description: "The primary software development framework that the service uses.",
				ForceNew:    false,
				Optional:    true,
			},
			"tier_alias": {
				Type:        schema.TypeString,
				Description: "The software tier that the service belongs to.",
				ForceNew:    false,
				Optional:    true,
			},
			"owner_alias": {
				Type:        schema.TypeString,
				Description: "The team that owns the service.",
				ForceNew:    false,
				Optional:    true,
			},
			"lifecycle_alias": {
				Type:        schema.TypeString,
				Description: "The lifecycle stage of the service.",
				ForceNew:    false,
				Optional:    true,
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "A list of human-friendly, unique identifiers for the service.",
				ForceNew:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:         schema.TypeList,
				Description:  "A list of tags applied to the service.",
				ForceNew:     false,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ValidateFunc: validateServiceTags,
			},
		},
	}
}

func validateServiceTags(i interface{}, k string) (warnings []string, errors []error) {
	data, ok := i.([]string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}
	for _, item := range data {
		key := strings.TrimSpace(strings.Split(item, ":")[0])
		if ok := TagKeyRegex.MatchString(key); !ok {
			return nil, []error{fmt.Errorf("'%s' - %s", key, TagKeyErrorMsg)}
		}
	}
	return nil, nil
}

func reconcileServiceAliases(d *schema.ResourceData, service *opslevel.Service, client *opslevel.Client) error {
	expectedAliases := getStringArray(d, "aliases")
	existingAliases := service.Aliases
	for _, existingAlias := range existingAliases {
		if stringInArray(existingAlias, expectedAliases) { continue }
		// Delete
		err := client.DeleteServiceAlias(existingAlias)
		if err != nil {
			return err
		}
	}
	for _, expectedAlias := range expectedAliases {
		if stringInArray(expectedAlias, existingAliases) { continue }
		// Add
		_, err := client.CreateAliases(service.Id, []string{expectedAlias})
		if err != nil {
			return err
		}
	}
	return nil
}

func reconcileTags(d *schema.ResourceData, service *opslevel.Service, client *opslevel.Client) error {
	tags := getStringArray(d, "tags")
	existingTags := []string{}
	for _, tag := range service.Tags.Nodes {
		flattenedTag := flattenTag(tag)
		existingTags = append(existingTags, flattenedTag)
		if stringInArray(flattenedTag, tags) {
			// Update
			continue
		}
		// Delete
		err := client.DeleteTag(tag.Id)
		if err != nil {
			return err
		}
	}
	tagInput := []opslevel.TagInput{}
	for _, tag := range tags {
		tagInput = append(tagInput, opslevel.TagInput{
			Key:   strings.TrimSpace(strings.Split(tag, ":")[0]),
			Value: strings.TrimSpace(strings.Split(tag, ":")[1]),
		})
	}
	_, err := client.AssignTags(opslevel.TagAssignInput{
		Id:   service.Id,
		Tags: tagInput,
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceServiceCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.ServiceCreateInput{
		Name:        d.Get("name").(string),
		Product:     d.Get("product").(string),
		Description: d.Get("description").(string),
		Language:    d.Get("language").(string),
		Framework:   d.Get("framework").(string),
		Tier:        d.Get("tier_alias").(string),
		Owner:       d.Get("owner_alias").(string),
		Lifecycle:   d.Get("lifecycle_alias").(string),
	}
	resource, err := client.CreateService(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	aliasesErr := reconcileServiceAliases(d, resource, client)
	if aliasesErr != nil {
		return aliasesErr
	}

	tagsErr := reconcileTags(d, resource, client)
	if tagsErr != nil {
		return tagsErr
	}

	return resourceServiceRead(d, client)
}

func resourceServiceRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetService(id)
	if err != nil {
		return err
	}

	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("product", resource.Product); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("language", resource.Language); err != nil {
		return err
	}
	if err := d.Set("framework", resource.Framework); err != nil {
		return err
	}
	if err := d.Set("tier_alias", resource.Tier.Alias); err != nil {
		return err
	}
	if err := d.Set("owner_alias", resource.Owner.Alias); err != nil {
		return err
	}
	if err := d.Set("lifecycle_alias", resource.Lifecycle.Alias); err != nil {
		return err
	}

	if err := d.Set("aliases", resource.Aliases); err != nil {
		return err
	}
	if err := d.Set("tags", flattenTagArray(resource.Tags.Nodes)); err != nil {
		return err
	}

	return nil
}

func resourceServiceUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	input := opslevel.ServiceUpdateInput{
		Id: id,
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("product") {
		input.Product = d.Get("product").(string)
	}
	if d.HasChange("description") {
		input.Description = d.Get("description").(string)
	}
	if d.HasChange("language") {
		input.Language = d.Get("language").(string)
	}
	if d.HasChange("framework") {
		input.Framework = d.Get("framework").(string)
	}
	if d.HasChange("tier_alias") {
		input.Tier = d.Get("tier_alias").(string)
	}
	if d.HasChange("owner_alias") {
		input.Owner = d.Get("owner_alias").(string)
	}
	if d.HasChange("lifecycle_alias") {
		input.Lifecycle = d.Get("lifecycle_alias").(string)
	}

	resource, err := client.UpdateService(input)
	if err != nil {
		return err
	}

	if d.HasChange("aliases") {
		tagsErr := reconcileServiceAliases(d, resource, client)
		if tagsErr != nil {
			return tagsErr
		}
	}

	if d.HasChange("tags") {
		tagsErr := reconcileTags(d, resource, client)
		if tagsErr != nil {
			return tagsErr
		}
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceServiceRead(d, client)
}

func resourceServiceDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteService(opslevel.ServiceDeleteInput{
		Id: id,
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
