package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ datasource.DataSourceWithConfigure = &UserDataSourcesAll{}

func NewUserDataSourcesAll() datasource.DataSource {
	return &UserDataSourcesAll{}
}

type UserDataSourcesAll struct {
	CommonDataSourceClient
}

type userDataSourcesAllModel struct {
	Users []userDataSourceModel `tfsdk:"users"`
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
			"users": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: userAttributes(map[string]schema.Attribute{}),
				},
				Description: "List of user data sources",
				Computed:    true,
			},
		},
	}
}

func (d *UserDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var planModel, stateModel userDataSourcesAllModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	users, err := d.client.ListUsers(nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list users, got error: %s", err))
		return
	}
	stateModel = newUserDataSourcesAllModel(users.Nodes)

	tflog.Trace(ctx, "listed all OpsLevel User data sources")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
