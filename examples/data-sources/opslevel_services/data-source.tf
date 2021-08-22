data "opslevel_tier" "tier1" {
    filter {
        field = "alias"
        value = "tier_1"
    }
}

data "opslevel_services" "all" {
}

data "opslevel_services" "tier1" {
  filter {
    field = "tier"
    value = data.opslevel_tier.tier1.alias
  }
}

data "opslevel_services" "frontend" {
  filter {
    field = "owner"
    value = "frontend"
  }
}

output "all_services" {
  value = data.opslevel_services.all.names
}

output "tier1_services" {
  value = data.opslevel_services.tier1.names
}

output "frontend_services" {
  value = data.opslevel_services.frontend.urls
}