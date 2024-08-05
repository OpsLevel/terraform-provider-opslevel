variables {
  name                 = "TF Test Filter with domain_id predicate"
  domain_id_predicates = setproduct(["domain_id"], var.predicate_types_equals_and_exists)
}

run "domain_module" {
  command = plan

  variables {
    name = "placeholder"
  }

  module {
    source = "./domain"
  }
}

run "resource_filter_with_domain_id_predicate_create" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.domain_id_predicates : "${pair[0]}_${pair[1]}" => {
        key = pair[0], type = pair[1], key_data = null, value = contains(var.predicate_types_exists, pair[1]) ? null : run.domain_module.first_domain.id
      }
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_does_not_equal"].predicate[0].key == "domain_id"
    error_message = format(
      "expected predicate key 'domain_id' got '%s'",
      opslevel_filter.all_predicates["domain_id_does_not_equal"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_does_not_equal"].predicate[0].type == "does_not_equal"
    error_message = format(
      "expected predicate type 'does_not_equal' got '%s'",
      opslevel_filter.all_predicates["domain_id_does_not_equal"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["domain_id_does_not_equal"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_does_not_equal"].predicate[0].value == run.domain_module.first_domain.id
    error_message = format(
      "expected predicate value '%s' got '%s'",
      run.domain_module.first_domain.id,
      opslevel_filter.all_predicates["domain_id_does_not_equal"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_equals"].predicate[0].key == "domain_id"
    error_message = format(
      "expected predicate key 'domain_id' got '%s'",
      opslevel_filter.all_predicates["domain_id_equals"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_equals"].predicate[0].type == "equals"
    error_message = format(
      "expected predicate type 'equals' got '%s'",
      opslevel_filter.all_predicates["domain_id_equals"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["domain_id_equals"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_equals"].predicate[0].value == run.domain_module.first_domain.id
    error_message = format(
      "expected predicate value '%s' got '%s'",
      run.domain_module.first_domain.id,
      opslevel_filter.all_predicates["domain_id_equals"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_does_not_exist"].predicate[0].key == "domain_id"
    error_message = format(
      "expected predicate key 'domain_id' got '%s'",
      opslevel_filter.all_predicates["domain_id_does_not_exist"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_does_not_exist"].predicate[0].type == "does_not_exist"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["domain_id_does_not_exist"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["domain_id_does_not_exist"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["domain_id_does_not_exist"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_exists"].predicate[0].key == "domain_id"
    error_message = format(
      "expected predicate key 'domain_id' got '%s'",
      opslevel_filter.all_predicates["domain_id_exists"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["domain_id_exists"].predicate[0].type == "exists"
    error_message = format(
      "expected predicate type 'exists' got '%s'",
      opslevel_filter.all_predicates["domain_id_exists"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["domain_id_exists"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["domain_id_exists"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

}
