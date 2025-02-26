---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opslevel_filter Resource - terraform-provider-opslevel"
subcategory: ""
description: |-
  Filter Resource
---

# opslevel_filter (Resource)

Filter Resource

## Example Usage

```terraform
resource "opslevel_filter" "tier1" {
  name = "foo"
  predicate {
    key   = "tier_index"
    type  = "equals"
    value = "1"
  }
  connective = "and"
}

resource "opslevel_filter" "tier2_alpha" {
  name = "foo"
  predicate {
    key   = "tier_index"
    type  = "equals"
    value = "1"
  }
  predicate {
    key   = "lifecycle_index"
    type  = "equals"
    value = "1"
  }
  connective = "and"
}

resource "opslevel_filter" "tier3" {
  name = "foo"
  predicate {
    key      = "tags"
    type     = "does_not_exist"
    key_data = "tier3"
  }
}

resource "opslevel_filter" "case_sensitive" {
  name = "foo"
  predicate {
    key            = "tags"
    type           = "equals"
    key_data       = "my-custom-tag"
    value          = "hello-world"
    case_sensitive = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The filter's display name.

### Optional

- `connective` (String) The logical operator to be used in conjunction with predicates. One of `and`, `or`
- `predicate` (Block List) (see [below for nested schema](#nestedblock--predicate))

### Read-Only

- `id` (String) The ID of the filter.

<a id="nestedblock--predicate"></a>
### Nested Schema for `predicate`

Required:

- `key` (String) The condition key used by the predicate. Valid values are `aliases`, `component_type_id`, `creation_source`, `domain_id`, `filter_id`, `framework`, `group_ids`, `language`, `lifecycle_index`, `name`, `owner_id`, `owner_ids`, `product`, `properties`, `repository_ids`, `system_id`, `tags`, `tier_index`
- `type` (String) The condition type used by the predicate. Valid values are `belongs_to`, `contains`, `does_not_contain`, `does_not_equal`, `does_not_exist`, `does_not_match`, `does_not_match_regex`, `ends_with`, `equals`, `exists`, `greater_than_or_equal_to`, `less_than_or_equal_to`, `matches`, `matches_regex`, `satisfies_jq_expression`, `satisfies_version_constraint`, `starts_with`

Optional:

- `case_insensitive` (Boolean, Deprecated) Option for determining whether to compare strings case-sensitively. Not settable for all predicate types.
- `case_sensitive` (Boolean) Option for determining whether to compare strings case-sensitively. Not settable for all predicate types.
- `key_data` (String) Additional data used by the predicate. This field is used by predicates with key = 'tags' to specify the tag key. For example, to create a predicate for services containing the tag 'db:mysql', set key_data = 'db' and value = 'mysql'.
- `value` (String) The condition value used by the predicate.

## Import

Import is supported using the following syntax:

```shell
terraform import opslevel_filter.example Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS82MDI0
```
