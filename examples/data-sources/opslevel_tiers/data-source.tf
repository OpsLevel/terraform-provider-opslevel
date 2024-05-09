data "opslevel_tiers" "all" {}

output "all" {
  value = data.opslevel_tiers.all.tiers
}

output "tier_aliases" {
  value = sort(data.opslevel_tiers.all.tiers[*].alias)
}

output "tier_names" {
  value = sort(data.opslevel_tiers.all.tiers[*].name)
}
