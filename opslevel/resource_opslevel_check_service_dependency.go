package opslevel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure   = &CheckServiceDependencyResource{}
	_ resource.ResourceWithImportState = &CheckServiceDependencyResource{}
)

func NewCheckServiceDependencyResource() resource.Resource {
	return &CheckServiceDependencyResource{}
}

// CheckServiceDependencyResource defines the resource implementation.
type CheckServiceDependencyResource struct {
	CommonResourceClient
}

type CheckServiceDependencyResourceModel struct {
	Category    types.String `tfsdk:"category"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	EnableOn    types.String `tfsdk:"enable_on"`
	Filter      types.String `tfsdk:"filter"`
	Id          types.String `tfsdk:"id"`
	Level       types.String `tfsdk:"level"`
	Name        types.String `tfsdk:"name"`
	Notes       types.String `tfsdk:"notes"`
	Owner       types.String `tfsdk:"owner"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func NewCheckServiceDependencyResourceModel(ctx context.Context, check opslevel.Check) CheckServiceDependencyResourceModel {
	var model CheckServiceDependencyResourceModel

	model.Category = types.StringValue(string(check.Category.Id))
	model.Enabled = types.BoolValue(check.Enabled)
	model.EnableOn = types.StringValue(check.EnableOn.Time.Format(time.RFC3339))
	model.Filter = types.StringValue(string(check.Filter.Id))
	model.Id = types.StringValue(string(check.Id))
	model.Level = types.StringValue(string(check.Level.Id))
	model.Name = types.StringValue(check.Name)
	model.Notes = types.StringValue(check.Notes)
	model.Owner = types.StringValue(string(check.Owner.Team.Id))
	model.LastUpdated = timeLastUpdated()

	return model
}

func (r *CheckServiceDependencyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_service_dependency"
}

func (r *CheckServiceDependencyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Service Dependency Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{}),
	}
}

func (r *CheckServiceDependencyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckServiceDependencyResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
	}
	input := opslevel.CheckServiceDependencyCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}

	data, err := r.client.CreateCheckServiceDependency(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_service_dependency, got error: %s", err))
		return
	}

	stateModel := NewCheckServiceDependencyResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a check service dependency resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckServiceDependencyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckServiceDependencyResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check service dependency, got error: %s", err))
		return
	}
	stateModel := NewCheckServiceDependencyResourceModel(ctx, *data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckServiceDependencyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckServiceDependencyResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
		return
	}
	input := opslevel.CheckServiceDependencyUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Id:         asID(planModel.Id),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}

	data, err := r.client.UpdateCheckServiceDependency(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_service_dependency, got error: %s", err))
		return
	}

	stateModel := NewCheckServiceDependencyResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a check service dependency resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckServiceDependencyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckServiceDependencyResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check service dependency, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check service dependency resource")
}

func (r *CheckServiceDependencyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
