package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2022"
)

func resourceCheckServiceOwnership() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service ownership check.",
		Create:      wrap(resourceCheckServiceOwnershipCreate),
		Read:        wrap(resourceCheckServiceOwnershipRead),
		Update:      wrap(resourceCheckServiceOwnershipUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"require_contact_method": {
				Type:        schema.TypeBool,
				Description: "True if a service's owner must have a contact method, False otherwise.",
				ForceNew:    false,
				Optional:    true,
			},
			"contact_method": {
				Type:         schema.TypeString,
				Description:  "The type of contact method that is required.",
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllServiceOwnershipCheckContactType(), true),
			},
		}),
	}
}

func resourceCheckServiceOwnershipCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServiceOwnershipCreateInput{}
	setCheckCreateInput(d, &input)

	if requireContactMethod, ok := d.GetOk("require_contact_method"); ok {
		input.RequireContactMethod = opslevel.Bool(requireContactMethod.(bool))
	}
	if value, ok := d.GetOk("contact_method"); ok {
		contactMethod := opslevel.ServiceOwnershipCheckContactType(value.(string))
		input.ContactMethod = &contactMethod
	}

	resource, err := client.CreateCheckServiceOwnership(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckServiceOwnershipRead(d, client)
}

func resourceCheckServiceOwnershipRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(id)
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}

	if _, ok := d.GetOk("require_contact_method"); ok {
		if err := d.Set("require_contact_method", resource.RequireContactMethod); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("contact_method"); ok {
		if err := d.Set("contact_method", resource.ContactMethod); err != nil {
			return err
		}
	}

	return nil
}

func resourceCheckServiceOwnershipUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckServiceOwnershipUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("require_contact_method") {
		input.RequireContactMethod = opslevel.Bool(d.Get("require_contact_method").(bool))
	}

	if d.HasChange("contact_method") {
		contactMethod := opslevel.ServiceOwnershipCheckContactType(d.Get("contact_method").(string))

		if contactMethod == "" {
			contactMethod = opslevel.ServiceOwnershipCheckContactType("ANY")
		}
		input.ContactMethod = &contactMethod
	}

	_, err := client.UpdateCheckServiceOwnership(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckServiceOwnershipRead(d, client)
}
