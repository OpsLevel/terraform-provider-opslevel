package provider

import (
	"context"
	"fmt"

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

var _ resource.Resource = &DomainResource{}

// _ resource.ResourceWithImportState = &DomainResource{}

func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

// DomainResource defines the resource implementation.
type DomainResource struct {
	client *opslevel.Client
}

// DomainResourceModel describes the resource data model.
type DomainResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Aliases     types.List   `tfsdk:"aliases"`
	Identifier  types.String `tfsdk:"identifier"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Note        types.String `tfsdk:"note"`
	Owner       types.String `tfsdk:"owner"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *DomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *DomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Domain Resource",

		Attributes: map[string]schema.Attribute{
			"last_updated": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"identifier": schema.StringAttribute{
				Description: "The id or alias of the domain to find.",
				Optional:    true,
			},
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The aliases of the domain.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the domain.",
				Computed:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The id of the team that owns the domain.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the domain.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*opslevel.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *opslevel.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.CreateDomain(opslevel.DomainInput{
		Name:        opslevel.NewString(data.Name.String()),
		Description: opslevel.NewString(data.Description.String()),
		OwnerId:     opslevel.NewID(data.Owner.String()),
		Note:        opslevel.NewString(data.Note.String()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create domain, got error: %s", err))
		return
	}
	data.Id = types.StringValue(string(resource.Id))
	// NOTE: copy common.go to get timeLastUpdated()

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a domain resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
	domainAliases, d := types.ListValueFrom(ctx, types.StringType, resource.Aliases)
	resp.Diagnostics.Append(d...)

	data.Aliases = domainAliases
	data.Description = types.StringValue(resource.Description)
	data.Name = types.StringValue(string(resource.Owner.Id()))
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

	resource, err := r.client.UpdateDomain(data.Identifier.String(), opslevel.DomainInput{
		Name:        opslevel.NewString(data.Name.String()),
		Description: opslevel.NewString(data.Description.String()),
		OwnerId:     opslevel.NewID(data.Owner.String()),
		Note:        opslevel.NewString(data.Note.String()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update domain, got error: %s", err))
		return
	}
	data.Id = types.StringValue(string(resource.Id))
	// NOTE: copy common.go to get timeLastUpdated()
	// data.LastUpdated - time

	// d.Set("last_updated", timeLastUpdated())
	// return resourceDomainRead(d, client)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// NOTE: REMOVE AFTER DONE
	// id := d.Id()
	err := r.client.DeleteDomain(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete domain, got error: %s", err))
		return
	}
}

func (r *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
