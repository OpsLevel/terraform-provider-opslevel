// package opslevel

// import (
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceUser() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages a User",
// 		Create:      wrap(resourceUserCreate),
// 		Read:        wrap(resourceUserRead),
// 		Update:      wrap(resourceUserUpdate),
// 		Delete:      wrap(resourceUserDelete),
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
// 				Description: "The name of the user.",
// 				ForceNew:    false,
// 				Required:    true,
// 			},
// 			"email": {
// 				Type:        schema.TypeString,
// 				Description: "The email address of the user.",
// 				ForceNew:    true,
// 				Required:    true,
// 			},
// 			"role": {
// 				Type:         schema.TypeString,
// 				Description:  "The access role (e.g. user vs admin) of the user.",
// 				ForceNew:     false,
// 				Optional:     true,
// 				ValidateFunc: validation.StringInSlice(opslevel.AllUserRole, false),
// 			},
// 			// There is no way to read this value from a User resource and no way
// 			// to set it in the terraform state. Do not remove it because some
// 			// customers rely on this feature regardless.
// 			"skip_welcome_email": {
// 				Type:        schema.TypeBool,
// 				Description: "Don't send an email welcoming the user to OpsLevel. Applies during creation only, this value cannot be read or updated.",
// 				Default:     false,
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 		},
// 	}
// }

// func resourceUserCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	email := d.Get("email").(string)
// 	input := opslevel.UserInput{
// 		Name:             opslevel.RefOf(d.Get("name").(string)),
// 		Role:             opslevel.RefOf(opslevel.UserRole(d.Get("role").(string))),
// 		SkipWelcomeEmail: opslevel.RefOf(d.Get("skip_welcome_email").(bool)),
// 	}
// 	resource, err := client.InviteUser(email, input)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(string(resource.Id))

// 	return resourceUserRead(d, client)
// }

// func resourceUserRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetUser(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("name", resource.Name); err != nil {
// 		return err
// 	}
// 	if err := d.Set("email", resource.Email); err != nil {
// 		return err
// 	}
// 	if err := d.Set("role", resource.Role); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceUserUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	input := opslevel.UserInput{}

// 	if d.HasChange("name") {
// 		input.Name = opslevel.RefOf(d.Get("name").(string))
// 	}
// 	if d.HasChange("role") {
// 		input.Role = opslevel.RefOf(opslevel.UserRole(d.Get("role").(string)))
// 	}

// 	_, err := client.UpdateUser(id, input)
// 	if err != nil {
// 		return err
// 	}

// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceUserRead(d, client)
// }

// func resourceUserDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteUser(id)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
