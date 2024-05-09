mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_tag" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_tag.example.id)
    error_message = "can use field id in opslevel_tag.example"
  }

  assert {
    condition     = opslevel_tag.example.resource_type == "Service"
    error_message = "value for field resource_type is bad opslevel_tag.example"

  }

  assert {
    condition     = opslevel_tag.example.resource_identifier == "test-service"
    error_message = "value for field resource_identifier is bad opslevel_tag.example"
  }

  assert {
    condition     = opslevel_tag.example.key == "yacht"
    error_message = "value for field key is bad opslevel_tag.example"
  }

  assert {
    condition     = opslevel_tag.example.value == "racing"
    error_message = "value for field value is bad opslevel_tag.example"
  }
}

