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

var _ datasource.DataSourceWithConfigure = &TeamPropertyDefinitionDataSource{}

func NewTeamPropertyDefinitionDataSource() datasource.DataSource {
	return &TeamPropertyDefinitionDataSource{}
}

type TeamPropertyDefinitionDataSource struct {
	CommonDataSourceClient
}

type teamPropertyDefinitionDataSourceModel struct {
	Alias        types.String `tfsdk:"alias"`
	Description  types.String `tfsdk:"description"`
	Id           types.String `tfsdk:"id"`
	Identifier   types.String `tfsdk:"identifier"`
	LockedStatus types.String `tfsdk:"locked_status"`
	Name         types.String `tfsdk:"name"`
	Schema       types.String `tfsdk:"schema"`
}

func NewTeamPropertyDefinitionDataSourceModel(definition opslevel.TeamPropertyDefinition, identifier string) teamPropertyDefinitionDataSourceModel {
	return teamPropertyDefinitionDataSourceModel{
		Alias:        ComputedStringValue(definition.Alias),
		Description:  ComputedStringValue(definition.Description),
		Id:           ComputedStringValue(string(definition.Id)),
		Identifier:   ComputedStringValue(identifier),
		LockedStatus: ComputedStringValue(string(definition.LockedStatus)),
		Name:         ComputedStringValue(definition.Name),
		Schema:       ComputedStringValue(definition.Schema.AsString()),
	}
}

var teamPropertyDefinitionSchemaAttrs = map[string]schema.Attribute{
	"alias": schema.StringAttribute{
		MarkdownDescription: "The human-friendly, unique identifier of the team property definition.",
		Computed:            true,
	},
	"description": schema.StringAttribute{
		MarkdownDescription: "The description of the team property definition.",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		MarkdownDescription: "The ID of this resource.",
		Computed:            true,
	},
	"locked_status": schema.StringAttribute{
		MarkdownDescription: "Restricts what sources are able to assign values to this property.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		MarkdownDescription: "The display name of the team property definition.",
		Computed:            true,
	},
	"schema": schema.StringAttribute{
		MarkdownDescription: "The schema of the team property definition.",
		Computed:            true,
	},
}

func teamPropertyDefinitionDataSourceAttrs(extra map[string]schema.Attribute) map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	for k, v := range teamPropertyDefinitionSchemaAttrs {
		attrs[k] = v
	}
	for k, v := range extra {
		attrs[k] = v
	}
	return attrs
}

func (d *TeamPropertyDefinitionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_property_definition"
}

func (d *TeamPropertyDefinitionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for a Team Property Definition",
		Attributes: teamPropertyDefinitionDataSourceAttrs(map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The id or alias of the team property definition.",
				Required:            true,
			},
		}),
	}
}

func (d *TeamPropertyDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	planModel := read[teamPropertyDefinitionDataSourceModel](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	definition, err := d.client.GetTeamPropertyDefinition(planModel.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team property definition, got error: %s", err))
		return
	}
	stateModel := NewTeamPropertyDefinitionDataSourceModel(*definition, planModel.Identifier.ValueString())

	tflog.Trace(ctx, "read an OpsLevel TeamPropertyDefinition data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
