package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceIntegrationAWS() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a AWS Integration",
		Create:      wrap(resourceIntegrationAWSCreate),
		Read:        wrap(resourceIntegrationAWSRead),
		Update:      wrap(resourceIntegrationAWSUpdate),
		Delete:      wrap(resourceIntegrationAWSDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the integration.",
				ForceNew:    false,
				Required:    true,
			},
			"iam_role": {
				Type:        schema.TypeString,
				Description: "The IAM role OpsLevel uses in order to access the AWS account.",
				ForceNew:    false,
				Required:    true,
			},
			"external_id": {
				Type:        schema.TypeString,
				Description: "The External ID defined in the trust relationship to ensure OpsLevel is the only third party assuming this role (See https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html for more details).",
				ForceNew:    false,
				Required:    true,
			},
			"ownership_tag_overrides": {
				Type:        schema.TypeBool,
				Description: "Allow tags imported from AWS to override ownership set in OpsLevel directly.",
				ForceNew:    false,
				Required:    true,
			},
			"ownership_tag_keys": {
				Type:        schema.TypeList,
				Description: "An Array of tag keys used to associate ownership from an integration. Max 5",
				ForceNew:    false,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, MaxItems: 5},
			},
		},
	}
}

func resourceIntegrationAWSCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.AWSIntegrationInput{
		Name:                 opslevel.NewString(d.Get("name").(string)),
		IAMRole:              opslevel.NewString(d.Get("iam_role").(string)),
		ExternalID:           opslevel.NewString(d.Get("external_id").(string)),
		OwnershipTagOverride: opslevel.Bool(d.Get("ownership_tag_overrides").(bool)),
	}

	input.OwnershipTagKeys = getStringArray(d, "ownership_tag_keys")

	resource, err := client.CreateIntegrationAWS(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceIntegrationAWSRead(d, client)
}

func resourceIntegrationAWSRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetIntegration(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := d.Set("iam_role", resource.IAMRole); err != nil {
		return err
	}
	if err := d.Set("external_id", resource.ExternalID); err != nil {
		return err
	}
	if err := d.Set("ownership_tag_overrides", resource.OwnershipTagOverride); err != nil {
		return err
	}
	if _, ok := d.GetOk("ownership_tag_keys"); ok {
		if err := d.Set("ownership_tag_keys", resource.OwnershipTagKeys); err != nil {
			return err
		}
	}

	return nil
}

func resourceIntegrationAWSUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.AWSIntegrationInput{
		Name:                 opslevel.NewString(d.Get("name").(string)),
		IAMRole:              opslevel.NewString(d.Get("iam_role").(string)),
		ExternalID:           opslevel.NewString(d.Get("external_id").(string)),
		OwnershipTagOverride: opslevel.Bool(d.Get("ownership_tag_overrides").(bool)),
	}

	input.OwnershipTagKeys = getStringArray(d, "ownership_tag_keys")

	_, err := client.UpdateIntegrationAWS(d.Id(), input)
	if err != nil {
		return err
	}

	d.Set("last_updated", timeLastUpdated())
	return resourceIntegrationAWSRead(d, client)
}

func resourceIntegrationAWSDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteIntegration(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
