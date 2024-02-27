package opslevel

// import (
// 	"github.com/opslevel/opslevel-go/v2024"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// )

// func datasourceRepository() *schema.Resource {
// 	return &schema.Resource{
// 		Read: wrap(datasourceRepositoryRead),
// 		Schema: map[string]*schema.Schema{
// 			"alias": {
// 				Type:        schema.TypeString,
// 				Description: "A human-friendly, unique identifier for the repository.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 			"id": {
// 				Type:        schema.TypeString,
// 				Description: "The id of the repository to find.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 		},
// 	}
// }

// func datasourceRepositoryRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	resource, err := findRepository("alias", "id", d, client)
// 	if err != nil {
// 		return err
// 	}

// 	d.SetId(string(resource.Id))

// 	return nil
// }
