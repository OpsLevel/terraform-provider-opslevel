variables {
  repository_one   = "opslevel_repository"
  repositories_all = "opslevel_repositories"

  # opslevel_repository fields
  identifier = "required"
  owner_id   = "optional"
}

run "datasource_repositories_all" {

  module {
    source = "./repository"
  }

  assert {
    condition     = can(data.opslevel_repositories.all.repositories)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.repositories_all)
  }

  assert {
    condition     = length(data.opslevel_repositories.all.repositories) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.repositories_all)
  }

}

run "datasource_repository_first" {

  module {
    source = "./repository"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_repository.first_repo_by_id.alias),
      can(data.opslevel_repository.first_repo_by_id.id),
      can(data.opslevel_repository.first_repo_by_id.languages),
      can(data.opslevel_repository.first_repo_by_id.name),
      can(data.opslevel_repository.first_repo_by_id.url),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.repository_one)
  }

  assert {
    condition     = data.opslevel_repository.first_repo_by_alias.alias == data.opslevel_repositories.all.repositories[0].alias
    error_message = replace(var.error_wrong_alias, "TYPE", var.repository_one)
  }

  assert {
    condition     = data.opslevel_repository.first_repo_by_id.id == data.opslevel_repositories.all.repositories[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.repository_one)
  }

}

# NOTE: "create" repository is really an update operation
# NOTE: not testing repository resource. No safe way to not overwrite actual data.

#run "resource_repository_update_set_all_fields" {
#
#  variables {
#    identifier = data.opslevel_repository.first_repo_by_id.id
#    owner_id   = run.get_owner_id.first_team.id # This would arbitrary overwrite actual owner data
#  }
#
#  module {
#    source = "./repository"
#  }
#
#  assert {
#    condition = alltrue([
#      can(opslevel_repository.test.id),
#      can(opslevel_repository.test.identifier),
#      can(opslevel_repository.test.owner),
#    ])
#    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.repository_one)
#  }
#
#}
