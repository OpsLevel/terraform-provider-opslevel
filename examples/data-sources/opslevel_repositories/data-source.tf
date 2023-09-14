data "opslevel_tier" "tier2" {
  filter {
    field = "alias"
    value = "tier_2"
  }
}

data "opslevel_repositories" "all" {
}

data "opslevel_repositories" "tier2" {
  filter {
    field = "tier"
    value = data.opslevel_tier.tier2.alias
  }
}

output "all" {
  value = data.opslevel_repositories.all.names
}

output "tier2" {
  value = data.opslevel_repositories.tier2.names
}