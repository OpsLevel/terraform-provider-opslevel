package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2022"
)

func resourceCheckRepositoryGrep() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a repository file check",
		Create:      wrap(resourceCheckRepositoryGrepCreate),
		Read:        wrap(resourceCheckRepositoryGrepRead),
		Update:      wrap(resourceCheckRepositoryGrepUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"directory_search": {
				Type:        schema.TypeBool,
				Description: "Whether the check looks for the existence of a directory instead of a file. Defaults to false",
				ForceNew:    false,
				Optional:    true,
			},
			"filepaths": {
				Type:        schema.TypeList,
				MinItems:    1,
				Description: "Restrict the search to certain file paths.",
				ForceNew:    false,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"file_contents_predicate": getPredicateInputSchema(false),
		}),
	}
}

func resourceCheckRepositoryGrepCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckRepositoryGrepCreateInput{}
	setCheckCreateInput(d, &input)

	input.DirectorySearch = d.Get("directory_search").(bool)
	input.Filepaths = getStringArray(d, "filepaths")
	input.FileContentsPredicate = expandPredicate(d, "file_contents_predicate")

	resource, err := client.CreateCheckRepositoryGrep(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckRepositoryGrepRead(d, client)
}

func resourceCheckRepositoryGrepRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(id)
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("directory_search", resource.RepositoryGrepCheckFragment.DirectorySearch); err != nil {
		return err
	}
	if err := d.Set("filepaths", resource.RepositoryGrepCheckFragment.Filepaths); err != nil {
		return err
	}
	if err := d.Set("file_contents_predicate", flattenPredicate(resource.RepositoryGrepCheckFragment.FileContentsPredicate)); err != nil {
		return err
	}
	return nil
}

func resourceCheckRepositoryGrepUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckRepositoryGrepUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("directory_search") {
		input.DirectorySearch = d.Get("directory_search").(bool)
	}

	if d.HasChange("filepaths") {
		input.Filepaths = getStringArray(d, "filepaths")
	}

	if d.HasChange("file_contents_predicate") {
		input.FileContentsPredicate = expandPredicate(d, "file_contents_predicate")
	}

	_, err := client.UpdateCheckRepositoryGrep(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckRepositoryGrepRead(d, client)
}
