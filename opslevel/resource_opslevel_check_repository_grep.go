package opslevel

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceCheckRepositoryGrep() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a repository file check",
// 		Create:      wrap(resourceCheckRepositoryGrepCreate),
// 		Read:        wrap(resourceCheckRepositoryGrepRead),
// 		Update:      wrap(resourceCheckRepositoryGrepUpdate),
// 		Delete:      wrap(resourceCheckDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: getCheckSchema(map[string]*schema.Schema{
// 			"directory_search": {
// 				Type:        schema.TypeBool,
// 				Description: "Whether the check looks for the existence of a directory instead of a file.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"filepaths": {
// 				Type:        schema.TypeList,
// 				MinItems:    1,
// 				Description: "Restrict the search to certain file paths.",
// 				ForceNew:    false,
// 				Required:    true,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 			},
// 			"file_contents_predicate": getPredicateInputSchema(false, "A condition that should be satisfied. Defaults to `exists` condition"),
// 		}),
// 	}
// }

// func resourceCheckRepositoryGrepCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	checkCreateInput := getCheckCreateInputFrom(d)
// 	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckRepositoryGrepCreateInput](checkCreateInput)
// 	input.DirectorySearch = opslevel.RefOf(d.Get("directory_search").(bool))
// 	input.FilePaths = getStringArray(d, "filepaths")
// 	fileContentsPredicate := expandPredicate(d, "file_contents_predicate")
// 	if fileContentsPredicate == nil {
// 		input.FileContentsPredicate = opslevel.PredicateInput{
// 			Type: opslevel.PredicateTypeEnumExists,
// 		}
// 	}

// 	resource, err := client.CreateCheckRepositoryGrep(*input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceCheckRepositoryGrepRead(d, client)
// }

// func resourceCheckRepositoryGrepRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetCheck(opslevel.ID(id))
// 	if err != nil {
// 		return err
// 	}

// 	if err := setCheckData(d, resource); err != nil {
// 		return err
// 	}
// 	if err := d.Set("directory_search", resource.RepositoryGrepCheckFragment.DirectorySearch); err != nil {
// 		return err
// 	}
// 	if err := d.Set("filepaths", resource.RepositoryGrepCheckFragment.Filepaths); err != nil {
// 		return err
// 	}
// 	if _, ok := d.GetOk("file_contents_predicate"); !ok {
// 		if err := d.Set("file_contents_predicate", nil); err != nil {
// 			return err
// 		}
// 	} else {
// 		if err := d.Set("file_contents_predicate", flattenPredicate(resource.RepositoryGrepCheckFragment.FileContentsPredicate)); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func resourceCheckRepositoryGrepUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	checkUpdateInput := getCheckUpdateInputFrom(d)
// 	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckRepositoryGrepUpdateInput](checkUpdateInput)
// 	input.DirectorySearch = opslevel.RefOf(d.Get("directory_search").(bool))

// 	if d.HasChange("filepaths") {
// 		input.FilePaths = opslevel.RefOf(getStringArray(d, "filepaths"))
// 	}
// 	if d.HasChange("file_contents_predicate") {
// 		input.FileContentsPredicate = expandPredicateUpdate(d, "file_contents_predicate")
// 	}

// 	_, err := client.UpdateCheckRepositoryGrep(*input)
// 	if err != nil {
// 		return err
// 	}
// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceCheckRepositoryGrepRead(d, client)
// }
