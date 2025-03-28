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

// Ensure SystemDataSource implements DataSourceWithConfigure interface
var _ datasource.DataSourceWithConfigure = &SystemDataSource{}

func NewSystemDataSource() datasource.DataSource {
	return &SystemDataSource{}
}

// SystemDataSource manages a System data source.
type SystemDataSource struct {
	CommonDataSourceClient
}

var systemDatasourceSchemaAttrs = map[string]schema.Attribute{
	"aliases": schema.ListAttribute{
		ElementType:         types.StringType,
		MarkdownDescription: "All of the aliases attached to the System.",
		Computed:            true,
	},
	"description": schema.StringAttribute{
		MarkdownDescription: "The description of the System.",
		Computed:            true,
	},
	"domain": schema.StringAttribute{
		MarkdownDescription: "ID of the parent domain of the System.",
		Computed:            true,
	},
	"id": schema.StringAttribute{
		Description: "The ID of this System.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the System.",
		Computed:    true,
	},
	"owner": schema.StringAttribute{
		Description: "The id of the team that owns the System.",
		Computed:    true,
	},
}

func systemAttributes(attrs map[string]schema.Attribute) map[string]schema.Attribute {
	for key, value := range systemDatasourceSchemaAttrs {
		attrs[key] = value
	}
	return attrs
}

type systemDataSourceModel struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Domain      types.String `tfsdk:"domain"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Owner       types.String `tfsdk:"owner"`
}

type systemDataSourceModelWithIdentifier struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Domain      types.String `tfsdk:"domain"`
	Id          types.String `tfsdk:"id"`
	Identifier  types.String `tfsdk:"identifier"`
	Name        types.String `tfsdk:"name"`
	Owner       types.String `tfsdk:"owner"`
}

func newSystemDataSourceModelWithIdentifier(system opslevel.System, identifier types.String) systemDataSourceModelWithIdentifier {
	aliases := OptionalStringListValue(system.Aliases)
	return systemDataSourceModelWithIdentifier{
		Aliases:     aliases,
		Description: ComputedStringValue(system.Description),
		Domain:      ComputedStringValue(string(system.Parent.Id)),
		Id:          ComputedStringValue(string(system.Id)),
		Identifier:  identifier,
		Name:        ComputedStringValue(system.Name),
		Owner:       ComputedStringValue(string(system.Owner.Id())),
	}
}

func newSystemDataSourceModel(system opslevel.System) systemDataSourceModel {
	aliases := OptionalStringListValue(system.Aliases)
	return systemDataSourceModel{
		Aliases:     aliases,
		Description: ComputedStringValue(system.Description),
		Domain:      ComputedStringValue(string(system.Parent.Id)),
		Id:          ComputedStringValue(string(system.Id)),
		Name:        ComputedStringValue(system.Name),
		Owner:       ComputedStringValue(string(system.Owner.Id())),
	}
}

func (sys *SystemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (sys *SystemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "System data source",

		Attributes: systemAttributes(map[string]schema.Attribute{
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the System.",
				Required:    true,
			},
		}),
	}
}

func (sys *SystemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := read[systemDataSourceModelWithIdentifier](ctx, &resp.Diagnostics, req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := sys.client.GetSystem(data.Identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("unable to read system, got error: %s", err))
		return
	}
	systemDataModel := newSystemDataSourceModelWithIdentifier(*system, data.Identifier)

	// Save data into Terraform state
	tflog.Trace(ctx, "read an OpsLevel System data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &systemDataModel)...)
}
