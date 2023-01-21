package opslevel

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
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
							ValidateFunc: validation.IsRFC3339Time,
						},
						"time_scale": {
							Type:         schema.TypeString,
							Description:  "The time scale type for the frequency.",
							ForceNew:     false,
							Required:     true,
							ValidateFunc: validation.StringInSlice(opslevel.AllFrequencyTimeScale(), false),
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

func expandUpdateFrequency(d *schema.ResourceData, key string) *opslevel.ManualCheckFrequencyInput {
	if _, ok := d.GetOk(key); !ok {
		return nil
	}
	return opslevel.NewManualCheckFrequencyInput(
		d.Get(fmt.Sprintf("%s.0.starting_data", key)).(string),
		opslevel.FrequencyTimeScale(d.Get(fmt.Sprintf("%s.0.time_scale", key)).(string)),
		d.Get(fmt.Sprintf("%s.0.value", key)).(int),
	)
}

func flattenUpdateFrequency(input *opslevel.ManualCheckFrequency) []map[string]interface{} {
	output := []map[string]interface{}{}
	if input != nil {
		output = append(output, map[string]interface{}{
			"starting_data": input.StartingDate.Format(time.RFC3339),
			"time_scale":    string(input.FrequencyTimeScale),
			"value":         input.FrequencyValue,
		})
	}
	return output
}

func resourceCheckManualCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckManualCreateInput{}
	setCheckCreateInput(d, &input)

	input.UpdateRequiresComment = d.Get("update_requires_comment").(bool)
	input.UpdateFrequency = expandUpdateFrequency(d, "update_frequency")

	resource, err := client.CreateCheckManual(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckManualRead(d, client)
}

func resourceCheckManualRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("update_frequency", flattenUpdateFrequency(resource.UpdateFrequency)); err != nil {
		return err
	}
	if err := d.Set("update_requires_comment", resource.UpdateRequiresComment); err != nil {
		return err
	}

	return nil
}

func resourceCheckManualUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckManualUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("update_frequency") {
		input.UpdateFrequency = expandUpdateFrequency(d, "update_frequency")
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
