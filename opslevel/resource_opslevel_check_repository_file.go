package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2022"
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
			"file_contents_predicate": getPredicateInputSchema(false, ""),
			"use_absolute_root": {
				Type:        schema.TypeBool,
				Description: "Whether the checks looks at the absolute root of a repo or the relative root (the directory specified when attached a repo to a service).",
				ForceNew:    false,
				Optional:    true,
			},
		}),
	}
}

func resourceCheckRepositoryFileCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckRepositoryFileCreateInput{}
	setCheckCreateInput(d, &input)

	input.DirectorySearch = d.Get("directory_search").(bool)
	input.Filepaths = getStringArray(d, "filepaths")
	input.FileContentsPredicate = expandPredicate(d, "file_contents_predicate")
	input.UseAbsoluteRoot = opslevel.Bool(d.Get("use_absolute_root").(bool))

	resource, err := client.CreateCheckRepositoryFile(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckRepositoryFileRead(d, client)
}

func resourceCheckRepositoryFileRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(id)
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
	input := opslevel.CheckRepositoryFileUpdateInput{}
	setCheckUpdateInput(d, &input)

	if d.HasChange("directory_search") {
		input.DirectorySearch = opslevel.Bool(d.Get("directory_search").(bool))
	}

	if d.HasChange("filepaths") {
		input.Filepaths = getStringArray(d, "filepaths")
	}

	if d.HasChange("file_contents_predicate") {
		input.FileContentsPredicate = expandPredicate(d, "file_contents_predicate")
	}

	if d.HasChange("use_absolute_root") {
		input.UseAbsoluteRoot = opslevel.Bool(d.Get("use_absolute_root").(bool))
	}

	_, err := client.UpdateCheckRepositoryFile(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckRepositoryFileRead(d, client)
}
