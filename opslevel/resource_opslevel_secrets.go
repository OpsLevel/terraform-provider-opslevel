package opslevel

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a secret",
		Create:      wrap(resourceSecretCreate),
		Read:        wrap(resourceSecretRead),
		Update:      wrap(resourceSecretUpdate),
		Delete:      wrap(resourceSecretDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alias": {
				Type:        schema.TypeString,
				Description: "The alias for this secret.",
				ForceNew:    true,
				Required:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The owner of this secret.",
				ForceNew:    false,
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "A sensitive value.",
				Sensitive:   true,
				ForceNew:    false,
				Required:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Timestamp of time created at.",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Timestamp of last update.",
				Computed:    true,
			},
		},
	}
}

func resourceSecretCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.SecretInput{
		Owner: opslevel.NewIdentifier(d.Get("owner").(string)),
		Value: opslevel.RefOf(d.Get("value").(string)),
	}
	alias := d.Get("alias").(string)
	resource, err := client.CreateSecret(alias, input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.ID))
	return resourceSecretRead(d, client)
}

func resourceSecretRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetSecret(id)
	if err != nil {
		return err
	}

	if opslevel.IsID(d.Get("owner").(string)) {
		if err := d.Set("owner", resource.Owner.Id); err != nil {
			return err
		}
	} else {
		if err := d.Set("owner", resource.Owner.Alias); err != nil {
			return err
		}
	}
	created_at := resource.Timestamps.CreatedAt.Local().Format(time.RFC850)
	if err := d.Set("created_at", created_at); err != nil {
		return err
	}
	updated_at := resource.Timestamps.UpdatedAt.Local().Format(time.RFC850)
	if err := d.Set("updated_at", updated_at); err != nil {
		return err
	}

	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.SecretInput{
		Owner: opslevel.NewIdentifier(d.Get("owner").(string)),
		Value: opslevel.RefOf(d.Get("value").(string)),
	}

	_, err := client.UpdateSecret(d.Id(), input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceSecretRead(d, client)
}

func resourceSecretDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteSecret(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
