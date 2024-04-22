package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure RepositoryDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &RepositoriesDataSourcesAll{}

func NewRepositoriesDataSourceAll() datasource.DataSource {
	return &RepositoriesDataSourcesAll{}
}

// RepositoryDataSource manages a Repository data source.
type RepositoriesDataSourcesAll struct {
	CommonDataSourceClient
}

// RepositoryDataSourceModel describes the data source data model.
type RepositoriesDataSourcesAllModel struct {
	Filter       *filterBlockModel           `tfsdk:"filter"`
	Repositories []RepositoryDataSourceModel `tfsdk:"repositories"`
}

func NewRepositoriesDataSourcesAllModel(repositories []opslevel.Repository) RepositoriesDataSourcesAllModel {
	repositoriesModels := []RepositoryDataSourceModel{}
	for _, repo := range repositories {
		repoModel := NewRepositoryDataSourceModel(repo)
		repositoriesModels = append(repositoriesModels, repoModel)
	}
	return RepositoriesDataSourcesAllModel{Repositories: repositoriesModels}
}

func (d *RepositoriesDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repositories"
}

func (d *RepositoriesDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	validFieldNames := []string{"tier"}
	resp.Schema = schema.Schema{
		MarkdownDescription: "List of all Repository data sources",

		Attributes: map[string]schema.Attribute{
			"filter": schema.SingleNestedAttribute{
				Description: fmt.Sprintf(
					"Used to filter repositories by one of '%s'",
					strings.Join(validFieldNames, "`, `"),
				),
				Optional:   true,
				Attributes: FilterAttrs(validFieldNames),
			},
			"repositories": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: repositoryDatasourceSchemaAttrs,
				},
				Description: "List of Repository data sources",
				Computed:    true,
			},
		},
	}
}

func (d *RepositoriesDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel RepositoriesDataSourcesAllModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repos *opslevel.RepositoryConnection
	var err error
	if planModel.Filter != nil && planModel.Filter.Field.ValueString() == "tier" {
		repos, err = d.client.ListRepositoriesWithTier(planModel.Filter.Value.ValueString(), nil)
	} else {
		repos, err = d.client.ListRepositories(nil)
	}
	if err != nil || repos == nil || repos.Nodes == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read OpsLevel Repositories data source, got error: %s", err))
		return
	}

	stateModel = NewRepositoriesDataSourcesAllModel(repos.Nodes)
	stateModel.Filter = planModel.Filter

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Repository data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
