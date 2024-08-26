output "first_service" {
  value = data.opslevel_service.first_service_by_id
}

output "last_service" {
  value = data.opslevel_service.last_service
}
