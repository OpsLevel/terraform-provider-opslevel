variables {
  name                  = "TF Test Filter with properties predicate"
  jq_expression         = ".[] | select(.name == fancy)"
  predicate_key_data    = "property_definition_id"
  properties_predicates = setproduct(["properties"], concat(var.predicate_types_exists, ["satisfies_jq_expression"]))
}

run "resource_filter_with_properties_predicate_exists" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.properties_predicates : "${pair[0]}_${pair[1]}" => {
        key = pair[0],
        type = pair[1],
        key_data = var.predicate_key_data,
        value = null
      }
      if contains(var.predicate_types_exists, pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_does_not_exist"].predicate[0].key == "properties"
    error_message = format(
      "expected predicate key 'properties' got '%s'",
      opslevel_filter.all_predicates["properties_does_not_exist"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_does_not_exist"].predicate[0].type == "does_not_exist"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["properties_does_not_exist"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_does_not_exist"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["properties_does_not_exist"].predicate[0].key_data
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["properties_does_not_exist"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_exists"].predicate[0].key == "properties"
    error_message = format(
      "expected predicate key 'properties' got '%s'",
      opslevel_filter.all_predicates["properties_exists"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_exists"].predicate[0].type == "exists"
    error_message = format(
      "expected predicate type 'exists' got '%s'",
      opslevel_filter.all_predicates["properties_exists"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_exists"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["properties_exists"].predicate[0].key_data
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["properties_exists"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

}

run "resource_filter_with_properties_predicate_satisfies_jq_expression" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.properties_predicates : "${pair[0]}_${pair[1]}" => {
        key = pair[0],
        type = pair[1],
        key_data = var.predicate_key_data,
        value = var.jq_expression
      }
      if "satisfies_jq_expression" == pair[1]
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_satisfies_jq_expression"].predicate[0].key == "properties"
    error_message = format(
      "expected predicate key 'properties' got '%s'",
      opslevel_filter.all_predicates["properties_satisfies_jq_expression"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_satisfies_jq_expression"].predicate[0].type == "satisfies_jq_expression"
    error_message = format(
      "expected predicate type 'satisfies_jq_expression' got '%s'",
      opslevel_filter.all_predicates["properties_satisfies_jq_expression"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_satisfies_jq_expression"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["properties_satisfies_jq_expression"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["properties_satisfies_jq_expression"].predicate[0].value == var.jq_expression
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.jq_expression,
      opslevel_filter.all_predicates["properties_satisfies_jq_expression"].predicate[0].type
    )
  }

}
