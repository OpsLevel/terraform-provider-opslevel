package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceCheckRepositorySearch() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a repository search check.",
		Create:      wrap(resourceCheckRepositorySearchCreate),
		Read:        wrap(resourceCheckRepositorySearchRead),
		Update:      wrap(resourceCheckRepositorySearchUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"file_extensions": {
				Type:        schema.TypeList,
				MinItems:    1,
				Description: "Restrict the search to files of given extensions. Extensions should contain only letters and numbers. For example: [\"py\", \"rb\"].",
				ForceNew:    false,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"file_contents_predicate": getPredicateInputSchema(true, DefaultPredicateDescription),
		}),
	}
}

func resourceCheckRepositorySearchCreate(d *schema.ResourceData, client *opslevel.Client) error {
	checkCreateInput := getCheckCreateInputFrom(d)
	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckRepositorySearchCreateInput](checkCreateInput)
	input.FileExtensions = opslevel.RefOf(getStringArray(d, "file_extensions"))
	input.FileContentsPredicate = *expandPredicate(d, "file_contents_predicate")

	resource, err := client.CreateCheckRepositorySearch(*input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckRepositorySearchRead(d, client)
}

func resourceCheckRepositorySearchRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("file_extensions", resource.FileExtensions); err != nil {
		return err
	}
	if err := d.Set("file_contents_predicate", flattenPredicate(&resource.RepositorySearchCheckFragment.FileContentsPredicate)); err != nil {
		return err
	}

	return nil
}

func resourceCheckRepositorySearchUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	checkUpdateInput := getCheckUpdateInputFrom(d)
	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckRepositorySearchUpdateInput](checkUpdateInput)

	if d.HasChange("file_extensions") {
		input.FileExtensions = opslevel.RefOf(getStringArray(d, "file_extensions"))
	}

	if d.HasChange("file_contents_predicate") {
		input.FileContentsPredicate = expandPredicateUpdate(d, "file_contents_predicate")
	}

	_, err := client.UpdateCheckRepositorySearch(*input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckRepositorySearchRead(d, client)
}
