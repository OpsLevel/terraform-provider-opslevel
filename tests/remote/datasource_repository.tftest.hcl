run "datasource_repositories_all" {

  assert {
    condition     = length(data.opslevel_repositories.all.repositories) > 0
    error_message = "zero Repositories found in data.opslevel_repositories"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_repositories.all.repositories[0].alias),
      can(data.opslevel_repositories.all.repositories[0].id),
      can(data.opslevel_repositories.all.repositories[0].languages),
      can(data.opslevel_repositories.all.repositories[0].name),
      can(data.opslevel_repositories.all.repositories[0].url),
    ])
    error_message = "cannot set all expected Repository datasource fields"
  }

}

run "datasource_repository_first" {

  assert {
    condition     = data.opslevel_repository.first_repo_by_alias.alias == data.opslevel_repositories.all.repositories[0].alias
    error_message = "wrong alias on first opslevel_repository"
  }

  assert {
    condition     = data.opslevel_repository.first_repo_by_id.id == data.opslevel_repositories.all.repositories[0].id
    error_message = "wrong ID on first opslevel_repository"
  }

}
