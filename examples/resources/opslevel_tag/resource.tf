data "opslevel_repository" "foo" {
  alias = "github.com:organization/example"
}

resource "opslevel_tag" "foo_repo" {
  resource_type       = "Repository"
  resource_identifier = data.opslevel_repository.foo.id

  key   = "type"
  value = "frontend"
}

resource "opslevel_tag" "bar_repo" {
  resource_type       = "Repository"
  resource_identifier = "github.com:organization/example-2"

  key   = "type"
  value = "backend"
}

resource "opslevel_tag" "foo_domain" {
  resource_type       = "Domain"
  resource_identifier = "test-team"

  key   = "space"
  value = "craft"
}

resource "opslevel_tag" "foo_service" {
  resource_type       = "Service"
  resource_identifier = "test-service"

  key   = "yacht"
  value = "racing"
}

resource "opslevel_tag" "foo_system" {
  resource_type       = "System"
  resource_identifier = "test-system"

  key   = "crisp"
  value = "audio"
}

resource "opslevel_tag" "foo_team" {
  resource_type       = "Team"
  resource_identifier = "platform"

  key   = "goals"
  value = "automation"
}

resource "opslevel_tag" "foo_infra" {
  resource_type       = "InfrastructureResource"
  resource_identifier = "foo-id"

  key   = "type"
  value = "storage"
}


resource "opslevel_tag" "foo_user" {
  resource_type       = "User"
  resource_identifier = "user-id"

  key   = "role"
  value = "dev"
}

