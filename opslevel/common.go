package opslevel

import (
	"context"
	"fmt"
	"slices"
	"strings"

	// "sort"
	"strconv"
	"time"

	// "strings"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	// "github.com/mitchellh/mapstructure"
	// "github.com/rs/zerolog/log"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2024"
)

// func cleanerString(s string) string {
// 	return strings.TrimSpace(strings.ToLower(s))
// }

// // interfacesMaps converts an interface{} into a []map[string]interface{}. This is a useful conversion for passing
// // schema.ResourceData objects from terraform into mapstructure.Decode to get actual struct types.
// func interfacesMaps(i interface{}) []map[string]interface{} {
// 	// interface{} 					to 		[]interface{}								segment into slices.
// 	interfaces := i.([]interface{})
// 	// interface{}					to		[]map[string]interface{}					convert each slice item into a map.
// 	mapStringInterfaces := make([]map[string]interface{}, len(interfaces))
// 	for i, item := range interfaces {
// 		mapStringInterfaces[i] = item.(map[string]interface{})
// 	}
// 	return mapStringInterfaces
// }

// var DefaultPredicateDescription = "A condition that should be satisfied."

const providerIssueUrl = "https://github.com/OpsLevel/terraform-provider-opslevel/issues"

type CommonResourceClient struct {
	client *opslevel.Client
}

// Configure sets up the OpsLevel client for datasources and resources
func (d *CommonResourceClient) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*opslevel.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("expected *opslevel.Client, got: %T please report this issue to the provider developers at %s.", req.ProviderData, providerIssueUrl),
		)

		return
	}

	d.client = client
}

type CommonDataSourceClient struct {
	client *opslevel.Client
}

// Configure sets up the OpsLevel client for datasources and resources
func (d *CommonDataSourceClient) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*opslevel.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("expected *opslevel.Client, got: %T please report this issue to the provider developers at %s.", req.ProviderData, providerIssueUrl),
		)

		return
	}

	d.client = client
}

func timeID() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func timeLastUpdated() basetypes.StringValue {
	return types.StringValue(time.Now().Format(time.RFC850))
}

// getValidOwner will compare the expected owner from the terraform plan OR state versus what is found in OpsLevel
// if the owner is not as expected it will return an error
func getValidOwner(client *opslevel.Client, resource opslevel.HasTeam, expectedOwner string) (types.String, error) {
	// validate that the resource does not have an owner set
	if expectedOwner == "" {
		if expectedOwner == string(resource.GetTeamId().Id) {
			return types.StringNull(), nil
		}
		return types.String{}, fmt.Errorf("expected no owner to be set, got owner with ID: '%s'", resource.GetTeamId().Id)
	}

	// validate that the resource owner ID is correct
	if opslevel.IsID(expectedOwner) {
		if expectedOwner == string(resource.GetTeamId().Id) {
			return types.StringValue(expectedOwner), nil
		}
		return types.String{}, fmt.Errorf("expected owner with ID '%s', got owner with ID '%s'", expectedOwner, resource.GetTeamId().Id)
	}

	// validate that the resource owner alias is correct
	if expectedOwner == resource.GetTeamId().Alias {
		return types.StringValue(expectedOwner), nil
	}
	// complex case - need to check through non-default aliases
	team, err := resource.GetTeam(client)
	if err != nil {
		return types.String{}, fmt.Errorf("error fetching owner team on resource: '%w'", err)
	}
	if team == nil || team.Id == "" {
		return types.String{}, fmt.Errorf("owner team on resource was not found")
	}
	if !slices.Contains(team.Aliases, expectedOwner) {
		return types.String{}, fmt.Errorf("owner team on resource does not have expected alias '%s'", expectedOwner)
	}
	return types.StringValue(expectedOwner), nil
}

// func wrap(handler func(data *schema.ResourceData, client *opslevel.Client) error) func(d *schema.ResourceData, meta interface{}) error {
// 	return func(data *schema.ResourceData, meta interface{}) error {
// 		client := meta.(*opslevel.Client)
// 		return handler(data, client)
// 	}
// }

// func getStringArray(d *schema.ResourceData, key string) []string {
// 	output := make([]string, 0)
// 	data, ok := d.GetOk(key)
// 	if !ok {
// 		return output
// 	}
// 	for _, item := range data.([]interface{}) {
// 		output = append(output, item.(string))
// 	}
// 	return output
// }

// func findService(aliasKey string, idKey string, d *schema.ResourceData, client *opslevel.Client) (*opslevel.Service, error) {
// 	alias := d.Get(aliasKey).(string)
// 	id := d.Get(idKey)
// 	if alias == "" && id == "" {
// 		return nil, fmt.Errorf("must provide one of `%s` or `%s` field to find by", aliasKey, idKey)
// 	}
// 	var resource *opslevel.Service
// 	if id == "" {
// 		found, err := client.GetServiceWithAlias(alias)
// 		if err != nil {
// 			return nil, err
// 		}
// 		resource = found
// 	} else {
// 		found, err := client.GetService(*opslevel.NewID(id.(string)))
// 		if err != nil {
// 			return nil, err
// 		}
// 		resource = found
// 	}
// 	if resource.Id == "" {
// 		return nil, fmt.Errorf("unable to find service with alias=`%s` or id=`%s`", alias, id.(string))
// 	}
// 	return resource, nil
// }

// func findRepository(aliasKey string, idKey string, d *schema.ResourceData, client *opslevel.Client) (*opslevel.Repository, error) {
// 	alias := d.Get(aliasKey).(string)
// 	id := d.Get(idKey)
// 	if alias == "" && id == "" {
// 		return nil, fmt.Errorf("must provide one of `%s` or `%s` field to find by", aliasKey, idKey)
// 	}
// 	var resource *opslevel.Repository
// 	if id == "" {
// 		found, err := client.GetRepositoryWithAlias(alias)
// 		if err != nil {
// 			return nil, err
// 		}
// 		resource = found
// 	} else {
// 		found, err := client.GetRepository(*opslevel.NewID(id.(string)))
// 		if err != nil {
// 			return nil, err
// 		}
// 		resource = found
// 	}
// 	if resource.Id == "" {
// 		return nil, fmt.Errorf("unable to find repository with alias=`%s` or id=`%s`", alias, id.(string))
// 	}
// 	return resource, nil
// }

// func findTeam(aliasKey string, idKey string, d *schema.ResourceData, client *opslevel.Client) (*opslevel.Team, error) {
// 	alias := d.Get(aliasKey).(string)
// 	id := d.Get(idKey)
// 	if alias == "" && id == "" {
// 		return nil, fmt.Errorf("must provide one of `%s` or `%s` field to find by", aliasKey, idKey)
// 	}
// 	var resource *opslevel.Team
// 	if id == "" {
// 		found, err := client.GetTeamWithAlias(alias)
// 		if err != nil {
// 			return nil, err
// 		}
// 		resource = found
// 	} else {
// 		found, err := client.GetTeam(*opslevel.NewID(id.(string)))
// 		if err != nil {
// 			return nil, err
// 		}
// 		resource = found
// 	}
// 	if resource.Id == "" {
// 		return nil, fmt.Errorf("unable to find service with alias=`%s` or id=`%s`", alias, id.(string))
// 	}
// 	return resource, nil
// }

// func getPredicateInputSchema(required bool, description string) *schema.Schema {
// 	output := &schema.Schema{
// 		Type:        schema.TypeList,
// 		MaxItems:    1,
// 		Description: "A condition that should be satisfied.",
// 		ForceNew:    false,
// 		Optional:    true,
// 		Elem: &schema.Resource{
// 			Schema: map[string]*schema.Schema{
// 				"type": {
// 					Type:         schema.TypeString,
// 					Description:  description,
// 					ForceNew:     false,
// 					Required:     true,
// 					ValidateFunc: validation.StringInSlice(opslevel.AllPredicateTypeEnum, false),
// 				},
// 				"value": {
// 					Type:        schema.TypeString,
// 					Description: "The condition value used by the predicate.",
// 					ForceNew:    false,
// 					Optional:    true,
// 				},
// 			},
// 		},
// 	}

// 	if required {
// 		output.Optional = false
// 		output.Required = true
// 	}
// 	return output
// }

// func expandPredicate(d *schema.ResourceData, key string) *opslevel.PredicateInput {
// 	if _, ok := d.GetOk(key); !ok {
// 		return nil
// 	}
// 	return &opslevel.PredicateInput{
// 		Type:  opslevel.PredicateTypeEnum(d.Get(fmt.Sprintf("%s.0.type", key)).(string)),
// 		Value: opslevel.RefOf(d.Get(fmt.Sprintf("%s.0.value", key)).(string)),
// 	}
// }

// func expandPredicateUpdate(d *schema.ResourceData, key string) *opslevel.PredicateUpdateInput {
// 	if _, ok := d.GetOk(key); !ok {
// 		return nil
// 	}
// 	return &opslevel.PredicateUpdateInput{
// 		Type:  opslevel.RefOf(opslevel.PredicateTypeEnum(d.Get(fmt.Sprintf("%s.0.type", key)).(string))),
// 		Value: opslevel.RefOf(d.Get(fmt.Sprintf("%s.0.value", key)).(string)),
// 	}
// }

// func flattenPredicate(input *opslevel.Predicate) []map[string]string {
// 	output := []map[string]string{}
// 	if input != nil {
// 		output = append(output, map[string]string{
// 			"type":  string(input.Type),
// 			"value": input.Value,
// 		})
// 	}
// 	return output
// }

// func expandFilterPredicateInputs(d interface{}) *[]opslevel.FilterPredicateInput {
// 	data := d.([]map[string]interface{})
// 	output := make([]opslevel.FilterPredicateInput, len(data))
// 	for i, item := range data {
// 		var predicate opslevel.FilterPredicateInput
// 		err := mapstructure.Decode(item, &predicate)
// 		if err != nil {
// 			log.Panic().Str("func", "expandFilterPredicateInputs").
// 				Str("item", fmt.Sprintf("%#v", item)).Err(err).
// 				Msg("mapstructure decoding error - please add a bug report https://github.com/OpsLevel/terraform-provider-opslevel/issues/new")
// 		}
// 		// special cases
// 		if item["key_data"] != nil {
// 			predicate.KeyData = opslevel.RefTo(item["key_data"].(string))
// 		} else {
// 			predicate.KeyData = nil
// 		}
// 		// all 4 cases of case_sensitive, case_insensitive need to be handled.
// 		// TODO: bug persists where we cannot unset predicate case_sensitive
// 		// value once it is set because opslevel-go cannot send null. Also
// 		// affects Predicate.key_data and Filter.connective.
// 		x := item["case_sensitive"] == true
// 		y := item["case_insensitive"] == true
// 		if x && y {
// 			// not possible because of input validation
// 		} else if x && !y {
// 			predicate.CaseSensitive = opslevel.RefTo(true)
// 		} else if !x && y {
// 			predicate.CaseSensitive = opslevel.RefTo(false)
// 		} else if !x && !y {
// 			predicate.CaseSensitive = nil
// 		}
// 		output[i] = predicate
// 	}
// 	return &output
// }

// func flattenFilterPredicates(input []opslevel.FilterPredicate) []map[string]any {
// 	output := make([]map[string]any, 0, len(input))
// 	for _, predicate := range input {
// 		o := map[string]any{
// 			"key":      string(predicate.Key),
// 			"key_data": predicate.KeyData,
// 			"type":     string(predicate.Type),
// 			"value":    predicate.Value,
// 		}
// 		// current terraform provider version cannot differentiate between nil and zero
// 		// this is the reverse of the 4 cases in the expand function
// 		if predicate.CaseSensitive == nil {
// 			o["case_sensitive"] = false
// 			o["case_insensitive"] = false
// 		} else if *predicate.CaseSensitive {
// 			o["case_sensitive"] = true
// 			o["case_insensitive"] = false
// 		} else {
// 			o["case_sensitive"] = false
// 			o["case_insensitive"] = true
// 		}
// 		output = append(output, o)
// 	}
// 	return output
// }

// filterBlockModel models data for a terraform block - used to filter resources
type filterBlockModel struct {
	Field types.String `tfsdk:"field"`
	Value types.String `tfsdk:"value"`
}

func NewFilterBlockModel(field string, value string) filterBlockModel {
	return filterBlockModel{
		Field: types.StringValue(string(field)),
		Value: types.StringValue(string(value)),
	}
}

func FilterAttrs(validFieldNames []string) map[string]schema.Attribute {
	filterAttrs := map[string]schema.Attribute{
		"field": schema.StringAttribute{
			Description: fmt.Sprintf(
				"The field of the target resource to filter upon. One of `%s`",
				strings.Join(validFieldNames, "`, `"),
			),
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(validFieldNames...),
			},
		},
		"value": schema.StringAttribute{
			Description: "The field value of the target resource to match.",
			Required:    true,
		},
	}
	return filterAttrs
}

// getDatasourceFilter originally had a "required" bool input parameter - no longer needed
func getDatasourceFilter(validFieldNames []string) schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		MarkdownDescription: "The filter of the target resource to filter upon.",
		Attributes:          FilterAttrs(validFieldNames),
	}
}

func flattenTag(tag opslevel.Tag) string {
	return fmt.Sprintf("%s:%s", tag.Key, tag.Value)
}

func flattenTagArray(tags []opslevel.Tag) []string {
	output := []string{}
	for _, tag := range tags {
		output = append(output, flattenTag(tag))
	}
	return output
}

func flattenServiceRepositoriesArray(repositories *opslevel.ServiceRepositoryConnection) []string {
	output := make([]string, len(repositories.Edges))
	for _, rep := range repositories.Edges {
		output = append(output, string(rep.Node.Id))
	}
	return output
}

// func flattenMembersArray(members *opslevel.UserConnection) []string {
// 	output := []string{}
// 	for _, member := range members.Nodes {
// 		output = append(output, member.Email)
// 	}
// 	return output
// }

// func mapMembershipsArray(members *opslevel.TeamMembershipConnection) []map[string]string {
// 	output := []map[string]string{}
// 	for _, membership := range members.Nodes {
// 		asMap := make(map[string]string)
// 		asMap["email"] = membership.User.Email
// 		asMap["role"] = membership.Role
// 		output = append(output, asMap)
// 	}
// 	return output
// }

// func mapServiceProperties(properties *opslevel.ServicePropertiesConnection) []map[string]any {
// 	output := []map[string]any{}
// 	for _, property := range properties.Nodes {
// 		asMap := make(map[string]any)
// 		asMap["definition"] = string(property.Definition.Id)
// 		asMap["owner"] = string(property.Owner.Id())
// 		if property.Value == nil {
// 			asMap["value"] = "null"
// 		} else {
// 			asMap["value"] = string(*property.Value)
// 		}
// 		output = append(output, asMap)
// 	}
// 	return output
// }

func flattenTeamsArray(teams *opslevel.TeamConnection) []string {
	output := []string{}
	for _, team := range teams.Nodes {
		output = append(output, team.Alias)
	}
	return output
}

// type (
// 	reconcileStringArrayAdd    func(v string) error
// 	reconcileStringArrayUpdate func(o string, n string) error
// 	reconcileStringArrayDelete func(v string) error
// )

// func reconcileStringArray(current []string, desired []string, add reconcileStringArrayAdd, update reconcileStringArrayUpdate, delete reconcileStringArrayDelete) error {
// 	errors := make([]string, 0)
// 	i_current := 0
// 	len_current := len(current)
// 	i_desired := 0
// 	len_desired := len(desired)
// 	sort.Strings(current)
// 	sort.Strings(desired)
// 	// fmt.Printf("Lengths: %v | %v\n", len_current, len_desired)
// 	if len_desired == 0 {
// 		// Delete All in current
// 		if delete == nil {
// 			return nil
// 		}
// 		for _, v := range current {
// 			if err := delete(v); err != nil {
// 				errors = append(errors, err.Error())
// 			}
// 		}
// 		return nil
// 	}
// 	if len_current == 0 {
// 		// Add All from desired
// 		if add == nil {
// 			return nil
// 		}
// 		for _, v := range desired {
// 			if err := add(v); err != nil {
// 				errors = append(errors, err.Error())
// 			}
// 		}

// 	} else {
// 		for i_current < len_current || i_desired < len_desired {
// 			// fmt.Printf("Step: %v | %v\n", i_current, i_desired)
// 			if i_desired >= len_desired {
// 				if delete != nil {
// 					if err := delete(current[i_current]); err != nil {
// 						errors = append(errors, err.Error())
// 					}
// 				}
// 				i_current++
// 				continue
// 			}

// 			if i_current >= len_current {
// 				if add != nil {
// 					if err := add(desired[i_desired]); err != nil {
// 						errors = append(errors, err.Error())
// 					}
// 				}
// 				i_desired++
// 				continue
// 			}
// 			a := current[i_current]
// 			b := desired[i_desired]
// 			if a == b {
// 				if update != nil {
// 					if err := update(a, b); err != nil {
// 						errors = append(errors, err.Error())
// 					}
// 				}
// 				i_current++
// 				i_desired++
// 				continue
// 			}
// 			if a > b {
// 				if add != nil {
// 					if err := add(b); err != nil {
// 						errors = append(errors, err.Error())
// 					}
// 				}
// 				i_desired++
// 				continue
// 			}
// 			if a < b {
// 				if delete != nil {
// 					if err := delete(a); err != nil {
// 						errors = append(errors, err.Error())
// 					}
// 				}
// 				i_current++
// 				continue
// 			}
// 		}
// 	}

// 	if len(errors) > 0 {
// 		return fmt.Errorf(strings.Join(errors, "\n"))
// 	}
// 	return nil
// }
