variables {
  name                    = "TF Test Filter with filter_id predicate"
  predicate_key           = "filter_id"
  predicate_types_matches = ["does_not_match", "matches"]
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_filters.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_filter_with_filter_id_predicate_matches" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_matches : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = null,
        value    = run.from_data_module.first_filter.id
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
      contains(var.predicate_types_matches, opslevel_filter.this.predicate[0].type),
      contains(var.predicate_types_matches, opslevel_filter.this.predicate[1].type)
    ])
    error_message = format(
      "expected predicate types to be one of '%v' got '%v'",
      var.predicate_types_matches,
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
      opslevel_filter.this.predicate[0].value == run.from_data_module.first_filter.id,
      opslevel_filter.this.predicate[1].value == run.from_data_module.first_filter.id
    ])
    error_message = format(
      "expected predicate values to all be '%v' got values '%v'",
      run.from_data_module.first_filter.id,
      distinct([opslevel_filter.this.predicate[0].value, opslevel_filter.this.predicate[1].value])
    )
  }

}
