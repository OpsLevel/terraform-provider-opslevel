package opslevel

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/rs/zerolog/log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2024"
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
			"owner": {
				Type:        schema.TypeString,
				Description: "The team that owns the service. ID or Alias my be used.",
				ForceNew:    false,
				Optional:    true,
			},
			"lifecycle_alias": {
				Type:        schema.TypeString,
				Description: "The lifecycle stage of the service.",
				ForceNew:    false,
				Optional:    true,
			},
			"api_document_path": {
				Type:        schema.TypeString,
				Description: "The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path.",
				ForceNew:    false,
				Optional:    true,
			},
			"preferred_api_document_source": {
				Type:         schema.TypeString,
				Description:  "The API document source (PUSH or PULL) used to determine the displayed document. If null, we use the order push and then pull.",
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllApiDocumentSourceEnum, false),
			},
			"aliases": {
				Type:        schema.TypeList,
				Description: "A list of human-friendly, unique identifiers for the service.",
				ForceNew:    false,
				Optional:    true,
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

func reconcileServiceAliases(d *schema.ResourceData, service *opslevel.Service, client *opslevel.Client) error {
	expectedAliases := getStringArray(d, "aliases")
	existingAliases := service.ManagedAliases
	for _, existingAlias := range existingAliases {
		if !slices.Contains(expectedAliases, existingAlias) {
			err := client.DeleteServiceAlias(existingAlias)
			if err != nil {
				return err
			}
		}
	}
	for _, expectedAlias := range expectedAliases {
		if !slices.Contains(existingAliases, expectedAlias) {
			_, err := client.CreateAliases(service.Id, []string{expectedAlias})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func reconcileTags(d *schema.ResourceData, service *opslevel.Service, client *opslevel.Client) error {
	tags := getStringArray(d, "tags")
	existingTags := make([]string, 0)
	for _, tag := range service.Tags.Nodes {
		flattenedTag := flattenTag(tag)
		existingTags = append(existingTags, flattenedTag)
		if !slices.Contains(tags, flattenedTag) {
			err := client.DeleteTag(tag.Id)
			if err != nil {
				return err
			}
		}
	}
	tagInput := map[string]string{}
	for _, tag := range tags {
		parts := strings.Split(tag, ":")
		if len(parts) != 2 {
			return fmt.Errorf("[%s] invalid tag, should be in format 'key:value' (only a single colon between the key and value, no spaces or special characters)", tag)
		}
		key := parts[0]
		value := parts[1]
		tagInput[key] = value
	}
	_, err := client.AssignTags(string(service.Id), tagInput)
	if err != nil {
		return err
	}

	return nil
}

func resourceServiceCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.ServiceCreateInput{
		Name:           d.Get("name").(string),
		Product:        opslevel.RefOf(d.Get("product").(string)),
		Description:    opslevel.RefOf(d.Get("description").(string)),
		Language:       opslevel.RefOf(d.Get("language").(string)),
		Framework:      opslevel.RefOf(d.Get("framework").(string)),
		TierAlias:      opslevel.RefOf(d.Get("tier_alias").(string)),
		LifecycleAlias: opslevel.RefOf(d.Get("lifecycle_alias").(string)),
	}
	if owner := d.Get("owner"); owner != "" {
		input.OwnerInput = opslevel.NewIdentifier(owner.(string))
	}

	resource, err := client.CreateService(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	err = reconcileServiceAliases(d, resource, client)
	if err != nil {
		return err
	}

	err = reconcileTags(d, resource, client)
	if err != nil {
		return err
	}

	docPath, ok1 := d.GetOk("api_document_path")
	docSource, ok2 := d.GetOk("preferred_api_document_source")
	if ok1 || ok2 {
		var source *opslevel.ApiDocumentSourceEnum = nil
		if ok2 {
			s := opslevel.ApiDocumentSourceEnum(docSource.(string))
			source = &s
		}
		_, err := client.ServiceApiDocSettingsUpdate(string(resource.Id), docPath.(string), source)
		if err != nil {
			log.Error().Err(err).Msgf("failed to update service '%s' api doc settings", resource.Aliases[0])
		}
	}

	return resourceServiceRead(d, client)
}

func resourceServiceRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetService(opslevel.ID(id))
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

	// only read in changes to optional fields if they have been set before
	// this will prevent HasChange() from detecting changes on update
	if owner, ok := d.GetOk("owner"); ok || owner != "" {
		var ownerValue string
		if opslevel.IsID(owner.(string)) {
			ownerValue = string(resource.Owner.Id)
		} else {
			ownerValue = string(resource.Owner.Alias)
		}

		if err := d.Set("owner", ownerValue); err != nil {
			return err
		}
	}

	if err := d.Set("lifecycle_alias", resource.Lifecycle.Alias); err != nil {
		return err
	}

	if err := d.Set("aliases", resource.ManagedAliases); err != nil {
		return err
	}
	if err := d.Set("tags", flattenTagArray(resource.Tags.Nodes)); err != nil {
		return err
	}

	if err := d.Set("api_document_path", resource.ApiDocumentPath); err != nil {
		return err
	}
	if err := d.Set("preferred_api_document_source", resource.PreferredApiDocumentSource); err != nil {
		return err
	}

	return nil
}

func resourceServiceUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	input := opslevel.ServiceUpdateInput{
		Id: opslevel.NewID(id),
	}

	if d.HasChange("name") {
		input.Name = opslevel.RefOf(d.Get("name").(string))
	}
	if d.HasChange("product") {
		input.Product = opslevel.RefOf(d.Get("product").(string))
	}
	if d.HasChange("description") {
		input.Description = opslevel.RefOf(d.Get("description").(string))
	}
	if d.HasChange("language") {
		input.Language = opslevel.RefOf(d.Get("language").(string))
	}
	if d.HasChange("framework") {
		input.Framework = opslevel.RefOf(d.Get("framework").(string))
	}
	if d.HasChange("tier_alias") {
		input.TierAlias = opslevel.RefOf(d.Get("tier_alias").(string))
	}
	if d.HasChange("owner") {
		if owner := d.Get("owner"); owner != "" {
			input.OwnerInput = opslevel.NewIdentifier(owner.(string))
		} else {
			input.OwnerInput = opslevel.NewIdentifier()
		}
	}
	if d.HasChange("lifecycle_alias") {
		input.LifecycleAlias = opslevel.RefOf(d.Get("lifecycle_alias").(string))
	}

	resource, err := client.UpdateService(input)
	if err != nil {
		return err
	}

	if d.HasChange("aliases") {
		err = reconcileServiceAliases(d, resource, client)
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		tagsErr := reconcileTags(d, resource, client)
		if tagsErr != nil {
			return tagsErr
		}
	}

	if d.HasChange("api_document_path") || d.HasChange("preferred_api_document_source") {
		var docPath string
		var docSource *opslevel.ApiDocumentSourceEnum
		if value, ok := d.GetOk("api_document_path"); ok {
			docPath = value.(string)
		} else {
			docPath = ""
		}
		if value, ok := d.GetOk("preferred_api_document_source"); ok {
			s := opslevel.ApiDocumentSourceEnum(value.(string))
			docSource = &s
		} else {
			docSource = nil
		}
		_, err := client.ServiceApiDocSettingsUpdate(string(resource.Id), docPath, docSource)
		if err != nil {
			log.Error().Err(err).Msgf("failed to update service '%s' api doc settings", resource.Aliases[0])
		}
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceServiceRead(d, client)
}

func resourceServiceDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteService(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
