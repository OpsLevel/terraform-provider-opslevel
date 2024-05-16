variables {
  tier_one  = "opslevel_tier"
  tiers_all = "opslevel_tiers"
}

run "datasource_tiers_all" {

  module {
    source = "./tier"
  }

  assert {
    condition     = can(data.opslevel_tiers.all.tiers)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.tiers_all)
  }

  assert {
    condition     = length(data.opslevel_tiers.all.tiers) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.tiers_all)
  }

}

run "datasource_tier_first" {

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
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.tier_one)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.alias == data.opslevel_tiers.all.tiers[0].alias
    error_message = replace(var.error_wrong_alias, "TYPE", var.tier_one)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.id == data.opslevel_tiers.all.tiers[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.tier_one)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.index == data.opslevel_tiers.all.tiers[0].index
    error_message = replace(var.error_wrong_index, "TYPE", var.tier_one)
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_name.name == data.opslevel_tiers.all.tiers[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.tier_one)
  }

}
