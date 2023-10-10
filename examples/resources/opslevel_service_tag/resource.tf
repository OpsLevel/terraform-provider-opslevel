data "opslevel_service" "foo" {
  alias = "foo"
}

resource "opslevel_tag" "foo" {
  type       = "Service"
  identifier = data.opslevel_service.foo.id

  key   = "environment"
  value = "foo"
}

resource "opslevel_tag" "bar" {
  type       = "Service"
  identifier = data.opslevel_service.foo.id

  key   = "environment"
  value = "bar"
}