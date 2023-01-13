package opslevel

import (
	"github.com/hasura/go-graphql-client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
)

func getCheckSchema(extras map[string]*schema.Schema) map[string]*schema.Schema {
	output := map[string]*schema.Schema{
		"last_updated": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The display name of the check.",
			ForceNew:    false,
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: `Whether the check is enabled or not.  Do not use this field in tandem with 'enable_on'.`,
			ForceNew:    false,
			Optional:    true,
		},
		"enable_on": {
			Type: schema.TypeString,
			Description: `The date when the check will be automatically enabled.
If you use this field you should add both 'enabled' and 'enable_on' to the lifecycle ignore_changes settings.
See example in opslevel_check_manual for proper configuration.
`,
			ForceNew:     false,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},
		"category": {
			Type:        schema.TypeString,
			Description: "The id of the category the check belongs to.",
			ForceNew:    false,
			Required:    true,
		},
		"level": {
			Type:        schema.TypeString,
			Description: "The id of the level the check belongs to.",
			ForceNew:    false,
			Required:    true,
		},
		"owner": {
			Type:        schema.TypeString,
			Description: "The id of the team that owns the check.",
			ForceNew:    false,
			Optional:    true,
		},
		"filter": {
			Type:        schema.TypeString,
			Description: "The id of the filter of the check.",
			ForceNew:    false,
			Optional:    true,
		},
		"notes": {
			Type:        schema.TypeString,
			Description: "Additional information about the check.",
			ForceNew:    false,
			Optional:    true,
		},
	}
	for k, v := range extras {
		output[k] = v
	}
	return output
}

func setCheckData(d *schema.ResourceData, resource *opslevel.Check) error {
	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("enabled", resource.Enabled); err != nil {
		return err
	}
	if _, ok := d.GetOk("enable_on"); ok {
		if err := d.Set("enable_on", resource.EnableOn.Format(time.RFC3339)); err != nil {
			return err
		}
	}
	if err := d.Set("category", resource.Category.Id); err != nil {
		return err
	}
	if err := d.Set("level", resource.Level.Id); err != nil {
		return err
	}
	if err := d.Set("owner", resource.Owner.Team.Id); err != nil {
		return err
	}
	if err := d.Set("filter", resource.Filter.Id); err != nil {
		return err
	}
	if err := d.Set("notes", resource.Notes); err != nil {
		return err
	}
	return nil
}

func setCheckCreateInput(d *schema.ResourceData, p opslevel.CheckCreateInputProvider) {
	input := p.GetCheckCreateInput()
	input.Name = d.Get("name").(string)
	input.Enabled = d.Get("enabled").(bool)
	if value, ok := d.GetOk("enable_on"); ok {
		enable_on := opslevel.NewISO8601Date(value.(string))
		input.EnableOn = &enable_on
	}
	input.Category = *getID(d, "category")
	input.Level = *getID(d, "level")
	input.Owner = getID(d, "owner")
	input.Filter = getID(d, "filter")
	input.Notes = d.Get("notes").(string)
}

func setCheckUpdateInput(d *schema.ResourceData, p opslevel.CheckUpdateInputProvider) {
	input := p.GetCheckUpdateInput()
	input.Id = graphql.ID(d.Id())

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("enabled") {
		value := d.Get("enabled").(bool)
		input.Enabled = &value
	}
	if d.HasChange("enable_on") {
		enable_on := opslevel.NewISO8601Date(d.Get("enable_on").(string))
		input.EnableOn = &enable_on
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
}

func resourceCheckDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteCheck(graphql.ID(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
