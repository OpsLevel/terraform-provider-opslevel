data "opslevel_repository" "foo" {
    alias = "github.com:organization/example"
}

resource "opslevel_tag" "foo" {
  type = "Repository"
  identifier = data.opslevel_repository.foo.id

  key = "type"
  value = "frontend"
}

resource "opslevel_tag" "bar" {
  type = "Repository"
  identifier = "github.com:organization/example"

  key = "type"
  value = "backend"
}
