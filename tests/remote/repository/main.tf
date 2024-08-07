data "opslevel_repositories" "all" {}

data "opslevel_repository" "first_repo_by_alias" {
  alias = data.opslevel_repositories.all.repositories[0].alias
}

data "opslevel_repository" "first_repo_by_id" {
  id = data.opslevel_repositories.all.repositories[0].id
}

# NOTE: placeholder until we can safely test this without overwrite actual data
#resource "opslevel_repository" "test" {
#  identifier = var.identifier
#  owner      = var.owner_id
#}
