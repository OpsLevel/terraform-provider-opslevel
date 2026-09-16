package opslevel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2026"
	"github.com/relvacode/iso8601"
)

// maxChecksPerCampaign mirrors the backend's own ceiling, enforced when checks
// are copied. Checking it up front turns a raw API failure part-way through an
// apply into a clear message.
const maxChecksPerCampaign = 25

var _ resource.ResourceWithConfigure = &CampaignCheckResource{}

var _ resource.ResourceWithImportState = &CampaignCheckResource{}

func NewCampaignCheckResource() resource.Resource {
	return &CampaignCheckResource{}
}

type CampaignCheckResource struct {
	CommonResourceClient
}

type CampaignCheckResourceModel struct {
	Id                types.String `tfsdk:"id"`
	CampaignId        types.String `tfsdk:"campaign_id"`
	CopiedFromCheckId types.String `tfsdk:"copied_from_check_id"`
	SourceCheckId     types.String `tfsdk:"source_check_id"`
	Name              types.String `tfsdk:"name"`
	Notes             types.String `tfsdk:"notes"`
	Owner             types.String `tfsdk:"owner"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	EnableOn          types.String `tfsdk:"enable_on"`
	Type              types.String `tfsdk:"type"`
	Description       types.String `tfsdk:"description"`
	Filter            types.String `tfsdk:"filter"`
}

// SelectCreatedCheck picks the single check a copy call created.
//
// The copy mutation reports exactly which checks it made, so this is a direct
// answer rather than an inference. Anything other than one check means we cannot
// say which check this resource owns, and guessing risks adopting - and later
// deleting - a check we did not create.
func SelectCreatedCheck(created []opslevel.Check) (*opslevel.Check, error) {
	switch len(created) {
	case 1:
		return &created[0], nil
	case 0:
		return nil, fmt.Errorf("the copy reported no created checks")
	default:
		return nil, fmt.Errorf("the copy reported %d created checks, expected exactly 1", len(created))
	}
}

// CampaignCheckWasDeleted reports whether a check the API cannot find has really
// been deleted, given the status of the campaign it belongs to.
//
// Once a campaign has ended, every check on it fails the API's active? gate and
// reads as missing even though it still exists and is still attached. Treating
// that as a deletion would drop the resource from state and make Terraform plan
// a re-copy, which is not the remedy.
func CampaignCheckWasDeleted(campaignStatus opslevel.CampaignStatusEnum) bool {
	return campaignStatus != opslevel.CampaignStatusEnumEnded
}

func NewCampaignCheckResourceModel(check opslevel.Check, givenModel CampaignCheckResourceModel) CampaignCheckResourceModel {
	model := CampaignCheckResourceModel{
		Id:                ComputedStringValue(string(check.Id)),
		CampaignId:        RequiredStringValue(string(check.Campaign.Id)),
		CopiedFromCheckId: givenModel.CopiedFromCheckId,
		SourceCheckId:     ComputedStringValue(string(check.SourceCheck.Id)),
		Name:              ComputedStringValue(check.Name),
		Notes:             StringValueFromResourceAndModelField(check.Notes, givenModel.Notes),
		Owner:             OptionalStringValue(string(check.Owner.Team.Id)),
		Enabled:           types.BoolValue(check.Enabled),
		Type:              ComputedStringValue(string(check.Type)),
		Description:       ComputedStringValue(check.Description),
		Filter:            ComputedStringValue(string(check.Filter.Id)),
	}
	// enable_on is write-only from the API's perspective once a check is enabled,
	// so the plan value is preserved rather than round-tripped, matching how the
	// rubric check resources handle it.
	if givenModel.EnableOn.IsNull() {
		model.EnableOn = types.StringNull()
	} else {
		model.EnableOn = givenModel.EnableOn
	}
	return model
}

func (r *CampaignCheckResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_campaign_check"
}

func (r *CampaignCheckResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `A rubric check copied onto a campaign.

Checks can only be added to a campaign by copying an existing rubric check. The copy is a
brand new, independent check with its own id - editing the rubric check it came from does
not change the copy.

Do not use this resource on a campaign that also sets ` + "`check_ids`" + ` on its
` + "`opslevel_campaign`" + ` resource. The two manage the same checks by different means and
will fight, deleting and re-copying checks on every apply.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the check that was created on the campaign.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"campaign_id": schema.StringAttribute{
				Description: "The id of the campaign to copy the check onto. Cannot be changed once set.",
				Required:    true,
				Validators:  []validator.String{IdStringValidator()},
				PlanModifiers: []planmodifier.String{
					ImmutableStringPlanModifier(),
				},
			},
			"copied_from_check_id": schema.StringAttribute{
				Description: "The id of the rubric check to copy. Cannot be changed once set, and cannot be " +
					"set on an imported check - the copy has already happened and there is nothing to record.",
				Optional:   true,
				Validators: []validator.String{IdStringValidator()},
				PlanModifiers: []planmodifier.String{
					ImmutableStringPlanModifier(),
				},
			},
			"source_check_id": schema.StringAttribute{
				Description: "The id of the check this check was copied from, as reported by the API. " +
					"Null if the source check has since been deleted, or if the copy was made before " +
					"OpsLevel began recording this.",
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the check, inherited from the check it was copied from.",
				Computed:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Additional information to display to the service owner about the check.",
				Optional:    true,
				Computed:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The id of the team that owns the check.",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{IdStringValidator()},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the check is enabled or not. Copies are always created disabled. " +
					"Do not use this field in tandem with 'enable_on'.",
				Optional:   true,
				Computed:   true,
				Default:    booldefault.StaticBool(false),
				Validators: []validator.Bool{boolvalidator.ConflictsWith(path.MatchRoot("enable_on"))},
			},
			"enable_on": schema.StringAttribute{
				Description: "The date when the check will be automatically enabled.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("enabled"))},
			},
			"type": schema.StringAttribute{
				Description: "The type of check that was copied.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the check type's purpose.",
				Computed:    true,
			},
			"filter": schema.StringAttribute{
				Description: "The id of the filter applied to this check. Read-only: a campaign check with no " +
					"filter of its own reports the campaign's filter, and the API cannot distinguish the two.",
				Computed: true,
			},
		},
	}
}

func (r *CampaignCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planModel := read[CampaignCheckResourceModel](ctx, &resp.Diagnostics, req.Plan)
	if resp.Diagnostics.HasError() {
		return
	}

	if planModel.CopiedFromCheckId.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("copied_from_check_id"),
			"Missing copied_from_check_id",
			"A campaign check can only be created by copying a rubric check, so copied_from_check_id "+
				"is required when creating one. It is optional only so that existing campaign checks, "+
				"which have no record of where they came from, can be imported.",
		)
		return
	}

	campaignId := opslevel.ID(planModel.CampaignId.ValueString())
	if !r.campaignHasRoomForAnotherCheck(ctx, &resp.Diagnostics, campaignId) {
		return
	}

	_, createdChecks, err := r.client.CopyChecksToCampaign(opslevel.ChecksCopyToCampaignInput{
		CampaignId: campaignId,
		CheckIds:   []opslevel.ID{opslevel.ID(planModel.CopiedFromCheckId.ValueString())},
	})
	if err != nil {
		title, detail := formatOpslevelError("copy check to campaign", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	check, err := SelectCreatedCheck(createdChecks)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not identify the copied check",
			fmt.Sprintf("The check was copied onto campaign %s but %s. The copy may have succeeded; "+
				"check the campaign in OpsLevel before retrying.", campaignId, err),
		)
		return
	}

	// The copy mutation accepts no field values, so anything the user asked for
	// has to be applied afterwards.
	check = r.applyEditableFields(ctx, &resp.Diagnostics, *check, planModel)
	if resp.Diagnostics.HasError() {
		return
	}

	createdModel := NewCampaignCheckResourceModel(*check, planModel)
	tflog.Trace(ctx, "created a campaign check resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdModel)...)
}

func (r *CampaignCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateModel := read[CampaignCheckResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	check, err := r.client.GetCheck(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil || check == nil || check.Id == "" {
		if !opslevel.IsOpsLevelApiError(err) && err != nil {
			title, detail := formatOpslevelError("read campaign check", err)
			resp.Diagnostics.AddError(title, detail)
			return
		}
		r.handleMissingCheck(ctx, resp, stateModel)
		return
	}

	readModel := NewCampaignCheckResourceModel(*check, stateModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, &readModel)...)
}

// handleMissingCheck decides what a not-found check means. See
// CampaignCheckWasDeleted.
func (r *CampaignCheckResource) handleMissingCheck(ctx context.Context, resp *resource.ReadResponse, stateModel CampaignCheckResourceModel) {
	campaignId := opslevel.ID(stateModel.CampaignId.ValueString())
	campaign, err := r.client.GetCampaign(campaignId)
	if err != nil || campaign == nil {
		// The campaign cannot be read either, so the check's absence cannot be
		// explained by the campaign having ended. Treat it as gone.
		resp.State.RemoveResource(ctx)
		return
	}

	if CampaignCheckWasDeleted(campaign.Status) {
		tflog.Info(ctx, "campaign check no longer exists, removing from state",
			map[string]any{"check_id": stateModel.Id.ValueString()})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.AddWarning(
		"Campaign has ended",
		fmt.Sprintf("Campaign %s has ended, so OpsLevel no longer reports its checks through the API. "+
			"Check %s is kept in state unchanged - it still exists and is still on the campaign, but "+
			"it cannot be read or updated while the campaign is ended.",
			campaignId, stateModel.Id.ValueString()),
	)
}

func (r *CampaignCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planModel := read[CampaignCheckResourceModel](ctx, &resp.Diagnostics, req.Plan)
	stateModel := read[CampaignCheckResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	check, err := r.client.GetCheck(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil {
		title, detail := formatOpslevelError("read campaign check for update", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}

	updated := r.applyEditableFields(ctx, &resp.Diagnostics, *check, planModel)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedModel := NewCampaignCheckResourceModel(*updated, planModel)
	tflog.Trace(ctx, "updated a campaign check resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedModel)...)
}

func (r *CampaignCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	stateModel := read[CampaignCheckResourceModel](ctx, &resp.Diagnostics, req.State)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCheck(opslevel.ID(stateModel.Id.ValueString())); err != nil {
		title, detail := formatOpslevelError("delete campaign check", err)
		resp.Diagnostics.AddError(title, detail)
		return
	}
	tflog.Trace(ctx, "deleted a campaign check resource")
}

// ImportState takes the copied check's own id. campaign_id is derived on read,
// and copied_from_check_id stays null - nothing records what an existing copy was
// made from.
func (r *CampaignCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// applyEditableFields pushes the model's editable values onto an existing check,
// returning the check unchanged when there is nothing to apply.
//
// A check's updatable fields depend on its concrete type, which is only known at
// runtime, so the shared fields are marshalled and handed to
// UnmarshalCheckUpdateInput to produce the right type-specific input.
func (r *CampaignCheckResource) applyEditableFields(ctx context.Context, diags *diag.Diagnostics, check opslevel.Check, model CampaignCheckResourceModel) *opslevel.Check {
	input := opslevel.CheckUpdateInput{Id: check.Id}
	changed := false

	if !model.Notes.IsNull() && !model.Notes.IsUnknown() && model.Notes.ValueString() != check.Notes {
		notes := model.Notes.ValueString()
		input.Notes = &notes
		changed = true
	}
	if !model.Owner.IsNull() && !model.Owner.IsUnknown() && model.Owner.ValueString() != string(check.Owner.Team.Id) {
		input.Owner = opslevel.NewNullableFrom(opslevel.ID(model.Owner.ValueString()))
		changed = true
	}
	if !model.Enabled.IsNull() && !model.Enabled.IsUnknown() && model.Enabled.ValueBool() != check.Enabled {
		input.Enabled = opslevel.NewNullableFrom(model.Enabled.ValueBool())
		changed = true
	}
	if !model.EnableOn.IsNull() && !model.EnableOn.IsUnknown() {
		when, err := iso8601.ParseString(model.EnableOn.ValueString())
		if err != nil {
			diags.AddAttributeError(
				path.Root("enable_on"),
				"Invalid enable_on",
				fmt.Sprintf("Could not parse %q as a date: %s", model.EnableOn.ValueString(), err),
			)
			return nil
		}
		input.EnableOn = opslevel.NewNullableFrom(iso8601.Time{Time: when})
		changed = true
	}

	if !changed {
		return &check
	}

	data, err := json.Marshal(input)
	if err != nil {
		diags.AddError("Could not build check update", err.Error())
		return nil
	}
	typedInput, err := opslevel.UnmarshalCheckUpdateInput(check.Type, data)
	if err != nil {
		diags.AddError(
			"Could not build check update",
			fmt.Sprintf("Check %s is of type %q, which this provider could not build an update for: %s",
				check.Id, check.Type, err),
		)
		return nil
	}

	updated, err := r.client.UpdateCheck(typedInput)
	if err != nil {
		title, detail := formatOpslevelError("update campaign check", err)
		diags.AddError(title, detail)
		return nil
	}
	return updated
}

// campaignHasRoomForAnotherCheck reports whether copying one more check onto the
// campaign would exceed the backend's ceiling. Terraform cannot catch this at
// plan time, because each resource instance plans independently and none of them
// can see how many others are aimed at the same campaign.
func (r *CampaignCheckResource) campaignHasRoomForAnotherCheck(ctx context.Context, diags *diag.Diagnostics, campaignId opslevel.ID) bool {
	checks, err := r.client.ListCampaignChecks(campaignId)
	if err != nil {
		title, detail := formatOpslevelError("list campaign checks", err)
		diags.AddError(title, detail)
		return false
	}
	if len(checks) >= maxChecksPerCampaign {
		diags.AddError(
			"Campaign check limit reached",
			fmt.Sprintf("Campaign %s already has %d checks, and OpsLevel allows at most %d. "+
				"Remove a check from the campaign before adding another.",
				campaignId, len(checks), maxChecksPerCampaign),
		)
		return false
	}
	return true
}
