package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go"
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
			"file_contents_predicate": getPredicateInputSchema(true),
		}),
	}
}

func resourceCheckRepositorySearchCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckRepositorySearchCreateInput{
		Name:     d.Get("name").(string),
		Enabled:  d.Get("enabled").(bool),
		Category: getID(d, "category"),
		Level:    getID(d, "level"),
		Owner:    getID(d, "owner"),
		Filter:   getID(d, "filter"),
		Notes:    d.Get("notes").(string),

		FileExtensions:        getStringArray(d, "file_extensions"),
		FileContentsPredicate: *getPredicateInput(d, "file_contents_predicate"),
	}
	resource, err := client.CreateCheckRepositorySearch(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceCheckRepositorySearchRead(d, client)
}

func resourceCheckRepositorySearchRead(d *schema.ResourceData, client *opslevel.Client) error {
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

func resourceCheckRepositorySearchUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckRepositorySearchUpdateInput{
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

	if d.HasChange("filepaths") {
		input.FileExtensions = getStringArray(d, "file_extensions")
	}

	if d.HasChange("file_contents_predicate") {
		input.FileContentsPredicate = getPredicateInput(d, "file_contents_predicate")
	}

	_, err := client.UpdateCheckRepositorySearch(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckRepositorySearchRead(d, client)
}
