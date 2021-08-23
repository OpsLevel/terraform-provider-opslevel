package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go"
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
				ValidateFunc: validation.StringInSlice(opslevel.GetToolCategoryTypes(), false),
			},
			"tool_name_predicate":   getPredicateInputSchema(false),
			"environment_predicate": getPredicateInputSchema(false),
		}),
	}
}

func resourceCheckToolUsageCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckToolUsageCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),

		ToolCategory:         opslevel.ToolCategory(d.Get("tool_category").(string)),
		ToolNamePredicate:    getPredicateInput(d, "tool_name_predicate"),
		EnvironmentPredicate: getPredicateInput(d, "environment_predicate"),
	}
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

	if err := resourceCheckRead(d, resource); err != nil {
		return err
	}

	return nil
}

func resourceCheckToolUsageUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckToolUsageUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("enabled") {
		value := d.Get("enabled").(bool)
		input.Enabled = &value
	}
	if d.HasChange("category") {
		input.Category = getID(d, "category")
	}
	if d.HasChange("level") {
		input.Level = getID(d, "level")
	}
	if d.HasChange("owner") {
		input.Owner = getID(d, "owner")
	}
	if d.HasChange("filter") {
		input.Filter = getID(d, "filter")
	}
	if d.HasChange("notes") {
		input.Notes = d.Get("notes").(string)
	}

	if d.HasChange("tool_category") {
		input.ToolCategory = opslevel.ToolCategory(d.Get("tool_category").(string))
	}
	if d.HasChange("tool_name_predicate") {
		input.ToolNamePredicate = getPredicateInput(d, "tool_name_predicate")
	}
	if d.HasChange("environment_predicate") {
		input.EnvironmentPredicate = getPredicateInput(d, "environment_predicate")
	}

	_, err := client.UpdateCheckToolUsage(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckToolUsageRead(d, client)
}
