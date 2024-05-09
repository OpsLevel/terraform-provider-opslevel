run "datasource_repositories_all" {

  variables {
    datasource_type = "opslevel_repositories"
  }

  assert {
    condition     = can(data.opslevel_repositories.all.repositories)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_repositories.all.repositories) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_repository_first" {

  variables {
    datasource_type = "opslevel_repository"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_repository.first_repo_by_id.alias),
      can(data.opslevel_repository.first_repo_by_id.id),
      can(data.opslevel_repository.first_repo_by_id.languages),
      can(data.opslevel_repository.first_repo_by_id.name),
      can(data.opslevel_repository.first_repo_by_id.url),
    ])
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_repository.first_repo_by_alias.alias == data.opslevel_repositories.all.repositories[0].alias
    error_message = replace(var.wrong_alias_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_repository.first_repo_by_id.id == data.opslevel_repositories.all.repositories[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

}
