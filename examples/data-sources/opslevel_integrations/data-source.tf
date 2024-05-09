data "opslevel_integrations" "all" {}

output "all" {
  value = data.opslevel_integrations.all.integrations
}

output "integration_names" {
  value = sort(data.opslevel_integrations.all.integrations[*].name)
}
