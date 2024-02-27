package opslevel

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceRepository() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages the ownership of a repository but does not create or delete the repository entity in OpsLevel.",
// 		Create:      wrap(resourceRepositoryCreate),
// 		Read:        wrap(resourceRepositoryRead),
// 		Update:      wrap(resourceRepositoryUpdate),
// 		Delete:      wrap(resourceRepositoryDelete),
// 		Importer: &schema.ResourceImporter{
// 			State: schema.ImportStatePassthrough,
// 		},
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"identifier": {
// 				Type:        schema.TypeString,
// 				Description: "The id or human-friendly, unique identifier for the repository.",
// 				ForceNew:    true,
// 				Optional:    true,
// 			},
// 			"owner": {
// 				Type:        schema.TypeString,
// 				Description: "The owner of the repository.",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 		},
// 	}
// }

// func resourceRepositoryCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	identifier := d.Get("identifier").(string)
// 	var repository *opslevel.Repository
// 	if opslevel.IsID(identifier) {
// 		resource, err := client.GetRepository(*opslevel.NewID(identifier))
// 		if err != nil {
// 			return err
// 		}
// 		repository = resource
// 	} else {
// 		resource, err := client.GetRepositoryWithAlias(identifier)
// 		if err != nil {
// 			return err
// 		}
// 		repository = resource
// 	}

// 	d.SetId(string(repository.Id))

// 	return resourceRepositoryUpdate(d, client)
// }

// func resourceRepositoryRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	repository, err := client.GetRepository(*opslevel.NewID(d.Id()))
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("owner", repository.Owner.Id); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceRepositoryUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	input := opslevel.RepositoryUpdateInput{
// 		Id: *opslevel.NewID(d.Id()),
// 	}

// 	if owner, ok := d.GetOk("owner"); ok {
// 		input.OwnerId = opslevel.NewID(owner.(string))
// 	} else {
// 		input.OwnerId = opslevel.NewID("")
// 	}

// 	_, err := client.UpdateRepository(input)
// 	if err != nil {
// 		return err
// 	}

// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceRepositoryRead(d, client)
// }

// func resourceRepositoryDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	// No API call to make because the repository is not able to be deleted
// 	d.SetId("")
// 	return nil
// }
