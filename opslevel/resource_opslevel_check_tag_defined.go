package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
)

func resourceCheckTagDefined() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a tag defined check",
		Create:      wrap(resourceCheckTagDefinedCreate),
		Read:        wrap(resourceCheckTagDefinedRead),
		Update:      wrap(resourceCheckTagDefinedUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"tag_key": {
				Type:        schema.TypeString,
				Description: "The tag key where the tag predicate should be applied.",
				ForceNew:    false,
				Required:    true,
			},
			"tag_predicate": getPredicateInputSchema(false),
		}),
	}
}

func resourceCheckTagDefinedCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckTagDefinedCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),

		TagKey:       d.Get("tag_key").(string),
		TagPredicate: getPredicateInput(d, "tag_predicate"),
	}
	resource, err := client.CreateCheckTagDefined(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckTagDefinedRead(d, client)
}

func resourceCheckTagDefinedRead(d *schema.ResourceData, client *opslevel.Client) error {
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

func resourceCheckTagDefinedUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckTagDefinedUpdateInput{
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

	if d.HasChange("tag_key") {
		input.TagKey = d.Get("tag_key").(string)
	}
	if d.HasChange("tag_predicate") {
		input.TagPredicate = getPredicateInput(d, "tag_predicate")
	}

	_, err := client.UpdateCheckTagDefined(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckTagDefinedRead(d, client)
}
