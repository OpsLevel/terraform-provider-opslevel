package opslevel

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &InfrastructureResource{}

var _ resource.ResourceWithImportState = &InfrastructureResource{}

func NewInfrastructureResource() resource.Resource {
	return &InfrastructureResource{}
}

// InfrastructureResource defines the resource implementation.
type InfrastructureResource struct {
	CommonResourceClient
}

type InfraProviderData struct {
	Account types.String `tfsdk:"account"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	Url     types.String `tfsdk:"url"`
}

func newInfraProviderData(infrastructure opslevel.InfrastructureResource) *InfraProviderData {
	return &InfraProviderData{
		Account: RequiredStringValue(infrastructure.ProviderData.AccountName),
		Name:    OptionalStringValue(infrastructure.ProviderData.ProviderName),
		Type:    OptionalStringValue(infrastructure.ProviderType),
		Url:     OptionalStringValue(infrastructure.ProviderData.ExternalURL),
	}
}

// InfrastructureResourceModel describes the Infrastructure managed resource.
type InfrastructureResourceModel struct {
	Aliases      types.List         `tfsdk:"aliases"`
	Data         types.String       `tfsdk:"data"`
	Id           types.String       `tfsdk:"id"`
	LastUpdated  types.String       `tfsdk:"last_updated"`
	ProviderData *InfraProviderData `tfsdk:"provider_data"`
	Owner        types.String       `tfsdk:"owner"`
	Schema       types.String       `tfsdk:"schema"`
}

func NewInfrastructureResourceModel(ctx context.Context, infrastructure opslevel.InfrastructureResource) (InfrastructureResourceModel, diag.Diagnostics) {
	var providerData *InfraProviderData
	if infrastructure.ProviderData.AccountName != "" {
		providerData = newInfraProviderData(infrastructure)
	}

	aliases, diags := OptionalStringListValue(ctx, infrastructure.Aliases)
	if diags != nil && diags.HasError() {
		return InfrastructureResourceModel{}, diags
	}

	return InfrastructureResourceModel{
		Aliases:      aliases,
		Data:         OptionalStringValue(infrastructure.Data.ToJSON()),
		Id:           ComputedStringValue(infrastructure.Id),
		ProviderData: providerData,
		Owner:        RequiredStringValue(string(infrastructure.Owner.Id())),
		Schema:       RequiredStringValue(infrastructure.Schema),
	}, diags
}

func (r *InfrastructureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infrastructure"
}

func (r *InfrastructureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Infrastructure Resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The aliases for the infrastructure resource.",
				Optional:    true,
			},
			"data": schema.StringAttribute{
				Description: "The data of the infrastructure resource in JSON format.",
				Required:    true,
				Validators: []validator.String{
					JsonStringValidator(),
					JsonHasNameKeyValidator(),
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the infrastructure.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"owner": schema.StringAttribute{
				Description: "The id of the team that owns the infrastructure resource. Does not support aliases!",
				Required:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"provider_data": schema.SingleNestedAttribute{
				Description: "The provider specific data for the infrastructure resource.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"account": schema.StringAttribute{
						Description: "The canonical account name for the provider of the infrastructure resource.",
						Required:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the provider of the infrastructure resource. (eg. AWS, GCP, Azure)",
						Optional:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type of the infrastructure resource as defined by its provider.",
						Optional:    true,
					},
					"url": schema.StringAttribute{
						Description: "The url for the provider of the infrastructure resource.",
						Optional:    true,
					},
				},
			},
			"schema": schema.StringAttribute{
				Description: "The schema of the infrastructure resource that determines its data specification.",
				Required:    true,
			},
		},
	}
}

func (r *InfrastructureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InfrastructureResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	infraInput, err := newInfraInput(data)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to create opslevel InfraInput, got error: %s", err))
		return
	}

	infrastructure, err := r.client.CreateInfrastructure(infraInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create infrastructure, got error: %s", err))
		return
	}

	givenAliases, diags := ListValueToStringSlice(ctx, data.Aliases)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given infrastructure aliases: '%s'", data.Aliases))
		return
	}
	if err = reconcileInfraAliases(*r.client, givenAliases, infrastructure); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile infrastructure aliases: '%s', got error: %s", givenAliases, err))
		return
	}

	createdInfrastructureResourceModel, diags := NewInfrastructureResourceModel(ctx, *infrastructure)
	resp.Diagnostics.Append(diags...)
	createdInfrastructureResourceModel.Aliases = data.Aliases
	createdInfrastructureResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a infrastructure resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdInfrastructureResourceModel)...)
}

func (r *InfrastructureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InfrastructureResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructure, err := r.client.GetInfrastructure(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read infrastructure, got error: %s", err))
		return
	}
	readInfrastructureResourceModel, diags := NewInfrastructureResourceModel(ctx, *infrastructure)
	resp.Diagnostics.Append(diags...)
	readInfrastructureResourceModel.Aliases = data.Aliases

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readInfrastructureResourceModel)...)
}

func (r *InfrastructureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InfrastructureResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	infraInput, err := newInfraInput(data)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create opslevel InfraInput, got error: %s", err))
		return
	}
	updatedInfrastructure, err := r.client.UpdateInfrastructure(data.Id.ValueString(), infraInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update infrastructure, got error: %s", err))
		return
	}

	givenAliases, diags := ListValueToStringSlice(ctx, data.Aliases)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given infrastructure aliases: '%s'", data.Aliases))
		return
	}
	if err = reconcileInfraAliases(*r.client, givenAliases, updatedInfrastructure); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile infrastructure aliases: '%s', got error: %s", givenAliases, err))
		return
	}

	updatedInfrastructureResourceModel, diags := NewInfrastructureResourceModel(ctx, *updatedInfrastructure)

	if data.ProviderData == nil && updatedInfrastructureResourceModel.ProviderData != nil {
		resp.Diagnostics.AddError("Known error", "Unable to unset 'provider_data' field for now. We have a planned fix for this.")
		return
	}

	resp.Diagnostics.Append(diags...)
	updatedInfrastructureResourceModel.Aliases = data.Aliases
	updatedInfrastructureResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a infrastructure resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedInfrastructureResourceModel)...)
}

func (r *InfrastructureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InfrastructureResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteInfrastructure(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete infrastructure, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a infrastructure resource")
}

func (r *InfrastructureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func newInfraInput(infraModel InfrastructureResourceModel) (opslevel.InfraInput, error) {
	infraInput := opslevel.InfraInput{Schema: infraModel.Schema.ValueString()}

	if infraModel.Owner.IsNull() {
		infraInput.Owner = opslevel.NewID("")
	} else {
		infraInput.Owner = opslevel.NewID(infraModel.Owner.ValueString())
	}

	if infraModel.Data.IsNull() {
		// Unsets this previously set field
		newJSON, err := opslevel.NewJSON("{}")
		if err != nil {
			return opslevel.InfraInput{}, err
		}
		infraInput.Data = newJSON
	} else if infraModel.Data.ValueString() != "" {
		newJSON, err := opslevel.NewJSON(infraModel.Data.ValueString())
		if err != nil {
			return opslevel.InfraInput{}, err
		}
		infraInput.Data = newJSON
	}

	if infraModel.ProviderData != nil {
		infraInput.Provider = expandInfraProviderData(*infraModel.ProviderData)
	}

	return infraInput, nil
}

func expandInfraProviderData(providerData InfraProviderData) *opslevel.InfraProviderInput {
	if providerData.Account.ValueString() == "" {
		return &opslevel.InfraProviderInput{Account: ""}
	}
	return &opslevel.InfraProviderInput{
		Account: providerData.Account.ValueString(),
		Name:    providerData.Name.ValueString(),
		Type:    providerData.Type.ValueString(),
		URL:     providerData.Url.ValueString(),
	}
}

func reconcileInfraAliases(client opslevel.Client, aliasesFromConfig []string, infra *opslevel.InfrastructureResource) error {
	// delete aliases found in infrastructure resource but not listed in Terraform config
	for _, alias := range infra.Aliases {
		if !slices.Contains(aliasesFromConfig, alias) {
			if err := client.DeleteInfraAlias(alias); err != nil {
				return err
			}
		}
	}

	// create aliases listed in Terraform config but not found in infrastructure resource
	newInfraAliases := []string{}
	for _, configAlias := range aliasesFromConfig {
		if !slices.Contains(infra.Aliases, configAlias) {
			newInfraAliases = append(newInfraAliases, configAlias)
		}
	}
	if len(newInfraAliases) > 0 {
		if _, err := client.CreateAliases(opslevel.ID(infra.Id), newInfraAliases); err != nil {
			return err
		}
	}
	infra.Aliases = aliasesFromConfig
	return nil
}
