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

// RepositoryDataSourceModel describes the data source data model.
type RepositoryDataSourceModel struct {
	Alias types.String `tfsdk:"alias"`
	Id    types.String `tfsdk:"id"`
}

func NewRepositoryDataSourceModel(ctx context.Context, repository opslevel.Repository) RepositoryDataSourceModel {
	return RepositoryDataSourceModel{
		Alias: types.StringValue(repository.DefaultAlias),
		Id:    types.StringValue(string(repository.Id)),
	}
}

func (d *RepositoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *RepositoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Repository data source",

		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				MarkdownDescription: "The human-friendly, unique identifier for the repository.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the repository.",
				Optional:            true,
			},
		},
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
		// use id
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

	repositoryDataModel := NewRepositoryDataSourceModel(ctx, *repository)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Repository data source")
	// resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &repositoryDataModel)...)
}
