
run "resource_domain" {

  module {
    source = "./domain"
  }

  assert {
    condition     = opslevel_domain.required_fields.name == "Test - name only"
    error_message = "wrong name for opslevel_domain"
  }

  assert {
    condition     = opslevel_domain.required_fields.id != null && opslevel_domain.required_fields.id != ""
    error_message = "opslevel_domain id should not be empty"
  }

  #assert {
  #  condition     = opslevel_domain.required_fields.owner == "Developers"
  #  error_message = "wrong owner of opslevel_domain resource"
  #}
}


