package opslevel

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	LastUpdated types.String `tfsdk:"last_updated"`
	Name        types.String `tfsdk:"name"`
	Note        types.String `tfsdk:"note"`
	Owner       types.String `tfsdk:"owner"`
}

func NewDomainResourceModel(ctx context.Context, domain opslevel.Domain) (DomainResourceModel, diag.Diagnostics) {
	var domainResourceModel DomainResourceModel

	domainAliases, diags := types.ListValueFrom(ctx, types.StringType, domain.Aliases)
	domainResourceModel.Aliases = domainAliases
	domainResourceModel.Description = types.StringValue(string(domain.Description))
	domainResourceModel.Id = types.StringValue(string(domain.Id))
	domainResourceModel.Name = types.StringValue(string(domain.Name))
	domainResourceModel.Note = types.StringValue(string(domain.Note))
	domainResourceModel.Owner = types.StringValue(string(domain.Owner.Id()))
	domainResourceModel.LastUpdated = types.StringValue(timeLastUpdated())

	return domainResourceModel, diags
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
			"last_updated": schema.StringAttribute{
				Optional: true,
				Computed: true,
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
				Description: "The id of the team that owns the domain.",
				Optional:    true,
			},
		},
	}
}

func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.CreateDomain(opslevel.DomainInput{
		Description: opslevel.RefOf(data.Description.ValueString()),
		Name:        opslevel.RefOf(data.Name.ValueString()),
		Note:        opslevel.RefOf(data.Note.ValueString()),
		OwnerId:     opslevel.NewID(data.Owner.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create domain, got error: %s", err))
		return
	}
	createdDomainResourceModel, diags := NewDomainResourceModel(ctx, *resource)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "created a domain resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdDomainResourceModel)...)
}

func (r *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.GetDomain(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read domain, got error: %s", err))
		return
	}
	domainAliases, d := types.ListValueFrom(ctx, types.StringType, resource.ManagedAliases)
	resp.Diagnostics.Append(d...)

	data.Aliases = domainAliases
	data.Description = types.StringValue(resource.Description)
	data.Name = types.StringValue(string(resource.Name))
	data.Note = types.StringValue(resource.Note)
	data.Owner = types.StringValue(string(resource.Owner.Id()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.UpdateDomain(data.Id.ValueString(), opslevel.DomainInput{
		Description: opslevel.RefOf(data.Description.ValueString()),
		Name:        opslevel.RefOf(data.Name.ValueString()),
		Note:        opslevel.RefOf(data.Note.ValueString()),
		OwnerId:     opslevel.NewID(data.Owner.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update domain, got error: %s", err))
		return
	}
	updatedDomainResourceModel, diags := NewDomainResourceModel(ctx, *resource)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "updated a domain resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedDomainResourceModel)...)
}

func (r *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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
