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
	Alias      types.String `tfsdk:"alias"`
	Id         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Name       types.String `tfsdk:"name"`
}

func NewRepositoryDataSourceModel(ctx context.Context, repository opslevel.Repository, identifier string) RepositoryDataSourceModel {
	return RepositoryDataSourceModel{
		Alias:      types.StringValue(repository.DefaultAlias),
		Id:         types.StringValue(string(repository.Id)),
		Identifier: types.StringValue(identifier),
		Name:       types.StringValue(repository.Name),
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
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the repository.",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The alias or id of the repository.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the repository.",
				Computed:    true,
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
	if opslevel.IsID(data.Identifier.ValueString()) {
		repository, err = d.client.GetRepository(opslevel.ID(data.Identifier.ValueString()))
	} else {
		repository, err = d.client.GetRepositoryWithAlias(data.Identifier.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to read repository datasource, got error: %s", err))
		return
	}

	repositoryDataModel := NewRepositoryDataSourceModel(ctx, *repository, data.Identifier.ValueString())

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Repository data source")
	// resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &repositoryDataModel)...)
}
