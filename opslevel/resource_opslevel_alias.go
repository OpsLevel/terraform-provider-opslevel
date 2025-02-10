package opslevel

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opslevel/opslevel-go/v2025"
)

var _ resource.ResourceWithConfigure = &AliasResource{}

func NewAliasResource() resource.Resource {
	return &AliasResource{}
}

// AliasResource defines the resource implementation for managing aliases.
type AliasResource struct {
	CommonResourceClient
}

func (r *AliasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alias"
}

type AliasResourceModel struct {
	ResourceType       types.String `tfsdk:"resource_type"`
	ResourceIdentifier types.String `tfsdk:"resource_identifier"`
	Aliases            types.Set    `tfsdk:"aliases"`

	Id types.String `tfsdk:"id"`
}

func (s AliasResourceModel) GetResource(d *diag.Diagnostics, client *opslevel.Client) opslevel.AliasableResourceInterface {
	resourceType := opslevel.AliasOwnerTypeEnum(s.ResourceType.ValueString())
	resourceIdentifier := s.ResourceIdentifier.ValueString()
	output, err := client.GetAliasableResource(resourceType, resourceIdentifier)
	if err != nil {
		d.AddError(
			"opslevel client error",
			fmt.Sprintf("Failed to find aliasable resource, %s", err),
		)
	}
	return output
}

func (s AliasResourceModel) GetAliases(ctx context.Context, d *diag.Diagnostics) []string {
	var output []string
	if !s.Aliases.IsNull() {
		d.Append(s.Aliases.ElementsAs(ctx, &output, true)...)
	}
	return output
}

func (r *AliasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Alias Resource",

		Attributes: map[string]schema.Attribute{
			"resource_identifier": schema.StringAttribute{
				Description: "The id or human-friendly, unique identifier of the resource this alias belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The resource type that the alias applies to. One of `%s`",
					strings.Join(opslevel.AllAliasOwnerTypeEnum, "`, `"),
				),
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllAliasOwnerTypeEnum...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"aliases": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "The unique set of aliases to ensure exist on the resource.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The id of the resource, maybe be duplicative of the 'resource_identifier' but in the case where that is an alias itself this is the identifier of what it found during lookup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AliasResource) createAlias(d *diag.Diagnostics, alias string, aliasable opslevel.AliasableResourceInterface) {
	input := opslevel.AliasCreateInput{
		Alias:   alias,
		OwnerId: aliasable.ResourceId(),
	}

	if _, err := r.client.CreateAlias(input); err != nil {
		d.AddError(
			"opslevel client error",
			fmt.Sprintf("Failed to create alias '%s', %s", alias, err),
		)
	}
}

func (r *AliasResource) deleteAlias(d *diag.Diagnostics, alias string, aliasable opslevel.AliasableResourceInterface) {
	input := opslevel.AliasDeleteInput{
		Alias:     alias,
		OwnerType: aliasable.AliasableType(),
	}
	if err := r.client.DeleteAlias(input); err != nil {
		// This allows locked slugs to be added and not cause a failure upon delete
		if strings.Contains(err.Error(), "slug is locked, it cannot be deleted") {
			return
		}
		d.AddError(
			"opslevel client error",
			fmt.Sprintf("Failed to delete alias '%s', %s", alias, err),
		)
	}
}

func (r *AliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[AliasResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	desiredAliases := planModel.GetAliases(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	aliasable := planModel.GetResource(&resp.Diagnostics, r.client)
	if resp.Diagnostics.HasError() {
		return
	}

	planModel.Id = types.StringValue(string(aliasable.ResourceId()))

	currentAliases := aliasable.GetAliases()
	for _, alias := range desiredAliases {
		if slices.Contains(currentAliases, alias) {
			continue
		}
		r.createAlias(&resp.Diagnostics, alias, aliasable)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planModel)...)
}

func (r *AliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	planModel := read[AliasResourceModel](ctx, &resp.Diagnostics, req.State)

	aliasable := planModel.GetResource(&resp.Diagnostics, r.client)
	if resp.Diagnostics.HasError() {
		return
	}

	planModel.Id = types.StringValue(string(aliasable.ResourceId()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &planModel)...)
}

func (r *AliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[AliasResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[AliasResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	desiredAliases := planModel.GetAliases(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	managedAliases := stateModel.GetAliases(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	aliasable := planModel.GetResource(&resp.Diagnostics, r.client)
	if resp.Diagnostics.HasError() {
		return
	}
	currentAliases := aliasable.GetAliases()

	for _, alias := range managedAliases {
		if slices.Contains(desiredAliases, alias) {
			continue
		}
		if slices.Contains(currentAliases, alias) {
			r.deleteAlias(&resp.Diagnostics, alias, aliasable)
		}
	}

	for _, alias := range desiredAliases {
		if slices.Contains(currentAliases, alias) {
			continue
		}
		r.createAlias(&resp.Diagnostics, alias, aliasable)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planModel)...)
}

func (r *AliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[AliasResourceModel](ctx, &resp.Diagnostics, req.State)

	managedAliases := stateModel.GetAliases(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	aliasable := stateModel.GetResource(&resp.Diagnostics, r.client)
	if resp.Diagnostics.HasError() {
		return
	}
	currentAliases := aliasable.GetAliases()

	for _, alias := range managedAliases {
		if slices.Contains(currentAliases, alias) {
			r.deleteAlias(&resp.Diagnostics, alias, aliasable)
		}
	}
}
