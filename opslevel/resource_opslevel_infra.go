package opslevel

// import (
// 	"slices"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
// 	"github.com/opslevel/opslevel-go/v2024"
// )

// func resourceInfrastructure() *schema.Resource {
// 	return &schema.Resource{
// 		Description: "Manages an infrastructure resource",
// 		Create:      wrap(resourceInfrastructureCreate),
// 		Read:        wrap(resourceInfrastructureRead),
// 		Update:      wrap(resourceInfrastructureUpdate),
// 		Delete:      wrap(resourceInfrastructureDelete),
// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"aliases": {
// 				Type:        schema.TypeList,
// 				Description: "The aliases of the infrastructure resource.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				Elem:        &schema.Schema{Type: schema.TypeString},
// 			},
// 			"schema": {
// 				Type:        schema.TypeString,
// 				Description: "The schema of the infrastructure resource that determines its data specification.",
// 				Required:    true,
// 				ForceNew:    true,
// 			},
// 			"owner": {
// 				Type:        schema.TypeString,
// 				Description: "The id of the team that owns the infrastructure resource. Does not support aliases!",
// 				ForceNew:    false,
// 				Optional:    true,
// 			},
// 			"provider_data": {
// 				Type:        schema.TypeList,
// 				Description: "The provider specific data for the infrastructure resource.",
// 				ForceNew:    false,
// 				Optional:    true,
// 				MaxItems:    1,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"name": {
// 							Type:        schema.TypeString,
// 							Description: "The name of the provider of the infrastructure resource. (eg. AWS, GCP, Azure)",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 						"type": {
// 							Type:        schema.TypeString,
// 							Description: "The type of the infrastructure resource as defined by its provider.",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 						"account": {
// 							Type:        schema.TypeString,
// 							Description: "The canonical account name for the provider of the infrastructure resource.",
// 							ForceNew:    false,
// 							Required:    true,
// 						},
// 						"url": {
// 							Type:        schema.TypeString,
// 							Description: "The url for the provider of the infrastructure resource.",
// 							ForceNew:    false,
// 							Optional:    true,
// 						},
// 					},
// 				},
// 			},
// 			"data": {
// 				Type:        schema.TypeString,
// 				Description: "The data of the infrastructure resource in JSON format.",
// 				Optional:    true,
// 			},
// 		},
// 	}
// }

// func reconcileInfraAliases(d *schema.ResourceData, resource *opslevel.InfrastructureResource, client *opslevel.Client) error {
// 	expectedAliases := getStringArray(d, "aliases")
// 	existingAliases := resource.Aliases
// 	for _, existingAlias := range existingAliases {
// 		if !slices.Contains(expectedAliases, existingAlias) {
// 			err := client.DeleteInfraAlias(existingAlias)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	for _, expectedAlias := range expectedAliases {
// 		if !slices.Contains(existingAliases, expectedAlias) {
// 			id := opslevel.NewID(resource.Id)
// 			_, err := client.CreateAliases(*id, []string{expectedAlias})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func flattenInfraProviderData(resource *opslevel.InfrastructureResource) []map[string]any {
// 	return []map[string]any{{
// 		"account": resource.ProviderData.AccountName,
// 		"name":    resource.ProviderData.ProviderName,
// 		"type":    resource.ProviderType,
// 		"url":     resource.ProviderData.ExternalURL,
// 	}}
// }

// func expandInfraProviderData(d *schema.ResourceData) *opslevel.InfraProviderInput {
// 	config := d.Get("provider_data").([]interface{})
// 	if len(config) > 0 {
// 		item := config[0].(map[string]interface{})
// 		return &opslevel.InfraProviderInput{
// 			Account: item["account"].(string),
// 			Name:    item["name"].(string),
// 			Type:    item["type"].(string),
// 			URL:     item["url"].(string),
// 		}
// 	}
// 	return &opslevel.InfraProviderInput{}
// }

// func resourceInfrastructureCreate(d *schema.ResourceData, client *opslevel.Client) error {
// 	newJSON, err := opslevel.NewJSON(d.Get("data").(string))
// 	if err != nil {
// 		return err
// 	}
// 	resource, err := client.CreateInfrastructure(opslevel.InfraInput{
// 		Schema:   d.Get("schema").(string),
// 		Owner:    opslevel.NewID(d.Get("owner").(string)),
// 		Provider: expandInfraProviderData(d),
// 		Data:     newJSON,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(resource.Id)

// 	err = reconcileInfraAliases(d, resource, client)
// 	if err != nil {
// 		return err
// 	}

// 	return resourceInfrastructureRead(d, client)
// }

// func resourceInfrastructureRead(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	resource, err := client.GetInfrastructure(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := d.Set("schema", resource.Schema); err != nil {
// 		return err
// 	}
// 	if err := d.Set("aliases", resource.Aliases); err != nil {
// 		return err
// 	}
// 	if err := d.Set("owner", resource.Owner.Id()); err != nil {
// 		return err
// 	}
// 	if err := d.Set("provider_data", flattenInfraProviderData(resource)); err != nil {
// 		return err
// 	}
// 	if err := d.Set("data", resource.Data.ToJSON()); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func resourceInfrastructureUpdate(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()

// 	newJSON, err := opslevel.NewJSON(d.Get("data").(string))
// 	if err != nil {
// 		return err
// 	}
// 	resource, err := client.UpdateInfrastructure(id, opslevel.InfraInput{
// 		Schema:   d.Get("schema").(string),
// 		Owner:    opslevel.NewID(d.Get("owner").(string)),
// 		Provider: expandInfraProviderData(d),
// 		Data:     newJSON,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if d.HasChange("aliases") {
// 		err = reconcileInfraAliases(d, resource, client)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	d.Set("last_updated", timeLastUpdated())
// 	return resourceInfrastructureRead(d, client)
// }

// func resourceInfrastructureDelete(d *schema.ResourceData, client *opslevel.Client) error {
// 	id := d.Id()
// 	err := client.DeleteInfrastructure(id)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
