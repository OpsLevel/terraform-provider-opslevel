package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2024"
)

func resourceCheckServiceProperty() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service property check.",
		Create:      wrap(resourceCheckServicePropertyCreate),
		Read:        wrap(resourceCheckServicePropertyRead),
		Update:      wrap(resourceCheckServicePropertyUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"property": {
				Type:         schema.TypeString,
				Description:  "The property of the service that the check will verify.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllServicePropertyTypeEnum, false),
			},
			"predicate": getPredicateInputSchema(false, DefaultPredicateDescription),
		}),
	}
}

func resourceCheckServicePropertyCreate(d *schema.ResourceData, client *opslevel.Client) error {
	checkCreateInput := getCheckCreateInputFrom(d)
	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckServicePropertyCreateInput](checkCreateInput)
	input.ServiceProperty = opslevel.ServicePropertyTypeEnum(d.Get("property").(string))
	input.PropertyValuePredicate = expandPredicate(d, "predicate")

	resource, err := client.CreateCheckServiceProperty(*input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckServicePropertyRead(d, client)
}

func resourceCheckServicePropertyRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("property", string(resource.Property)); err != nil {
		return err
	}
	if err := d.Set("predicate", flattenPredicate(resource.Predicate)); err != nil {
		return err
	}

	return nil
}

func resourceCheckServicePropertyUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	checkUpdateInput := getCheckUpdateInputFrom(d)
	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckServicePropertyUpdateInput](checkUpdateInput)

	if d.HasChange("property") {
		input.ServiceProperty = opslevel.RefOf(opslevel.ServicePropertyTypeEnum(d.Get("property").(string)))
	}
	if d.HasChange("predicate") {
		input.PropertyValuePredicate = expandPredicateUpdate(d, "predicate")
	}

	_, err := client.UpdateCheckServiceProperty(*input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckServicePropertyRead(d, client)
}
