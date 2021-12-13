package opslevel

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opslevel/opslevel-go"
)

func resourceFilter() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a filter",
		Create:      wrap(resourceFilterCreate),
		Read:        wrap(resourceFilterRead),
		Update:      wrap(resourceFilterUpdate),
		Delete:      wrap(resourceFilterDelete),
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
				Description: "The team's display name.",
				ForceNew:    false,
				Required:    true,
			},
			"predicate": {
				Type:        schema.TypeList,
				Description: "The list of predicates used to select which services apply to the filter.",
				ForceNew:    false,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Description:  "The condition type used by the predicate. Valid values are `contains`, `does_not_contain`, `does_not_equal`, `does_not_exist`, `ends_with`, `equals`, `exists`, `greater_than_or_equal_to`, `less_than_or_equal_to`, `starts_with`, `satisfies_version_constraint`, `matches_regex`, `satisfies_jq_expression`",
							ForceNew:     false,
							Required:     true,
							ValidateFunc: validation.StringInSlice(opslevel.AllPredicateTypeEnum(), false),
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The condition value used by the predicate.",
							ForceNew:    false,
							Optional:    true,
						},
						"key": {
							Type:         schema.TypeString,
							Description:  "The condition key used by the predicate.",
							ForceNew:     false,
							Required:     true,
							ValidateFunc: validation.StringInSlice(opslevel.AllPredicateKeyEnum(), false),
						},
						"key_data": {
							Type:        schema.TypeString,
							Description: "Additional data used by the predicate. This field is used by predicates with key = 'tags' to specify the tag key. For example, to create a predicate for services containing the tag 'db:mysql', set key_data = 'db' and value = 'mysql'.",
							ForceNew:    false,
							Optional:    true,
						},
					},
				},
			},
			"connective": {
				Type:         schema.TypeString,
				Description:  "The logical operator to be used in conjunction with predicates.",
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(append(opslevel.AllConnectiveEnum(), ""), false),
			},
		},
	}
}

func resourceFilterCreate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.FilterCreateInput{
		Name:       d.Get("name").(string),
		Predicates: expandFilterPredicates(d),
		Connective: opslevel.ConnectiveEnum(d.Get("connective").(string)),
	}

	resource, err := client.CreateFilter(input)
	if err != nil {
		return err
	}
	d.SetId(resource.Id.(string))

	return resourceFilterRead(d, client)
}

func resourceFilterRead(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()

	resource, err := client.GetFilter(id)
	if err != nil {
		return err
	}

	if err := d.Set("name", resource.Name); err != nil {
		return err
	}

	if err := d.Set("connective", string(resource.Connective)); err != nil {
		return err
	}

	if err := d.Set("predicate", flattenFilterPredicates(resource.Predicates)); err != nil {
		return err
	}

	return nil
}

func resourceFilterUpdate(d *schema.ResourceData, client *opslevel.Client) error {
	input := opslevel.FilterUpdateInput{
		Id: d.Id(),
	}

	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("predicate") {
		input.Predicates = expandFilterPredicates(d)
	}
	if d.HasChange("connective") {
		input.Connective = opslevel.ConnectiveEnum(d.Get("connective").(string))
	}

	_, err := client.UpdateFilter(input)
	if err != nil {
		return err
	}
	d.Set("last_updated", timeLastUpdated())
	return resourceFilterRead(d, client)
}

func resourceFilterDelete(d *schema.ResourceData, client *opslevel.Client) error {
	id := d.Id()
	err := client.DeleteFilter(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
