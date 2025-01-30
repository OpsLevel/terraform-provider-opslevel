package opslevel

import (
	"context"
	"fmt"
	"strings"

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
	CheckCodeBaseResourceModel

	DocumentType    types.String `tfsdk:"document_type"`
	DocumentSubtype types.String `tfsdk:"document_subtype"`
}

func NewCheckHasDocumentationResourceModel(ctx context.Context, check opslevel.Check, planModel CheckHasDocumentationResourceModel) CheckHasDocumentationResourceModel {
	var stateModel CheckHasDocumentationResourceModel

	stateModel.Category = RequiredStringValue(string(check.Category.Id))
	stateModel.Description = ComputedStringValue(check.Description)
	if planModel.Enabled.IsNull() {
		stateModel.Enabled = types.BoolValue(false)
	} else {
		stateModel.Enabled = OptionalBoolValue(&check.Enabled)
	}
	if planModel.EnableOn.IsNull() {
		stateModel.EnableOn = types.StringNull()
	} else {
		// We pass through the plan value because of time formatting issue to ensure the state gets the exact value the customer specified
		stateModel.EnableOn = planModel.EnableOn
	}
	stateModel.Filter = OptionalStringValue(string(check.Filter.Id))
	stateModel.Id = ComputedStringValue(string(check.Id))
	stateModel.Level = RequiredStringValue(string(check.Level.Id))
	stateModel.Name = RequiredStringValue(check.Name)
	stateModel.Notes = OptionalStringValue(check.Notes)
	stateModel.Owner = OptionalStringValue(string(check.Owner.Team.Id))

	stateModel.DocumentType = types.StringValue(string(check.DocumentType))
	stateModel.DocumentSubtype = types.StringValue(string(check.DocumentSubtype))

	return stateModel
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
				Description: fmt.Sprintf(
					"The type of the document. One of `%s`",
					strings.Join(opslevel.AllHasDocumentationTypeEnum, "`, `"),
				),
				Required:   true,
				Validators: []validator.String{stringvalidator.OneOf(opslevel.AllHasDocumentationTypeEnum...)},
			},
			"document_subtype": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The subtype of the document. One of `%s`",
					strings.Join(opslevel.AllHasDocumentationSubtypeEnum, "`, `"),
				),
				Required:   true,
				Validators: []validator.String{stringvalidator.OneOf(opslevel.AllHasDocumentationSubtypeEnum...)},
			},
		}),
	}
}

func (r *CheckHasDocumentationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CheckHasDocumentationResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckHasDocumentationCreateInput{
		CategoryId: asID(planModel.Category),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      nullable(planModel.Notes.ValueStringPointer()),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	input.DocumentType = opslevel.HasDocumentationTypeEnum(planModel.DocumentType.ValueString())
	input.DocumentSubtype = opslevel.HasDocumentationSubtypeEnum(planModel.DocumentSubtype.ValueString())

	data, err := r.client.CreateCheckHasDocumentation(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_has_documentation, got error: %s", err))
		return
	}

	stateModel := NewCheckHasDocumentationResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check has documentation resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckHasDocumentationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CheckHasDocumentationResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check has documentation, got error: %s", err))
		return
	}
	verifiedStateModel := NewCheckHasDocumentationResourceModel(ctx, *data, stateModel)
	verifiedStateModel.EnableOn = stateModel.EnableOn

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *CheckHasDocumentationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckHasDocumentationResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckHasDocumentationUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   opslevel.RefOf(asID(planModel.Filter)),
		Id:         asID(planModel.Id),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      nullable(planModel.Notes.ValueStringPointer()),
		OwnerId:    opslevel.RefOf(asID(planModel.Owner)),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}

	input.DocumentType = asEnum[opslevel.HasDocumentationTypeEnum](planModel.DocumentType.ValueString())
	input.DocumentSubtype = asEnum[opslevel.HasDocumentationSubtypeEnum](planModel.DocumentSubtype.ValueString())

	data, err := r.client.UpdateCheckHasDocumentation(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_has_documentation, got error: %s", err))
		return
	}

	stateModel := NewCheckHasDocumentationResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check has documentation resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckHasDocumentationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CheckHasDocumentationResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check has documentation, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check has documentation resource")
}

func (r *CheckHasDocumentationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
