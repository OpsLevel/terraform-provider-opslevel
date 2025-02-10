package opslevel

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/opslevel/opslevel-go/v2025"
)

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

// return strings not in both slices
func diffBetweenStringSlices(sliceOne, sliceTwo []string) []string {
	var diffValues []string

	// collect values that are in sliceOne but not in sliceTwo
	for _, value := range sliceOne {
		if !slices.Contains(sliceTwo, value) {
			diffValues = append(diffValues, value)
		}
	}

	// collect values that are in sliceTwo but not in sliceOne
	for _, value := range sliceTwo {
		if !slices.Contains(sliceOne, value) {
			diffValues = append(diffValues, value)
		}
	}
	return diffValues
}

// getDatasourceFilter originally had a "required" bool input parameter - no longer needed
func getDatasourceFilter(validFieldNames []string) schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		MarkdownDescription: "The filter of the target resource to filter upon.",
		Attributes:          FilterAttrs(validFieldNames),
	}
}

// Temp wrapper until opslevel-go is updated
func getService(client *opslevel.Client, serviceIdentifier string) (*opslevel.Service, error) {
	var err error
	var service *opslevel.Service

	if opslevel.IsID(serviceIdentifier) {
		service, err = client.GetService(serviceIdentifier)
	} else {
		service, err = client.GetServiceWithAlias(serviceIdentifier)
	}
	if err == nil && (service == nil || string(service.Id) == "") {
		err = fmt.Errorf("service %s not found", serviceIdentifier)
	}

	return service, err
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

func flattenTeamsArray(teams *opslevel.TeamConnection) []string {
	output := []string{}
	for _, team := range teams.Nodes {
		output = append(output, team.Alias)
	}
	return output
}
