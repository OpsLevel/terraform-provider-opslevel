variables {
  name                 = "TF Test Filter with filter_id predicate"
  filter_id_predicates = setproduct(["filter_id"], ["does_not_match", "matches"])
}

run "filter_module" {
  command = plan

  variables {
    name = "placeholder"
  }

  module {
    source = "./filter"
  }
}

run "resource_filter_with_filter_id_predicate_create" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.filter_id_predicates : "${pair[0]}_${pair[1]}" => {
        key = pair[0], type = pair[1], key_data = null, value = run.filter_module.first_filter.id
      }
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["filter_id_does_not_match"].predicate[0].key == "filter_id"
    error_message = format(
      "expected predicate key 'filter_id' got '%s'",
      opslevel_filter.all_predicates["filter_id_does_not_match"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["filter_id_does_not_match"].predicate[0].type == "does_not_match"
    error_message = format(
      "expected predicate type 'does_not_match' got '%s'",
      opslevel_filter.all_predicates["filter_id_does_not_match"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["filter_id_does_not_match"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["filter_id_does_not_match"].predicate[0].value == run.filter_module.first_filter.id
    error_message = format(
      "expected predicate value '%s' got '%s'",
      run.filter_module.first_filter.id,
      opslevel_filter.all_predicates["filter_id_does_not_match"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["filter_id_matches"].predicate[0].key == "filter_id"
    error_message = format(
      "expected predicate key 'filter_id' got '%s'",
      opslevel_filter.all_predicates["filter_id_matches"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["filter_id_matches"].predicate[0].type == "matches"
    error_message = format(
      "expected predicate type 'matches' got '%s'",
      opslevel_filter.all_predicates["filter_id_matches"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["filter_id_matches"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["filter_id_matches"].predicate[0].value == run.filter_module.first_filter.id
    error_message = format(
      "expected predicate value '%s' got '%s'",
      run.filter_module.first_filter.id,
      opslevel_filter.all_predicates["filter_id_matches"].predicate[0].type
    )
  }

}
