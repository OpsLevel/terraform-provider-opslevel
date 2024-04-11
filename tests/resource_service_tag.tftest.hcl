mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_service_tag_using_service_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service_tag.using_service_id.key == "hello_with_id"
    error_message = "wrong service tag key"
  }

  assert {
    condition     = opslevel_service_tag.using_service_id.value == "world_with_id"
    error_message = "wrong service tag value"
  }

  assert {
    condition     = opslevel_service_tag.using_service_id.service == "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNzQxMg"
    error_message = "expected service identifier to be an id"
  }

  assert {
    condition     = can(opslevel_service_tag.using_service_id.id)
    error_message = "expected service tag to have an ID"
  }
}

run "resource_service_tag_using_service_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_service_tag.using_service_alias.key == "hello_with_alias"
    error_message = "wrong service tag key"
  }

  assert {
    condition     = opslevel_service_tag.using_service_alias.value == "world_with_alias"
    error_message = "wrong service tag value"
  }

  assert {
    condition     = opslevel_service_tag.using_service_alias.service_alias == "cart"
    error_message = "expected service identifier to be an alias"
  }

  assert {
    condition     = can(opslevel_service_tag.using_service_alias.id)
    error_message = "expected service tag to have an ID"
  }
}