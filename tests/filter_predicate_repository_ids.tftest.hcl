variables {
  name          = "TF Test Filter with repository_ids predicate"
  predicate_key = "repository_ids"
}

run "resource_filter_with_repository_ids_predicate_exists" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_exists : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = null,
        value    = null
      }
    ]
  }

  module {
    source = "./opslevel_modules/modules/filter"
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].key == var.predicate_key,
      opslevel_filter.this.predicate[1].key == var.predicate_key
    ])
    error_message = format(
      "expected predicate keys to all be '%v' got keys '%v'",
      var.predicate_key,
      [opslevel_filter.this.predicate[0].key, opslevel_filter.this.predicate[1].key]
    )
  }

  assert {
    condition = alltrue([
      contains(var.predicate_types_exists, opslevel_filter.this.predicate[0].type),
      contains(var.predicate_types_exists, opslevel_filter.this.predicate[1].type)
    ])
    error_message = format(
      "expected predicate types to be one of '%v' got '%v'",
      var.predicate_types_exists,
      tolist([opslevel_filter.this.predicate[0].type, opslevel_filter.this.predicate[1].type])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].key_data == null,
      opslevel_filter.this.predicate[1].key_data == null
    ])
    error_message = var.error_expected_null_field
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].value == null,
      opslevel_filter.this.predicate[1].value == null
    ])
    error_message = var.error_expected_null_field
  }

}
