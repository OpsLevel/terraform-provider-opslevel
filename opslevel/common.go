package opslevel

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2024"
)

func cleanerString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// interfacesMaps converts an interface{} into a []map[string]interface{}. This is a useful conversion for passing
// schema.ResourceData objects from terraform into mapstructure.Decode to get actual struct types.
func interfacesMaps(i interface{}) []map[string]interface{} {
	// interface{} 					to 		[]interface{}								segment into slices.
	interfaces := i.([]interface{})
	// interface{}					to		[]map[string]interface{}					convert each slice item into a map.
	mapStringInterfaces := make([]map[string]interface{}, len(interfaces))
	for i, item := range interfaces {
		mapStringInterfaces[i] = item.(map[string]interface{})
	}
	return mapStringInterfaces
}

var DefaultPredicateDescription = "A condition that should be satisfied."

func timeID() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func timeLastUpdated() string {
	return time.Now().Format(time.RFC850)
}

func wrap(handler func(data *schema.ResourceData, client *opslevel.Client) error) func(d *schema.ResourceData, meta interface{}) error {
	return func(data *schema.ResourceData, meta interface{}) error {
		client := meta.(*opslevel.Client)
		return handler(data, client)
	}
}

func stringInArray(term string, search []string) bool {
	for _, item := range search {
		if term == item {
			return true
		}
	}
	return false
}

func expandStringArray(m []interface{}) []string {
	result := make([]string, 0)
	for _, v := range m {
		result = append(result, v.(string))
	}
	return result
}

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func getStringArray(d *schema.ResourceData, key string) []string {
	output := make([]string, 0)
	data, ok := d.GetOk(key)
	if !ok {
		return output
	}
	for _, item := range data.([]interface{}) {
		output = append(output, item.(string))
	}
	return output
}

func findService(aliasKey string, idKey string, d *schema.ResourceData, client *opslevel.Client) (*opslevel.Service, error) {
	alias := d.Get(aliasKey).(string)
	id := d.Get(idKey)
	if alias == "" && id == "" {
		return nil, fmt.Errorf("must provide one of `%s` or `%s` field to find by", aliasKey, idKey)
	}
	var resource *opslevel.Service
	if id == "" {
		found, err := client.GetServiceWithAlias(alias)
		if err != nil {
			return nil, err
		}
		resource = found
	} else {
		found, err := client.GetService(*opslevel.NewID(id.(string)))
		if err != nil {
			return nil, err
		}
		resource = found
	}
	if resource.Id == "" {
		return nil, fmt.Errorf("unable to find service with alias=`%s` or id=`%s`", alias, id.(string))
	}
	return resource, nil
}

func findRepository(aliasKey string, idKey string, d *schema.ResourceData, client *opslevel.Client) (*opslevel.Repository, error) {
	alias := d.Get(aliasKey).(string)
	id := d.Get(idKey)
	if alias == "" && id == "" {
		return nil, fmt.Errorf("must provide one of `%s` or `%s` field to find by", aliasKey, idKey)
	}
	var resource *opslevel.Repository
	if id == "" {
		found, err := client.GetRepositoryWithAlias(alias)
		if err != nil {
			return nil, err
		}
		resource = found
	} else {
		found, err := client.GetRepository(*opslevel.NewID(id.(string)))
		if err != nil {
			return nil, err
		}
		resource = found
	}
	if resource.Id == "" {
		return nil, fmt.Errorf("unable to find repository with alias=`%s` or id=`%s`", alias, id.(string))
	}
	return resource, nil
}

func findTeam(aliasKey string, idKey string, d *schema.ResourceData, client *opslevel.Client) (*opslevel.Team, error) {
	alias := d.Get(aliasKey).(string)
	id := d.Get(idKey)
	if alias == "" && id == "" {
		return nil, fmt.Errorf("must provide one of `%s` or `%s` field to find by", aliasKey, idKey)
	}
	var resource *opslevel.Team
	if id == "" {
		found, err := client.GetTeamWithAlias(alias)
		if err != nil {
			return nil, err
		}
		resource = found
	} else {
		found, err := client.GetTeam(*opslevel.NewID(id.(string)))
		if err != nil {
			return nil, err
		}
		resource = found
	}
	if resource.Id == "" {
		return nil, fmt.Errorf("unable to find service with alias=`%s` or id=`%s`", alias, id.(string))
	}
	return resource, nil
}

func getPredicateInputSchema(required bool, description string) *schema.Schema {
	output := &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Description: "A condition that should be satisfied.",
		ForceNew:    false,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Description:  description,
					ForceNew:     false,
					Required:     true,
					ValidateFunc: validation.StringInSlice(opslevel.AllPredicateTypeEnum, false),
				},
				"value": {
					Type:        schema.TypeString,
					Description: "The condition value used by the predicate.",
					ForceNew:    false,
					Optional:    true,
				},
			},
		},
	}

	if required {
		output.Optional = false
		output.Required = true
	}
	return output
}

func expandPredicate(d *schema.ResourceData, key string) *opslevel.PredicateInput {
	if _, ok := d.GetOk(key); !ok {
		return nil
	}
	return &opslevel.PredicateInput{
		Type:  opslevel.PredicateTypeEnum(d.Get(fmt.Sprintf("%s.0.type", key)).(string)),
		Value: opslevel.RefOf(d.Get(fmt.Sprintf("%s.0.value", key)).(string)),
	}
}

func expandPredicateUpdate(d *schema.ResourceData, key string) *opslevel.PredicateUpdateInput {
	if _, ok := d.GetOk(key); !ok {
		return nil
	}
	return &opslevel.PredicateUpdateInput{
		Type:  opslevel.RefOf(opslevel.PredicateTypeEnum(d.Get(fmt.Sprintf("%s.0.type", key)).(string))),
		Value: opslevel.RefOf(d.Get(fmt.Sprintf("%s.0.value", key)).(string)),
	}
}

func flattenPredicate(input *opslevel.Predicate) []map[string]string {
	output := []map[string]string{}
	if input != nil {
		output = append(output, map[string]string{
			"type":  string(input.Type),
			"value": input.Value,
		})
	}
	return output
}

func expandFilterPredicateInputs(d interface{}) *[]opslevel.FilterPredicateInput {
	data := d.([]map[string]interface{})
	output := make([]opslevel.FilterPredicateInput, len(data))
	for i, item := range data {
		var predicate opslevel.FilterPredicateInput
		err := mapstructure.Decode(item, &predicate)
		if err != nil {
			log.Panic().Str("func", "expandFilterPredicateInputs").
				Str("item", fmt.Sprintf("%#v", item)).Err(err).
				Msg("mapstructure decoding error - please add a bug report https://github.com/OpsLevel/terraform-provider-opslevel/issues/new")
		}
		// special cases
		if item["key_data"] != nil {
			predicate.KeyData = opslevel.RefTo(item["key_data"].(string))
		} else {
			predicate.KeyData = nil
		}
		if item["case_sensitive"] == "true" && item["case_insensitive"] == "true" {
			log.Panic().Str("func", "expandFilterPredicateInputs").
				Str("item", fmt.Sprintf("%#v", item)).Err(err).Msg("can't set both case_insensitive and case_sensitive at the same time.")
		} else if item["case_sensitive"] == "true" {
			predicate.CaseSensitive = opslevel.RefTo(true)
		} else if item["case_insensitive"] == "true" {
			predicate.CaseSensitive = opslevel.RefTo(false)
		}
		output[i] = predicate
	}
	return &output
}

func flattenFilterPredicates(input []opslevel.FilterPredicate) []map[string]any {
	output := make([]map[string]any, 0, len(input))
	for _, predicate := range input {
		o := map[string]any{
			"key":      string(predicate.Key),
			"key_data": predicate.KeyData,
			"type":     string(predicate.Type),
			"value":    predicate.Value,
		}
		// special cases
		if predicate.CaseSensitive == nil {
		} else if *predicate.CaseSensitive == true {
			o["case_sensitive"] = true
		} else if *predicate.CaseSensitive == false {
			o["case_insensitive"] = true
		}
		output = append(output, o)
	}
	return output
}

func getDatasourceFilter(required bool, validFieldNames []string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		ForceNew: true,
		Required: required,
		Optional: !required,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field": {
					Type:         schema.TypeString,
					Description:  "The field of the target resource to filter upon.",
					ForceNew:     true,
					Required:     true,
					ValidateFunc: validation.StringInSlice(validFieldNames, false),
				},
				"value": {
					Type:        schema.TypeString,
					Description: "The field value of the target resource to match.",
					ForceNew:    true,
					Optional:    true,
				},
			},
		},
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
	output := []string{}
	for _, rep := range repositories.Edges {
		output = append(output, string(rep.Node.Id))
	}
	return output
}

func flattenMembersArray(members *opslevel.UserConnection) []string {
	output := []string{}
	for _, member := range members.Nodes {
		output = append(output, member.Email)
	}
	return output
}

func mapMembershipsArray(members *opslevel.TeamMembershipConnection) []map[string]string {
	output := []map[string]string{}
	for _, membership := range members.Nodes {
		asMap := make(map[string]string)
		asMap["email"] = membership.User.Email
		asMap["role"] = membership.Role
		output = append(output, asMap)
	}
	return output
}

func mapServiceProperties(properties *opslevel.ServicePropertiesConnection) []map[string]any {
	output := []map[string]any{}
	for _, property := range properties.Nodes {
		asMap := make(map[string]any)
		asMap["definition"] = string(property.Definition.Id)
		asMap["owner"] = string(property.Owner.Id())
		if property.Value == nil {
			asMap["value"] = "null"
		} else {
			asMap["value"] = string(*property.Value)
		}
		output = append(output, asMap)
	}
	return output
}

func flattenTeamsArray(teams *opslevel.TeamConnection) []string {
	output := []string{}
	for _, team := range teams.Nodes {
		output = append(output, team.Alias)
	}
	return output
}

type (
	reconcileStringArrayAdd    func(v string) error
	reconcileStringArrayUpdate func(o string, n string) error
	reconcileStringArrayDelete func(v string) error
)

func reconcileStringArray(current []string, desired []string, add reconcileStringArrayAdd, update reconcileStringArrayUpdate, delete reconcileStringArrayDelete) error {
	errors := make([]string, 0)
	i_current := 0
	len_current := len(current)
	i_desired := 0
	len_desired := len(desired)
	sort.Strings(current)
	sort.Strings(desired)
	// fmt.Printf("Lengths: %v | %v\n", len_current, len_desired)
	if len_desired == 0 {
		// Delete All in current
		if delete == nil {
			return nil
		}
		for _, v := range current {
			if err := delete(v); err != nil {
				errors = append(errors, err.Error())
			}
		}
		return nil
	}
	if len_current == 0 {
		// Add All from desired
		if add == nil {
			return nil
		}
		for _, v := range desired {
			if err := add(v); err != nil {
				errors = append(errors, err.Error())
			}
		}

	} else {
		for i_current < len_current || i_desired < len_desired {
			// fmt.Printf("Step: %v | %v\n", i_current, i_desired)
			if i_desired >= len_desired {
				if delete != nil {
					if err := delete(current[i_current]); err != nil {
						errors = append(errors, err.Error())
					}
				}
				i_current++
				continue
			}

			if i_current >= len_current {
				if add != nil {
					if err := add(desired[i_desired]); err != nil {
						errors = append(errors, err.Error())
					}
				}
				i_desired++
				continue
			}
			a := current[i_current]
			b := desired[i_desired]
			if a == b {
				if update != nil {
					if err := update(a, b); err != nil {
						errors = append(errors, err.Error())
					}
				}
				i_current++
				i_desired++
				continue
			}
			if a > b {
				if add != nil {
					if err := add(b); err != nil {
						errors = append(errors, err.Error())
					}
				}
				i_desired++
				continue
			}
			if a < b {
				if delete != nil {
					if err := delete(a); err != nil {
						errors = append(errors, err.Error())
					}
				}
				i_current++
				continue
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}
