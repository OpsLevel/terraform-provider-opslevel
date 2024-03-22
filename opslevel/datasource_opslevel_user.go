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

// Ensure UserDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource manages a User data source.
type UserDataSource struct {
	CommonDataSourceClient
}

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
	Email      types.String `tfsdk:"email"`
	Id         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Name       types.String `tfsdk:"name"`
	Role       types.String `tfsdk:"role"`
}

func NewUserDataSourceModel(ctx context.Context, user opslevel.User, identifier string) UserDataSourceModel {
	return UserDataSourceModel{
		Email:      types.StringValue(user.Email),
		Id:         types.StringValue(string(user.Id)),
		Identifier: types.StringValue(identifier),
		Name:       types.StringValue(user.Name),
		Role:       types.StringValue(string(user.Role)),
	}
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User data source",

		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "The email of the user.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier for the user.",
				Computed:    true,
			},
			"identifier": schema.StringAttribute{
				Description: "The id or email of the user to find.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the user.",
				Computed:    true,
			},
			"role": schema.StringAttribute{
				Description: "The user's assigned role.",
				Computed:    true,
			},
		},
	}
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.GetUser(data.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}
	userDataModel := NewUserDataSourceModel(ctx, *user, data.Identifier.ValueString())

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel User data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &userDataModel)...)
}
