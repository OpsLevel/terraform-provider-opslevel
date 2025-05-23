package opslevel

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
)

// Ensure LevelDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &LevelDataSource{}

func NewLevelDataSource() datasource.DataSource {
	return &LevelDataSource{}
}

// LevelDataSource manages a Level data source.
type LevelDataSource struct {
	CommonDataSourceClient
}

// levelDataSourceWithFilterModel describes the data source data model.
type levelDataSourceWithFilterModel struct {
	Alias  types.String     `tfsdk:"alias"`
	Filter filterBlockModel `tfsdk:"filter"`
	Id     types.String     `tfsdk:"id"`
	Index  types.Int64      `tfsdk:"index"`
	Name   types.String     `tfsdk:"name"`
}

func NewLevelDataSourceWithFilterModel(ctx context.Context, level opslevel.Level, filter filterBlockModel) levelDataSourceWithFilterModel {
	return levelDataSourceWithFilterModel{
		Alias:  ComputedStringValue(level.Alias),
		Filter: filter,
		Id:     ComputedStringValue(string(level.Id)),
		Index:  types.Int64Value(int64(level.Index)),
		Name:   ComputedStringValue(level.Name),
	}
}

var rubricLevelSchemaAttrs = map[string]schema.Attribute{
	"alias": schema.StringAttribute{
		MarkdownDescription: "An alias of the rubric level to find by.",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this resource.",
		Computed:            true,
	},
	"index": schema.Int64Attribute{
		MarkdownDescription: "An integer allowing this level to be inserted between others.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description: "The display name of the rubric level.",
		Computed:    true,
	},
}

func (d *LevelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rubric_level"
}

func (d *LevelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"alias", "id", "index", "name"}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rubric Level data source",

		Attributes: rubricLevelSchemaAttrs,
		Blocks: map[string]schema.Block{
			"filter": getDatasourceFilter(validFieldNames),
		},
	}
}

func (d *LevelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	configModel := read[levelDataSourceWithFilterModel](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	levels, err := d.client.ListLevels(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rubric_level datasource, got error: %s", err))
		return
	}

	level, err := filterRubricLevels(levels.Nodes, configModel.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to filter rubric_level datasource, got error: %s", err))
		return
	}

	stateModel := NewLevelDataSourceWithFilterModel(ctx, *level, configModel.Filter)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Rubric Level data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func filterRubricLevels(levels []opslevel.Level, filter filterBlockModel) (*opslevel.Level, error) {
	if filter.Value.Equal(types.StringValue("")) {
		return nil, fmt.Errorf("Please provide a non-empty value for filter's value")
	}
	for _, level := range levels {
		switch filter.Field.ValueString() {
		case "alias":
			if filter.Value.Equal(types.StringValue(level.Alias)) {
				return &level, nil
			}
		case "id":
			if filter.Value.Equal(types.StringValue(string(level.Id))) {
				return &level, nil
			}
		case "index":
			index := strconv.Itoa(int(level.Index))
			if filter.Value.Equal(types.StringValue(index)) {
				return &level, nil
			}
		case "name":
			if filter.Value.Equal(types.StringValue(level.Name)) {
				return &level, nil
			}

		}
	}

	return nil, fmt.Errorf("Unable to find rubric level with: %s==%s", filter.Field, filter.Value)
}
