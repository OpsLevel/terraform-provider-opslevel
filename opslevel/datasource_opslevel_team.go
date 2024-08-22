package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure TeamDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &TeamDataSource{}

func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

// TeamDataSource manages a Team data source.
type TeamDataSource struct {
	CommonDataSourceClient
}

// teamDataSourceModel describes the data source data model.
type teamDataSourceModel struct {
	Alias       types.String      `tfsdk:"alias"`
	Contacts    types.List        `tfsdk:"contacts"`
	Id          types.String      `tfsdk:"id"`
	Members     []teamMemberModel `tfsdk:"members"`
	Name        types.String      `tfsdk:"name"`
	ParentAlias types.String      `tfsdk:"parent_alias"`
	ParentId    types.String      `tfsdk:"parent_id"`
}

type teamContactModel struct {
	Address     types.String `tfsdk:"address"`
	DisplayName types.String `tfsdk:"display_name"`
	DisplayType types.String `tfsdk:"display_yype"`
	ExternalId  types.String `tfsdk:"external_id"`
	Id          types.String `tfsdk:"id"`
	IsDefault   types.Bool   `tfsdk:"is_default"`
	Type        types.String `tfsdk:"type"`
}

var teamContactsNestedSchemaAttrs = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Description: "The contact address. Examples: 'support@company.com' for type email, 'https://opslevel.com' for type web.",
		Computed:    true,
	},
	"display_name": schema.StringAttribute{
		Description: "The name shown in the UI for the contact.",
		Computed:    true,
	},
	"display_type": schema.StringAttribute{
		Description: "The type shown in the UI for the contact.",
		Computed:    true,
	},
	"external_id": schema.StringAttribute{
		Description: "The remote identifier of the contact method.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The unique identifier for the contact.",
		Computed:    true,
	},
	"is_default": schema.BoolAttribute{
		Description: "Indicates if this address is a team's default for the given type.",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: fmt.Sprintf("The method of contact. One of [`%s`].",
			strings.Join(opslevel.AllContactType, "`, `"),
		),
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf(opslevel.AllContactType...),
		},
	},
}

func newTeamContactModel(contact opslevel.Contact) teamContactModel {
	return teamContactModel{
		Address:     ComputedStringValue(contact.Address),
		DisplayName: ComputedStringValue(contact.DisplayName),
		DisplayType: ComputedStringValue(contact.DisplayType),
		ExternalId:  ComputedStringValue(contact.ExternalId),
		Id:          ComputedStringValue(string(contact.Id)),
		IsDefault:   types.BoolValue(contact.IsDefault),
		Type:        ComputedStringValue(string(contact.Type)),
	}
}

func (tcm teamContactModel) asObjectValue() basetypes.ObjectValue {
	attrValues := map[string]attr.Value{
		"address":      tcm.Address,
		"display_name": tcm.DisplayName,
		"display_type": tcm.DisplayType,
		"external_id":  tcm.ExternalId,
		"id":           tcm.Id,
		"is_default":   tcm.IsDefault,
		"type":         tcm.Type,
	}
	return types.ObjectValueMust(teamContactAttrs(), attrValues)
}

func teamContactAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"address":      types.StringType,
		"display_name": types.StringType,
		"display_type": types.StringType,
		"external_id":  types.StringType,
		"id":           types.StringType,
		"is_default":   types.BoolType,
		"type":         types.StringType,
	}
}

var teamDatasourceSchemaAttrs = map[string]schema.Attribute{
	"contacts": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: teamContactsNestedSchemaAttrs,
		},
		Description: "The contacts for the team.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the Team.",
		Computed:    true,
	},
	"members": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: memberNestedSchemaAttrs,
		},
		Description: "List of team members on the team with email address and role.",
		Computed:    true,
	},
	"parent_alias": schema.StringAttribute{
		Description: "The alias of the parent team.",
		Computed:    true,
	},
	"parent_id": schema.StringAttribute{
		Description: "The id of the parent team.",
		Computed:    true,
	},
}

func teamAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range teamDatasourceSchemaAttrs {
		attrs[key] = value
	}
	return attrs
}

var memberNestedSchemaAttrs = map[string]schema.Attribute{
	"email": schema.StringAttribute{
		MarkdownDescription: "The email address of the team member.",
		Computed:            true,
	},
	"role": schema.StringAttribute{
		MarkdownDescription: "The role of the team member.",
		Computed:            true,
	},
}

type teamMemberModel struct {
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

func newTeamMemberModel(member opslevel.TeamMembership) teamMemberModel {
	return teamMemberModel{
		Email: ComputedStringValue(member.User.Email),
		Role:  ComputedStringValue(member.Role),
	}
}

func newTeamMembersAllModel(members []opslevel.TeamMembership) []teamMemberModel {
	membersModel := make([]teamMemberModel, 0)
	for _, member := range members {
		membersModel = append(membersModel, newTeamMemberModel(member))
	}
	return membersModel
}

func newTeamDataSourceModel(team opslevel.Team) teamDataSourceModel {
	teamDataSourceModel := teamDataSourceModel{
		Alias:       ComputedStringValue(team.Alias),
		Id:          ComputedStringValue(string(team.Id)),
		Name:        ComputedStringValue(team.Name),
		ParentAlias: ComputedStringValue(team.ParentTeam.Alias),
		ParentId:    ComputedStringValue(string(team.ParentTeam.Id)),
	}
	teamContactAttrTypes := teamContactAttrs()
	if len(team.Contacts) == 0 {
		teamDataSourceModel.Contacts = types.ListNull(types.ObjectType{AttrTypes: teamContactAttrTypes})
	} else {
		teamContactModels := []attr.Value{}
		for _, contact := range team.Contacts {
			teamContactModels = append(teamContactModels, newTeamContactModel(contact).asObjectValue())
		}
		teamDataSourceModel.Contacts = types.ListValueMust(types.ObjectType{AttrTypes: teamContactAttrTypes}, teamContactModels)
	}
	if team.Memberships != nil {
		teamDataSourceModel.Members = newTeamMembersAllModel(team.Memberships.Nodes)
	}
	return teamDataSourceModel
}

func (teamDataSource *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (teamDataSource *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Team data source",

		Attributes: teamAttributes(map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				MarkdownDescription: "The alias attached to the Team.",
				Computed:            true,
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of this Team.",
				Computed:    true,
				Optional:    true,
			},
		}),
	}
}

func (teamDataSource *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data teamDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var team *opslevel.Team
	if data.Alias.ValueString() != "" {
		team, err = teamDataSource.client.GetTeamWithAlias(data.Alias.ValueString())
	} else if opslevel.IsID(data.Id.ValueString()) {
		team, err = teamDataSource.client.GetTeam(opslevel.ID(data.Id.ValueString()))
	} else {
		resp.Diagnostics.AddError("Config Error", "'alias' or 'id' for opslevel_team datasource must be set")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to read team, got error: %s", err))
		return
	}
	if team == nil || team.Id == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to find team with alias=`%s` or id=`%s`", data.Alias.ValueString(), data.Id.ValueString()))
		return
	}

	teamDataModel := newTeamDataSourceModel(*team)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel Team data source")
	resp.Diagnostics.Append(resp.Diagnostics...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &teamDataModel)...)
}
