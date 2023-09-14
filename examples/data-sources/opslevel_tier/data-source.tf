data "opslevel_tier" "tier1" {
  filter {
    field = "alias"
    value = "tier_1"
  }
}

data "opslevel_tier" "tier3" {
  filter {
    field = "index"
    value = "3"
  }
}