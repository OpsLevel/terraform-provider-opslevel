package opslevel

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
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
		// all 4 cases of case_sensitive, case_insensitive need to be handled.
		// TODO: bug persists where we cannot unset predicate case_sensitive
		// value once it is set because opslevel-go cannot send null. Also
		// affects Predicate.key_data and Filter.connective.
		x := item["case_sensitive"] == true
		y := item["case_insensitive"] == true
		if x && y {
			// not possible because of input validation
		} else if x && !y {
			predicate.CaseSensitive = opslevel.RefTo(true)
		} else if !x && y {
			predicate.CaseSensitive = opslevel.RefTo(false)
		} else if !x && !y {
			predicate.CaseSensitive = nil
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
		// current terraform provider version cannot differentiate between nil and zero
		// this is the reverse of the 4 cases in the expand function
		if predicate.CaseSensitive == nil {
			o["case_sensitive"] = false
			o["case_insensitive"] = false
		} else if *predicate.CaseSensitive == true {
			o["case_sensitive"] = true
			o["case_insensitive"] = false
		} else if *predicate.CaseSensitive == false {
			o["case_sensitive"] = false
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
