package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &InfrastructureResource{}

var _ resource.ResourceWithImportState = &InfrastructureResource{}

func NewInfrastructureResource() resource.Resource {
	return &InfrastructureResource{}
}

// InfrastructureResource defines the resource implementation.
type InfrastructureResource struct {
	CommonResourceClient
}

// WIP: pick up here
var infraProviderAttrTypes = map[string]attr.Type{
	"account": types.StringType,
	"name":    types.StringType,
	"type":    types.StringType,
	"url":     types.StringType,
}

// NOTE: old approach, might drop
var infraProviderDataObjectType = types.ObjectType{
	AttrTypes: infraProviderAttrTypes,
}

func newInfraProviderData(infrastructure opslevel.InfrastructureResource) (basetypes.ObjectValue, diag.Diagnostics) {
	infraAttrs := make(map[string]attr.Value)
	infraAttrs["account"] = RequiredStringValue(infrastructure.ProviderData.AccountName)
	infraAttrs["name"] = OptionalStringValue(infrastructure.ProviderData.ProviderName)
	infraAttrs["type"] = OptionalStringValue(infrastructure.ProviderType)
	infraAttrs["url"] = OptionalStringValue(infrastructure.ProviderData.ExternalURL)
	return types.ObjectValue(identifierObjectType.AttrTypes, infraAttrs)
}

// InfrastructureResourceModel describes the Infrastructure managed resource.
type InfrastructureResourceModel struct {
	Aliases      types.List   `tfsdk:"aliases"`
	Data         types.String `tfsdk:"data"`
	Id           types.String `tfsdk:"id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	ProviderData types.Object `tfsdk:"provider_data"`
	Owner        types.String `tfsdk:"owner"`
	Schema       types.String `tfsdk:"schema"`
}

func NewInfrastructureResourceModel(ctx context.Context, infrastructure opslevel.InfrastructureResource) (InfrastructureResourceModel, diag.Diagnostics) {
	providerData, diags := newInfraProviderData(infrastructure)

	aliases, diags := OptionalStringListValue(ctx, infrastructure.Aliases)
	return InfrastructureResourceModel{
		Aliases:      aliases,
		Data:         OptionalStringValue(infrastructure.Data.ToJSON()),
		Id:           ComputedStringValue(infrastructure.Id),
		ProviderData: providerData,
		Owner:        OptionalStringValue(string(infrastructure.Owner.Id())),
		Schema:       RequiredStringValue(infrastructure.Schema),
	}, diags
}

func (r *InfrastructureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infrastructure"
}

func (r *InfrastructureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Infrastructure Resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The aliases for the infrastructure resource.",
				Optional:    true,
			},
			"data": schema.StringAttribute{
				Description: "The data of the infrastructure resource in JSON format.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the infrastructure.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"owner": schema.StringAttribute{
				Description: "The id of the team that owns the infrastructure resource. Does not support aliases!",
				Optional:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"provider_data": schema.SingleNestedAttribute{
				Description: "The provider specific data for the infrastructure resource.",
				Optional:    true,
				Default:     objectdefault.StaticValue(types.ObjectNull(infraProviderAttrTypes)),
				Attributes: map[string]schema.Attribute{
					"account": schema.StringAttribute{
						Description: "The canonical account name for the provider of the infrastructure resource.",
						Required:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the provider of the infrastructure resource. (eg. AWS, GCP, Azure)",
						Optional:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type of the infrastructure resource as defined by its provider.",
						Optional:    true,
					},
					"url": schema.StringAttribute{
						Description: "The url for the provider of the infrastructure resource.",
						Optional:    true,
					},
				},
			},
			"schema": schema.StringAttribute{
				Description: "The schema of the infrastructure resource that determines its data specification.",
				Required:    true,
			},
		},
	}
}

func (r *InfrastructureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InfrastructureResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newJSON, err := opslevel.NewJSON(data.Data.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to parse infrastructure resource 'data' into JSON, got error: %s", err))
		return
	}

	infrastructure, err := r.client.CreateInfrastructure(opslevel.InfraInput{
		Schema: data.Schema.ValueString(),
		Owner:  opslevel.NewID(data.Owner.ValueString()),
		Data:   newJSON,
		// Provider: expandInfraProviderData(data.ProviderData), // TODO: fix this
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create infrastructure, got error: %s", err))
		return
	}

	createdInfrastructureResourceModel, diags := NewInfrastructureResourceModel(ctx, *infrastructure)
	resp.Diagnostics.Append(diags...)
	if data.Aliases.IsNull() {
		createdInfrastructureResourceModel.Aliases = data.Aliases
	}
	createdInfrastructureResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a infrastructure resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdInfrastructureResourceModel)...)
}

func (r *InfrastructureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InfrastructureResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructure, err := r.client.GetInfrastructure(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read infrastructure, got error: %s", err))
		return
	}
	readInfrastructureResourceModel, diags := NewInfrastructureResourceModel(ctx, *infrastructure)
	resp.Diagnostics.Append(diags...)
	if data.Aliases.IsNull() {
		readInfrastructureResourceModel.Aliases = data.Aliases
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readInfrastructureResourceModel)...)
}

func (r *InfrastructureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InfrastructureResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newJSON, err := opslevel.NewJSON(data.Data.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to parse infrastructure resource 'data' into JSON, got error: %s", err))
		return
	}
	updatedInfrastructure, err := r.client.UpdateInfrastructure(data.Id.ValueString(), opslevel.InfraInput{
		Schema: data.Schema.ValueString(),
		Owner:  opslevel.NewID(data.Owner.ValueString()),
		Data:   newJSON,
		// Provider: expandInfraProviderData(data.ProviderData), // TODO: fix this
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update infrastructure, got error: %s", err))
		return
	}
	updatedInfrastructureResourceModel, diags := NewInfrastructureResourceModel(ctx, *updatedInfrastructure)
	resp.Diagnostics.Append(diags...)
	if data.Aliases.IsNull() {
		updatedInfrastructureResourceModel.Aliases = data.Aliases
	}
	updatedInfrastructureResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a infrastructure resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedInfrastructureResourceModel)...)
}

func (r *InfrastructureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InfrastructureResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteInfrastructure(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete infrastructure, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a infrastructure resource")
}

func (r *InfrastructureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// TODO: fix this
// func expandInfraProviderData(providerData map[string]attr.Type) *opslevel.InfraProviderInput {
// 	var blank infraProviderData
// 	if providerData == blank {
// 		return &opslevel.InfraProviderInput{}
// 	}
// 	return &opslevel.InfraProviderInput{
// 		Account: providerData.Account.ValueString(),
// 		Name:    providerData.Name.ValueString(),
// 		Type:    providerData.Type.ValueString(),
// 		URL:     providerData.Url.ValueString(),
// 	}
// }

// import (
// 	"time"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceInfrastructure() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a infrastructure",
// 		Create:      wrap(resourceInfrastructureCreate),
// 		Read:        wrap(resourceInfrastructureRead),
// 		Update:      wrap(resourceInfrastructureUpdate),
// 		Delete:      wrap(resourceInfrastructureDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"alias": {
// 				Type:        schema.TypeString,
// 				Description: "The alias for this infrastructure.",
// 				ForceNew:    true,
// 				Required:    true,
// 			},
// 			"owner": {
// 				Type:        schema.TypeString,
// 				Description: "The owner of this infrastructure.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"value": {
// 				Type:        schema.TypeString,
// 				Description: "A sensitive value.",
// 				Sensitive:   true,
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"created_at": {
// 				Type:        schema.TypeString,
// 				Description: "Timestamp of time created at.",
// 				Computed:    true,
// 			},
// 			"updated_at": {
// 				Type:        schema.TypeString,
// 				Description: "Timestamp of last update.",
// 				Computed:    true,
// 			},
// 		},
// 	}
// }

// func resourceInfrastructureCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.InfrastructureInput{
// 		Owner: opslevel.NewIdentifier(d.Get("owner").(string)),
// 		Value: opslevel.RefOf(d.Get("value").(string)),
// 	}
// 	alias := d.Get("alias").(string)
// 	resource, err := client.CreateInfrastructure(alias, input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.ID))
// 	return resourceInfrastructureRead(d, client)
// }

// func resourceInfrastructureRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetInfrastructure(id)
// 	if err != nil {
// 		return err
// 	}

// 	if opslevel.IsID(d.Get("owner").(string)) {
// 		if err := d.Set("owner", resource.Owner.Id); err != nil {
// 			return err
// 		}
// 	} else {
// 		if err := d.Set("owner", resource.Owner.Alias); err != nil {
// 			return err
// 		}
// 	}
// 	created_at := resource.Timestamps.CreatedAt.Local().Format(time.RFC850)
// 	if err := d.Set("created_at", created_at); err != nil {
// 		return err
// 	}
// 	updated_at := resource.Timestamps.UpdatedAt.Local().Format(time.RFC850)
// 	if err := d.Set("updated_at", updated_at); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceInfrastructureUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.InfrastructureInput{
// 		Owner: opslevel.NewIdentifier(d.Get("owner").(string)),
// 		Value: opslevel.RefOf(d.Get("value").(string)),
// 	}

// 	_, err := client.UpdateInfrastructure(d.Id(), input)
// 	if err != nil {
// 		return err
// 	}

// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceInfrastructureRead(d, client)
// }

// func resourceInfrastructureDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteInfrastructure(id)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }

// import (
// 	"slices"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceInfrastructure() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages an infrastructure resource",
// 		Create:      wrap(resourceInfrastructureCreate),
// 		Read:        wrap(resourceInfrastructureRead),
// 		Update:      wrap(resourceInfrastructureUpdate),
// 		Delete:      wrap(resourceInfrastructureDelete),
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"aliases": {
// 				Type:        schema.TypeList,
// 				Description: "The aliases of the infrastructure resource.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				Elem:        &schema.Schema{Type: schema.TypeString},
// 			},
// 			"schema": {
// 				Type:        schema.TypeString,
// 				Description: "The schema of the infrastructure resource that determines its data specification.",
// 				Required:    true,
// 				ForceNew:    true,
// 			},
// 			"owner": {
// 				Type:        schema.TypeString,
// 				Description: "The id of the team that owns the infrastructure resource. Does not support aliases!",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"provider_data": {
// 				Type:        schema.TypeList,
// 				Description: "The provider specific data for the infrastructure resource.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				MaxItems:    1,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"name": {
// 							Type:        schema.TypeString,
// 							Description: "The name of the provider of the infrastructure resource. (eg. AWS, GCP, Azure)",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 						"type": {
// 							Type:        schema.TypeString,
// 							Description: "The type of the infrastructure resource as defined by its provider.",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 						"account": {
// 							Type:        schema.TypeString,
// 							Description: "The canonical account name for the provider of the infrastructure resource.",
// 							ForceNew:    false,
// 							Required:    true,
// 						},
// 						"url": {
// 							Type:        schema.TypeString,
// 							Description: "The url for the provider of the infrastructure resource.",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 					},
// 				},
// 			},
// 			"data": {
// 				Type:        schema.TypeString,
// 				Description: "The data of the infrastructure resource in JSON format.",
// 				Optional:    true,
// 			},
// 		},
// 	}
// }

// func reconcileInfraAliases(d *schema.ResourceData, resource *opslevel.InfrastructureResource, client *opslevel.Client) error {
// 	expectedAliases := getStringArray(d, "aliases")
// 	existingAliases := resource.Aliases
// 	for _, existingAlias := range existingAliases {
// 		if !slices.Contains(expectedAliases, existingAlias) {
// 			err := client.DeleteInfraAlias(existingAlias)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	for _, expectedAlias := range expectedAliases {
// 		if !slices.Contains(existingAliases, expectedAlias) {
// 			id := opslevel.NewID(resource.Id)
// 			_, err := client.CreateAliases(*id, []string{expectedAlias})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func flattenInfraProviderData(resource *opslevel.InfrastructureResource) []map[string]any {
// 	return []map[string]any{{
// 		"account": resource.ProviderData.AccountName,
// 		"name":    resource.ProviderData.ProviderName,
// 		"type":    resource.ProviderType,
// 		"url":     resource.ProviderData.ExternalURL,
// 	}}
// }

// func expandInfraProviderData(d *schema.ResourceData) *opslevel.InfraProviderInput {
// 	config := d.Get("provider_data").([]interface{})
// 	if len(config) > 0 {
// 		item := config[0].(map[string]interface{})
// 		return &opslevel.InfraProviderInput{
// 			Account: item["account"].(string),
// 			Name:    item["name"].(string),
// 			Type:    item["type"].(string),
// 			URL:     item["url"].(string),
// 		}
// 	}
// 	return &opslevel.InfraProviderInput{}
// }

// func resourceInfrastructureCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	newJSON, err := opslevel.NewJSON(d.Get("data").(string))
// 	if err != nil {
// 		return err
// 	}
// 	resource, err := client.CreateInfrastructure(opslevel.InfraInput{
// 		Schema:   d.Get("schema").(string),
// 		Owner:    opslevel.NewID(d.Get("owner").(string)),
// 		Provider: expandInfraProviderData(d),
// 		Data:     newJSON,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(resource.Id)

// 	err = reconcileInfraAliases(d, resource, client)
// 	if err != nil {
// 		return err
// 	}

// 	return resourceInfrastructureRead(d, client)
// }

// func resourceInfrastructureRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetInfrastructure(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("schema", resource.Schema); err != nil {
// 		return err
// 	}
// 	if err := d.Set("aliases", resource.Aliases); err != nil {
// 		return err
// 	}
// 	if err := d.Set("owner", resource.Owner.Id()); err != nil {
// 		return err
// 	}
// 	if err := d.Set("provider_data", flattenInfraProviderData(resource)); err != nil {
// 		return err
// 	}
// 	if err := d.Set("data", resource.Data.ToJSON()); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceInfrastructureUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	newJSON, err := opslevel.NewJSON(d.Get("data").(string))
// 	if err != nil {
// 		return err
// 	}
// 	resource, err := client.UpdateInfrastructure(id, opslevel.InfraInput{
// 		Schema:   d.Get("schema").(string),
// 		Owner:    opslevel.NewID(d.Get("owner").(string)),
// 		Provider: expandInfraProviderData(d),
// 		Data:     newJSON,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if d.HasChange("aliases") {
// 		err = reconcileInfraAliases(d, resource, client)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceInfrastructureRead(d, client)
// }

// func resourceInfrastructureDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteInfrastructure(id)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
