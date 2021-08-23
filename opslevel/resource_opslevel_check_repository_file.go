package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
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
			"file_contents_predicate": getPredicateInputSchema(false),
		}),
	}
}

func resourceCheckRepositoryFileCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckRepositoryFileCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),

		DirectorySearch:       d.Get("directory_search").(bool),
		Filepaths:             getStringArray(d, "filepaths"),
		FileContentsPredicate: getPredicateInput(d, "file_contents_predicate"),
	}
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

	if err := resourceCheckRead(d, resource); err != nil {
		return err
	}

	return nil
}

func resourceCheckRepositoryFileUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckRepositoryFileUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("enabled") {
		value := d.Get("enabled").(bool)
		input.Enabled = &value
	}
	if d.HasChange("category") {
		input.Category = getID(d, "category")
	}
	if d.HasChange("level") {
		input.Level = getID(d, "level")
	}
	if d.HasChange("owner") {
		input.Owner = getID(d, "owner")
	}
	if d.HasChange("filter") {
		input.Filter = getID(d, "filter")
	}
	if d.HasChange("notes") {
		input.Notes = d.Get("notes").(string)
	}

	if d.HasChange("directory_search") {
		input.DirectorySearch = d.Get("directory_search").(bool)
	}

	if d.HasChange("filepaths") {
		input.Filepaths = getStringArray(d, "filepaths")
	}

	if d.HasChange("file_contents_predicate") {
		input.FileContentsPredicate = getPredicateInput(d, "file_contents_predicate")
	}

	_, err := client.UpdateCheckRepositoryFile(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckRepositoryFileRead(d, client)
}
