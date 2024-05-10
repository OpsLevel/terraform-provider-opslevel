data "opslevel_lifecycles" "all" {}

data "opslevel_lifecycle" "first_lifecycle_by_alias" {
  filter {
    field = "alias"
    value = data.opslevel_lifecycles.all.lifecycles[0].alias
  }
}

data "opslevel_lifecycle" "first_lifecycle_by_id" {
  filter {
    field = "id"
    value = data.opslevel_lifecycles.all.lifecycles[0].id
  }
}

data "opslevel_lifecycle" "first_lifecycle_by_index" {
  filter {
    field = "index"
    value = data.opslevel_lifecycles.all.lifecycles[0].index
  }
}

data "opslevel_lifecycle" "first_lifecycle_by_name" {
  filter {
    field = "name"
    value = data.opslevel_lifecycles.all.lifecycles[0].name
  }
}
