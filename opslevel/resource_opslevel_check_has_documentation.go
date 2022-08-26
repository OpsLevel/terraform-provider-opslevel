package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2022"
)

func resourceCheckHasDocumentation() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a has documentation check",
		Create:      wrap(resourceCheckHasDocumentationCreate),
		Read:        wrap(resourceCheckHasDocumentationRead),
		Update:      wrap(resourceCheckHasDocumentationUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"document_type": {
				Type:         schema.TypeString,
				Description:  "The type of the document.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllHasDocumentationTypeEnum(), false),
			},
			"document_subtype": {
				Type:         schema.TypeString,
				Description:  "The subtype of the document.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllHasDocumentationSubtypeEnum(), false),
			},
		}),
	}
}

func resourceCheckHasDocumentationCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckHasDocumentationCreateInput{}
	setCheckCreateInput(d, &input)

	input.DocumentType = opslevel.HasDocumentationTypeEnum(d.Get("document_type").(string))
	input.DocumentSubtype = opslevel.HasDocumentationSubtypeEnum(d.Get("document_subtype").(string))

	resource, err := client.CreateCheckHasDocumentation(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckHasDocumentationRead(d, client)
}

func resourceCheckHasDocumentationRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(id)
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("document_type", string(resource.DocumentType)); err != nil {
		return err
	}
	if err := d.Set("document_subtype", string(resource.DocumentSubtype)); err != nil {
		return err
	}


	return nil
}

func resourceCheckHasDocumentationUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckHasDocumentationUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("document_type") {
		input.DocumentType = opslevel.HasDocumentationTypeEnum(d.Get("document_type").(string))
	}
	if d.HasChange("document_subtype") {
		input.DocumentSubtype = opslevel.HasDocumentationSubtypeEnum(d.Get("document_subtype").(string))
	}

	_, err := client.UpdateCheckHasDocumentation(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckHasDocumentationRead(d, client)
}
