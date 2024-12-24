package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var _ resource.ResourceWithConfigure = &DomainResource{}

var _ resource.ResourceWithImportState = &DomainResource{}

func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

// DomainResource defines the resource implementation.
type DomainResource struct {
	CommonResourceClient
}

// DomainResourceModel describes the Domain managed resource.
type DomainResourceModel struct {
	Aliases     types.List   `tfsdk:"aliases"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Note        types.String `tfsdk:"note"`
	Owner       types.String `tfsdk:"owner"`
}

func NewDomainResourceModel(ctx context.Context, domain opslevel.Domain, givenModel DomainResourceModel) DomainResourceModel {
	domainResourceModel := DomainResourceModel{
		Aliases:     OptionalStringListValue(domain.Aliases),
		Description: StringValueFromResourceAndModelField(domain.Description, givenModel.Description),
		Id:          ComputedStringValue(string(domain.Id)),
		Name:        RequiredStringValue(domain.Name),
		Note:        StringValueFromResourceAndModelField(domain.Note, givenModel.Note),
		Owner:       StringValueFromResourceAndModelField(string(domain.Owner.Id()), givenModel.Owner),
	}

	return domainResourceModel
}

func (r *DomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *DomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Domain Resource",

		Attributes: map[string]schema.Attribute{
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The aliases of the domain.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the domain.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the domain.",
				Required:    true,
			},
			"note": schema.StringAttribute{
				Description: "Additional information about the domain.",
				Optional:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The id or alias of the team that owns the domain.",
				Optional:    true,
			},
		},
	}
}

func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[DomainResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.CreateDomain(opslevel.DomainInput{
		Description: planModel.Description.ValueStringPointer(),
		Name:        opslevel.RefOf(planModel.Name.ValueString()),
		Note:        planModel.Note.ValueStringPointer(),
		OwnerId:     GetTeamID(&resp.Diagnostics, r.client, planModel.Owner.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create domain, got error: %s", err))
		return
	}
	createdDomainResourceModel := NewDomainResourceModel(ctx, *resource, planModel)

	tflog.Trace(ctx, "created a domain resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdDomainResourceModel)...)
}

func (r *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[DomainResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.GetDomain(stateModel.Id.ValueString())
	if err != nil {
		if (resource == nil || resource.Id == "") && opslevel.IsOpsLevelApiError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read domain, got error: %s", err))
		return
	}
	readDomainResourceModel := NewDomainResourceModel(ctx, *resource, stateModel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readDomainResourceModel)...)
}

func (r *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[DomainResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.UpdateDomain(planModel.Id.ValueString(), opslevel.DomainInput{
		Description: opslevel.RefOf(planModel.Description.ValueString()),
		Name:        opslevel.RefOf(planModel.Name.ValueString()),
		Note:        opslevel.RefOf(planModel.Note.ValueString()),
		OwnerId:     GetTeamID(&resp.Diagnostics, r.client, planModel.Owner.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update domain, got error: %s", err))
		return
	}
	updatedDomainResourceModel := NewDomainResourceModel(ctx, *resource, planModel)

	tflog.Trace(ctx, "updated a domain resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedDomainResourceModel)...)
}

func (r *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := read[DomainResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDomain(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete domain, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a domain resource")
}

func (r *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
