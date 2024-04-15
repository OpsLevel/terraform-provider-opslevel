package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

var _ resource.ResourceWithConfigure = &FilterResource{}

var _ resource.ResourceWithImportState = &FilterResource{}

func NewFilterResource() resource.Resource {
	return &FilterResource{}
}

// FilterResource defines the resource implementation.
type FilterResource struct {
	CommonResourceClient
}

// FilterResourceModel describes the Filter managed resource.
type FilterResourceModel struct {
	Connective  types.String      `tfsdk:"connective"`
	Id          types.String      `tfsdk:"id"`
	LastUpdated types.String      `tfsdk:"last_updated"`
	Name        types.String      `tfsdk:"name"`
	Predicate   []filterPredicate `tfsdk:"predicate"`
}

type filterPredicate struct {
	CaseInsensitive types.Bool   `tfsdk:"case_insensitive"`
	CaseSensitive   types.Bool   `tfsdk:"case_sensitive"`
	Key             types.String `tfsdk:"key"`
	KeyData         types.String `tfsdk:"key_data"`
	Type            types.String `tfsdk:"type"`
	Value           types.String `tfsdk:"value"`
}

func convertPredicate(predicate opslevel.FilterPredicate) filterPredicate {
	convertedFilterPredicate := filterPredicate{
		CaseSensitive:   types.BoolNull(),
		CaseInsensitive: types.BoolNull(),
		Key:             RequiredStringValue(string(predicate.Key)),
		KeyData:         OptionalStringValue(predicate.KeyData),
		Type:            RequiredStringValue(string(predicate.Type)),
		Value:           OptionalStringValue(predicate.Value),
	}
	if predicate.CaseSensitive != nil {
		isCaseSensitive := *predicate.CaseSensitive
		if isCaseSensitive {
			convertedFilterPredicate.CaseSensitive = types.BoolValue(true)
			convertedFilterPredicate.CaseInsensitive = types.BoolValue(false)
		} else {
			convertedFilterPredicate.CaseSensitive = types.BoolValue(false)
			convertedFilterPredicate.CaseInsensitive = types.BoolValue(true)
		}
	}
	return convertedFilterPredicate
}

func NewFilterResourceModel(filter opslevel.Filter) FilterResourceModel {
	filterPredicates := []filterPredicate{}
	for _, predicate := range filter.Predicates {
		filterPredicates = append(filterPredicates, convertPredicate(predicate))
	}
	return FilterResourceModel{
		Connective: OptionalStringValue(string(filter.Connective)),
		Id:         ComputedStringValue(string(filter.Id)),
		Name:       RequiredStringValue(filter.Name),
		Predicate:  filterPredicates,
	}
}

func (r *FilterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filter"
}

func (r *FilterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Filter Resource",

		Attributes: map[string]schema.Attribute{
			"connective": schema.StringAttribute{
				Description: fmt.Sprintf(
					"The logical operator to be used in conjunction with predicates. Valid values are `%s`",
					strings.Join(opslevel.AllConnectiveEnum, "`, `"),
				),
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(opslevel.AllConnectiveEnum...),
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the filter.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The filter's display name.",
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"predicate": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"case_insensitive": schema.BoolAttribute{
							Description: "Option for determining whether to compare strings case-sensitively. Not settable for all predicate types.",
							Optional:    true,
							Computed:    true,
						},
						"case_sensitive": schema.BoolAttribute{
							Description: "Option for determining whether to compare strings case-sensitively. Not settable for all predicate types.",
							Optional:    true,
							Computed:    true,
						},
						"key": schema.StringAttribute{
							Description: fmt.Sprintf(
								"The condition key used by the predicate. Valid values are `%s`",
								strings.Join(opslevel.AllPredicateKeyEnum, "`, `"),
							),
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(opslevel.AllPredicateKeyEnum...),
							},
						},
						"key_data": schema.StringAttribute{
							Description: "Additional data used by the predicate. This field is used by predicates with key = 'tags' to specify the tag key. For example, to create a predicate for services containing the tag 'db:mysql', set key_data = 'db' and value = 'mysql'.",
							Optional:    true,
						},
						"type": schema.StringAttribute{
							Description: fmt.Sprintf(
								"The condition type used by the predicate. Valid values are `%s`",
								strings.Join(opslevel.AllPredicateTypeEnum, "`, `"),
							),
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(opslevel.AllPredicateTypeEnum...),
							},
						},
						"value": schema.StringAttribute{
							Description: "The condition value used by the predicate.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *FilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel FilterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}
	predicates, err := getFilterPredicates(planModel.Predicate)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("misconfigured filter predicate, got error: %s", err))
		return
	}

	filter, err := r.client.CreateFilter(opslevel.FilterCreateInput{
		Name:       planModel.Name.ValueString(),
		Predicates: predicates,
		Connective: getConnectiveEnum(planModel.Connective.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create filter, got error: %s", err))
		return
	}
	stateModel := NewFilterResourceModel(*filter)
	stateModel.Connective = planModel.Connective
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a filter resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *FilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var planModel FilterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filter, err := r.client.GetFilter(opslevel.ID(planModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read filter, got error: %s", err))
		return
	}
	stateModel := NewFilterResourceModel(*filter)
	stateModel.Connective = planModel.Connective

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *FilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel FilterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	predicates, err := getFilterPredicates(planModel.Predicate)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("misconfigured filter predicate, got error: %s", err))
		return
	}

	updatedFilter, err := r.client.UpdateFilter(opslevel.FilterUpdateInput{
		Id:         opslevel.ID(planModel.Id.ValueString()),
		Name:       opslevel.RefOf(planModel.Name.ValueString()),
		Predicates: predicates,
		Connective: getConnectiveEnum(planModel.Connective.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update filter, got error: %s", err))
		return
	}
	if planModel.Connective.ValueString() != "" {
		connectiveEnum := getConnectiveEnum(planModel.Connective.ValueString())
		updatedFilter.Connective = *connectiveEnum
	}

	stateModel := NewFilterResourceModel(*updatedFilter)
	stateModel.Connective = planModel.Connective
	stateModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a filter resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *FilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var planModel FilterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFilter(opslevel.ID(planModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete filter, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "deleted a filter resource")
}

func (r *FilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getConnectiveEnum(connective string) *opslevel.ConnectiveEnum {
	switch connective {
	case "and":
		return opslevel.RefTo(opslevel.ConnectiveEnumAnd)
	case "or":
		return opslevel.RefTo(opslevel.ConnectiveEnumOr)
	default:
		return nil
	}
}

func getFilterPredicates(predicates []filterPredicate) (*[]opslevel.FilterPredicateInput, error) {
	filterPredicateInputs := []opslevel.FilterPredicateInput{}
	for _, predicate := range predicates {
		tmpPredicateInput := opslevel.FilterPredicateInput{
			Key:     opslevel.PredicateKeyEnum(predicate.Key.ValueString()),
			KeyData: predicate.KeyData.ValueStringPointer(),
			Type:    opslevel.PredicateTypeEnum(predicate.Type.ValueString()),
			Value:   predicate.Value.ValueStringPointer(),
		}

		if predicate.CaseSensitive.ValueBool() && predicate.CaseInsensitive.ValueBool() {
			return nil, fmt.Errorf("a predicate should not have 'case_sensitive' and 'case_insensitive' set at the same time")
		} else if predicate.CaseSensitive.ValueBool() {
			tmpPredicateInput.CaseSensitive = opslevel.RefOf(true)
		} else if predicate.CaseInsensitive.ValueBool() {
			tmpPredicateInput.CaseSensitive = opslevel.RefOf(false)
		}

		filterPredicateInputs = append(filterPredicateInputs, tmpPredicateInput)
	}
	if len(filterPredicateInputs) > 0 {
		return &filterPredicateInputs, nil
	}
	return nil, nil
}
