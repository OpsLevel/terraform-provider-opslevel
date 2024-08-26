resource "opslevel_secret" "test" {
  alias = var.alias
  owner = var.owner
  value = var.value
}
