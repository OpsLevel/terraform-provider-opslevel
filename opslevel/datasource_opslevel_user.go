package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
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

type userWithIdentifierDataSourceModel struct {
	Email      types.String `tfsdk:"email"`
	Id         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Name       types.String `tfsdk:"name"`
	Role       types.String `tfsdk:"role"`
}

func newUserWithIdentifierDataSourceModel(user opslevel.User, identifier string) userWithIdentifierDataSourceModel {
	return userWithIdentifierDataSourceModel{
		Email:      types.StringValue(user.Email),
		Id:         types.StringValue(string(user.Id)),
		Identifier: types.StringValue(identifier),
		Name:       types.StringValue(user.Name),
		Role:       types.StringValue(string(user.Role)),
	}
}

type userDataSourceModel struct {
	Email types.String `tfsdk:"email"`
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Role  types.String `tfsdk:"role"`
}

func newUserDataSourceModel(user opslevel.User) userDataSourceModel {
	return userDataSourceModel{
		Email: types.StringValue(user.Email),
		Id:    types.StringValue(string(user.Id)),
		Name:  types.StringValue(user.Name),
		Role:  types.StringValue(string(user.Role)),
	}
}

var userDatasourceSchemaAttrs = map[string]schema.Attribute{
	"email": schema.StringAttribute{
		Description: "The email of the user.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The unique identifier for the user.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the user.",
		Computed:    true,
	},
	"role": schema.StringAttribute{
		Description: "The user's assigned role.",
		Computed:    true,
	},
}

func userAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range userDatasourceSchemaAttrs {
		attrs[key] = value
	}
	return attrs
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User data source",

		Attributes: userAttributes(map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Description: "The id or email of the user to find.",
				Required:    true,
			},
		}),
	}
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := read[userWithIdentifierDataSourceModel](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.GetUser(data.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}
	userDataModel := newUserWithIdentifierDataSourceModel(*user, data.Identifier.ValueString())

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel User data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &userDataModel)...)
}
