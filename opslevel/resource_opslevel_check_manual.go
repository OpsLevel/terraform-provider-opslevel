package opslevel

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &CheckManualResource{}
var _ resource.ResourceWithImportState = &CheckManualResource{}

func NewCheckManualResource() resource.Resource {
	return &CheckManualResource{}
}

// CheckManualResource defines the resource implementation.
type CheckManualResource struct {
	CommonResourceClient
}

type CheckUpdateFrequency struct {
	StartingDate timetypes.RFC3339 `tfsdk:"starting_date"`
	TimeScale    types.String      `tfsdk:"time_scale"`
	Value        types.Int64       `tfsdk:"value"`
}

type CheckManualResourceModel struct {
	CheckBaseModel
	UpdateFrequency       CheckUpdateFrequency `tfsdk:"update_frequency"`
	UpdateRequiresComment types.Bool           `tfsdk:"update_requires_comment"`
}

func (r *CheckManualResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_manual"
}

func (r *CheckManualResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Manual Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"update_requires_comment": schema.BoolAttribute{
				Description: "Whether the check requires a comment or not.",
				Required:    true,
			},
		}),
		Blocks: map[string]schema.Block{
			"update_frequency": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"starting_date": schema.StringAttribute{
						Description: "The date that the check will start to evaluate.",
						Required:    true,
					},
					"time_scale": schema.StringAttribute{
						Description: "The time scale type for the frequency.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf(opslevel.AllFrequencyTimeScale...),
						},
					},
					"value": schema.Int64Attribute{
						Description: "The value to be used together with the frequency time_scale.",
						Required:    true,
					},
				},
			},
		},
	}
}

func (r *CheckManualResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

func (r *CheckManualResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *CheckManualResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *CheckManualResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *CheckManualResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// import (
// 	"fmt"
// 	"time"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceCheckManual() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a manual check.",
// 		Create:      wrap(resourceCheckManualCreate),
// 		Read:        wrap(resourceCheckManualRead),
// 		Update:      wrap(resourceCheckManualUpdate),
// 		Delete:      wrap(resourceCheckDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: getCheckSchema(map[string]*schema.Schema{
// 			"update_frequency": {
// 				Type:        schema.TypeList,
// 				MaxItems:    1,
// 				Description: "Defines the minimum frequency of the updates.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"starting_data": {
// 							Type:         schema.TypeString,
// 							Description:  "The date that the check will start to evaluate.",
// 							ForceNew:     false,
// 							Required:     true,
// 							ValidateFunc: validation.IsRFC3339Time,
// 						},
// 						"time_scale": {
// 							Type:         schema.TypeString,
// 							Description:  "The time scale type for the frequency.",
// 							ForceNew:     false,
// 							Required:     true,
// 							ValidateFunc: validation.StringInSlice(opslevel.AllFrequencyTimeScale, false),
// 						},
// 						"value": {
// 							Type:        schema.TypeInt,
// 							Description: "The value to be used together with the frequency scale.",
// 							ForceNew:    false,
// 							Required:    true,
// 						},
// 					},
// 				},
// 			},
// 			"update_requires_comment": {
// 				Type:        schema.TypeBool,
// 				Description: "Whether the check requires a comment or not.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 		}),
// 	}
// }

// func expandUpdateFrequencyOnCreate(d *schema.ResourceData, key string) *opslevel.ManualCheckFrequencyInput {
// 	if _, ok := d.GetOk(key); !ok {
// 		return nil
// 	}
// 	return opslevel.NewManualCheckFrequencyInput(
// 		d.Get(fmt.Sprintf("%s.0.starting_data", key)).(string),
// 		opslevel.FrequencyTimeScale(d.Get(fmt.Sprintf("%s.0.time_scale", key)).(string)),
// 		d.Get(fmt.Sprintf("%s.0.value", key)).(int),
// 	)
// }

// func expandUpdateFrequencyOnUpdate(d *schema.ResourceData, key string) *opslevel.ManualCheckFrequencyUpdateInput {
// 	if _, ok := d.GetOk(key); !ok {
// 		return nil
// 	}
// 	return opslevel.NewManualCheckFrequencyUpdateInput(
// 		d.Get(fmt.Sprintf("%s.0.starting_data", key)).(string),
// 		opslevel.FrequencyTimeScale(d.Get(fmt.Sprintf("%s.0.time_scale", key)).(string)),
// 		d.Get(fmt.Sprintf("%s.0.value", key)).(int),
// 	)
// }

// func flattenUpdateFrequency(input *opslevel.ManualCheckFrequency) []map[string]interface{} {
// 	output := []map[string]interface{}{}
// 	if input != nil {
// 		output = append(output, map[string]interface{}{
// 			"starting_data": input.StartingDate.Format(time.RFC3339),
// 			"time_scale":    string(input.FrequencyTimeScale),
// 			"value":         input.FrequencyValue,
// 		})
// 	}
// 	return output
// }

// func resourceCheckManualCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	checkCreateInput := getCheckCreateInputFrom(d)
// 	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckManualCreateInput](checkCreateInput)
// 	input.UpdateRequiresComment = d.Get("update_requires_comment").(bool)
// 	input.UpdateFrequency = expandUpdateFrequencyOnCreate(d, "update_frequency")

// 	resource, err := client.CreateCheckManual(*input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceCheckManualRead(d, client)
// }

// func resourceCheckManualRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetCheck(opslevel.ID(id))
// 	if err != nil {
// 		return err
// 	}

// 	if err := setCheckData(d, resource); err != nil {
// 		return err
// 	}
// 	if err := d.Set("update_frequency", flattenUpdateFrequency(resource.UpdateFrequency)); err != nil {
// 		return err
// 	}
// 	if err := d.Set("update_requires_comment", resource.UpdateRequiresComment); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceCheckManualUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	checkUpdateInput := getCheckUpdateInputFrom(d)
// 	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckManualUpdateInput](checkUpdateInput)
// 	if d.HasChange("update_frequency") {
// 		input.UpdateFrequency = expandUpdateFrequencyOnUpdate(d, "update_frequency")
// 	}
// 	input.UpdateRequiresComment = opslevel.RefOf(d.Get("update_requires_comment").(bool))

// 	_, err := client.UpdateCheckManual(*input)
// 	if err != nil {
// 		return err
// 	}
// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceCheckManualRead(d, client)
// }
