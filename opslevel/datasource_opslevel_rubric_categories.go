package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure CategoryDataSourcesAll implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &CategoryDataSourcesAll{}

func NewCategoryDataSourcesAll() datasource.DataSource {
	return &CategoryDataSourcesAll{}
}

// CategoryDataSourcesAll manages a Category data source.
type CategoryDataSourcesAll struct {
	CommonDataSourceClient
}

// rubricCategoryObjectType is derived from DomainDataSourceModel, needed for lists
var rubricCategoryObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	},
}

// CategoryDataSourcesModel describes the data source data model.
type CategoryDataSourcesModel struct {
	RubricCategories types.List `tfsdk:"rubric_categories"`
}

func (d *CategoryDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rubric_categories"
}

func (d *CategoryDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rubric Category data sources",

		Attributes: map[string]schema.Attribute{
			"rubric_categories": schema.ListAttribute{
				ElementType: rubricCategoryObjectType,
				Description: "List of Rubric Category data sources",
				Computed:    true,
			},
		},
	}
}

func (d *CategoryDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CategoryDataSourcesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	categories, err := d.client.ListCategories(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list rubric_categories datasource, got error: %s", err))
		return
	}
	categoriesListValue, diags := allCategoriesToListValue(ctx, categories.Nodes)
	resp.Diagnostics.Append(diags...)

	data.RubricCategories = categoriesListValue

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all rubric_categories data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func allCategoriesToListValue(ctx context.Context, opslevelCategories []opslevel.Category) (basetypes.ListValue, diag.Diagnostics) {
	categories := make([]attr.Value, len(opslevelCategories))

	for idx, category := range opslevelCategories {
		categoryObject, diags := categoryToObject(ctx, category)
		if diags != nil && diags.HasError() {
			return basetypes.NewListNull(rubricCategoryObjectType), diags
		}
		categories[idx] = categoryObject
	}

	result, diags := basetypes.NewListValue(
		rubricCategoryObjectType,
		categories,
	)
	if diags != nil && diags.HasError() {
		return basetypes.NewListNull(rubricCategoryObjectType), diags
	}

	return result, nil
}

// categoryToObject converts an opslevel.Category to a basetypes.ObjectValue
func categoryToObject(ctx context.Context, opslevelCategory opslevel.Category) (basetypes.ObjectValue, diag.Diagnostics) {
	categoryObject := NewCategoryDataSourceModel(ctx, opslevelCategory)

	categoryModel := make(map[string]attr.Value)
	categoryModel["id"] = categoryObject.Id
	categoryModel["name"] = categoryObject.Name

	parsedCategory, diags := types.ObjectValue(rubricCategoryObjectType.AttrTypes, categoryModel)
	if diags != nil && diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	return parsedCategory, nil
}
