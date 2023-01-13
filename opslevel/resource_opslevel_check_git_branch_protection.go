package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hasura/go-graphql-client"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceCheckGitBranchProtection() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a git branch protection check",
		Create:      wrap(resourceCheckGitBranchProtectionCreate),
		Read:        wrap(resourceCheckGitBranchProtectionRead),
		Update:      wrap(resourceCheckGitBranchProtectionUpdate),
		Delete:      wrap(resourceCheckDelete),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCheckSchema(map[string]*schema.Schema{}),
	}
}

func resourceCheckGitBranchProtectionCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckGitBranchProtectionCreateInput{}
	setCheckCreateInput(d, &input)

	resource, err := client.CreateCheckGitBranchProtection(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceCheckGitBranchProtectionRead(d, client)
}

func resourceCheckGitBranchProtectionRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetCheck(graphql.ID(id))
	if err != nil {
		return err
	}

	if err := setCheckData(d, resource); err != nil {
		return err
	}

	return nil
}

func resourceCheckGitBranchProtectionUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.CheckGitBranchProtectionUpdateInput{}
	setCheckUpdateInput(d, &input)

	_, err := client.UpdateCheckGitBranchProtection(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceCheckGitBranchProtectionRead(d, client)
}
