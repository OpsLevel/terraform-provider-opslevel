variables {
  name = "TF Test Filter with owner_id predicate"
  owner_id_predicates = setproduct(
    ["owner_id"],
    concat(var.predicate_types_equals, var.predicate_types_exists)
  )
}

run "get_team" {
  command = plan

  variables {
    name = "placeholder"
  }

  module {
    source = "./team"
  }
}

run "resource_filter_with_owner_id_predicate_create" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.owner_id_predicates : "${pair[0]}_${pair[1]}" => {
        key = pair[0],
        type = pair[1],
        key_data = null,
        value = run.get_team.first_team.id
      }
      if contains(var.predicate_types_equals, pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_does_not_equal"].predicate[0].key == "owner_id"
    error_message = format(
      "expected predicate key 'owner_id' got '%s'",
      opslevel_filter.all_predicates["owner_id_does_not_equal"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_does_not_equal"].predicate[0].type == "does_not_equal"
    error_message = format(
      "expected predicate type 'does_not_equal' got '%s'",
      opslevel_filter.all_predicates["owner_id_does_not_equal"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["owner_id_does_not_equal"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_does_not_equal"].predicate[0].value == run.get_team.first_team.id
    error_message = format(
      "expected predicate value '%s' got '%s'",
      run.get_team.first_team.id,
      opslevel_filter.all_predicates["owner_id_does_not_equal"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_equals"].predicate[0].key == "owner_id"
    error_message = format(
      "expected predicate key 'owner_id' got '%s'",
      opslevel_filter.all_predicates["owner_id_equals"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_equals"].predicate[0].type == "equals"
    error_message = format(
      "expected predicate type 'equals' got '%s'",
      opslevel_filter.all_predicates["owner_id_equals"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["owner_id_equals"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_equals"].predicate[0].value == run.get_team.first_team.id
    error_message = format(
      "expected predicate value '%s' got '%s'",
      run.get_team.first_team.id,
      opslevel_filter.all_predicates["owner_id_equals"].predicate[0].type
    )
  }

}

run "resource_filter_with_owner_id_predicate_exists" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.owner_id_predicates : "${pair[0]}_${pair[1]}" => {
        key = pair[0],
        type = pair[1],
        key_data = null,
        value = null
      }
      if contains(var.predicate_types_exists, pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_does_not_exist"].predicate[0].key == "owner_id"
    error_message = format(
      "expected predicate key 'owner_id' got '%s'",
      opslevel_filter.all_predicates["owner_id_does_not_exist"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_does_not_exist"].predicate[0].type == "does_not_exist"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["owner_id_does_not_exist"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["owner_id_does_not_exist"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["owner_id_does_not_exist"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_exists"].predicate[0].key == "owner_id"
    error_message = format(
      "expected predicate key 'owner_id' got '%s'",
      opslevel_filter.all_predicates["owner_id_exists"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["owner_id_exists"].predicate[0].type == "exists"
    error_message = format(
      "expected predicate type 'exists' got '%s'",
      opslevel_filter.all_predicates["owner_id_exists"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["owner_id_exists"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["owner_id_exists"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

}
