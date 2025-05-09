package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/opslevel/opslevel-go/v2025"
)

var _ datasource.DataSourceWithConfigure = (*TFDataSourceSingle[any, any])(nil)

type TFDataSourceSingle[TData any, TModel any] struct {
	datasource.DataSource
	client *opslevel.Client

	Name        string
	Description string
	Attributes  map[string]schema.Attribute

	ReadFn  func(ctx context.Context, client *opslevel.Client, identifier string) (TData, error)
	ToModel func(ctx context.Context, identifier string, data TData) (TModel, error)
}

func (s *TFDataSourceSingle[TData, TModel]) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*opslevel.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("expected *opslevel.Client, got: %T please report this issue to the provider developers at %s.", req.ProviderData, providerIssueUrl),
		)

		return
	}

	s.client = client
}

func (s *TFDataSourceSingle[TData, TModel]) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + s.Name
}

func (s *TFDataSourceSingle[TData, TModel]) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: s.Description,
		Attributes: MergeMaps(s.Attributes,
			map[string]schema.Attribute{
				"identifier": schema.StringAttribute{
					MarkdownDescription: "The identifier (id or alias) of the component type to lookup.",
					Required:            true,
				},
			},
		),
	}
}

func (s *TFDataSourceSingle[TData, TModel]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	id := ""
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("identifier"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := s.ReadFn(ctx, s.client, id)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to get resource, got error: %s", err))
		return
	}

	model, err := s.ToModel(ctx, id, item)
	if err != nil {
		resp.Diagnostics.AddError("error", fmt.Sprintf("unable to build model, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

var _ datasource.DataSourceWithConfigure = (*TFDataSourceMulti[any, any])(nil)

type TFDataSourceMulti[TData any, TModel any] struct {
	datasource.DataSource
	client *opslevel.Client

	Name        string
	Description string
	Attributes  map[string]schema.Attribute

	ReadFn  func(ctx context.Context, client *opslevel.Client) ([]TData, error)
	ToModel func(ctx context.Context, data TData) (TModel, error)
}

func (s *TFDataSourceMulti[TData, TModel]) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*opslevel.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("expected *opslevel.Client, got: %T please report this issue to the provider developers at %s.", req.ProviderData, providerIssueUrl),
		)

		return
	}

	s.client = client
}

func (s *TFDataSourceMulti[TData, TModel]) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + s.Name
}

func (s *TFDataSourceMulti[TData, TModel]) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: s.Description,
		Attributes: map[string]schema.Attribute{
			"all": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: s.Attributes,
				},
				Computed: true,
			},
		},
	}
}

func (s *TFDataSourceMulti[TData, TModel]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, err := s.ReadFn(ctx, s.client)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("unable to list resource, got error: %s", err))
		return
	}
	var items []TModel
	for _, item := range data {
		model, err := s.ToModel(ctx, item)
		if err != nil {
			resp.Diagnostics.AddError("error", fmt.Sprintf("unable to build model, got error: %s", err))
			continue
		}
		items = append(items, model)
	}
	var state struct {
		All []TModel `tfsdk:"all"`
	}
	state.All = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
