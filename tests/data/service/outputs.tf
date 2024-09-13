output "all" {
  value = data.opslevel_services.all
}

output "first" {
  value = data.opslevel_services.all.services[0]
}
