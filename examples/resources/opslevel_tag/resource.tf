data "opslevel_repository" "foo" {
  alias = "github.com:organization/example"
}

resource "opslevel_tag" "foo" {
  resource_type = "Repository"
  resource_id   = data.opslevel_repository.foo.id

  key   = "type"
  value = "frontend"
}

resource "opslevel_tag" "bar" {
  resource_type = "Repository"
  resource_id   = "github.com:organization/example-2"

  key   = "type"
  value = "backend"
}
