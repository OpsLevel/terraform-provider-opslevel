package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &SystemResource{}

var _ resource.ResourceWithImportState = &SystemResource{}

func NewSystemResource() resource.Resource {
	return &SystemResource{}
}

// SystemResource defines the resource implementation.
type SystemResource struct {
	CommonResourceClient
}

// SystemResourceModel describes the System managed resource.
type SystemResourceModel struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Domain      types.String `tfsdk:"domain"`
	Id          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Name        types.String `tfsdk:"name"`
	Note        types.String `tfsdk:"note"`
	Owner       types.String `tfsdk:"owner"`
}

func NewSystemResourceModel(ctx context.Context, system opslevel.System) (SystemResourceModel, diag.Diagnostics) {
	aliases, diags := OptionalStringListValue(ctx, system.Aliases)
	systemDataSourceModel := SystemResourceModel{
		Aliases:     aliases,
		Description: types.StringValue(system.Description),
		Domain:      types.StringValue(string(system.Parent.Id)),
		Id:          types.StringValue(string(system.Id)),
		Name:        types.StringValue(system.Name),
		Note:        types.StringValue(system.Note),
		Owner:       types.StringValue(string(system.Owner.Id())),
	}
	return systemDataSourceModel, diags
}

func (r *SystemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (r *SystemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "System Resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The aliases of the system.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description for the system.",
				Optional:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The id of the parent domain this system is a child for.",
				Optional:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the system.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The name for the system.",
				Required:    true,
			},
			"note": schema.StringAttribute{
				Description: "Additional information about the system.",
				Optional:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The id of the team that owns the system.",
				Optional:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
		},
	}
}

func (r *SystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel SystemResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	system, err := r.client.CreateSystem(opslevel.SystemInput{
		Name:        planModel.Name.ValueStringPointer(),
		Description: planModel.Description.ValueStringPointer(),
		OwnerId:     opslevel.NewID(planModel.Owner.ValueString()),
		Parent:      opslevel.NewIdentifier(planModel.Domain.ValueString()),
		Note:        planModel.Note.ValueStringPointer(),
	})
	if err != nil || system == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create system, got error: %s", err))
		return
	}
	stateModel, diags := NewSystemResourceModel(ctx, *system)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a system resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *SystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel SystemResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	readSystem, err := r.client.GetSystem(planModel.Id.ValueString())
	if err != nil || readSystem == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read system, got error: %s", err))
		return
	}
	stateModel, diags := NewSystemResourceModel(ctx, *readSystem)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *SystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel SystemResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := r.client.UpdateSystem(planModel.Id.ValueString(), opslevel.SystemInput{
		Name:        planModel.Name.ValueStringPointer(),
		Description: planModel.Description.ValueStringPointer(),
		OwnerId:     opslevel.NewID(planModel.Owner.ValueString()),
		Parent:      opslevel.NewIdentifier(planModel.Domain.ValueString()),
		Note:        planModel.Note.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update system, got error: %s", err))
		return
	}
	stateModel, diags := NewSystemResourceModel(ctx, *system)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a system resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *SystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel SystemResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSystem(planModel.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete system, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a system resource")
}

func (r *SystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
