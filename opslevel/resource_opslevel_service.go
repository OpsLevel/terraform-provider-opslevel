package opslevel

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
)

var _ resource.ResourceWithConfigure = &ServiceResource{}

var _ resource.ResourceWithImportState = &ServiceResource{}

func NewServiceResource() resource.Resource {
	return &ServiceResource{}
}

// ServiceResource defines the resource implementation.
type ServiceResource struct {
	CommonResourceClient
}

// ServiceResourceModel describes the Service managed resource.
type ServiceResourceModel struct {
	Aliases                    types.Set    `tfsdk:"aliases"`
	ApiDocumentPath            types.String `tfsdk:"api_document_path"`
	Description                types.String `tfsdk:"description"`
	Framework                  types.String `tfsdk:"framework"`
	Id                         types.String `tfsdk:"id"`
	Language                   types.String `tfsdk:"language"`
	LifecycleAlias             types.String `tfsdk:"lifecycle_alias"`
	Name                       types.String `tfsdk:"name"`
	Note                       types.String `tfsdk:"note"`
	Owner                      types.String `tfsdk:"owner"`
	Parent                     types.String `tfsdk:"parent"`
	PreferredApiDocumentSource types.String `tfsdk:"preferred_api_document_source"`
	Product                    types.String `tfsdk:"product"`
	Tags                       types.Set    `tfsdk:"tags"`
	TierAlias                  types.String `tfsdk:"tier_alias"`
	Type                       types.String `tfsdk:"type"`
}

func newServiceResourceModel(ctx context.Context, service opslevel.Service, givenModel ServiceResourceModel) (ServiceResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	serviceResourceModel := ServiceResourceModel{
		ApiDocumentPath: OptionalStringValue(service.ApiDocumentPath),
		Description:     OptionalStringValue(service.Description),
		Framework:       OptionalStringValue(service.Framework),
		Id:              ComputedStringValue(string(service.Id)),
		Language:        OptionalStringValue(service.Language),
		LifecycleAlias:  OptionalStringValue(service.Lifecycle.Alias),
		Name:            RequiredStringValue(service.Name),
		Note:            OptionalStringValue(service.Note),
		Product:         OptionalStringValue(service.Product),
		TierAlias:       OptionalStringValue(service.Tier.Alias),
		Type:            OptionalStringValue(givenModel.Type.ValueString()),
	}

	if string(service.Owner.Id) == "" {
		serviceResourceModel.Owner = types.StringNull()
	} else if string(service.Owner.Id) == givenModel.Owner.ValueString() || service.Owner.Alias == givenModel.Owner.ValueString() {
		serviceResourceModel.Owner = givenModel.Owner
	} else {
		serviceResourceModel.Owner = OptionalStringValue(string(service.Owner.Id))
	}

	if service.Parent == nil {
		serviceResourceModel.Parent = types.StringNull()
	} else if string(service.Parent.Id) == givenModel.Parent.ValueString() || slices.Contains(service.Parent.Aliases, givenModel.Parent.ValueString()) {
		serviceResourceModel.Parent = givenModel.Parent
	} else {
		serviceResourceModel.Parent = OptionalStringValue(string(service.Parent.Id))
	}

	if givenModel.Aliases.IsNull() {
		serviceResourceModel.Aliases = types.SetNull(types.StringType)
	} else {
		serviceResourceModel.Aliases = givenModel.Aliases
	}

	if givenModel.Tags.IsNull() && (service.Tags != nil || len(service.Tags.Nodes) == 0) {
		serviceResourceModel.Tags = types.SetNull(types.StringType)
	} else {
		serviceResourceModel.Tags, diags = types.SetValueFrom(ctx, types.StringType, flattenTagArray(service.Tags.Nodes))
		if diags.HasError() {
			return serviceResourceModel, diags
		}
	}

	if service.PreferredApiDocumentSource != nil {
		apiDocSource := service.PreferredApiDocumentSource
		serviceResourceModel.PreferredApiDocumentSource = types.StringValue(string(*apiDocSource))
	}

	return serviceResourceModel, diags
}

// updateServiceResourceModelWithPlan mutates the input ServiceResourceModel based on the current terraform plan
func updateServiceResourceModelWithPlan(serviceResourceModel *ServiceResourceModel, planModel ServiceResourceModel) {
	// handle edge case of a basic field being null vs ""
	if serviceResourceModel.Description.IsNull() && !planModel.Description.IsNull() && planModel.Description.ValueString() == "" {
		serviceResourceModel.Description = types.StringValue("")
	}
	if serviceResourceModel.Framework.IsNull() && !planModel.Framework.IsNull() && planModel.Framework.ValueString() == "" {
		serviceResourceModel.Framework = types.StringValue("")
	}
	if serviceResourceModel.Language.IsNull() && !planModel.Language.IsNull() && planModel.Language.ValueString() == "" {
		serviceResourceModel.Language = types.StringValue("")
	}
	if serviceResourceModel.Product.IsNull() && !planModel.Product.IsNull() && planModel.Product.ValueString() == "" {
		serviceResourceModel.Product = types.StringValue("")
	}
}

func (r *ServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *ServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Service Resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "A list of human-friendly, unique identifiers for the service.",
				Optional:    true,
			},
			"api_document_path": schema.StringAttribute{
				Description: "The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.+\.(json|ya?ml)$`),
						"ends with '.json', '.yml', or '.yaml'",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: "A brief description of the service.",
				Optional:    true,
			},
			"framework": schema.StringAttribute{
				Description: "The primary software development framework that the service uses.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "The id of the service to find",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"language": schema.StringAttribute{
				Description: "The primary programming language that the service is written in.",
				Optional:    true,
			},
			"lifecycle_alias": schema.StringAttribute{
				Description: "The lifecycle stage of the service.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the service.",
				Required:    true,
			},
			"note": schema.StringAttribute{
				Description: "Additional information about the service.",
				Optional:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The team that owns the service. ID or Alias may be used.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
			"parent": schema.StringAttribute{
				Description: "The id or alias of the parent system of this service",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
			"preferred_api_document_source": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The API document source (%s) used to determine the displayed document. If null, defaults to PUSH.",
					strings.Join(opslevel.AllApiDocumentSourceEnum, " or "),
				),
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllApiDocumentSourceEnum...),
				},
			},
			"product": schema.StringAttribute{
				Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
				Optional:    true,
			},
			"tags": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "A list of tags applied to the service.",
				Optional:    true,
				Validators:  []validator.Set{TagFormatValidator()},
			},
			"tier_alias": schema.StringAttribute{
				Description: "The software tier that the service belongs to.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "The component type of the service.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.NoneOf(""),
				},
			},
		},
	}
}

func (r *ServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[ServiceResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.ServiceCreateInput{
		Description:    nullable(planModel.Description.ValueStringPointer()),
		Framework:      nullable(planModel.Framework.ValueStringPointer()),
		Language:       nullable(planModel.Language.ValueStringPointer()),
		LifecycleAlias: nullable(planModel.LifecycleAlias.ValueStringPointer()),
		Name:           planModel.Name.ValueString(),
		OwnerInput:     opslevel.NewIdentifier(),
		Parent:         opslevel.NewIdentifier(),
		Product:        nullable(planModel.Product.ValueStringPointer()),
		TierAlias:      nullable(planModel.TierAlias.ValueStringPointer()),
	}

	if planModel.Owner.ValueString() != "" {
		input.OwnerInput = opslevel.NewIdentifier(planModel.Owner.ValueString())
	}

	if planModel.Parent.ValueString() != "" {
		input.Parent = opslevel.NewIdentifier(planModel.Parent.ValueString())
	}

	if !planModel.Type.IsNull() {
		input.Type = opslevel.NewIdentifier(planModel.Type.ValueString())
	}

	service, err := r.client.CreateService(input)
	if err != nil || service == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create service, got error: %s", err))
		return
	}

	// TODO: the post create/update steps are the same and can be extracted into a function so we repeat less code
	if len(planModel.Aliases.Elements()) > 0 {
		aliases, diags := SetValueToStringSlice(ctx, planModel.Aliases)
		if diags != nil && diags.HasError() {
			resp.Diagnostics.Append(diags...)
			resp.Diagnostics.AddAttributeError(path.Root("aliases"), "Config error", "unable to handle given service aliases")
			return
		}
		// add "unique identifiers" (OpsLevel created aliases) before reconciling.
		// this ensures that we don't try to create an alias that already exists
		aliases = append(aliases, service.UniqueIdentifiers()...)
		if err = service.ReconcileAliases(r.client, aliases); err != nil {
			resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service aliases: '%s'\n%s", aliases, err))

			// delete newly created team to avoid dupliate team creation on next 'terraform apply'
			if err := r.client.DeleteService(string(service.Id)); err != nil {
				resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("failed to delete incorrectly created service '%s' following aliases error:\n%s", service.Name, err))
			}
			return
		}
	}

	if _, err := updateServiceNote(*r.client, *service, planModel); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update service note, got error: %s", err))
		return
	}

	givenTags, diags := TagSetValueToTagSlice(ctx, planModel.Tags)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service tags: '%s'", planModel.Tags))
		return
	}
	if err = r.client.ReconcileTags(service, givenTags); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service tags '%s', got error: %s", givenTags, err))
		return
	}

	if planModel.ApiDocumentPath.ValueString() != "" {
		apiDocPath := planModel.ApiDocumentPath.ValueString()
		if planModel.PreferredApiDocumentSource.IsNull() {
			if _, err := r.client.ServiceApiDocSettingsUpdate(string(service.Id), apiDocPath, nil); err != nil {
				resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to set provided 'api_document_path' %s for service. error: %s", apiDocPath, err))
				return
			}
		} else {
			sourceEnum := opslevel.ApiDocumentSourceEnum(planModel.PreferredApiDocumentSource.ValueString())
			if _, err := r.client.ServiceApiDocSettingsUpdate(string(service.Id), apiDocPath, &sourceEnum); err != nil {
				resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to set provided 'api_document_path' %s with doc source '%s' for service. error: %s", apiDocPath, sourceEnum, err))
				return
			}
		}
	}

	// fetch the service again, since other mutations are performed after the create/update step
	service, err = r.client.GetService(string(service.Id))
	if err != nil {
		if (service == nil || service.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get service after creation, got error: %s", err))
		return
	}

	newStateModel, diags := newServiceResourceModel(ctx, *service, planModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateServiceResourceModelWithPlan(&newStateModel, planModel)

	tflog.Trace(ctx, "created a service resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func (r *ServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[ServiceResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	service, err := r.client.GetService(stateModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	if service.Id == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	newStateModel, diags := newServiceResourceModel(ctx, *service, stateModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func unsetStringHelper(plan, state basetypes.StringValue) *opslevel.Nullable[string] {
	if !plan.IsNull() {
		return nullable(plan.ValueStringPointer())
	} else if !state.IsNull() { // Unset
		return &opslevel.Nullable[string]{SetNull: true}
	}
	return nil
}

func unsetIDHelper(plan, state basetypes.StringValue) *opslevel.Nullable[opslevel.ID] {
	if !plan.IsNull() {
		return nullable(opslevel.NewID(plan.ValueString()))
	} else if !state.IsNull() { // Unset
		return &opslevel.Nullable[opslevel.ID]{SetNull: true}
	}
	return nil
}

func unsetIdentifierHelper(plan, state basetypes.StringValue) *opslevel.IdentifierInput {
	if !plan.IsNull() {
		return opslevel.NewIdentifier(plan.ValueString())
	} else if !state.IsNull() { // Unset
		return opslevel.NewIdentifier()
	}
	return nil
}

func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[ServiceResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[ServiceResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceUpdateInput := opslevel.ServiceUpdateInput{
		Description:    unsetStringHelper(planModel.Description, stateModel.Description),
		Framework:      unsetStringHelper(planModel.Framework, stateModel.Framework),
		Id:             nullable(opslevel.NewID(planModel.Id.ValueString())),
		Language:       unsetStringHelper(planModel.Language, stateModel.Language),
		LifecycleAlias: unsetStringHelper(planModel.LifecycleAlias, stateModel.LifecycleAlias),
		Name:           unsetStringHelper(planModel.Name, stateModel.Name),
		Product:        unsetStringHelper(planModel.Product, stateModel.Product),
		TierAlias:      unsetStringHelper(planModel.TierAlias, stateModel.TierAlias),

		OwnerInput: unsetIdentifierHelper(planModel.Owner, stateModel.Owner),
		Parent:     unsetIdentifierHelper(planModel.Parent, stateModel.Parent),
		Type:       unsetIdentifierHelper(planModel.Type, stateModel.Type),
	}

	service, err := r.client.UpdateService(serviceUpdateInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update service, got error: %s", err))
		return
	}

	aliases, diags := SetValueToStringSlice(ctx, planModel.Aliases)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddAttributeError(path.Root("aliases"), "Config error", "unable to handle given service aliases")
		return
	}

	// Try deleting uniqueIdentifiers (aka default alias) if not declared in Terraform config
	// Deleting this alias may fail because its locked but that's ok
	uniqueIdentifiers := service.UniqueIdentifiers()
	for _, uniqueIdentifier := range uniqueIdentifiers {
		if !slices.Contains(aliases, uniqueIdentifier) {
			_ = r.client.DeleteAlias(opslevel.AliasDeleteInput{
				Alias:     uniqueIdentifier,
				OwnerType: opslevel.AliasOwnerTypeEnumService,
			})
		}
	}
	// add "unique identifiers" (OpsLevel created aliases) before reconciling.
	// this ensures that we don't try to create an alias that already exists
	aliases = append(aliases, uniqueIdentifiers...)
	if err = service.ReconcileAliases(r.client, aliases); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service aliases: '%s'\n%s", aliases, err))
	}

	// update service note only if known to plan and/or state
	if !planModel.Note.IsNull() || !stateModel.Note.IsNull() {
		if _, err := updateServiceNote(*r.client, *service, planModel); err != nil {
			resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update service note, got error: %s", err))
			return
		}
	}

	givenTags, diags := TagSetValueToTagSlice(ctx, planModel.Tags)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddAttributeError(path.Root("tags"), "Config error", "unable to handle given service tags")
		return
	}

	if !stateModel.Tags.IsNull() || !planModel.Tags.IsNull() {
		if err = r.client.ReconcileTags(service, givenTags); err != nil {
			resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service tags '%s', got error: %s", givenTags, err))
			return
		}
	}

	if planModel.ApiDocumentPath.IsNull() {
		if _, err := r.client.ServiceApiDocSettingsUpdate(string(service.Id), "", nil); err != nil {
			resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to unset 'api_document_path' for service %s. error: %s", service.Name, err))
			return
		}
	} else {
		apiDocPath := planModel.ApiDocumentPath.ValueString()
		if planModel.PreferredApiDocumentSource.IsNull() {
			if _, err := r.client.ServiceApiDocSettingsUpdate(string(service.Id), apiDocPath, nil); err != nil {
				resp.Diagnostics.AddError("opslevel client error",
					fmt.Sprintf(
						"Unable to set provided 'api_document_path' %s for service. error: %s",
						apiDocPath, err),
				)
				return
			}
		} else {
			sourceEnum := opslevel.ApiDocumentSourceEnum(planModel.PreferredApiDocumentSource.ValueString())
			if _, err := r.client.ServiceApiDocSettingsUpdate(string(service.Id), apiDocPath, &sourceEnum); err != nil {
				resp.Diagnostics.AddError("opslevel client error",
					fmt.Sprintf(
						"Unable to set provided 'api_document_path' %s with doc source '%s' for service. error: %s",
						apiDocPath, sourceEnum, err),
				)
				return
			}
		}
	}

	// fetch the service again, since other mutations are performed after the create/update step
	service, err = r.client.GetService(string(service.Id))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get service after update, got error: %s", err))
		return
	}

	newStateModel, diags := newServiceResourceModel(ctx, *service, planModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateServiceResourceModelWithPlan(&newStateModel, planModel)

	tflog.Trace(ctx, "updated a service resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateModel)...)
}

func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[ServiceResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteService(stateModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a service resource")
}

func (r *ServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func updateServiceNote(client opslevel.Client, service opslevel.Service, planModel ServiceResourceModel) (*opslevel.Service, error) {
	if planModel.Note.ValueString() == service.Note {
		return nil, nil
	}

	serviceNoteUpdateInput := opslevel.ServiceNoteUpdateInput{
		Service: *opslevel.NewIdentifier(string(service.Id)),
	}

	if service.Note != "" && planModel.Note.IsNull() {
		// unset the previously set note
		serviceNoteUpdateInput.Note = opslevel.RefOf("")
	} else {
		serviceNoteUpdateInput.Note = nullable(planModel.Note.ValueStringPointer())
	}

	return client.UpdateServiceNote(serviceNoteUpdateInput)
}
