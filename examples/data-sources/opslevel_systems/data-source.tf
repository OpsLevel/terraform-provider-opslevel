data "opslevel_systems" "all" {}

output "all" {
  value = data.opslevel_systems.all.systems
}

output "system_names" {
  value = sort(data.opslevel_systems.all.systems[*].name)
}
