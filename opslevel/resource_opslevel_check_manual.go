package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go"
)

func resourceCheckManual() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a manual check.",
		Create:      wrap(resourceCheckManualCreate),
		Read:        wrap(resourceCheckManualRead),
		Update:      wrap(resourceCheckManualUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"update_frequency": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Defines the minimum frequency of the updates.",
				ForceNew:    false,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"starting_data": {
							Type:         schema.TypeString,
							Description:  "The date that the check will start to evaluate.",
							ForceNew:     false,
							Required:     true,
							ValidateFunc: validation.ValidateRFC3339TimeString,
						},
						"time_scale": {
							Type:         schema.TypeString,
							Description:  "The time scale type for the frequency.",
							ForceNew:     false,
							Required:     true,
							ValidateFunc: validation.StringInSlice(opslevel.GetFrequencyTimeScales(), false),
						},
						"value": {
							Type:        schema.TypeInt,
							Description: "The value to be used together with the frequency scale.",
							ForceNew:    false,
							Required:    true,
						},
					},
				},
			},
			"update_requires_comment": {
				Type:        schema.TypeBool,
				Description: "Whether the check requires a comment or not.",
				ForceNew:    false,
				Required:    true,
			},
		}),
	}
}

func resourceCheckManualCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckManualCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),

		UpdateRequiresComment: d.Get("update_requires_comment").(bool),
	}
	if _, ok := d.GetOk("update_frequency"); ok {
		input.UpdateFrequency = opslevel.NewManualCheckFrequencyInput(
			d.Get("update_frequency.0.starting_data").(string),
			opslevel.FrequencyTimeScale(d.Get("update_frequency.0.time_scale").(string)),
			d.Get("update_frequency.0.value").(int),
		)
	}

	resource, err := client.CreateCheckManual(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckManualRead(d, client)
}

func resourceCheckManualRead(d *schema.ResourceData, client *opslevel.Client) error {
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

func resourceCheckManualUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckManualUpdateInput{
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

	if d.HasChange("update_frequency") {
		input.UpdateFrequency = opslevel.NewManualCheckFrequencyInput(
			d.Get("update_frequency.0.starting_data").(string),
			opslevel.FrequencyTimeScale(d.Get("update_frequency.0.time_scale").(string)),
			d.Get("update_frequency.0.value").(int),
		)
	}
	if d.HasChange("update_requires_comment") {
		input.UpdateRequiresComment = d.Get("update_requires_comment").(bool)
	}

	_, err := client.UpdateCheckManual(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckManualRead(d, client)
}
