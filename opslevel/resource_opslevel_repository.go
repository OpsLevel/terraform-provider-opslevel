package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2025"
)

var _ resource.ResourceWithConfigure = &RepositoryResource{}

var _ resource.ResourceWithImportState = &RepositoryResource{}

func NewRepositoryResource() resource.Resource {
	return &RepositoryResource{}
}

// RepositoryResource defines the resource implementation.
type RepositoryResource struct {
	CommonResourceClient
}

// RepositoryResourceModel describes the Repository managed resource.
type RepositoryResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Owner      types.String `tfsdk:"owner"`
}

func NewRepositoryResourceModel(ctx context.Context, repository opslevel.Repository) RepositoryResourceModel {
	return RepositoryResourceModel{
		Id:    ComputedStringValue(string(repository.Id)),
		Owner: OptionalStringValue(string(repository.Owner.Id)),
	}
}

func (r *RepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (r *RepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Repository Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identifier": schema.StringAttribute{
				Description: "The id or human-friendly, unique identifier for the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "The ID of the owner of the repository.",
				Optional:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
		},
	}
}

func (r *RepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[RepositoryResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	var repository *opslevel.Repository

	identifier := planModel.Identifier.ValueString()
	if opslevel.IsID(identifier) {
		repository, err = r.client.GetRepository(*opslevel.NewID(identifier))
	} else {
		repository, err = r.client.GetRepositoryWithAlias(identifier)
	}
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to get repository, got error: %s", err))
		return
	}
	input := opslevel.RepositoryUpdateInput{
		Id: repository.Id,
	}
	if !planModel.Owner.IsNull() {
		input.OwnerId = nullable(opslevel.NewID(planModel.Owner.ValueString()))
	}

	updatedRepository, err := r.client.UpdateRepository(input)
	if err != nil || updatedRepository == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update repository, got error: %s", err))
		return
	}
	stateModel := NewRepositoryResourceModel(ctx, *updatedRepository)

	// Identifier from plan can be an id or alias
	switch planModel.Identifier.ValueString() {
	case string(updatedRepository.Id), updatedRepository.DefaultAlias:
		stateModel.Identifier = planModel.Identifier
	default:
		resp.Diagnostics.AddError(
			"opslevel client error",
			fmt.Sprintf("given repository identifier '%s' did not match found repository's id '%s' or alias '%s'",
				planModel.Identifier.ValueString(),
				string(updatedRepository.Id),
				updatedRepository.DefaultAlias,
			),
		)
		return
	}

	tflog.Trace(ctx, "created a repository resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *RepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[RepositoryResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	readRepository, err := r.client.GetRepository(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil || readRepository == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read repository, got error: %s", err))
		return
	}
	verifiedStateModel := NewRepositoryResourceModel(ctx, *readRepository)

	// Identifier from plan can be an id or alias
	switch stateModel.Identifier.ValueString() {
	case string(readRepository.Id), readRepository.DefaultAlias:
		verifiedStateModel.Identifier = stateModel.Identifier
	default:
		resp.Diagnostics.AddError(
			"opslevel client error",
			fmt.Sprintf("given repository identifier '%s' did not match found repository's id '%s' or alias '%s'",
				stateModel.Identifier.ValueString(),
				string(readRepository.Id),
				readRepository.DefaultAlias,
			),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *RepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[RepositoryResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	input := opslevel.RepositoryUpdateInput{
		Id: opslevel.ID(planModel.Id.ValueString()),
	}
	if !planModel.Owner.IsNull() {
		input.OwnerId = nullable(opslevel.NewID(planModel.Owner.ValueString()))
	}

	updatedRepository, err := r.client.UpdateRepository(input)
	if err != nil || updatedRepository == nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update repository, got error: %s", err))
		return
	}
	stateModel := NewRepositoryResourceModel(ctx, *updatedRepository)

	// Identifier from plan can be an id or alias
	switch planModel.Identifier.ValueString() {
	case string(updatedRepository.Id), updatedRepository.DefaultAlias:
		stateModel.Identifier = planModel.Identifier
	default:
		resp.Diagnostics.AddError(
			"opslevel client error",
			fmt.Sprintf("given repository identifier '%s' did not match found repository's id '%s' or alias '%s'",
				planModel.Identifier.ValueString(),
				string(updatedRepository.Id),
				updatedRepository.DefaultAlias,
			),
		)
		return
	}

	// Owner from plan can be an id or alias
	switch planModel.Owner.ValueString() {
	case string(updatedRepository.Owner.Id), updatedRepository.Owner.Alias:
		stateModel.Owner = planModel.Owner
	case "":
		stateModel.Owner = types.StringNull()
	default:
		resp.Diagnostics.AddError(
			"opslevel client error",
			fmt.Sprintf("repository owner found '%s' did not match given owner '%s'",
				stateModel.Owner.ValueString(),
				planModel.Owner.ValueString(),
			),
		)
		return
	}

	tflog.Trace(ctx, "updated a repository resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *RepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "unset a repository resource, actual repository not deleted")
}

func (r *RepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
