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
				Description: "A list of human-friendly, unique identifiers for the service",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "A list of tags applied to the service.",
				ForceNew:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func parseTags(d *schema.ResourceData) []string {
	output := []string{}
	for _, entry := range d.Get("tags").([]interface{}) {
		output = append(output, entry.(string))
	}
	return output
}

func reconcileTags(d *schema.ResourceData, service *opslevel.Service, client *opslevel.Client) error {
	tags := parseTags(d)
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
	for i, configuredTag := range tags {
		if stringInArray(configuredTag, existingTags) {
			continue
		}
		// Add
		parts := strings.Split(configuredTag, ":")
		keyName := strings.TrimSpace(parts[0])
		if !tagKeyRegex.MatchString(keyName) {
			return fmt.Errorf("field 'tags.%d.key' == '%s' - %s", i, keyName, tagKeyRegexErrorMsg)
		}
		keyValue := strings.TrimSpace(parts[1])
		_, err := client.CreateTag(opslevel.TagCreateInput{
			Id:    service.Id,
			Key:   keyName,
			Value: keyValue,
		})
		if err != nil {
			return err
		}
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

	// aliases := expandStringArray(d.Get("aliases").([]interface{}))
	// if len(aliases) > 0 {
	// 	for _, alias := range aliases {
	// 		if resource.HasAlias(alias) == false {
	// 			_, aliasErr := client.CreateAlias(opslevel.AliasCreateInput{
	// 				OwnerId: resource.Id,
	// 				Alias:   alias,
	// 			})
	// 			if aliasErr != nil {
	// 				return aliasErr
	// 			}
	// 		}
	// 	}
	// }

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

	// if d.HasChange("aliases") {
	// 	aliases := expandStringArray(d.Get("aliases").([]interface{}))
	// 	doAdd := func(v string) error {
	// 		_, aliasErr := client.CreateAlias(opslevel.AliasCreateInput{
	// 			OwnerId: resource.Id,
	// 			Alias:   v,
	// 		})
	// 		if aliasErr != nil {
	// 			return aliasErr
	// 		}
	// 		return nil
	// 	}
	// 	// TODO: DeleteAlias
	// 	// doDelete := func(v string) error {
	// 	// 	_, aliasErr := client.DeleteAlias(opslevel.AliasDeleteInput{
	// 	// 		OwnerId: resource.Id,
	// 	// 		Alias:   v,
	// 	// 	})
	// 	// 	if aliasErr != nil {
	// 	// 		return aliasErr
	// 	// 	}
	// 	// 	return nil
	// 	// }
	// 	reconcileStringArray(aliases, resource.Aliases, doAdd, nil, nil) // doDelete
	// }

	if d.HasChange("tags") {
		tagsErr := reconcileTags(d, resource, client)
		if tagsErr != nil {
			return tagsErr
		}
	}

	tags := map[string]string{}
	for i, entry := range d.Get("tags").([]interface{}) {
		parts := strings.Split(strings.TrimSpace(entry.(string)), ":")
		key := strings.TrimSpace(parts[0])
		if !tagKeyRegex.MatchString(key) {
			return fmt.Errorf("tag.%d.key '%s' - %s", i, key, tagKeyRegexErrorMsg)
		}
		tags[key] = strings.TrimSpace(parts[1])
	}
	_, tagsErr := client.AssignTagsForId(resource.Id, tags)
	if tagsErr != nil {
		return tagsErr
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
