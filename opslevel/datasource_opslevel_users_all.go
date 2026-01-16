package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
)

var _ datasource.DataSourceWithConfigure = &UserDataSourcesAll{}

func NewUserDataSourcesAll() datasource.DataSource {
	return &UserDataSourcesAll{}
}

type UserDataSourcesAll struct {
	CommonDataSourceClient
}

type userDataSourcesAllModel struct {
	IgnoreDeactivated types.Bool            `tfsdk:"ignore_deactivated"`
	Users             []userDataSourceModel `tfsdk:"users"`
}

func newUserDataSourcesAllModel(users []opslevel.User) userDataSourcesAllModel {
	userModels := make([]userDataSourceModel, 0)
	for _, user := range users {
		userModel := newUserDataSourceModel(user)
		userModels = append(userModels, userModel)
	}
	return userDataSourcesAllModel{Users: userModels}
}

func (d *UserDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *UserDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List of all User data sources",

		Attributes: map[string]schema.Attribute{
			"ignore_deactivated": schema.BoolAttribute{
				Description: "Do not list deactivated users if set.",
				Optional:    true,
			},
			"users": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: userDatasourceSchemaAttrs,
				},
				Description: "List of user data sources",
				Computed:    true,
			},
		},
	}
}

func (d *UserDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var users *opslevel.UserConnection
	var err error

	planModel := read[userDataSourcesAllModel](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	if planModel.IgnoreDeactivated.ValueBool() {
		withoutDeactivedUsers := d.client.InitialPageVariablesPointer().WithoutDeactivedUsers()
		users, err = d.client.ListUsers(withoutDeactivedUsers)
	} else {
		users, err = d.client.ListUsers(nil)
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list users, got error: %s", err))
		return
	}
	stateModel := newUserDataSourcesAllModel(users.Nodes)
	stateModel.IgnoreDeactivated = planModel.IgnoreDeactivated

	tflog.Trace(ctx, "listed all OpsLevel User data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
