package opslevel

import (
	"log"

	"github.com/shurcooL/graphql"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kr/pretty"
	"github.com/opslevel/opslevel-go"
)

func resourceOpsLevelService() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsLevelServiceCreate,
		Read:   resourceOpsLevelServiceRead,
		Update: resourceOpsLevelServiceUpdate,
		Delete: resourceOpsLevelServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: serviceSchema,
	}
}

func resourceOpsLevelServiceCreate(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)

	svcCreate := opslevel.ServiceCreateInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Framework:   d.Get("framework").(string),
		Language:    d.Get("language").(string),
		Product:     d.Get("product").(string),
	}

	svc, err := p.client.CreateService(svcCreate)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] created OpsLevel service ID:", svc.Id)
	d.SetId(svc.Id.(string))

	// add Aliases
	aliases := expandServiceAliases(d.Get("aliases").([]interface{}))
	if len(aliases) > 0 {
		p.client.CreateAliases(svc.Id, aliases)
	}

	// add tags
	if v, ok := d.Get("tags").(map[string]interface{}); ok && len(v) > 0 {
		p.client.AssignTagsForId(svc.Id, expandStringMap(d.Get("tags").(map[string]interface{})))
	}

	return resourceOpsLevelServiceRead(d, meta)
}

func resourceOpsLevelServiceRead(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)
	id := d.Id()

	log.Println("[DEBUG] querying OpsLevel for service with ID:", id)

	svc, err := GetServiceWithId(p.client, id)
	if err != nil {
		log.Printf("[DEBUG] query service error: %s", pretty.Sprint(err))

		return err
	}
	if svc.Id == nil || svc.Id.(string) == "" {
		log.Println("[DEBUG] service no longer exists")
		d.SetId("")
		return nil
	}

	flattened := flattenService(svc)
	for k, v := range flattened {
		d.Set(k, v)
	}

	return nil
}

func resourceOpsLevelServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	//p := meta.(provider)
	//id := d.Id()

	return resourceOpsLevelServiceRead(d, meta)
}

func resourceOpsLevelServiceDelete(d *schema.ResourceData, meta interface{}) error {
	p := meta.(provider)
	id := d.Id()

	delInput := opslevel.ServiceDeleteInput{
		Id: graphql.ID(id),
	}

	err := p.client.DeleteService(delInput)
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
