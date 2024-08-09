data "opslevel_service" "foo" {
  alias = "foo"
}

data "opslevel_service" "bar" {
  id = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS84Njcw"
}

data "opslevel_service_dependencies" "by_alias" {
  service = data.opslevel_service.foo.alias
}

data "opslevel_service_dependencies" "by_id" {
  service = data.opslevel_service.bar.id
}
