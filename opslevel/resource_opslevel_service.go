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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
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
		Note:            givenModel.Note,
		Owner:           OptionalStringValue(givenModel.Owner.ValueString()),
		Parent:          OptionalStringValue(givenModel.Parent.ValueString()),
		Product:         OptionalStringValue(service.Product),
		TierAlias:       OptionalStringValue(service.Tier.Alias),
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
		},
	}
}

func (r *ServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel ServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceCreateInput := opslevel.ServiceCreateInput{
		Description:    planModel.Description.ValueStringPointer(),
		Framework:      planModel.Framework.ValueStringPointer(),
		Language:       planModel.Language.ValueStringPointer(),
		LifecycleAlias: planModel.LifecycleAlias.ValueStringPointer(),
		Name:           planModel.Name.ValueString(),
		OwnerInput:     opslevel.NewIdentifier(),
		Parent:         opslevel.NewIdentifier(),
		Product:        planModel.Product.ValueStringPointer(),
		TierAlias:      planModel.TierAlias.ValueStringPointer(),
	}

	if planModel.Owner.ValueString() != "" {
		serviceCreateInput.OwnerInput = opslevel.NewIdentifier(planModel.Owner.ValueString())
	}

	if planModel.Parent.ValueString() != "" {
		serviceCreateInput.Parent = opslevel.NewIdentifier(planModel.Parent.ValueString())
	}

	service, err := r.client.CreateService(serviceCreateInput)
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
	service, err = r.client.GetService(service.Id)
	if err != nil {
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
	var stateModel ServiceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	service, err := r.client.GetService(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
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

func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel, stateModel ServiceResourceModel
	var ownerFromState, parentFromState types.String

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("owner"), &ownerFromState)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("parent"), &parentFromState)...)
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceUpdateInput := opslevel.ServiceUpdateInputV2{
		Description:    NullableStringConfigValue(planModel.Description),
		Framework:      NullableStringConfigValue(planModel.Framework),
		Id:             opslevel.NewID(planModel.Id.ValueString()),
		Language:       NullableStringConfigValue(planModel.Language),
		LifecycleAlias: NullableStringConfigValue(planModel.LifecycleAlias),
		Name:           opslevel.NewNullableFrom(planModel.Name.ValueString()),
		Product:        NullableStringConfigValue(planModel.Product),
		TierAlias:      NullableStringConfigValue(planModel.TierAlias),
	}
	if planModel.Owner.ValueString() != "" {
		serviceUpdateInput.OwnerInput = opslevel.NewIdentifier(planModel.Owner.ValueString())
	} else if !ownerFromState.IsNull() && planModel.Owner.ValueString() == "" {
		// unset owner field only if it's set in state and owner field in plan is not set
		serviceUpdateInput.OwnerInput = opslevel.NewIdentifier()
	}

	if planModel.Parent.ValueString() != "" {
		serviceUpdateInput.Parent = opslevel.NewIdentifier(planModel.Parent.ValueString())
	} else if !parentFromState.IsNull() && planModel.Parent.ValueString() == "" {
		// unset parent field only if it's set in state and parent field in plan is not set
		serviceUpdateInput.Parent = opslevel.NewIdentifier()
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
	if err = r.client.ReconcileTags(service, givenTags); err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service tags '%s', got error: %s", givenTags, err))
		return
	}
	apiDocPath := ""
	if !planModel.ApiDocumentPath.IsNull() {
		apiDocPath = planModel.ApiDocumentPath.ValueString()
	}
	var apiDocSource *opslevel.ApiDocumentSourceEnum
	if !planModel.PreferredApiDocumentSource.IsNull() {
		apiDocSource = opslevel.RefOf(opslevel.ApiDocumentSourceEnum(planModel.PreferredApiDocumentSource.ValueString()))
	}
	service, err = r.client.ServiceApiDocSettingsUpdate(string(service.Id), apiDocPath, apiDocSource)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update api document settings for service %s. error: %s", service.Name, err))
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
	var stateModel ServiceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
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
		Service: opslevel.IdentifierInput{Id: opslevel.RefOf(service.Id)},
	}

	if service.Note != "" && planModel.Note.IsNull() {
		// unset the previously set note
		serviceNoteUpdateInput.Note = opslevel.RefOf("")
	} else {
		serviceNoteUpdateInput.Note = planModel.Note.ValueStringPointer()
	}

	return client.UpdateServiceNote(serviceNoteUpdateInput)
}
