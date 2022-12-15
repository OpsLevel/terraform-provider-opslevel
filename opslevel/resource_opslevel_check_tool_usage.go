package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2022"
)

func resourceCheckToolUsage() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a tool usage check",
		Create:      wrap(resourceCheckToolUsageCreate),
		Read:        wrap(resourceCheckToolUsageRead),
		Update:      wrap(resourceCheckToolUsageUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"tool_category": {
				Type:         schema.TypeString,
				Description:  "The category that the tool belongs to.",
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice(opslevel.AllToolCategory(), false),
			},
			"tool_name_predicate":   getPredicateInputSchema(false, DefaultPredicateDescription),
			"tool_url_predicate":    getPredicateInputSchema(false, DefaultPredicateDescription),
			"environment_predicate": getPredicateInputSchema(false, DefaultPredicateDescription),
		}),
	}
}

func resourceCheckToolUsageCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckToolUsageCreateInput{}
	setCheckCreateInput(d, &input)

	input.ToolCategory = opslevel.ToolCategory(d.Get("tool_category").(string))
	input.ToolNamePredicate = expandPredicate(d, "tool_name_predicate")
	input.ToolUrlPredicate = expandPredicate(d, "tool_url_predicate")
	input.EnvironmentPredicate = expandPredicate(d, "environment_predicate")

	resource, err := client.CreateCheckToolUsage(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckToolUsageRead(d, client)
}

func resourceCheckToolUsageRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(id)
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("tool_category", string(resource.ToolCategory)); err != nil {
		return err
	}
	if err := d.Set("tool_name_predicate", flattenPredicate(resource.ToolNamePredicate)); err != nil {
		return err
	}
	if err := d.Set("tool_url_predicate", flattenPredicate(resource.ToolUrlPredicate)); err != nil {
		return err
	}
	if err := d.Set("environment_predicate", flattenPredicate(resource.EnvironmentPredicate)); err != nil {
		return err
	}

	return nil
}

func resourceCheckToolUsageUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckToolUsageUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("tool_category") {
		input.ToolCategory = opslevel.ToolCategory(d.Get("tool_category").(string))
	}
	if d.HasChange("tool_name_predicate") {
		input.ToolNamePredicate = expandPredicate(d, "tool_name_predicate")
	}
	if d.HasChange("tool_url_predicate") {
		input.ToolUrlPredicate = expandPredicate(d, "tool_url_predicate")
	}
	if d.HasChange("environment_predicate") {
		input.EnvironmentPredicate = expandPredicate(d, "environment_predicate")
	}

	_, err := client.UpdateCheckToolUsage(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckToolUsageRead(d, client)
}
