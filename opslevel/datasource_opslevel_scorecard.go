package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure ScorecardDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &ScorecardDataSource{}

func NewScorecardDataSource() datasource.DataSource {
	return &ScorecardDataSource{}
}

// ScorecardDataSource manages a Scorecard data source.
type ScorecardDataSource struct {
	CommonDataSourceClient
}

var scorecardSchemaAttrs = map[string]schema.Attribute{
	"affects_overall_service_levels": schema.BoolAttribute{
		Description: "Specifies whether the checks on this scorecard affect services' overall maturity level.",
		Computed:    true,
	},
	"aliases": schema.ListAttribute{
		ElementType: types.StringType,
		Description: "The scorecard's aliases.",
		Computed:    true,
	},
	"categories": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: rubricCategorySchemaAttrs,
		},
		Description: "The scorecard's rubric categories.",
		Computed:    true,
	},
	"description": schema.StringAttribute{
		Description: "The scorecard's description.",
		Computed:    true,
	},
	"filter_id": schema.StringAttribute{
		MarkdownDescription: "The scorecard's filter.",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this resource.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description: "The scorecard's name.",
		Computed:    true,
	},
	"owner_id": schema.StringAttribute{
		Description: "The scorecard's owner id.",
		Computed:    true,
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
}

func ScorecardAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range scorecardSchemaAttrs {
		attrs[key] = value
	}
	return attrs
}

// scorecardDataSourceWithIdentifierModel describes the data source data model.
type scorecardDataSourceWithIdentifierModel struct {
	AffectsOverallServiceLevels types.Bool                `tfsdk:"affects_overall_service_levels"`
	Aliases                     types.List                `tfsdk:"aliases"`
	Categories                  []categoryDataSourceModel `tfsdk:"categories"`
	Description                 types.String              `tfsdk:"description"`
	FilterId                    types.String              `tfsdk:"filter_id"`
	Id                          types.String              `tfsdk:"id"`
	Identifier                  types.String              `tfsdk:"identifier"`
	Name                        types.String              `tfsdk:"name"`
	OwnerId                     types.String              `tfsdk:"owner_id"`
	PassingChecks               types.Int64               `tfsdk:"passing_checks"`
	ServiceCount                types.Int64               `tfsdk:"service_count"`
	TotalChecks                 types.Int64               `tfsdk:"total_checks"`
}

func NewScorecardDataSourceWithIdentifierModel(
	scorecard opslevel.Scorecard,
	identifier string,
	categoriesModel []categoryDataSourceModel,
) scorecardDataSourceWithIdentifierModel {
	scorecardAliases := OptionalStringListValue(scorecard.Aliases)
	return scorecardDataSourceWithIdentifierModel{
		AffectsOverallServiceLevels: types.BoolValue(scorecard.AffectsOverallServiceLevels),
		Aliases:                     scorecardAliases,
		Categories:                  categoriesModel,
		Description:                 ComputedStringValue(scorecard.Description),
		FilterId:                    ComputedStringValue(string(scorecard.Filter.Id)),
		Id:                          ComputedStringValue(string(scorecard.Id)),
		Identifier:                  ComputedStringValue(identifier),
		Name:                        ComputedStringValue(scorecard.Name),
		OwnerId:                     ComputedStringValue(string(scorecard.Owner.Id())),
		PassingChecks:               types.Int64Value(int64(scorecard.PassingChecks)),
		ServiceCount:                types.Int64Value(int64(scorecard.ServiceCount)),
		TotalChecks:                 types.Int64Value(int64(scorecard.ChecksCount)),
	}
}

func (d *ScorecardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scorecard"
}

func (d *ScorecardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Scorecard data source",

		Attributes: ScorecardAttributes(map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the scorecard to find.",
				Required:    true,
			},
		}),
	}
}

func (d *ScorecardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel scorecardDataSourceWithIdentifierModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scorecard, err := d.client.GetScorecard(planModel.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scorecard datasource, got error: %s", err))
		return
	}
	categoriesModel, diags := getCategoriesModelFromScorecard(d.client, scorecard)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	stateModel = NewScorecardDataSourceWithIdentifierModel(*scorecard, planModel.Identifier.ValueString(), categoriesModel)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Scorecard data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func getCategoriesModelFromScorecard(client *opslevel.Client, scorecard *opslevel.Scorecard) ([]categoryDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	categories, err := scorecard.ListCategories(client, nil)
	if err != nil || categories == nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to list categories from scorecard with id '%s', got error: %s", scorecard.Id, err))
	}
	categoriesModel := NewCategoryDataSourcesAllModel(categories.Nodes)
	return categoriesModel.RubricCategories, diags
}
