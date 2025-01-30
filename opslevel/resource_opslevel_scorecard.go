package opslevel

import (
	"context"
	"fmt"

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
	CategoryIds                 types.List   `tfsdk:"categories"`
	Description                 types.String `tfsdk:"description"`
	FilterId                    types.String `tfsdk:"filter_id"`
	Id                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	OwnerId                     types.String `tfsdk:"owner_id"`
	PassingChecks               types.Int64  `tfsdk:"passing_checks"`
	ServiceCount                types.Int64  `tfsdk:"service_count"`
	TotalChecks                 types.Int64  `tfsdk:"total_checks"`
}

func NewScorecardResourceModel(ctx context.Context, scorecard opslevel.Scorecard, categoryIds []string, givenModel ScorecardResourceModel) ScorecardResourceModel {
	scorecardDataSourceModel := ScorecardResourceModel{
		AffectsOverallServiceLevels: types.BoolValue(scorecard.AffectsOverallServiceLevels),
		Description:                 StringValueFromResourceAndModelField(scorecard.Description, givenModel.Description),
		FilterId:                    OptionalStringValue(string(scorecard.Filter.Id)),
		Id:                          ComputedStringValue(string(scorecard.Id)),
		Name:                        RequiredStringValue(scorecard.Name),
		OwnerId:                     RequiredStringValue(string(scorecard.Owner.Id())),
		PassingChecks:               types.Int64Value(int64(scorecard.PassingChecks)),
		ServiceCount:                types.Int64Value(int64(scorecard.ServiceCount)),
		TotalChecks:                 types.Int64Value(int64(scorecard.TotalChecks)),
	}

	scorecardDataSourceModel.CategoryIds = OptionalStringListValue(categoryIds)
	scorecardDataSourceModel.Aliases = OptionalStringListValue(scorecard.Aliases)

	return scorecardDataSourceModel
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
			"categories": schema.ListAttribute{
				Description: "The ids of the categories on this scorecard.",
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
	planModel := read[ScorecardResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	scorecard, err := r.client.CreateScorecard(opslevel.ScorecardInput{
		AffectsOverallServiceLevels: nullable(planModel.AffectsOverallServiceLevels.ValueBoolPointer()),
		Description:                 nullable(planModel.Description.ValueStringPointer()),
		FilterId:                    nullable(opslevel.NewID(planModel.FilterId.ValueString())),
		Name:                        planModel.Name.ValueString(),
		OwnerId:                     opslevel.ID(planModel.OwnerId.ValueString()),
	})
	if err != nil || scorecard == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create scorecard, got error: %s", err))
		return
	}
	categoryIds, err := getScorecardCategoyIds(r.client, *scorecard)
	if err != nil {
		resp.Diagnostics.AddWarning("opslevel client error", fmt.Sprintf("Unable to retrieve category ids from scorecard, got error: %s", err))
	}
	createdScorecardResourceModel := NewScorecardResourceModel(ctx, *scorecard, categoryIds, planModel)

	tflog.Trace(ctx, "created a scorecard resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdScorecardResourceModel)...)
}

func getScorecardCategoyIds(client *opslevel.Client, scorecard opslevel.Scorecard) ([]string, error) {
	var categoryIds []string

	result, err := scorecard.ListCategories(client, nil)
	if err != nil {
		return categoryIds, err
	}
	for _, category := range result.Nodes {
		categoryIds = append(categoryIds, string(category.Id))
	}

	return categoryIds, nil
}

func (r *ScorecardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[ScorecardResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	readScorecard, err := r.client.GetScorecard(stateModel.Id.ValueString())
	if err != nil || readScorecard == nil {
		if (readScorecard == nil || readScorecard.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read scorecard, got error: %s", err))
		return
	}
	categoryIds, err := getScorecardCategoyIds(r.client, *readScorecard)
	if err != nil {
		resp.Diagnostics.AddWarning("opslevel client error", fmt.Sprintf("Unable to retrieve category ids from scorecard, got error: %s", err))
	}
	readScorecardResourceModel := NewScorecardResourceModel(ctx, *readScorecard, categoryIds, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readScorecardResourceModel)...)
}

func (r *ScorecardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[ScorecardResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[ScorecardResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.ScorecardInput{
		AffectsOverallServiceLevels: nullable(planModel.AffectsOverallServiceLevels.ValueBoolPointer()),
		Description:                 nullable(planModel.Description.ValueStringPointer()),
		Name:                        planModel.Name.ValueString(),
		OwnerId:                     opslevel.ID(planModel.OwnerId.ValueString()),
	}
	if !planModel.FilterId.IsNull() {
		input.FilterId = nullable(opslevel.NewID(planModel.FilterId.ValueString()))
	} else if !stateModel.FilterId.IsNull() { // Then unset
		input.FilterId = nullable(opslevel.NewID(""))
	}

	scorecard, err := r.client.UpdateScorecard(planModel.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update scorecard, got error: %s", err))
		return
	}
	categoryIds, err := getScorecardCategoyIds(r.client, *scorecard)
	if err != nil {
		resp.Diagnostics.AddWarning("opslevel client error", fmt.Sprintf("Unable to retrieve category ids from scorecard, got error: %s", err))
	}
	updatedScorecardResourceModel := NewScorecardResourceModel(ctx, *scorecard, categoryIds, planModel)

	tflog.Trace(ctx, "updated a scorecard resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedScorecardResourceModel)...)
}

func (r *ScorecardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[ScorecardResourceModel](ctx, &resp.Diagnostics, req.State)
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
