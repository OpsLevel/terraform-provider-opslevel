package opslevel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure   = &CheckHasDocumentationResource{}
	_ resource.ResourceWithImportState = &CheckHasDocumentationResource{}
)

func NewCheckHasDocumentationResource() resource.Resource {
	return &CheckHasDocumentationResource{}
}

// CheckHasDocumentationResource defines the resource implementation.
type CheckHasDocumentationResource struct {
	CommonResourceClient
}

type CheckHasDocumentationResourceModel struct {
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

	DocumentType    types.String `tfsdk:"document_type"`
	DocumentSubtype types.String `tfsdk:"document_subtype"`
}

func NewCheckHasDocumentationResourceModel(ctx context.Context, check opslevel.Check) CheckHasDocumentationResourceModel {
	var model CheckHasDocumentationResourceModel

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

	model.DocumentType = types.StringValue(string(check.DocumentType))
	model.DocumentSubtype = types.StringValue(string(check.DocumentSubtype))

	return model
}

func (r *CheckHasDocumentationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_has_documentation"
}

func (r *CheckHasDocumentationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Has Documentation Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"document_type": schema.StringAttribute{
				Description: "The type of the document.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllHasDocumentationTypeEnum...)},
			},
			"document_subtype": schema.StringAttribute{
				Description: "The subtype of the document.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllHasDocumentationSubtypeEnum...)},
			},
		}),
	}
}

func (r *CheckHasDocumentationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel CheckHasDocumentationResourceModel

	// Read Terraform plan data into the planModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error", err.Error())
	}
	input := opslevel.CheckHasDocumentationCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    planModel.Enabled.ValueBoolPointer(),
		EnableOn:   &iso8601.Time{Time: enabledOn},
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      planModel.Notes.ValueStringPointer(),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}

	input.DocumentType = opslevel.HasDocumentationTypeEnum(planModel.DocumentType.ValueString())
	input.DocumentSubtype = opslevel.HasDocumentationSubtypeEnum(planModel.DocumentSubtype.ValueString())

	data, err := r.client.CreateCheckHasDocumentation(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_has_documentation, got error: %s", err))
		return
	}

	stateModel := NewCheckHasDocumentationResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a check has documentation resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckHasDocumentationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel CheckHasDocumentationResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check has documentation, got error: %s", err))
		return
	}
	stateModel := NewCheckHasDocumentationResourceModel(ctx, *data)
	stateModel.EnableOn = planModel.EnableOn

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckHasDocumentationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel CheckHasDocumentationResourceModel

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
	input := opslevel.CheckHasDocumentationUpdateInput{
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
	input.DocumentType = opslevel.RefOf(opslevel.HasDocumentationTypeEnum(planModel.DocumentType.ValueString()))
	input.DocumentSubtype = opslevel.RefOf(opslevel.HasDocumentationSubtypeEnum(planModel.DocumentSubtype.ValueString()))

	data, err := r.client.UpdateCheckHasDocumentation(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_has_documentation, got error: %s", err))
		return
	}

	stateModel := NewCheckHasDocumentationResourceModel(ctx, *data)

	tflog.Trace(ctx, "updated a check has documentation resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckHasDocumentationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel CheckHasDocumentationResourceModel

	// Read Terraform prior state data into the planModel
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(planModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check has documentation, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check has documentation resource")
}

func (r *CheckHasDocumentationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
