package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

type ICRUD interface {
	GetTitle() string
	GetSchema() schema.Schema
	Create(ctx context.Context, client *opslevel.Client, req resource.CreateRequest, resp *resource.CreateResponse)
	Read(ctx context.Context, client *opslevel.Client, req resource.ReadRequest, resp *resource.ReadResponse)
	Update(ctx context.Context, client *opslevel.Client, req resource.UpdateRequest, resp *resource.UpdateResponse)
	Delete(ctx context.Context, client *opslevel.Client, req resource.DeleteRequest, resp *resource.DeleteResponse)
}

type CRUD[TRes any, TModel any] struct {
	Title      string
	Schema     schema.Schema
	BuildModel func(ctx context.Context, res TRes) (TModel, diag.Diagnostics)
	DoCreate   func(client *opslevel.Client, model TModel) (TRes, diag.Diagnostics)
	DoRead     func(client *opslevel.Client, model TModel) (TRes, diag.Diagnostics)
	DoUpdate   func(client *opslevel.Client, model TModel) (TRes, diag.Diagnostics)
	DoDelete   func(client *opslevel.Client, model TModel) error
}

func (s CRUD[TRes, TModel]) GetTitle() string {
	return s.Title
}

func (s CRUD[TRes, TModel]) GetSchema() schema.Schema {
	return s.Schema
}

func (s CRUD[TRes, TModel]) Create(ctx context.Context, client *opslevel.Client, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model TModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, diags := s.DoCreate(client, model)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	created, diags := s.BuildModel(ctx, res)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, fmt.Sprintf("created a '%s' resource", s.Title))
	resp.Diagnostics.Append(resp.State.Set(ctx, &created)...)
}

func (s CRUD[TRes, TModel]) Read(ctx context.Context, client *opslevel.Client, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model TModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := s.DoRead(client, model)
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read check manual, got error: %s", err))
		return
	}
	created, diags := s.BuildModel(ctx, res)
	resp.Diagnostics.Append(diags...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &created)...)
}

func (s CRUD[TRes, TModel]) Update(ctx context.Context, client *opslevel.Client, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model TModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, diags := s.DoUpdate(client, model)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	created, diags := s.BuildModel(ctx, res)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, fmt.Sprintf("updated a '%s' resource", s.Title))
	resp.Diagnostics.Append(resp.State.Set(ctx, &created)...)
}

func (s CRUD[TRes, TModel]) Delete(ctx context.Context, client *opslevel.Client, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model TModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := s.DoDelete(client, model)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete check manual, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a check manual resource")
}
