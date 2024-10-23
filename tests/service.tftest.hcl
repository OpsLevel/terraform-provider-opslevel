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
  note                          = "TF Test Service Note"
  owner                         = null
  preferred_api_document_source = "PUSH"
  product                       = "widgets"
  tags                          = toset(["key1:value1", "key2:value2"])
  tier_alias                    = null
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_lifecycles.all,
      data.opslevel_teams.all,
      data.opslevel_tiers.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_service_create_with_all_fields" {

  variables {
    aliases                       = var.aliases
    api_document_path             = var.api_document_path
    description                   = var.description
    framework                     = var.framework
    language                      = var.language
    lifecycle_alias               = run.from_data_module.first_lifecycle.alias
    name                          = var.name
    note                          = var.note
    owner                         = run.from_data_module.first_team.id
    preferred_api_document_source = var.preferred_api_document_source
    product                       = var.product
    tags                          = var.tags
    tier_alias                    = run.from_data_module.first_tier.alias
  }

  module {
    source = "./opslevel_modules/modules/service"
  }


  assert {
    condition = alltrue([
      can(opslevel_service.this.aliases),
      can(opslevel_service.this.api_document_path),
      can(opslevel_service.this.description),
      can(opslevel_service.this.framework),
      can(opslevel_service.this.id),
      can(opslevel_service.this.language),
      can(opslevel_service.this.lifecycle_alias),
      can(opslevel_service.this.name),
      can(opslevel_service.this.note),
      can(opslevel_service.this.owner),
      can(opslevel_service.this.preferred_api_document_source),
      can(opslevel_service.this.product),
      can(opslevel_service.this.tags),
      can(opslevel_service.this.tier_alias),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.service_one)
  }

  assert {
    condition     = opslevel_service.this.aliases == toset(var.aliases)
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_service.this.aliases,
    )
  }

  assert {
    condition     = opslevel_service.this.api_document_path == var.api_document_path
    error_message = format(
      "expected '%v' but got '%v'",
      var.api_document_path,
      opslevel_service.this.api_document_path,
    )
  }

  assert {
    condition     = opslevel_service.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_service.this.description,
    )
  }

  assert {
    condition     = opslevel_service.this.framework == var.framework
    error_message = format(
      "expected '%v' but got '%v'",
      var.framework,
      opslevel_service.this.framework,
    )
  }

  assert {
    condition     = startswith(opslevel_service.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.service_one)
  }

  assert {
    condition     = opslevel_service.this.language == var.language
    error_message = format(
      "expected '%v' but got '%v'",
      var.language,
      opslevel_service.this.language,
    )
  }

  assert {
    condition     = opslevel_service.this.lifecycle_alias == var.lifecycle_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.lifecycle_alias,
      opslevel_service.this.lifecycle_alias,
    )
  }

  assert {
    condition     = opslevel_service.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_service.this.name,
    )
  }

  assert {
    condition     = opslevel_service.this.note == var.note
    error_message = format(
      "expected '%v' but got '%v'",
      var.note,
      opslevel_service.this.note,
    )
  }

  assert {
    condition     = opslevel_service.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_service.this.owner,
    )
  }

  assert {
    condition     = opslevel_service.this.preferred_api_document_source == var.preferred_api_document_source
    error_message = format(
      "expected '%v' but got '%v'",
      var.preferred_api_document_source,
      opslevel_service.this.preferred_api_document_source,
    )
  }

  assert {
    condition     = opslevel_service.this.product == var.product
    error_message = format(
      "expected '%v' but got '%v'",
      var.product,
      opslevel_service.this.product,
    )
  }

  assert {
    condition     = opslevel_service.this.tags == var.tags
    error_message = format(
      "expected '%v' but got '%v'",
      var.tags,
      opslevel_service.this.tags,
    )
  }

  assert {
    condition     = opslevel_service.this.tier_alias == var.tier_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.tier_alias,
      opslevel_service.this.tier_alias,
    )
  }

}

run "resource_service_create_with_empty_optional_fields" {

  variables {
    description = ""
    framework   = ""
    language    = ""
    name        = "New ${var.name} with empty fields"
    product     = ""
  }

  module {
    source = "./opslevel_modules/modules/service"
  }

  assert {
    condition     = opslevel_service.this.description == ""
    error_message = var.error_expected_empty_string
  }

  assert {
    condition     = opslevel_service.this.framework == ""
    error_message = var.error_expected_empty_string
  }

  assert {
    condition     = opslevel_service.this.language == ""
    error_message = var.error_expected_empty_string
  }

  assert {
    condition     = opslevel_service.this.product == ""
    error_message = var.error_expected_empty_string
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
    note                          = null
    owner                         = null
    preferred_api_document_source = null
    product                       = null
    tags                          = null
    tier_alias                    = null
  }

  module {
    source = "./opslevel_modules/modules/service"
  }

  assert {
    condition     = opslevel_service.this.aliases == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.api_document_path == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.framework == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.language == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.lifecycle_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.note == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.preferred_api_document_source == null
    error_message = "expected 'PUSH' default for preferred_api_document_source in opslevel_service resource"
  }

  assert {
    condition     = opslevel_service.this.product == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.tags == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_service.this.tier_alias == null
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
    lifecycle_alias               = run.from_data_module.first_lifecycle.alias
    name                          = "${var.name} updated"
    note                          = "${var.note} updated"
    owner                         = run.from_data_module.first_team.id
    preferred_api_document_source = var.preferred_api_document_source
    product                       = var.product
    tags                          = setunion(var.tags, ["key3:value3"])
    tier_alias                    = run.from_data_module.first_tier.alias
  }

  module {
    source = "./opslevel_modules/modules/service"
  }

  assert {
    condition     = opslevel_service.this.aliases == toset(var.aliases)
    error_message = format(
      "expected '%v' but got '%v'",
      var.aliases,
      opslevel_service.this.aliases,
    )
  }

  assert {
    condition     = opslevel_service.this.api_document_path == var.api_document_path
    error_message = format(
      "expected '%v' but got '%v'",
      var.api_document_path,
      opslevel_service.this.api_document_path,
    )
  }

  assert {
    condition     = opslevel_service.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_service.this.description,
    )
  }

  assert {
    condition     = opslevel_service.this.framework == var.framework
    error_message = format(
      "expected '%v' but got '%v'",
      var.framework,
      opslevel_service.this.framework,
    )
  }

  assert {
    condition     = opslevel_service.this.language == var.language
    error_message = format(
      "expected '%v' but got '%v'",
      var.language,
      opslevel_service.this.language,
    )
  }

  assert {
    condition     = opslevel_service.this.lifecycle_alias == var.lifecycle_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.lifecycle_alias,
      opslevel_service.this.lifecycle_alias,
    )
  }

  assert {
    condition     = opslevel_service.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_service.this.name,
    )
  }

  assert {
    condition     = opslevel_service.this.note == var.note
    error_message = format(
      "expected '%v' but got '%v'",
      var.note,
      opslevel_service.this.note,
    )
  }

  assert {
    condition     = opslevel_service.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_service.this.owner,
    )
  }

  assert {
    condition     = opslevel_service.this.preferred_api_document_source == var.preferred_api_document_source
    error_message = format(
      "expected '%v' but got '%v'",
      var.preferred_api_document_source,
      opslevel_service.this.preferred_api_document_source,
    )
  }

  assert {
    condition     = opslevel_service.this.product == var.product
    error_message = format(
      "expected '%v' but got '%v'",
      var.product,
      opslevel_service.this.product,
    )
  }

  assert {
    condition     = opslevel_service.this.tags == var.tags
    error_message = format(
      "expected '%v' but got '%v'",
      var.tags,
      opslevel_service.this.tags,
    )
  }

  assert {
    condition     = opslevel_service.this.tier_alias == var.tier_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.tier_alias,
      opslevel_service.this.tier_alias,
    )
  }

}
