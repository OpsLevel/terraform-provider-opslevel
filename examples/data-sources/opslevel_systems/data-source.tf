data "opslevel_systems" "all" {
}

output "found" {
  value = data.opslevel_systems.all
}