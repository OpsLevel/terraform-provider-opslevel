package opslevel

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/relvacode/iso8601"
)

var (
	_ resource.ResourceWithConfigure   = &CheckCodeIssueResource{}
	_ resource.ResourceWithImportState = &CheckCodeIssueResource{}
)

func NewCheckCodeIssueResource() resource.Resource {
	return &CheckCodeIssueResource{}
}

// CheckCodeIssueResource defines the resource implementation.
type CheckCodeIssueResource struct {
	CommonResourceClient
}

var resolutionTimeType = map[string]attr.Type{
	"unit":  types.StringType,
	"value": types.Int64Type,
}

type CheckCodeIssueResourceModel struct {
	CheckCodeBaseResourceModel

	Constraint     types.String `tfsdk:"constraint"`
	IssueName      types.String `tfsdk:"issue_name"`
	IssueType      types.List   `tfsdk:"issue_type"`
	MaxAllowed     types.Int64  `tfsdk:"max_allowed"`
	ResolutionTime types.Object `tfsdk:"resolution_time"`
	Severity       types.List   `tfsdk:"severity"`
}

func NewCheckCodeIssueResourceModel(ctx context.Context, check opslevel.Check, givenModel CheckCodeIssueResourceModel) CheckCodeIssueResourceModel {
	var stateModel CheckCodeIssueResourceModel

	stateModel.Category = RequiredStringValue(string(check.Category.Id))
	stateModel.Description = ComputedStringValue(check.Description)
	if givenModel.Enabled.IsNull() {
		stateModel.Enabled = types.BoolValue(false)
	} else {
		stateModel.Enabled = OptionalBoolValue(&check.Enabled)
	}
	if givenModel.EnableOn.IsNull() {
		stateModel.EnableOn = types.StringNull()
	} else {
		// We pass through the plan value because of time formatting issue to ensure the state gets the exact value the customer specified
		stateModel.EnableOn = givenModel.EnableOn
	}
	stateModel.Filter = OptionalStringValue(string(check.Filter.Id))
	stateModel.Id = ComputedStringValue(string(check.Id))
	stateModel.Level = RequiredStringValue(string(check.Level.Id))
	stateModel.Name = RequiredStringValue(check.Name)
	stateModel.Notes = OptionalStringValue(check.Notes)
	stateModel.Owner = OptionalStringValue(string(check.Owner.Team.Id))

	stateModel.Constraint = RequiredStringValue(string(check.Constraint))
	stateModel.IssueName = OptionalStringValue(check.IssueName)
	stateModel.IssueType = OptionalStringListValue(check.IssueType)
	// NOTE: API prevents MaxAllowed from being zero
	if check.MaxAllowed > 0 {
		stateModel.MaxAllowed = types.Int64Value(int64(check.MaxAllowed))
	}
	emptyResolutionTime := opslevel.CodeIssueResolutionTime{}
	if check.ResolutionTime == emptyResolutionTime {
		stateModel.ResolutionTime = types.ObjectNull(resolutionTimeType)
	} else {
		resolutionTime := check.ResolutionTime
		resolutionTimeAttrValues := map[string]attr.Value{
			"unit":  types.StringValue(string(resolutionTime.Unit)),
			"value": types.Int64Value(int64(resolutionTime.Value)),
		}
		stateModel.ResolutionTime = types.ObjectValueMust(resolutionTimeType, resolutionTimeAttrValues)
	}
	stateModel.Severity = OptionalStringListValue(check.Severity)

	return stateModel
}

func (r *CheckCodeIssueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_code_issue"
}

func (r *CheckCodeIssueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 1,
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Check Code Issue Resource",

		Attributes: CheckBaseAttributes(map[string]schema.Attribute{
			"constraint": schema.StringAttribute{
				Description: "The type of constraint used in evaluation the code issues check.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllCheckCodeIssueConstraintEnum...)},
			},
			"issue_name": schema.StringAttribute{
				Description: "The issue name used for code issue lookup.",
				Optional:    true,
			},
			"issue_type": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The types of code issues to consider.",
				Optional:    true,
			},
			"max_allowed": schema.Int64Attribute{
				Description: "The threshold count of code issues beyond which the check starts failing.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"resolution_time": schema.SingleNestedAttribute{
				Description: "Defines the minimum frequency of the updates.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"unit": schema.StringAttribute{
						Description: "The name of duration of time.",
						Required:    true,
						Validators:  []validator.String{stringvalidator.OneOf(opslevel.AllCodeIssueResolutionTimeUnitEnum...)},
					},
					"value": schema.Int64Attribute{
						Description: "The count value of the specified unit.",
						Required:    true,
					},
				},
			},
			"severity": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The severity levels of the issue.",
				Optional:    true,
			},
		}),
	}
}

func (r *CheckCodeIssueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CheckCodeIssueResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckCodeIssueCreateInput{
		CategoryId: asID(planModel.Category),
		Constraint: opslevel.CheckCodeIssueConstraintEnum(planModel.Constraint.ValueString()),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   nullableID(planModel.Filter.ValueStringPointer()),
		IssueName:  nullable(planModel.IssueName.ValueStringPointer()),
		LevelId:    asID(planModel.Level),
		Name:       planModel.Name.ValueString(),
		Notes:      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:    nullableID(planModel.Owner.ValueStringPointer()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}
	if !planModel.IssueType.IsNull() {
		issueType, _ := ListValueToStringSlice(ctx, planModel.IssueType)
		input.IssueType = opslevel.RefOf(issueType)
	}
	if !planModel.MaxAllowed.IsNull() {
		input.MaxAllowed = refOf(int(planModel.MaxAllowed.ValueInt64()))
	}
	if !planModel.ResolutionTime.IsNull() {
		attrs := planModel.ResolutionTime.Attributes()
		input.ResolutionTime = refOf(opslevel.CodeIssueResolutionTimeInput{
			Unit:  opslevel.CodeIssueResolutionTimeUnitEnum(attrs["unit"].(basetypes.StringValue).ValueString()),
			Value: int(attrs["value"].(basetypes.Int64Value).ValueInt64()),
		})
	}
	if !planModel.Severity.IsNull() {
		severity, _ := ListValueToStringSlice(ctx, planModel.Severity)
		input.Severity = opslevel.RefOf(severity)
	}

	data, err := r.client.CreateCheckCodeIssue(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create check_code_issue, got error: %s", err))
		return
	}

	stateModel := NewCheckCodeIssueResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "created a check_code_issue resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckCodeIssueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CheckCodeIssueResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetCheck(asID(stateModel.Id))
	if err != nil {
		if (data == nil || data.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check_code_issue, got error: %s", err))
		return
	}
	stateModel = NewCheckCodeIssueResourceModel(ctx, *data, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckCodeIssueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CheckCodeIssueResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.CheckCodeIssueUpdateInput{
		CategoryId: opslevel.RefOf(asID(planModel.Category)),
		Constraint: opslevel.CheckCodeIssueConstraintEnum(planModel.Constraint.ValueString()),
		Enabled:    nullable(planModel.Enabled.ValueBoolPointer()),
		FilterId:   nullableID(planModel.Filter.ValueStringPointer()),
		Id:         asID(planModel.Id),
		IssueName:  nullable(planModel.IssueName.ValueStringPointer()),
		LevelId:    opslevel.RefOf(asID(planModel.Level)),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Notes:      opslevel.NewString(planModel.Notes.ValueString()),
		OwnerId:    nullableID(planModel.Owner.ValueStringPointer()),
	}
	if !planModel.EnableOn.IsNull() {
		enabledOn, err := iso8601.ParseString(planModel.EnableOn.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error", err.Error())
		}
		input.EnableOn = opslevel.RefOf(iso8601.Time{Time: enabledOn})
	}
	if planModel.IssueType.IsNull() {
		input.IssueType = opslevel.NewNullOf[[]string]()
	} else {
		issueType, _ := ListValueToStringSlice(ctx, planModel.IssueType)
		input.IssueType = opslevel.NewNullableFrom(issueType)
	}
	if !planModel.MaxAllowed.IsNull() {
		input.MaxAllowed = refOf(int(planModel.MaxAllowed.ValueInt64()))
	}
	if planModel.ResolutionTime.IsNull() {
		input.ResolutionTime = (*opslevel.CodeIssueResolutionTimeInput)(nil)
	} else {
		attrs := planModel.ResolutionTime.Attributes()
		input.ResolutionTime = refOf(opslevel.CodeIssueResolutionTimeInput{
			Unit:  opslevel.CodeIssueResolutionTimeUnitEnum(attrs["unit"].(basetypes.StringValue).ValueString()),
			Value: int(attrs["value"].(basetypes.Int64Value).ValueInt64()),
		})
	}
	if planModel.Severity.IsNull() {
		input.Severity = opslevel.NewNullOf[[]string]()
	} else {
		severity, _ := ListValueToStringSlice(ctx, planModel.Severity)
		input.Severity = opslevel.NewNullableFrom(severity)
	}

	data, err := r.client.UpdateCheckCodeIssue(input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update check_code_issue, got error: %s", err))
		return
	}

	stateModel := NewCheckCodeIssueResourceModel(ctx, *data, planModel)

	tflog.Trace(ctx, "updated a check_code_issue resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *CheckCodeIssueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CheckCodeIssueResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCheck(asID(stateModel.Id))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check_code_issue, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check_code_issue resource")
}

func (r *CheckCodeIssueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
