output "all" {
  value = data.opslevel_repositories.all
}

output "first" {
  value = data.opslevel_repositories.all.repositories[0]
}
