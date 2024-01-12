package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceCheckRepositoryFile() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a repository file check",
		Create:      wrap(resourceCheckRepositoryFileCreate),
		Read:        wrap(resourceCheckRepositoryFileRead),
		Update:      wrap(resourceCheckRepositoryFileUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{
			"directory_search": {
				Type:        schema.TypeBool,
				Description: "Whether the check looks for the existence of a directory instead of a file.",
				ForceNew:    false,
				Required:    true,
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
			"file_contents_predicate": getPredicateInputSchema(false, DefaultPredicateDescription),
			"use_absolute_root": {
				Type:        schema.TypeBool,
				Description: "Whether the checks looks at the absolute root of a repo or the relative root (the directory specified when attached a repo to a service).",
				ForceNew:    false,
				Required:    true,
			},
		}),
	}
}

func resourceCheckRepositoryFileCreate(d *schema.ResourceData, client *opslevel.Client) error {
	checkCreateInput := getCheckCreateInputFrom(d)
	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckRepositoryFileCreateInput](checkCreateInput)

	input.DirectorySearch = opslevel.RefOf(d.Get("directory_search").(bool))
	input.FilePaths = getStringArray(d, "filepaths")
	input.FileContentsPredicate = expandPredicate(d, "file_contents_predicate")
	input.UseAbsoluteRoot = opslevel.RefOf(d.Get("use_absolute_root").(bool))

	resource, err := client.CreateCheckRepositoryFile(*input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckRepositoryFileRead(d, client)
}

func resourceCheckRepositoryFileRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}
	if err := d.Set("directory_search", resource.RepositoryFileCheckFragment.DirectorySearch); err != nil {
		return err
	}
	if err := d.Set("filepaths", resource.RepositoryFileCheckFragment.Filepaths); err != nil {
		return err
	}
	if err := d.Set("file_contents_predicate", flattenPredicate(resource.RepositoryFileCheckFragment.FileContentsPredicate)); err != nil {
		return err
	}
	if err := d.Set("use_absolute_root", resource.UseAbsoluteRoot); err != nil {
		return err
	}
	return nil
}

func resourceCheckRepositoryFileUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	checkUpdateInput := getCheckUpdateInputFrom(d)
	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckRepositoryFileUpdateInput](checkUpdateInput)
	input.DirectorySearch = opslevel.RefOf(d.Get("directory_search").(bool))
	input.UseAbsoluteRoot = opslevel.RefOf(d.Get("use_absolute_root").(bool))

	if d.HasChange("filepaths") {
		input.FilePaths = opslevel.RefOf(getStringArray(d, "filepaths"))
	}

	if d.HasChange("file_contents_predicate") {
		input.FileContentsPredicate = expandPredicateUpdate(d, "file_contents_predicate")
	}

	_, err := client.UpdateCheckRepositoryFile(*input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckRepositoryFileRead(d, client)
}
