variables {
  filter_one  = "opslevel_filter"
  filters_all = "opslevel_filters"

  # required fields
  name = "TF Test Filter"

  # optional fields
  connective = "and"
}

run "resource_filter_create_with_all_fields" {

  variables {
    connective = var.connective
    name       = var.name
  }

  module {
    source = "./opslevel_modules/modules/filter"
  }

  assert {
    condition = alltrue([
      can(opslevel_filter.this.connective),
      can(opslevel_filter.this.id),
      can(opslevel_filter.this.name),
      can(opslevel_filter.this.predicate),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.filter_one)
  }

  assert {
    condition     = opslevel_filter.this.connective == var.connective
    error_message = "wrong connective of opslevel_filter resource"
  }

  assert {
    condition     = opslevel_filter.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.filter_one)
  }

}

#run "resource_filter_update_unset_optional_fields" {
#
#  variables {
#    predicate_list = null
#  }
#
#  module {
#    source = "./opslevel_modules/modules/filter"
#  }
#
#  assert {
#    condition     = opslevel_filter.this.predicate_list == null
#    error_message = var.error_expected_null_field
#  }
#
#}

run "resource_filter_update_set_all_fields" {

  variables {
    connective = var.connective == "and" ? "or" : "and"
    name       = "${var.name} updated"
    # predicate_list = var.predicate_list
  }

  module {
    source = "./opslevel_modules/modules/filter"
  }

  assert {
    condition     = opslevel_filter.this.connective == var.connective
    error_message = "wrong connective of opslevel_filter resource"
  }

  assert {
    condition     = opslevel_filter.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.filter_one)
  }

  #assert {
  #  condition     = opslevel_filter.this.predicate_list == var.predicate_list
  #  error_message = "wrong predicate_list of opslevel_filter resource"
  #}

}
