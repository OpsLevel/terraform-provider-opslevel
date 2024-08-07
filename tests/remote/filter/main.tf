data "opslevel_filters" "all" {}

data "opslevel_filter" "first_filter_by_name" {
  filter {
    field = "name"
    value = data.opslevel_filters.all.filters[0].name
  }
}

data "opslevel_filter" "first_filter_by_id" {
  filter {
    field = "id"
    value = data.opslevel_filters.all.filters[0].id
  }
}

resource "opslevel_filter" "test" {
  name       = var.name
  connective = var.connective
}

resource "opslevel_filter" "all_predicates" {
  for_each = var.predicates
  name     = each.key
  predicate {
    key      = each.value.key
    key_data = each.value.key_data
    type     = each.value.type
    value    = each.value.value
  }
}
