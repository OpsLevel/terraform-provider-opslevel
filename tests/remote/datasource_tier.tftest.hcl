run "datasource_tiers_all" {

  assert {
    condition     = length(data.opslevel_tiers.all.tiers) > 0
    error_message = "zero tiers found in data.opslevel_tiers"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_tiers.all.tiers[0].alias),
      can(data.opslevel_tiers.all.tiers[0].id),
      can(data.opslevel_tiers.all.tiers[0].index),
      can(data.opslevel_tiers.all.tiers[0].name),
    ])
    error_message = "cannot set all expected tier datasource fields"
  }

}

run "datasource_tier_first" {

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.alias == data.opslevel_tiers.all.tiers[0].alias
    error_message = "wrong alias on opslevel_tier"
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.id == data.opslevel_tiers.all.tiers[0].id
    error_message = "wrong ID on opslevel_tier"
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_id.index == data.opslevel_tiers.all.tiers[0].index
    error_message = "wrong index on opslevel_tier"
  }

  assert {
    condition     = data.opslevel_tier.first_tier_by_name.name == data.opslevel_tiers.all.tiers[0].name
    error_message = "wrong name on opslevel_tier"
  }

}

