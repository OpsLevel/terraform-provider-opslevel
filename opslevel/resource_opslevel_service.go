package opslevel

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	Aliases                    types.List   `tfsdk:"aliases"`
	ApiDocumentPath            types.String `tfsdk:"api_document_path"`
	Description                types.String `tfsdk:"description"`
	Framework                  types.String `tfsdk:"framework"`
	Id                         types.String `tfsdk:"id"`
	Language                   types.String `tfsdk:"language"`
	LastUpdated                types.String `tfsdk:"last_updated"`
	LifecycleAlias             types.String `tfsdk:"lifecycle_alias"`
	Name                       types.String `tfsdk:"name"`
	Owner                      types.String `tfsdk:"owner"`
	PreferredApiDocumentSource types.String `tfsdk:"preferred_api_document_source"`
	Product                    types.String `tfsdk:"product"`
	Tags                       types.List   `tfsdk:"tags"`
	TierAlias                  types.String `tfsdk:"tier_alias"`
}

func NewServiceResourceModel(ctx context.Context, service opslevel.Service, planModel ServiceResourceModel, setLastUpdated bool) (ServiceResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	serviceResourceModel := ServiceResourceModel{
		ApiDocumentPath: OptionalStringValue(service.ApiDocumentPath),
		Description:     OptionalStringValue(service.Description),
		Framework:       OptionalStringValue(service.Framework),
		Id:              ComputedStringValue(string(service.Id)),
		Language:        OptionalStringValue(service.Language),
		LifecycleAlias:  OptionalStringValue(service.Lifecycle.Alias),
		Name:            RequiredStringValue(service.Name),
		Product:         OptionalStringValue(service.Product),
		TierAlias:       OptionalStringValue(service.Tier.Alias),
	}
	if setLastUpdated {
		serviceResourceModel.LastUpdated = timeLastUpdated()
	}

	if len(service.ManagedAliases) == 0 {
		serviceResourceModel.Aliases = types.ListNull(types.StringType)
	} else {
		serviceResourceModel.Aliases, diags = types.ListValueFrom(ctx, types.StringType, service.ManagedAliases)
		if diags.HasError() {
			return serviceResourceModel, diags
		}
	}

	if service.Tags != nil && len(service.Tags.Nodes) > 0 {
		serviceResourceModel.Tags, diags = types.ListValueFrom(ctx, types.StringType, flattenTagArray(service.Tags.Nodes))
		if diags.HasError() {
			return serviceResourceModel, diags
		}
	} else {
		serviceResourceModel.Tags = types.ListNull(types.StringType)
	}

	if service.PreferredApiDocumentSource != nil {
		apiDocSource := service.PreferredApiDocumentSource
		serviceResourceModel.PreferredApiDocumentSource = types.StringValue(string(*apiDocSource))
	}

	// after creating resource model, differentiate between a basic field being unset (null) vs empty string ("")
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

	// after creating resource model, set the owner to alias/ID or to null
	switch planModel.Owner.ValueString() {
	case string(service.Owner.Id), service.Owner.Alias:
		serviceResourceModel.Owner = planModel.Owner
	case "":
		serviceResourceModel.Owner = types.StringNull()
	default:
		diags.AddError(
			"opslevel client error",
			fmt.Sprintf("service owner found '%s' did not match given owner '%s'",
				serviceResourceModel.Owner.ValueString(),
				planModel.Owner.ValueString(),
			),
		)
		return serviceResourceModel, diags
	}

	return serviceResourceModel, diags
}

func (r *ServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *ServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Service Resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "A list of human-friendly, unique identifiers for the service.",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
				},
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
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"lifecycle_alias": schema.StringAttribute{
				Description: "The lifecycle stage of the service.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the service.",
				Required:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The team that owns the service. ID or Alias may be used.",
				Optional:    true,
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
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "A list of tags applied to the service.",
				Optional:    true,
				Validators:  []validator.List{TagFormatValidator()},
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

	// Read Terraform plan data into the model
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
		Product:        planModel.Product.ValueStringPointer(),
		TierAlias:      planModel.TierAlias.ValueStringPointer(),
	}

	if planModel.Owner.ValueString() != "" {
		serviceCreateInput.OwnerInput = opslevel.NewIdentifier(planModel.Owner.ValueString())
	}

	service, err := r.client.CreateService(serviceCreateInput)
	if err != nil || service == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create service, got error: %s", err))
		return
	}

	// TODO: the post create/update steps are the same and can be extracted into a function so we repeat less code
	givenAliases, diags := ListValueToStringSlice(ctx, planModel.Aliases)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service aliases: '%s'", planModel.Aliases))
		return
	}
	err = reconcileServiceAliases(*r.client, givenAliases, service)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service aliases: '%s', got error: %s", givenAliases, err))
		return
	}

	givenTags, diags := ListValueToStringSlice(ctx, planModel.Tags)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service tags: '%s'", planModel.Tags))
		return
	}
	err = reconcileTags(*r.client, givenTags, service)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service tags: '%s', got error: %s", givenTags, err))
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

	stateModel, diags := NewServiceResourceModel(ctx, *service, planModel, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created a service resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel ServiceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	service, err := r.client.GetService(opslevel.ID(planModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}

	stateModel, diags := NewServiceResourceModel(ctx, *service, planModel, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel ServiceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
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
	}

	service, err := r.client.UpdateService(serviceUpdateInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update service, got error: %s", err))
		return
	}

	givenAliases, diags := ListValueToStringSlice(ctx, planModel.Aliases)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service aliases: '%s'", planModel.Aliases))
		return
	}
	err = reconcileServiceAliases(*r.client, givenAliases, service)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service aliases '%s', go error: %s", givenAliases, err))
		return
	}

	givenTags, diags := ListValueToStringSlice(ctx, planModel.Tags)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service tags: '%s'", planModel.Tags))
		return
	}
	err = reconcileTags(*r.client, givenTags, service)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to reconcile service tags '%s', got error: %s", givenTags, err))
		return
	}
	if planModel.ApiDocumentPath.ValueString() != "" {
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
	service, err = r.client.GetService(service.Id)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get service after update, got error: %s", err))
		return
	}

	stateModel, diags := NewServiceResourceModel(ctx, *service, planModel, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updated a service resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel ServiceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteService(planModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a service resource")
}

func (r *ServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Assigns new aliases from terraform config to service and deletes aliases not in config
func reconcileServiceAliases(client opslevel.Client, aliasesFromConfig []string, service *opslevel.Service) error {
	// delete service aliases found in service but not listed in Terraform config
	for _, managedAlias := range service.ManagedAliases {
		if !slices.Contains(aliasesFromConfig, managedAlias) {
			if err := client.DeleteServiceAlias(managedAlias); err != nil {
				return err
			}
		}
	}

	// create aliases listed in Terraform config but not found in service
	newServiceAliases := []string{}
	for _, configAlias := range aliasesFromConfig {
		if !slices.Contains(service.ManagedAliases, configAlias) {
			newServiceAliases = append(newServiceAliases, configAlias)
		}
	}
	if len(newServiceAliases) > 0 {
		if _, err := client.CreateAliases(service.Id, newServiceAliases); err != nil {
			return err
		}
	}
	service.ManagedAliases = aliasesFromConfig
	return nil
}

// Assigns new tags from terraform config to service and deletes tags not in config
func reconcileTags(client opslevel.Client, tagsFromConfig []string, service *opslevel.Service) error {
	// delete service tags found in service but not listed in Terraform config
	existingTags := []string{}
	for _, tag := range service.Tags.Nodes {
		flattenedTag := flattenTag(tag)
		existingTags = append(existingTags, flattenedTag)
		if !slices.Contains(tagsFromConfig, flattenedTag) {
			if err := client.DeleteTag(tag.Id); err != nil {
				return err
			}
		}
	}

	// format tags listed in Terraform config but not found in service
	tagInput := map[string]string{}
	for _, tag := range tagsFromConfig {
		parts := strings.Split(tag, ":")
		if len(parts) != 2 {
			return fmt.Errorf("[%s] invalid tag, should be in format 'key:value' (only a single colon between the key and value, no spaces or special characters)", tag)
		}
		key := parts[0]
		value := parts[1]
		tagInput[key] = value
	}
	// assign tags listed in Terraform config but not found in service
	assignedTags, err := client.AssignTags(string(service.Id), tagInput)
	if err != nil {
		return err
	}

	service.Tags.Nodes = assignedTags
	return nil
}
