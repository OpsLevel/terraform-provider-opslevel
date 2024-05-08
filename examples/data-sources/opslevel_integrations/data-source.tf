data "opslevel_integrations" "all" {}

output "found" {
  value = data.opslevel_integrations.all.id[0]
}
