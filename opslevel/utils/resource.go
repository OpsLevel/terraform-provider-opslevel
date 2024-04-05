package utils

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/opslevel/opslevel-go/v2024"
)

type BaseResource struct {
	Client *opslevel.Client
	Helper ICRUD
}

func (r *BaseResource) SetHelper(helper ICRUD) resource.Resource {
	r.Helper = helper
	return r
}

func (r *BaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.Helper.GetTitle()
}

func (r *BaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.Helper.GetSchema()
}

func (r *BaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.Helper.Create(ctx, r.Client, req, resp)
}

func (r *BaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.Helper.Read(ctx, r.Client, req, resp)
}

func (r *BaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.Helper.Update(ctx, r.Client, req, resp)
}

func (r *BaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.Helper.Delete(ctx, r.Client, req, resp)
}

func (r *BaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
