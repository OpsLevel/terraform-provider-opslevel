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

// Ensure CategoryDataSourcesAll implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &CategoryDataSourcesAll{}

func NewCategoryDataSourcesAll() datasource.DataSource {
	return &CategoryDataSourcesAll{}
}

// CategoryDataSourcesAll manages a Category data source.
type CategoryDataSourcesAll struct {
	CommonDataSourceClient
}

// categoryDataSourceModel describes the data source data model.
type categoryDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// CategoryDataSourcesModel describes the data source data model.
type CategoryDataSourcesModel struct {
	RubricCategories []categoryDataSourceModel `tfsdk:"rubric_categories"`
}

func NewCategoryDataSourcesModel(categories []opslevel.Category) CategoryDataSourcesModel {
	rubricCategories := []categoryDataSourceModel{}
	for _, category := range categories {
		rubricCategory := categoryDataSourceModel{
			Id:   ComputedStringValue(string(category.Id)),
			Name: ComputedStringValue(category.Name),
		}
		rubricCategories = append(rubricCategories, rubricCategory)
	}
	return CategoryDataSourcesModel{RubricCategories: rubricCategories}
}

func (d *CategoryDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rubric_categories"
}

func (d *CategoryDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rubric Category data sources",

		Attributes: map[string]schema.Attribute{
			"rubric_categories": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: rubricCategorySchemaAttrs,
				},
				Description: "List of Rubric Category data sources",
				Computed:    true,
			},
		},
	}
}

func (d *CategoryDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel CategoryDataSourcesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	categories, err := d.client.ListCategories(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list rubric_categories datasource, got error: %s", err))
		return
	}
	stateModel = NewCategoryDataSourcesModel(categories.Nodes)

	// Save data into Terraform state
	tflog.Trace(ctx, "listed all rubric_categories data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
