data "opslevel_services" "all" {}

data "opslevel_tier" "tier1" {
  filter {
    field = "alias"
    value = "tier_1"
  }
}

data "opslevel_services" "tier1" {
  filter = {
    field = "tier"
    value = data.opslevel_tier.tier1.alias
  }
}

data "opslevel_services" "frontend" {
  filter = {
    field = "owner"
    value = "frontend"
  }
}

output "all_services" {
  value = data.opslevel_services.all.services
}

output "all_service_names" {
  value = sort(data.opslevel_services.all.services[*].name)
}

output "tier1_services" {
  value = data.opslevel_services.tier1.services
}

output "frontend_services" {
  value = data.opslevel_services.frontend.services
}


output "frontend_services_urls" {
  value = sort(data.opslevel_services.frontend.services[*].url)
}
