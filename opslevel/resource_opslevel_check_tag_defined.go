package opslevel

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceCheckTagDefined() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a tag defined check",
// 		Create:      wrap(resourceCheckTagDefinedCreate),
// 		Read:        wrap(resourceCheckTagDefinedRead),
// 		Update:      wrap(resourceCheckTagDefinedUpdate),
// 		Delete:      wrap(resourceCheckDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: getCheckSchema(map[string]*schema.Schema{
// 			"tag_key": {
// 				Type:        schema.TypeString,
// 				Description: "The tag key where the tag predicate should be applied.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"tag_predicate": getPredicateInputSchema(false, DefaultPredicateDescription),
// 		}),
// 	}
// }

// func resourceCheckTagDefinedCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	checkCreateInput := getCheckCreateInputFrom(d)
// 	input := opslevel.NewCheckCreateInputTypeOf[opslevel.CheckTagDefinedCreateInput](checkCreateInput)
// 	input.TagKey = d.Get("tag_key").(string)
// 	input.TagPredicate = expandPredicate(d, "tag_predicate")

// 	resource, err := client.CreateCheckTagDefined(*input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceCheckTagDefinedRead(d, client)
// }

// func resourceCheckTagDefinedRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetCheck(opslevel.ID(id))
// 	if err != nil {
// 		return err
// 	}

// 	if err := setCheckData(d, resource); err != nil {
// 		return err
// 	}
// 	if err := d.Set("tag_key", resource.TagKey); err != nil {
// 		return err
// 	}
// 	if err := d.Set("tag_predicate", flattenPredicate(resource.TagPredicate)); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceCheckTagDefinedUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	checkUpdateInput := getCheckUpdateInputFrom(d)
// 	input := opslevel.NewCheckUpdateInputTypeOf[opslevel.CheckTagDefinedUpdateInput](checkUpdateInput)

// 	if d.HasChange("tag_key") {
// 		input.TagKey = opslevel.RefOf(d.Get("tag_key").(string))
// 	}
// 	if d.HasChange("tag_predicate") {
// 		input.TagPredicate = expandPredicateUpdate(d, "tag_predicate")
// 	}

// 	_, err := client.UpdateCheckTagDefined(*input)
// 	if err != nil {
// 		return err
// 	}
// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceCheckTagDefinedRead(d, client)
// }
