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

var _ resource.ResourceWithConfigure = &ScorecardResource{}

var _ resource.ResourceWithImportState = &ScorecardResource{}

func NewScorecardResource() resource.Resource {
	return &ScorecardResource{}
}

// ScorecardResource defines the resource implementation.
type ScorecardResource struct {
	CommonResourceClient
}

// ScorecardResourceModel describes the Scorecard managed resource.
type ScorecardResourceModel struct {
	AffectsOverallServiceLevels types.Bool   `tfsdk:"affects_overall_service_levels"`
	Aliases                     types.List   `tfsdk:"aliases"`
	Description                 types.String `tfsdk:"description"`
	FilterId                    types.String `tfsdk:"filter_id"`
	Id                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	OwnerId                     types.String `tfsdk:"owner_id"`
	PassingChecks               types.Int64  `tfsdk:"passing_checks"`
	ServiceCount                types.Int64  `tfsdk:"service_count"`
	TotalChecks                 types.Int64  `tfsdk:"total_checks"`
}

func NewScorecardResourceModel(ctx context.Context, scorecard opslevel.Scorecard, givenModel ScorecardResourceModel) (ScorecardResourceModel, diag.Diagnostics) {
	scorecardDataSourceModel := ScorecardResourceModel{
		AffectsOverallServiceLevels: types.BoolValue(scorecard.AffectsOverallServiceLevels),
		Description:                 StringValueFromResourceAndModelField(scorecard.Description, givenModel.Description),
		FilterId:                    OptionalStringValue(string(scorecard.Filter.Id)),
		Id:                          ComputedStringValue(string(scorecard.Id)),
		Name:                        RequiredStringValue(scorecard.Name),
		OwnerId:                     RequiredStringValue(string(scorecard.Owner.Id())),
		PassingChecks:               types.Int64Value(int64(scorecard.PassingChecks)),
		ServiceCount:                types.Int64Value(int64(scorecard.ServiceCount)),
		TotalChecks:                 types.Int64Value(int64(scorecard.ChecksCount)),
	}

	scorecardAliases, diags := types.ListValueFrom(ctx, types.StringType, scorecard.Aliases)
	scorecardDataSourceModel.Aliases = scorecardAliases

	return scorecardDataSourceModel, diags
}

func (r *ScorecardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard"
}

func (r *ScorecardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Scorecard Resource",

		Attributes: map[string]schema.Attribute{
			"affects_overall_service_levels": schema.BoolAttribute{
				Description: "Specifies whether the checks on this scorecard affect services' overall maturity level.",
				Required:    true,
			},
			"aliases": schema.ListAttribute{
				Description: "The scorecard's aliases.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "The scorecard's description.",
				Optional:    true,
			},
			"filter_id": schema.StringAttribute{
				Description: "The scorecard's filter.",
				Optional:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the scorecard.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The scorecard's name.",
				Required:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The scorecard's owner.",
				Required:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"passing_checks": schema.Int64Attribute{
				Description: "The scorecard's number of checks that are passing.",
				Computed:    true,
			},
			"service_count": schema.Int64Attribute{
				Description: "The scorecard's number of services matched.",
				Computed:    true,
			},
			"total_checks": schema.Int64Attribute{
				Description: "The scorecard's total number of checks.",
				Computed:    true,
			},
		},
	}
}

func (r *ScorecardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel ScorecardResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	scorecard, err := r.client.CreateScorecard(opslevel.ScorecardInput{
		AffectsOverallServiceLevels: planModel.AffectsOverallServiceLevels.ValueBoolPointer(),
		Description:                 planModel.Description.ValueStringPointer(),
		FilterId:                    opslevel.NewID(planModel.FilterId.ValueString()),
		Name:                        planModel.Name.ValueString(),
		OwnerId:                     opslevel.ID(planModel.OwnerId.ValueString()),
	})
	if err != nil || scorecard == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create scorecard, got error: %s", err))
		return
	}
	createdScorecardResourceModel, diags := NewScorecardResourceModel(ctx, *scorecard, planModel)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Trace(ctx, "created a scorecard resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdScorecardResourceModel)...)
}

func (r *ScorecardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel ScorecardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	readScorecard, err := r.client.GetScorecard(stateModel.Id.ValueString())
	if err != nil || readScorecard == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read scorecard, got error: %s", err))
		return
	}
	readScorecardResourceModel, diags := NewScorecardResourceModel(ctx, *readScorecard, stateModel)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readScorecardResourceModel)...)
}

func (r *ScorecardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel ScorecardResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scorecard, err := r.client.UpdateScorecard(planModel.Id.ValueString(), opslevel.ScorecardInput{
		AffectsOverallServiceLevels: planModel.AffectsOverallServiceLevels.ValueBoolPointer(),
		Description:                 opslevel.RefOf(planModel.Description.ValueString()),
		FilterId:                    opslevel.NewID(planModel.FilterId.ValueString()),
		Name:                        planModel.Name.ValueString(),
		OwnerId:                     opslevel.ID(planModel.OwnerId.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update scorecard, got error: %s", err))
		return
	}
	updatedScorecardResourceModel, diags := NewScorecardResourceModel(ctx, *scorecard, planModel)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Trace(ctx, "updated a scorecard resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedScorecardResourceModel)...)
}

func (r *ScorecardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ScorecardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteScorecard(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete scorecard, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a scorecard resource")
}

func (r *ScorecardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
