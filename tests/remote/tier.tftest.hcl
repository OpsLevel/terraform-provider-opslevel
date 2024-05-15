run "datasource_tiers_all" {

  variables {
    datasource_type = "opslevel_tiers"
  }

  module {
    source = "./tier"
  }

  assert {
    condition     = can(data.opslevel_tiers.all.tiers)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_tiers.all.tiers) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

}

run "datasource_tier_first" {

  variables {
    datasource_type = "opslevel_tier"
  }

  module {
    source = "./tier"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_tier.first_tier_by_id.alias),
      can(data.opslevel_tier.first_tier_by_id.id),
      can(data.opslevel_tier.first_tier_by_id.index),
      can(data.opslevel_tier.first_tier_by_id.name),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.alias == data.opslevel_tiers.all.tiers[0].alias
    error_message = replace(var.error_wrong_alias, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.id == data.opslevel_tiers.all.tiers[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.index == data.opslevel_tiers.all.tiers[0].index
    error_message = replace(var.error_wrong_index, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_name.name == data.opslevel_tiers.all.tiers[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.datasource_type)
  }

}

run "resource_tier_create_with_all_fields" {

  variables {
  }

  module {
    source = "./tier"
  }

}

run "resource_tier_update_unset_optional_fields" {

  variables {
  }

  module {
    source = "./tier"
  }

}

run "resource_tier_update_set_optional_fields" {

  variables {
  }

  module {
    source = "./tier"
  }

}
