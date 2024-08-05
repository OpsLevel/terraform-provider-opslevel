variables {
  name                      = "TF Test Filter with repository_ids predicate"
  repository_ids_predicates = setproduct(["repository_ids"], var.predicate_types_exists)
}

run "resource_filter_with_repository_ids_predicate_exists" {

  variables {
    predicates = tomap({
      for pair in var.repository_ids_predicates : "${pair[0]}_${pair[1]}" => {
        key = pair[0],
        type = pair[1],
        key_data = null,
        value = null
      }
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["repository_ids_does_not_exist"].predicate[0].key == "repository_ids"
    error_message = format(
      "expected predicate key 'repository_ids' got '%s'",
      opslevel_filter.all_predicates["repository_ids_does_not_exist"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["repository_ids_does_not_exist"].predicate[0].type == "does_not_exist"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["repository_ids_does_not_exist"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["repository_ids_does_not_exist"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["repository_ids_does_not_exist"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["repository_ids_exists"].predicate[0].key == "repository_ids"
    error_message = format(
      "expected predicate key 'repository_ids' got '%s'",
      opslevel_filter.all_predicates["repository_ids_exists"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["repository_ids_exists"].predicate[0].type == "exists"
    error_message = format(
      "expected predicate type 'exists' got '%s'",
      opslevel_filter.all_predicates["repository_ids_exists"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["repository_ids_exists"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["repository_ids_exists"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

}
