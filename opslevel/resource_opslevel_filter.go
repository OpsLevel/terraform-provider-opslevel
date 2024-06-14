package opslevel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

var (
	_ resource.ResourceWithConfigure      = &FilterResource{}
	_ resource.ResourceWithImportState    = &FilterResource{}
	_ resource.ResourceWithValidateConfig = &FilterResource{}
)

func NewFilterResource() resource.Resource {
	return &FilterResource{}
}

// FilterResource defines the resource implementation.
type FilterResource struct {
	CommonResourceClient
}

// FilterResourceModel describes the Filter managed resource.
type FilterResourceModel struct {
	Connective types.String `tfsdk:"connective"`
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Predicate  types.List   `tfsdk:"predicate"`
}

type filterPredicateModel struct {
	CaseInsensitive types.Bool   `tfsdk:"case_insensitive"`
	CaseSensitive   types.Bool   `tfsdk:"case_sensitive"`
	Key             types.String `tfsdk:"key"`
	KeyData         types.String `tfsdk:"key_data"`
	Type            types.String `tfsdk:"type"`
	Value           types.String `tfsdk:"value"`
}

var filterPredicateType = map[string]attr.Type{
	"case_insensitive": types.BoolType,
	"case_sensitive":   types.BoolType,
	"key":              types.StringType,
	"key_data":         types.StringType,
	"type":             types.StringType,
	"value":            types.StringType,
}

func (fp filterPredicateModel) Validate() error {
	opslevelFilterPredicate := opslevel.FilterPredicate{
		CaseSensitive: fp.CaseSensitive.ValueBoolPointer(),
		Key:           opslevel.PredicateKeyEnum(fp.Key.ValueString()),
		KeyData:       fp.KeyData.ValueString(),
		Type:          opslevel.PredicateTypeEnum(fp.Type.ValueString()),
		Value:         fp.Value.ValueString(),
	}
	if opslevelFilterPredicate.CaseSensitive == nil && !fp.CaseInsensitive.IsNull() {
		opslevelFilterPredicate.CaseSensitive = fp.CaseInsensitive.ValueBoolPointer()
	}
	return opslevelFilterPredicate.Validate()
}

func NewFilterResourceModel(ctx context.Context, filter opslevel.Filter, givenModel FilterResourceModel) (FilterResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var filterPredicateAttrs []attr.Value
	var filterPredicatesListValue basetypes.ListValue
	var givenPredicateModels []filterPredicateModel

	// Convert predicates from plan to slice of models
	givenModel.Predicate.ElementsAs(ctx, &givenPredicateModels, false)

	for _, opslevelPredicate := range filter.Predicates {
		predicateObj := OpslevelFilterPredicateToObjectValue(nil, &opslevelPredicate)
		attrs := predicateObj.Attributes()

		// find predicate from plan that matches predicateObj from API
		foundPlanPredModel, extractDiags := ExtractFilterPredicateModel(ctx, attrs, givenPredicateModels)
		diags.Append(extractDiags...)
		if diags.HasError() {
			return FilterResourceModel{}, diags
		}

		if !foundPlanPredModel.CaseSensitive.IsNull() && !foundPlanPredModel.CaseSensitive.IsUnknown() {
			attrs["case_sensitive"] = types.BoolValue(*opslevelPredicate.CaseSensitive)
		} else {
			attrs["case_sensitive"] = types.BoolNull()
		}
		if !foundPlanPredModel.CaseInsensitive.IsNull() && !foundPlanPredModel.CaseInsensitive.IsUnknown() {
			attrs["case_insensitive"] = types.BoolValue(*opslevelPredicate.CaseSensitive)
		} else {
			attrs["case_insensitive"] = types.BoolNull()
		}

		predicateObj = types.ObjectValueMust(filterPredicateType, attrs)
		filterPredicateAttrs = append(filterPredicateAttrs, predicateObj)
	}
	if len(filterPredicateAttrs) == 0 {
		filterPredicatesListValue = types.ListNull(types.ObjectType{AttrTypes: filterPredicateType})
	} else {
		filterPredicatesListValue = types.ListValueMust(types.ObjectType{AttrTypes: filterPredicateType}, filterPredicateAttrs)
	}

	return FilterResourceModel{
		Connective: givenModel.Connective,
		Id:         ComputedStringValue(string(filter.Id)),
		Name:       RequiredStringValue(filter.Name),
		Predicate:  filterPredicatesListValue,
	}, diags
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
					"The logical operator to be used in conjunction with predicates. One of `%s`",
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
							Validators:  []validator.Bool{boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("case_sensitive"))},
						},
						"case_sensitive": schema.BoolAttribute{
							Description: "Option for determining whether to compare strings case-sensitively. Not settable for all predicate types.",
							Optional:    true,
							Computed:    true,
							Validators:  []validator.Bool{boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("case_insensitive"))},
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

func (r *FilterResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var configModel FilterResourceModel
	var predicateModels []filterPredicateModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &configModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(configModel.Predicate.ElementsAs(ctx, &predicateModels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, filterPredicate := range predicateModels {
		if err := filterPredicate.Validate(); err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("predicate"), "Invalid Attribute Configuration", err.Error())
		}
	}
}

func (r *FilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel FilterResourceModel
	var predicateModels []filterPredicateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(planModel.Predicate.ElementsAs(ctx, &predicateModels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	predicates, err := getFilterPredicates(predicateModels)
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
	stateModel, diags := NewFilterResourceModel(ctx, *filter, planModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created a filter resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
}

func (r *FilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateModel FilterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filter, err := r.client.GetFilter(opslevel.ID(stateModel.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read filter, got error: %s", err))
		return
	}
	verifiedStateModel, diags := NewFilterResourceModel(ctx, *filter, stateModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &verifiedStateModel)...)
}

func (r *FilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel FilterResourceModel
	var predicateModels []filterPredicateModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(planModel.Predicate.ElementsAs(ctx, &predicateModels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	predicates, err := getFilterPredicates(predicateModels)
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

	stateModel, diags := NewFilterResourceModel(ctx, *updatedFilter, planModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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

func getFilterPredicates(predicates []filterPredicateModel) (*[]opslevel.FilterPredicateInput, error) {
	filterPredicateInputs := []opslevel.FilterPredicateInput{}

	for _, predicate := range predicates {
		tmpPredicateInput := opslevel.FilterPredicateInput{
			Key:     opslevel.PredicateKeyEnum(predicate.Key.ValueString()),
			KeyData: predicate.KeyData.ValueStringPointer(),
			Type:    opslevel.PredicateTypeEnum(predicate.Type.ValueString()),
			Value:   predicate.Value.ValueStringPointer(),
		}
		isCaseSensitiveSet := !predicate.CaseSensitive.IsNull() && !predicate.CaseSensitive.IsUnknown()
		isCaseInsensitiveSet := !predicate.CaseInsensitive.IsNull() && !predicate.CaseInsensitive.IsUnknown()

		if isCaseSensitiveSet && isCaseInsensitiveSet {
			return nil, fmt.Errorf("a predicate should not have 'case_sensitive' and 'case_insensitive' set at the same time")
		}
		if isCaseSensitiveSet {
			tmpPredicateInput.CaseSensitive = opslevel.RefOf(predicate.CaseSensitive.ValueBool())
		} else if isCaseInsensitiveSet {
			tmpPredicateInput.CaseSensitive = opslevel.RefOf(predicate.CaseInsensitive.ValueBool())
		}

		filterPredicateInputs = append(filterPredicateInputs, tmpPredicateInput)
	}

	return &filterPredicateInputs, nil
}
