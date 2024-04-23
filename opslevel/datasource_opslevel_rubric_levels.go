package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure LevelDataSourcesAll implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &LevelDataSourcesAll{}

func NewLevelDataSourcesAll() datasource.DataSource {
	return &LevelDataSourcesAll{}
}

// LevelDataSourcesAll manages a Level data source.
type LevelDataSourcesAll struct {
	CommonDataSourceClient
}

// levelDataSourceModel describes the data source data model.
type levelDataSourceModel struct {
	Alias types.String `tfsdk:"alias"`
	Id    types.String `tfsdk:"id"`
	Index types.Int64  `tfsdk:"index"`
	Name  types.String `tfsdk:"name"`
}

// LevelDataSourcesAllModel describes the data source data model.
type LevelDataSourcesAllModel struct {
	RubricLevels []levelDataSourceModel `tfsdk:"rubric_levels"`
}

func NewLevelDataSourcesAllModel(levels []opslevel.Level) LevelDataSourcesAllModel {
	rubricLevels := []levelDataSourceModel{}
	for _, level := range levels {
		rubricLevel := levelDataSourceModel{
			Alias: ComputedStringValue(level.Alias),
			Id:    ComputedStringValue(string(level.Id)),
			Index: types.Int64Value(int64(level.Index)),
			Name:  ComputedStringValue(level.Name),
		}
		rubricLevels = append(rubricLevels, rubricLevel)
	}
	return LevelDataSourcesAllModel{RubricLevels: rubricLevels}
}

func (d *LevelDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rubric_levels"
}

func (d *LevelDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rubric Level data sources",

		Attributes: map[string]schema.Attribute{
			"rubric_levels": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: rubricLevelSchemaAttrs,
				},
				Description: "List of Rubric Level data sources",
				Computed:    true,
			},
		},
	}
}

func (d *LevelDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel LevelDataSourcesAllModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	levels, err := d.client.ListLevels()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list rubric_levels datasource, got error: %s", err))
		return
	}
	stateModel = NewLevelDataSourcesAllModel(levels)

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all rubric_levels data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
