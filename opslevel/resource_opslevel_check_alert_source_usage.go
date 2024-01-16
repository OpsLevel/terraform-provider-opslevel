package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2024"
)

func resourceCheckAlertSourceUsage() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an alert source usage check",
		Create:      wrap(resourceCheckAlertSourceUsageCreate),
		Read:        wrap(resourceCheckAlertSourceUsageRead),
		Update:      wrap(resourceCheckAlertSourceUsageUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"alert_type": {
				Type:         schema.TypeString,
				Description:  "The type of the alert source.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllAlertSourceTypeEnum, false),
			},
			"alert_name_predicate": getPredicateInputSchema(false, DefaultPredicateDescription),
		}),
	}
}

func resourceCheckAlertSourceUsageCreate(d *schema.ResourceData, client *opslevel.Client) error {
	checkCreateInput := getCheckCreateInputFrom(d)
	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckAlertSourceUsageCreateInput](checkCreateInput)
	input.AlertSourceType = opslevel.RefOf(opslevel.AlertSourceTypeEnum(d.Get("alert_type").(string)))
	input.AlertSourceNamePredicate = expandPredicate(d, "alert_name_predicate")

	resource, err := client.CreateCheckAlertSourceUsage(*input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckAlertSourceUsageRead(d, client)
}

func resourceCheckAlertSourceUsageRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("alert_type", string(resource.AlertSourceType)); err != nil {
		return err
	}
	if _, ok := d.GetOk("alert_name_predicate"); ok {
		if err := d.Set("alert_name_predicate", flattenPredicate(&resource.AlertSourceNamePredicate)); err != nil {
			return err
		}
	}
	return nil
}

func resourceCheckAlertSourceUsageUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	checkUpdateInput := getCheckUpdateInputFrom(d)
	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckAlertSourceUsageUpdateInput](checkUpdateInput)

	if d.HasChange("alert_type") {
		input.AlertSourceType = opslevel.RefOf(opslevel.AlertSourceTypeEnum(d.Get("alert_type").(string)))
	}
	if d.HasChange("alert_name_predicate") {
		input.AlertSourceNamePredicate = expandPredicateUpdate(d, "alert_name_predicate")
	}

	_, err := client.UpdateCheckAlertSourceUsage(*input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckAlertSourceUsageRead(d, client)
}
