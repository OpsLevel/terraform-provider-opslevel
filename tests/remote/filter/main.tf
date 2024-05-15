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
  #predicate {
  #  case_insensitive = var.predicate.case_insensitive
  #  case_sensitive   = var.predicate.case_sensitive
  #  key              = var.predicate.key
  #  key_data         = var.predicate.key_data
  #  type             = var.predicate.type
  #  value            = var.predicate.value
  #}
}
