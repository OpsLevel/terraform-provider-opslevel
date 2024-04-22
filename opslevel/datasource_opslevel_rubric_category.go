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

// Ensure CategoryDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &CategoryDataSource{}

func NewCategoryDataSource() datasource.DataSource {
	return &CategoryDataSource{}
}

// CategoryDataSource manages a Category data source.
type CategoryDataSource struct {
	CommonDataSourceClient
}

// categoryDataSourceWithFilterModel contains categoryDataSourceModel fields and a filterBlockModel
type categoryDataSourceWithFilterModel struct {
	Filter filterBlockModel `tfsdk:"filter"`
	Id     types.String     `tfsdk:"id"`
	Name   types.String     `tfsdk:"name"`
}

func NewCategoryDataSourceWithFilterModel(category opslevel.Category, filter filterBlockModel) categoryDataSourceWithFilterModel {
	return categoryDataSourceWithFilterModel{
		Id:     ComputedStringValue(string(category.Id)),
		Name:   ComputedStringValue(category.Name),
		Filter: filter,
	}
}

func (d *CategoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rubric_category"
}

var rubricCategorySchemaAttrs = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this resource.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the rubric category.",
		Computed:    true,
	},
}

func (d *CategoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"id", "name"}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rubric Category data source",

		Attributes: rubricCategorySchemaAttrs,
		Blocks: map[string]schema.Block{
			"filter": getDatasourceFilter(validFieldNames),
		},
	}
}

func (d *CategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel categoryDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	categories, err := d.client.ListCategories(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rubric_category datasource, got error: %s", err))
		return
	}

	category, err := filterRubricCategories(categories.Nodes, planModel.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to filter rubric_category datasource, got error: %s", err))
		return
	}

	stateModel = NewCategoryDataSourceWithFilterModel(*category, planModel.Filter)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Rubric Category data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func filterRubricCategories(categories []opslevel.Category, filter filterBlockModel) (*opslevel.Category, error) {
	if filter.Value.Equal(types.StringValue("")) {
		return nil, fmt.Errorf("please provide a non-empty value for filter's value")
	}
	for _, category := range categories {
		switch filter.Field.ValueString() {
		case "id":
			if filter.Value.Equal(types.StringValue(string(category.Id))) {
				return &category, nil
			}
		case "name":
			if filter.Value.Equal(types.StringValue(category.Name)) {
				return &category, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find category with: %s==%s", filter.Field, filter.Value)
}
