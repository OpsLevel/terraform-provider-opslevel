data "opslevel_repository" "foo" {
  alias = "github.com:organization/example"
}

resource "opslevel_tag" "foo_repo" {
  resource_type = "Repository"
  resource_id   = data.opslevel_repository.foo.id

  key   = "type"
  value = "frontend"
}

resource "opslevel_tag" "bar_repo" {
  resource_type = "Repository"
  resource_id   = "github.com:organization/example-2"

  key   = "type"
  value = "backend"
}

resource "opslevel_tag" "foo_domain" {
  resource_type = "Domain"
  resource_id   = "test-team"

  key   = "space"
  value = "craft"
}

resource "opslevel_tag" "foo_service" {
  resource_type = "Service"
  resource_id   = "test-service"

  key   = "yacht"
  value = "racing"
}

resource "opslevel_tag" "foo_system" {
  resource_type = "System"
  resource_id   = "test-system"

  key   = "crisp"
  value = "audio"
}

resource "opslevel_tag" "foo_team" {
  resource_type = "Team"
  resource_id   = "platform"

  key   = "goals"
  value = "automation"
}
