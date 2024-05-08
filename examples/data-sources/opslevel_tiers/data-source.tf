data "opslevel_tiers" "all" {}

output "found" {
  value = data.opslevel_tiers.all.aliases[0]
}
