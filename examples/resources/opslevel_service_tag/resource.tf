data "opslevel_service" "foo" {
  alias = "foo"
}

resource "opslevel_service_tag" "foo_environment" {
  service = data.opslevel_service.foo.id

  key   = "environment"
  value = "production"
}

resource "opslevel_service_tool" "bar_environment" {
  service_alias = "bar"

  key   = "environment"
  value = "production"
}