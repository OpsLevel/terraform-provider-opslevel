data "opslevel_services" "all" {}

data "opslevel_service" "first_service_by_id" {
  id = data.opslevel_services.all.services[0].id
}

data "opslevel_service" "last_service" {
  id = element(data.opslevel_services.all.services, length(data.opslevel_services.all.services) - 1).id
}

resource "opslevel_service" "test" {
  aliases                       = var.aliases
  api_document_path             = var.api_document_path
  description                   = var.description
  framework                     = var.framework
  language                      = var.language
  lifecycle_alias               = var.lifecycle_alias
  name                          = var.name
  owner                         = var.owner
  preferred_api_document_source = var.preferred_api_document_source
  product                       = var.product
  tags                          = var.tags
  tier_alias                    = var.tier_alias
}
