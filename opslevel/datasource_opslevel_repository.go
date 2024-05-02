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

// Ensure RepositoryDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &RepositoryDataSource{}

func NewRepositoryDataSource() datasource.DataSource {
	return &RepositoryDataSource{}
}

// RepositoryDataSource manages a Repository data source.
type RepositoryDataSource struct {
	CommonDataSourceClient
}

// LanguagesModel describes the model for the Languages data of the repository.
type LanguagesModel struct {
	Name  types.String  `tfsdk:"name"`
	Usage types.Float64 `tfsdk:"usage"`
}

// RepositoryDataSourceModel describes the data source data model.
type RepositoryDataSourceModel struct {
	Alias     types.String     `tfsdk:"alias"`
	Id        types.String     `tfsdk:"id"`
	Name      types.String     `tfsdk:"name"`
	Url       types.String     `tfsdk:"url"`
	Languages []LanguagesModel `tfsdk:"languages"`
}

// LanguagesValue function converts the raw opslevel data to terraform friendly format
func LanguagesValue(value []opslevel.Language) []LanguagesModel {
	var languages []LanguagesModel

	if len(value) == 0 {
		return []LanguagesModel{}
	}

	for _, lang := range value {
		language := LanguagesModel{
			Name:  types.StringValue(lang.Name),
			Usage: types.Float64Value(lang.Usage),
		}

		languages = append(languages, language)
	}
	return languages
}

func NewRepositoryDataSourceModel(repository opslevel.Repository) RepositoryDataSourceModel {
	return RepositoryDataSourceModel{
		Alias:     OptionalStringValue(repository.DefaultAlias),
		Id:        OptionalStringValue(string(repository.Id)),
		Name:      ComputedStringValue(repository.Name),
		Url:       ComputedStringValue(repository.Url),
		Languages: LanguagesValue(repository.Languages),
	}
}

var repositoryDatasourceSchemaAttrs = map[string]schema.Attribute{
	"alias": schema.StringAttribute{
		MarkdownDescription: "The human-friendly, unique identifier for the repository.",
		Optional:            true,
		Computed:            true,
	},
	"id": schema.StringAttribute{
		Description: "The unique identifier for the repository.",
		Optional:    true,
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The display name of the repository.",
		Computed:    true,
	},
	"url": schema.StringAttribute{
		Description: "The url of the the repository.",
		Computed:    true,
	},
	"languages": schema.ListNestedAttribute{
		Description: "The list of programming languages used in the repository.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Optional: true,
				},
				"usage": schema.Float64Attribute{
					Optional: true,
				},
			},
		},
	},
}

func (d *RepositoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *RepositoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Repository data source",

		Attributes: repositoryDatasourceSchemaAttrs,
	}
}

func (d *RepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RepositoryDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var repository *opslevel.Repository

	if data.Alias.ValueString() != "" {
		repository, err = d.client.GetRepositoryWithAlias(data.Alias.ValueString())
	} else if data.Id.ValueString() != "" {
		repository, err = d.client.GetRepository(opslevel.ID(data.Id.ValueString()))
	} else {
		resp.Diagnostics.AddError("Config Error", "please provide a value for alias or id")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to read repository datasource, got error: %s", err))
		return
	}
	if repository == nil || repository.Id == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to find repository with alias=`%s` or id=`%s`", data.Alias.ValueString(), data.Id.ValueString()))
		return
	}

	repositoryDataModel := NewRepositoryDataSourceModel(*repository)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Repository data source")
	resp.Diagnostics.Append(resp.Diagnostics...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &repositoryDataModel)...)
}
