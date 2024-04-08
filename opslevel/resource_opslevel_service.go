package opslevel

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	OwnerId                    types.String `tfsdk:"owner_id"`
	PreferredApiDocumentSource types.String `tfsdk:"preferred_api_document_source"`
	Product                    types.String `tfsdk:"product"`
	Tags                       types.List   `tfsdk:"tags"`
	TierAlias                  types.String `tfsdk:"tier_alias"`
}

func NewServiceResourceModel(ctx context.Context, service opslevel.Service) (ServiceResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	serviceResourceModel := ServiceResourceModel{
		ApiDocumentPath: OptionalStringValue(service.ApiDocumentPath),
		Description:     OptionalStringValue(service.Description),
		Framework:       OptionalStringValue(service.Framework),
		Id:              ComputedStringValue(string(service.Id)),
		Language:        OptionalStringValue(service.Language),
		LifecycleAlias:  OptionalStringValue(service.Lifecycle.Alias),
		Name:            RequiredStringValue(service.Name),
		Owner:           OptionalStringValue(service.Owner.Alias),
		OwnerId:         OptionalStringValue(string(service.Owner.Id)),
		Product:         OptionalStringValue(service.Product),
		TierAlias:       OptionalStringValue(service.Tier.Alias),
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
		apiDocSource := *&service.PreferredApiDocumentSource
		serviceResourceModel.PreferredApiDocumentSource = types.StringValue(string(*apiDocSource))
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
			},
			"api_document_path": schema.StringAttribute{
				Description: "The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path.",
				Optional:    true,
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
				Optional: true,
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
				Description: "The team that owns the service. ID or Alias my be used.",
				Optional:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The team ID that owns the service.",
				Optional:    true,
			},
			"preferred_api_document_source": schema.StringAttribute{
				Description: "The API document source (PUSH or PULL) used to determine the displayed document. If null, we use the order push and then pull.",
				Optional:    true,
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
			},
			"tier_alias": schema.StringAttribute{
				Description: "The software tier that the service belongs to.",
				Optional:    true,
			},
		},
	}
}

func (r *ServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	serviceCreateInput := opslevel.ServiceCreateInput{
		Description:    data.Description.ValueStringPointer(),
		Framework:      data.Framework.ValueStringPointer(),
		Language:       data.Language.ValueStringPointer(),
		LifecycleAlias: data.LifecycleAlias.ValueStringPointer(),
		Name:           data.Name.ValueString(),
		Product:        data.Product.ValueStringPointer(),
		TierAlias:      data.TierAlias.ValueStringPointer(),
	}

	if data.Owner.ValueString() != "" {
		serviceCreateInput.OwnerInput = opslevel.NewIdentifier(data.Owner.ValueString())
	}

	service, err := r.client.CreateService(serviceCreateInput)
	if err != nil || service == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create service, got error: %s", err))
		return
	}

	givenAliases, diags := ListValueToStringSlice(ctx, data.Aliases)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service aliases: '%s'", data.Aliases))
		return
	}
	err = reconcileServiceAliases(*r.client, givenAliases, service)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create service aliases: '%s'", givenAliases))
		return
	}

	givenTags, diags := ListValueToStringSlice(ctx, data.Tags)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service tags: '%s'", data.Tags))
		return
	}
	err = reconcileTags(*r.client, givenTags, service)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create service tags: '%s'", givenTags))
		return
	}

	createdServiceResourceModel, diags := NewServiceResourceModel(ctx, *service)
	createdServiceResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a service resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdServiceResourceModel)...)
}

func (r *ServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServiceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	service, err := r.client.GetService(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read service, got error: %s", err))
		return
	}
	readServiceResourceModel, diags := NewServiceResourceModel(ctx, *service)
	resp.Diagnostics.Append(diags...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readServiceResourceModel)...)
}

func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServiceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceUpdateInput := opslevel.ServiceUpdateInput{
		Description:    data.Description.ValueStringPointer(),
		Framework:      data.Framework.ValueStringPointer(),
		Id:             opslevel.NewID(data.Id.ValueString()),
		Language:       data.Language.ValueStringPointer(),
		LifecycleAlias: data.LifecycleAlias.ValueStringPointer(),
		Name:           data.Name.ValueStringPointer(),
		Product:        data.Product.ValueStringPointer(),
		TierAlias:      data.TierAlias.ValueStringPointer(),
	}
	service, err := r.client.UpdateService(serviceUpdateInput)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update service, got error: %s", err))
		return
	}

	givenAliases, diags := ListValueToStringSlice(ctx, data.Aliases)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service aliases: '%s'", data.Aliases))
		return
	}
	err = reconcileServiceAliases(*r.client, givenAliases, service)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create service aliases: '%s'", givenAliases))
		return
	}

	givenTags, diags := ListValueToStringSlice(ctx, data.Tags)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("Unable to handle given service tags: '%s'", data.Tags))
		return
	}
	err = reconcileTags(*r.client, givenTags, service)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create service tags: '%s'", givenTags))
		return
	}

	updatedServiceResourceModel, diags := NewServiceResourceModel(ctx, *service)
	updatedServiceResourceModel.LastUpdated = timeLastUpdated()
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "updated a service resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedServiceResourceModel)...)
}

func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServiceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteService(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a service resource")
}

func (r *ServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// import (
// 	"fmt"
// 	"slices"
// 	"strings"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
// 	"github.com/rs/zerolog/log"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceService() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a service",
// 		Create:      wrap(resourceServiceCreate),
// 		Read:        wrap(resourceServiceRead),
// 		Update:      wrap(resourceServiceUpdate),
// 		Delete:      wrap(resourceServiceDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The display name of the service.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"product": {
// 				Type:        schema.TypeString,
// 				Description: "A product is an application that your end user interacts with. Multiple services can work together to power a single product.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"description": {
// 				Type:        schema.TypeString,
// 				Description: "A brief description of the service.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"language": {
// 				Type:        schema.TypeString,
// 				Description: "The primary programming language that the service is written in.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"framework": {
// 				Type:        schema.TypeString,
// 				Description: "The primary software development framework that the service uses.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"tier_alias": {
// 				Type:        schema.TypeString,
// 				Description: "The software tier that the service belongs to.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"owner": {
// 				Type:        schema.TypeString,
// 				Description: "The team that owns the service. ID or Alias my be used.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"lifecycle_alias": {
// 				Type:        schema.TypeString,
// 				Description: "The lifecycle stage of the service.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"api_document_path": {
// 				Type:        schema.TypeString,
// 				Description: "The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"preferred_api_document_source": {
// 				Type:         schema.TypeString,
// 				Description:  "The API document source (PUSH or PULL) used to determine the displayed document. If null, we use the order push and then pull.",
// 				ForceNew:     false,
// 				Optional:     true,
// 				ValidateFunc: validation.StringInSlice(opslevel.AllApiDocumentSourceEnum, false),
// 			},
// 			"aliases": {
// 				Type:        schema.TypeList,
// 				Description: "A list of human-friendly, unique identifiers for the service.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				Elem:        &schema.Schema{Type: schema.TypeString},
// 			},
// 			"tags": {
// 				Type:        schema.TypeList,
// 				Description: "A list of tags applied to the service.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				Elem:        &schema.Schema{Type: schema.TypeString},
// 			},
// 		},
// 	}
// }

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

// func reconcileServiceAliasesOld(d *schema.ResourceData, service *opslevel.Service, client *opslevel.Client) error {
// 	expectedAliases := getStringArray(d, "aliases")
// 	existingAliases := service.ManagedAliases
// 	for _, existingAlias := range existingAliases {
// 		if !slices.Contains(expectedAliases, existingAlias) {
// 			err := client.DeleteServiceAlias(existingAlias)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	for _, expectedAlias := range expectedAliases {
// 		if !slices.Contains(existingAliases, expectedAlias) {
// 			_, err := client.CreateAliases(service.Id, []string{expectedAlias})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

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

// func reconcileTagsOld(d *schema.ResourceData, service *opslevel.Service, client *opslevel.Client) error {
// 	tags := getStringArray(d, "tags")
// 	existingTags := make([]string, 0)
// 	for _, tag := range service.Tags.Nodes {
// 		flattenedTag := flattenTag(tag)
// 		existingTags = append(existingTags, flattenedTag)
// 		if !slices.Contains(tags, flattenedTag) {
// 			err := client.DeleteTag(tag.Id)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	tagInput := map[string]string{}
// 	for _, tag := range tags {
// 		parts := strings.Split(tag, ":")
// 		if len(parts) != 2 {
// 			return fmt.Errorf("[%s] invalid tag, should be in format 'key:value' (only a single colon between the key and value, no spaces or special characters)", tag)
// 		}
// 		key := parts[0]
// 		value := parts[1]
// 		tagInput[key] = value
// 	}
// 	_, err := client.AssignTags(string(service.Id), tagInput)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceServiceCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.ServiceCreateInput{
// 		Name:           d.Get("name").(string),
// 		Product:        opslevel.RefOf(d.Get("product").(string)),
// 		Description:    opslevel.RefOf(d.Get("description").(string)),
// 		Language:       opslevel.RefOf(d.Get("language").(string)),
// 		Framework:      opslevel.RefOf(d.Get("framework").(string)),
// 		TierAlias:      opslevel.RefOf(d.Get("tier_alias").(string)),
// 		LifecycleAlias: opslevel.RefOf(d.Get("lifecycle_alias").(string)),
// 	}
// 	if owner := d.Get("owner"); owner != "" {
// 		input.OwnerInput = opslevel.NewIdentifier(owner.(string))
// 	}

// 	resource, err := client.CreateService(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	err = reconcileServiceAliases(d, resource, client)
// 	if err != nil {
// 		return err
// 	}

// 	err = reconcileTags(d, resource, client)
// 	if err != nil {
// 		return err
// 	}

// 	docPath, ok1 := d.GetOk("api_document_path")
// 	docSource, ok2 := d.GetOk("preferred_api_document_source")
// 	if ok1 || ok2 {
// 		var source *opslevel.ApiDocumentSourceEnum = nil
// 		if ok2 {
// 			s := opslevel.ApiDocumentSourceEnum(docSource.(string))
// 			source = &s
// 		}
// 		_, err := client.ServiceApiDocSettingsUpdate(string(resource.Id), docPath.(string), source)
// 		if err != nil {
// 			log.Error().Err(err).Msgf("failed to update service '%s' api doc settings", resource.Aliases[0])
// 		}
// 	}

// 	return resourceServiceRead(d, client)
// }

// func resourceServiceRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetService(opslevel.ID(id))
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("name", resource.Name); err != nil {
// 		return err
// 	}
// 	if err := d.Set("product", resource.Product); err != nil {
// 		return err
// 	}
// 	if err := d.Set("description", resource.Description); err != nil {
// 		return err
// 	}
// 	if err := d.Set("language", resource.Language); err != nil {
// 		return err
// 	}
// 	if err := d.Set("framework", resource.Framework); err != nil {
// 		return err
// 	}
// 	if err := d.Set("tier_alias", resource.Tier.Alias); err != nil {
// 		return err
// 	}

// 	// only read in changes to optional fields if they have been set before
// 	// this will prevent HasChange() from detecting changes on update
// 	if owner, ok := d.GetOk("owner"); ok || owner != "" {
// 		var ownerValue string
// 		if opslevel.IsID(owner.(string)) {
// 			ownerValue = string(resource.Owner.Id)
// 		} else {
// 			ownerValue = string(resource.Owner.Alias)
// 		}

// 		if err := d.Set("owner", ownerValue); err != nil {
// 			return err
// 		}
// 	}

// 	if err := d.Set("lifecycle_alias", resource.Lifecycle.Alias); err != nil {
// 		return err
// 	}

// 	if err := d.Set("aliases", resource.ManagedAliases); err != nil {
// 		return err
// 	}
// 	if err := d.Set("tags", flattenTagArray(resource.Tags.Nodes)); err != nil {
// 		return err
// 	}

// 	if err := d.Set("api_document_path", resource.ApiDocumentPath); err != nil {
// 		return err
// 	}
// 	if err := d.Set("preferred_api_document_source", resource.PreferredApiDocumentSource); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceServiceUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	input := opslevel.ServiceUpdateInput{
// 		Id: opslevel.NewID(id),
// 	}

// 	if d.HasChange("name") {
// 		input.Name = opslevel.RefOf(d.Get("name").(string))
// 	}
// 	if d.HasChange("product") {
// 		input.Product = opslevel.RefOf(d.Get("product").(string))
// 	}
// 	if d.HasChange("description") {
// 		input.Description = opslevel.RefOf(d.Get("description").(string))
// 	}
// 	if d.HasChange("language") {
// 		input.Language = opslevel.RefOf(d.Get("language").(string))
// 	}
// 	if d.HasChange("framework") {
// 		input.Framework = opslevel.RefOf(d.Get("framework").(string))
// 	}
// 	if d.HasChange("tier_alias") {
// 		input.TierAlias = opslevel.RefOf(d.Get("tier_alias").(string))
// 	}
// 	if d.HasChange("owner") {
// 		if owner := d.Get("owner"); owner != "" {
// 			input.OwnerInput = opslevel.NewIdentifier(owner.(string))
// 		} else {
// 			input.OwnerInput = opslevel.NewIdentifier()
// 		}
// 	}
// 	if d.HasChange("lifecycle_alias") {
// 		input.LifecycleAlias = opslevel.RefOf(d.Get("lifecycle_alias").(string))
// 	}

// 	resource, err := client.UpdateService(input)
// 	if err != nil {
// 		return err
// 	}

// 	if d.HasChange("aliases") {
// 		err = reconcileServiceAliases(d, resource, client)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if d.HasChange("tags") {
// 		tagsErr := reconcileTags(d, resource, client)
// 		if tagsErr != nil {
// 			return tagsErr
// 		}
// 	}

// 	if d.HasChange("api_document_path") || d.HasChange("preferred_api_document_source") {
// 		var docPath string
// 		var docSource *opslevel.ApiDocumentSourceEnum
// 		if value, ok := d.GetOk("api_document_path"); ok {
// 			docPath = value.(string)
// 		} else {
// 			docPath = ""
// 		}
// 		if value, ok := d.GetOk("preferred_api_document_source"); ok {
// 			s := opslevel.ApiDocumentSourceEnum(value.(string))
// 			docSource = &s
// 		} else {
// 			docSource = nil
// 		}
// 		_, err := client.ServiceApiDocSettingsUpdate(string(resource.Id), docPath, docSource)
// 		if err != nil {
// 			log.Error().Err(err).Msgf("failed to update service '%s' api doc settings", resource.Aliases[0])
// 		}
// 	}

// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceServiceRead(d, client)
// }

// func resourceServiceDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteService(id)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
