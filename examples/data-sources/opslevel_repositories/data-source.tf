data "opslevel_repositories" "all" {}

data "opslevel_tier" "tier2" {
  filter {
    field = "alias"
    value = "tier_2"
  }
}

data "opslevel_repositories" "tier2" {
  filter = {
    field = "tier"
    value = data.opslevel_tier.tier2.alias
  }
}

output "all" {
  value = data.opslevel_repositories.all.repositories
}

output "tier2_repositories" {
  value = data.opslevel_repositories.tier2.repositories
}

output "all_repository_names" {
  value = sort(data.opslevel_repositories.all.repositories[*].name)
}

