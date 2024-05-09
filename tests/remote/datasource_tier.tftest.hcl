run "datasource_tiers_all" {

  variables {
    datasource_type = "opslevel_tiers"
  }

  assert {
    condition     = can(data.opslevel_tiers.all.tiers)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_tiers.all.tiers) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_tier_first" {

  variables {
    datasource_type = "opslevel_tier"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_tier.first_tier_by_id.alias),
      can(data.opslevel_tier.first_tier_by_id.id),
      can(data.opslevel_tier.first_tier_by_id.index),
      can(data.opslevel_tier.first_tier_by_id.name),
    ])
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.alias == data.opslevel_tiers.all.tiers[0].alias
    error_message = replace(var.wrong_alias_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.id == data.opslevel_tiers.all.tiers[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.index == data.opslevel_tiers.all.tiers[0].index
    error_message = replace(var.wrong_index_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_name.name == data.opslevel_tiers.all.tiers[0].name
    error_message = replace(var.wrong_name_error, "TYPE", var.datasource_type)
  }

}

