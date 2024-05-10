data "opslevel_services" "all" {}

data "opslevel_service" "first_service_by_id" {
  id = data.opslevel_services.all.services[0].id
}
