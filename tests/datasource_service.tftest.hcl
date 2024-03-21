mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_service_given_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_alias.alias == "mock-service-alias"
    error_message = "alias in opslevel_service mock was not set"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_alias.id == null
    error_message = "expected null id in opslevel_service mock"
  }

}

run "datasource_service_given_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.alias == null
    error_message = "expected null alias in opslevel_service mock"
  }

  assert {
    condition     = startswith(data.opslevel_service.mock_service_with_id.id, "Z2lk")
    error_message = "wrong id prefix in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.id == "Z2lkOi8vmock123"
    error_message = "wrong id in opslevel_service mock"
  }

}

run "datasource_service_defaults" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.aliases == tolist(["alias-one", "alias-two"])
    error_message = "wrong aliases in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.api_document_path ==  "mock-api-document-path"
    error_message = "wrong api_document_path in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.description == "mock-service-description"
    error_message = "wrong description in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.framework == "mock-framework"
    error_message = "wrong framework in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.language == "mock-language"
    error_message = "wrong language in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.lifecycle_alias == "alpha"
    error_message = "wrong lifecycle_alias in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.name == "mock-service-name"
    error_message = "wrong name in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.owner == "mock-team"
    error_message = "wrong owner in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.owner_id == "Z2lkOi8vmockowner123"
    error_message = "wrong owner_id in opslevel_service mock"
  }

  assert {
    condition     = contains(["PUSH", "PULL"], data.opslevel_service.mock_service_with_id.preferred_api_document_source)
    error_message = "wrong preferred_api_document_source in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.properties == tolist(["prop-one", "prop-two"])
    error_message = "wrong properties list in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.repositories == tolist(["repo-one", "repo-two"])
    error_message = "wrong repositories list in opslevel_service mock"
  }

  assert {
    condition     = data.opslevel_service.mock_service_with_id.tags == tolist(["key1:value2", "key2:value2"])
    error_message = "wrong tags list in opslevel_service mock"
  }
}
