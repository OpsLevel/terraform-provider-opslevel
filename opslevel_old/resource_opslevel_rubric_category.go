// package opslevel

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceRubricCategory() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a rubric category",
// 		Create:      wrap(resourceRubricCategoryCreate),
// 		Read:        wrap(resourceRubricCategoryRead),
// 		Update:      wrap(resourceRubricCategoryUpdate),
// 		Delete:      wrap(resourceRubricCategoryDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "The display name of the category.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 		},
// 	}
// }

// func resourceRubricCategoryCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.CategoryCreateInput{
// 		Name: d.Get("name").(string),
// 	}
// 	resource, err := client.CreateCategory(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceRubricCategoryRead(d, client)
// }

// func resourceRubricCategoryRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetCategory(opslevel.ID(id))
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("name", resource.Name); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceRubricCategoryUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.CategoryUpdateInput{
// 		Id: opslevel.ID(d.Id()),
// 	}

// 	if d.HasChange("name") {
// 		input.Name = opslevel.RefOf(d.Get("name").(string))
// 	}

// 	_, err := client.UpdateCategory(input)
// 	if err != nil {
// 		return err
// 	}
// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceRubricCategoryRead(d, client)
// }

// func resourceRubricCategoryDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteCategory(opslevel.ID(id))
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
