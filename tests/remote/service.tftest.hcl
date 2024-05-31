variables {
  service_one  = "opslevel_service"
  services_all = "opslevel_services"

  # required fields
  name = "TF Test Service"

  # lifecycle_alias fields
  aliases                       = ["service_one"]
  api_document_path             = "test.json"
  description                   = "Service description"
  framework                     = "Ruby on Rails"
  language                      = "ruby"
  lifecycle_alias               = null
  owner                         = null
  preferred_api_document_source = "PUSH"
  product                       = "widgets"
  tags                          = tolist(["key1:value1", "key2:value2"])
  tier_alias                    = null
}

run "from_lifecycle_get_lifecycle_alias" {
  command = plan

  module {
    source = "./lifecycle"
  }
}

run "from_team_get_owner_id" {
  command = plan

  variables {
    aliases          = null
    name             = ""
    parent           = null
    responsibilities = null
  }

  module {
    source = "./team"
  }
}

run "from_tier_get_tier_alias" {
  command = plan

  module {
    source = "./tier"
  }
}

run "resource_service_create_with_all_fields" {

  variables {
    aliases                       = var.aliases
    api_document_path             = var.api_document_path
    description                   = var.description
    framework                     = var.framework
    language                      = var.language
    lifecycle_alias               = run.from_lifecycle_get_lifecycle_alias.first_lifecycle.alias
    name                          = var.name
    owner                         = run.from_team_get_owner_id.first_team.id
    preferred_api_document_source = var.preferred_api_document_source
    product                       = var.product
    tags                          = var.tags
    tier_alias                    = run.from_tier_get_tier_alias.first_tier.alias
  }

  module {
    source = "./service"
  }


  assert {
    condition = alltrue([
      can(opslevel_service.test.aliases),
      can(opslevel_service.test.api_document_path),
      can(opslevel_service.test.description),
      can(opslevel_service.test.framework),
      can(opslevel_service.test.id),
      can(opslevel_service.test.language),
      can(opslevel_service.test.lifecycle_alias),
      can(opslevel_service.test.owner),
      can(opslevel_service.test.preferred_api_document_source),
      can(opslevel_service.test.product),
      can(opslevel_service.test.tags),
      can(opslevel_service.test.tier_alias),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.service_one)
  }

  assert {
    condition     = opslevel_service.test.aliases == toset(var.aliases)
    error_message = "wrong aliases of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.api_document_path == var.api_document_path
    error_message = "wrong api_document_path of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.service_one)
  }

  assert {
    condition     = opslevel_service.test.framework == var.framework
    error_message = "wrong framework of opslevel_service resource"
  }

  assert {
    condition     = startswith(opslevel_service.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.service_one)
  }

  assert {
    condition     = opslevel_service.test.language == var.language
    error_message = "wrong language of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.lifecycle_alias == var.lifecycle_alias
    error_message = "wrong lifecycle_alias of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.owner == var.owner
    error_message = "wrong owner of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.preferred_api_document_source == var.preferred_api_document_source
    error_message = "wrong preferred_api_document_source of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.product == var.product
    error_message = "wrong product of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.tags == var.tags
    error_message = "wrong tags of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.tier_alias == var.tier_alias
    error_message = "wrong tier_alias of opslevel_service resource"
  }

}

run "resource_service_update_unset_optional_fields" {

  variables {
    aliases                       = null
    api_document_path             = null
    description                   = null
    framework                     = null
    language                      = null
    lifecycle_alias               = null
    owner                         = null
    preferred_api_document_source = null
    product                       = null
    tags                          = null
    tier_alias                    = null
  }

  module {
    source = "./service"
  }

  assert {
    condition     = opslevel_service.test.aliases == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.api_document_path == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.framework == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.language == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.lifecycle_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.preferred_api_document_source == null
    error_message = "expected 'PUSH' default for preferred_api_document_source in opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.product == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.tags == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.test.tier_alias == null
    error_message = var.error_expected_null_field
  }

}

run "resource_service_update_set_all_fields" {

  variables {
    aliases                       = setunion(var.aliases, ["test_alias"])
    api_document_path             = var.api_document_path
    description                   = "${var.description} updated"
    framework                     = upper(var.framework)
    language                      = upper(var.language)
    lifecycle_alias               = run.from_lifecycle_get_lifecycle_alias.first_lifecycle.alias
    name                          = "${var.name} updated"
    owner                         = run.from_team_get_owner_id.first_team.id
    preferred_api_document_source = var.preferred_api_document_source
    product                       = var.product
    tags                          = concat(var.tags, ["key3:value3"])
    tier_alias                    = run.from_tier_get_tier_alias.first_tier.alias
  }

  module {
    source = "./service"
  }

  assert {
    condition     = opslevel_service.test.aliases == toset(var.aliases)
    error_message = "wrong aliases of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.api_document_path == var.api_document_path
    error_message = "wrong api_document_path of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.service_one)
  }

  assert {
    condition     = opslevel_service.test.framework == var.framework
    error_message = "wrong framework of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.language == var.language
    error_message = "wrong language of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.lifecycle_alias == var.lifecycle_alias
    error_message = "wrong lifecycle_alias of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.owner == var.owner
    error_message = "wrong owner of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.preferred_api_document_source == var.preferred_api_document_source
    error_message = "wrong preferred_api_document_source of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.product == var.product
    error_message = "wrong product of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.tags == var.tags
    error_message = "wrong tags of opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.test.tier_alias == var.tier_alias
    error_message = "wrong tier_alias of opslevel_service resource"
  }

}

run "datasource_services_all" {

  module {
    source = "./service"
  }

  assert {
    condition     = can(data.opslevel_services.all.services)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.services_all)
  }

  assert {
    condition     = length(data.opslevel_services.all.services) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.services_all)
  }

  assert {
    condition = alltrue([
      can(data.opslevel_services.all.services[0].id),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.services_all)
  }

}

run "datasource_service_first" {

  module {
    source = "./service"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_service.first_service_by_id.alias),
      can(data.opslevel_service.first_service_by_id.aliases),
      can(data.opslevel_service.first_service_by_id.api_document_path),
      can(data.opslevel_service.first_service_by_id.description),
      can(data.opslevel_service.first_service_by_id.framework),
      can(data.opslevel_service.first_service_by_id.id),
      can(data.opslevel_service.first_service_by_id.language),
      can(data.opslevel_service.first_service_by_id.lifecycle_alias),
      can(data.opslevel_service.first_service_by_id.name),
      can(data.opslevel_service.first_service_by_id.owner),
      can(data.opslevel_service.first_service_by_id.owner_id),
      can(data.opslevel_service.first_service_by_id.preferred_api_document_source),
      can(data.opslevel_service.first_service_by_id.product),
      can(data.opslevel_service.first_service_by_id.properties),
      can(data.opslevel_service.first_service_by_id.repositories),
      can(data.opslevel_service.first_service_by_id.tags),
      can(data.opslevel_service.first_service_by_id.tier_alias),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.service_one)
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.id == data.opslevel_services.all.services[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.service_one)
  }

}
