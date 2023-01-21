package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go/v2023"
)

func resourceRubricLevel() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a rubric level",
		Create:      wrap(resourceRubricLevelCreate),
		Read:        wrap(resourceRubricLevelRead),
		Update:      wrap(resourceRubricLevelUpdate),
		Delete:      wrap(resourceRubricLevelDelete),
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
				Description: "The display name of the category.",
				ForceNew:    false,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the category.",
				ForceNew:    false,
				Optional:    true,
			},
			"index": {
				Type:         schema.TypeInt,
				Description:  "An integer allowing this level to be inserted between others. Must be unique per Rubric.",
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3, 4, 5, 6}),
			},
		},
	}
}

func resourceRubricLevelCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.LevelCreateInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	if v, ok := d.GetOk("index"); ok {
		index := v.(int)
		input.Index = &index
	}
	resource, err := client.CreateLevel(input)
	if err != nil {
		return err
	}
	d.SetId(string(resource.Id))

	return resourceRubricLevelRead(d, client)
}

func resourceRubricLevelRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetLevel(opslevel.ID(id))
	if err != nil {
		return err
	}

	if err := d.Set("name", resource.Name); err != nil {
		return err
	}
	if err := d.Set("description", resource.Description); err != nil {
		return err
	}
	if err := d.Set("index", resource.Index); err != nil {
		return err
	}

	return nil
}

func resourceRubricLevelUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.LevelUpdateInput{
		Id: opslevel.ID(d.Id()),
	}

	if d.HasChange("name") {
		input.Name = *opslevel.NewString(d.Get("name").(string))
	}
	if d.HasChange("description") {
		input.Description = opslevel.NewString(d.Get("description").(string))
	}

	_, err := client.UpdateLevel(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceRubricLevelRead(d, client)
}

func resourceRubricLevelDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteLevel(opslevel.ID(id))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
