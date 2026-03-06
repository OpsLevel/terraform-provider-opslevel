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

var _ datasource.DataSourceWithConfigure = &TeamPropertyDefinitionDataSourcesAll{}

func NewTeamPropertyDefinitionDataSourcesAll() datasource.DataSource {
	return &TeamPropertyDefinitionDataSourcesAll{}
}

type TeamPropertyDefinitionDataSourcesAll struct {
	CommonDataSourceClient
}

type teamPropertyDefinitionDataSourcesAllModel struct {
	TeamPropertyDefinitions []teamPropertyDefinitionDataSourceItemModel `tfsdk:"team_property_definitions"`
}

type teamPropertyDefinitionDataSourceItemModel struct {
	Alias        types.String `tfsdk:"alias"`
	Description  types.String `tfsdk:"description"`
	Id           types.String `tfsdk:"id"`
	LockedStatus types.String `tfsdk:"locked_status"`
	Name         types.String `tfsdk:"name"`
	Schema       types.String `tfsdk:"schema"`
}

func NewTeamPropertyDefinitionDataSourcesAllModel(definitions []opslevel.TeamPropertyDefinition) teamPropertyDefinitionDataSourcesAllModel {
	items := make([]teamPropertyDefinitionDataSourceItemModel, len(definitions))
	for i, d := range definitions {
		items[i] = teamPropertyDefinitionDataSourceItemModel{
			Alias:        ComputedStringValue(d.Alias),
			Description:  ComputedStringValue(d.Description),
			Id:           ComputedStringValue(string(d.Id)),
			LockedStatus: ComputedStringValue(string(d.LockedStatus)),
			Name:         ComputedStringValue(d.Name),
			Schema:       ComputedStringValue(d.Schema.AsString()),
		}
	}
	return teamPropertyDefinitionDataSourcesAllModel{TeamPropertyDefinitions: items}
}

func (d *TeamPropertyDefinitionDataSourcesAll) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_property_definitions"
}

func (d *TeamPropertyDefinitionDataSourcesAll) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for all Team Property Definitions",
		Attributes: map[string]schema.Attribute{
			"team_property_definitions": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: teamPropertyDefinitionSchemaAttrs,
				},
				Description: "List of Team Property Definition data sources",
				Computed:    true,
			},
		},
	}
}

func (d *TeamPropertyDefinitionDataSourcesAll) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	definitions, err := d.client.ListTeamPropertyDefinitions(nil)
	if err != nil || definitions == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team property definitions, got error: %s", err))
		return
	}
	stateModel := NewTeamPropertyDefinitionDataSourcesAllModel(definitions.Nodes)

	tflog.Trace(ctx, "read OpsLevel TeamPropertyDefinitions data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}
