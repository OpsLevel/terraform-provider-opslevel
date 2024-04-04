package opslevel

import (
	"context"
	"fmt"
	"strings"

	// "github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
		Id:         types.StringValue(string(filter.Id)),
		Name:       types.StringValue(filter.Name),
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
			},
			"last_updated": schema.StringAttribute{
				Optional: true,
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
							// NOTE: uncomment later if we want to enforce one field or the other
							// Validators: []validator.Bool{
							// 	boolvalidator.ConflictsWith(path.MatchRelative().AtName("case_sensitive")),
							// },
						},
						"case_sensitive": schema.BoolAttribute{
							Description: "Option for determining whether to compare strings case-sensitively. Not settable for all predicate types.",
							Optional:    true,
							// NOTE: uncomment later if we want to enforce one field or the other
							// Validators: []validator.Bool{
							// 	boolvalidator.ConflictsWith(path.MatchRelative().AtName("case_insensitive")),
							// },
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
	var data FilterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	predicates, err := getFilterPredicates(data.Predicate)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("misconfigured filter predicate, got error: %s", err))
		return
	}

	filter, err := r.client.CreateFilter(opslevel.FilterCreateInput{
		Name:       data.Name.ValueString(),
		Predicates: predicates,
		Connective: getConnectiveEnum(data.Connective.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to create filter, got error: %s", err))
		return
	}
	createdFilterResourceModel := NewFilterResourceModel(*filter)
	createdFilterResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "created a filter resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &createdFilterResourceModel)...)
}

func (r *FilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FilterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filter, err := r.client.GetFilter(opslevel.ID(data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to read filter, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("predicate 0 key: %s", filter.Predicates[0].Key))
	tflog.Info(ctx, fmt.Sprintf("predicate 0 key data: %s", filter.Predicates[0].KeyData))
	tflog.Info(ctx, fmt.Sprintf("predicate 0 type: %s", filter.Predicates[0].Type))
	tflog.Info(ctx, fmt.Sprintf("predicate 0 value: %s", filter.Predicates[0].Value))
	readFilterResourceModel := NewFilterResourceModel(*filter)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readFilterResourceModel)...)
}

func (r *FilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FilterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	predicates, err := getFilterPredicates(data.Predicate)
	if err != nil {
		resp.Diagnostics.AddError("Config error", fmt.Sprintf("misconfigured filter predicate, got error: %s", err))
		return
	}

	updatedFilter, err := r.client.UpdateFilter(opslevel.FilterUpdateInput{
		Id:         opslevel.ID(data.Id.ValueString()),
		Name:       data.Name.ValueStringPointer(),
		Predicates: predicates,
		Connective: getConnectiveEnum(data.Connective.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("opslevel client error", fmt.Sprintf("Unable to update filter, got error: %s", err))
		return
	}

	updatedFilterResourceModel := NewFilterResourceModel(*updatedFilter)
	updatedFilterResourceModel.LastUpdated = timeLastUpdated()

	tflog.Trace(ctx, "updated a filter resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedFilterResourceModel)...)
}

func (r *FilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FilterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFilter(opslevel.ID(data.Id.ValueString()))
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

// import (
// 	"errors"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceFilter() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a filter",
// 		Create:      wrap(resourceFilterCreate),
// 		Read:        wrap(resourceFilterRead),
// 		Update:      wrap(resourceFilterUpdate),
// 		Delete:      wrap(resourceFilterDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The filter's display name.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"predicate": {
// 				Type:        schema.TypeList,
// 				Description: "The list of predicates used to select which services apply to the filter.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"type": {
// 							Type:         schema.TypeString,
// 							Description:  "The condition type used by the predicate. Valid values are `contains`, `does_not_contain`, `does_not_equal`, `does_not_exist`, `ends_with`, `equals`, `exists`, `greater_than_or_equal_to`, `less_than_or_equal_to`, `starts_with`, `satisfies_version_constraint`, `matches_regex`, `matches`, `does_not_match`, `satisfies_jq_expression`",
// 							ForceNew:     false,
// 							Required:     true,
// 							ValidateFunc: validation.StringInSlice(opslevel.AllPredicateTypeEnum, false),
// 						},
// 						"value": {
// 							Type:        schema.TypeString,
// 							Description: "The condition value used by the predicate.",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 						"key": {
// 							Type:         schema.TypeString,
// 							Description:  "The condition key used by the predicate.",
// 							ForceNew:     false,
// 							Required:     true,
// 							ValidateFunc: validation.StringInSlice(opslevel.AllPredicateKeyEnum, false),
// 						},
// 						"key_data": {
// 							Type:        schema.TypeString,
// 							Description: "Additional data used by the predicate. This field is used by predicates with key = 'tags' to specify the tag key. For example, to create a predicate for services containing the tag 'db:mysql', set key_data = 'db' and value = 'mysql'.",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 						"case_insensitive": {
// 							Type:        schema.TypeBool,
// 							Description: "Option for determining whether to compare strings case-sensitively. Not settable for all predicate types.",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 						"case_sensitive": {
// 							Type:        schema.TypeBool,
// 							Description: "Option for determining whether to compare strings case-sensitively. Not settable for all predicate types.",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 					},
// 				},
// 			},
// 			"connective": {
// 				Type:        schema.TypeString,
// 				Description: "The logical operator to be used in conjunction with predicates.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 		},
// 	}
// }

// func getConnectiveEnum(d *schema.ResourceData) *opslevel.ConnectiveEnum {
// 	switch cleanerString(d.Get("connective").(string)) {
// 	case "or":
// 		return opslevel.RefTo(opslevel.ConnectiveEnumOr)
// 	case "and":
// 		return opslevel.RefTo(opslevel.ConnectiveEnumAnd)
// 	default:
// 		return nil
// 	}
// }

// func getFilterPredicates(d *schema.ResourceData) (*[]opslevel.FilterPredicateInput, error) {
// 	predicates := interfacesMaps(d.Get("predicate"))
// 	for _, pred := range predicates {
// 		if pred["case_sensitive"] == true && pred["case_insensitive"] == true {
// 			return nil, errors.New("a predicate should not have 'case_sensitive' and 'case_insensitive' set at the same time")
// 		}
// 	}
// 	return expandFilterPredicateInputs(predicates), nil
// }

// func resourceFilterCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	var err error
// 	predicates, err := getFilterPredicates(d)
// 	if err != nil {
// 		return err
// 	}
// 	input := opslevel.FilterCreateInput{
// 		Name:       d.Get("name").(string),
// 		Predicates: predicates,
// 		Connective: getConnectiveEnum(d),
// 	}
// 	if input.Predicates != nil && len(*input.Predicates) > 1 && input.Connective == nil {
// 		return errors.New("if there is more than 1 'predicate' then 'connective' must be set")
// 	}

// 	resource, err := client.CreateFilter(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceFilterRead(d, client)
// }

// func resourceFilterRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := opslevel.NewID(d.Id())

// 	resource, err := client.GetFilter(*id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("name", resource.Name); err != nil {
// 		return err
// 	}

// 	if err := d.Set("connective", string(resource.Connective)); err != nil {
// 		return err
// 	}

// 	if err := d.Set("predicate", flattenFilterPredicates(resource.Predicates)); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceFilterUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.FilterUpdateInput{
// 		Id:   *opslevel.NewID(d.Id()),
// 		Name: opslevel.RefOf(d.Get("name").(string)),
// 	}

// 	predicates, err := getFilterPredicates(d)
// 	if err != nil {
// 		return err
// 	}
// 	input.Predicates = predicates
// 	input.Connective = getConnectiveEnum(d)
// 	if input.Predicates != nil && len(*input.Predicates) > 1 && input.Connective == nil {
// 		return errors.New("if there is more than 1 'predicate' then 'connective' must be set")
// 	}

// 	_, err = client.UpdateFilter(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceFilterRead(d, client)
// }

// func resourceFilterDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := opslevel.NewID(d.Id())
// 	err := client.DeleteFilter(*id)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
