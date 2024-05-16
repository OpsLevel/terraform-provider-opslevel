data "opslevel_tiers" "all" {}

# data "opslevel_tier" "first_tier_by_alias" {
#   filter {
#     field = "alias"
#     value = data.opslevel_tiers.all.tiers[0].alias
#   }
# }

data "opslevel_tier" "first_tier_by_id" {
  filter {
    field = "id"
    value = data.opslevel_tiers.all.tiers[0].id
  }
}

data "opslevel_tier" "first_tier_by_index" {
  filter {
    field = "index"
    value = data.opslevel_tiers.all.tiers[0].index
  }
}

data "opslevel_tier" "first_tier_by_name" {
  filter {
    field = "name"
    value = data.opslevel_tiers.all.tiers[0].name
  }
}
