data "opslevel_rubric_levels" "all" {}

data "opslevel_rubric_level" "first_level_by_alias" {
  filter {
    field = "alias"
    value = data.opslevel_rubric_levels.all.rubric_levels[0].alias
  }
}

data "opslevel_rubric_level" "first_level_by_id" {
  filter {
    field = "id"
    value = data.opslevel_rubric_levels.all.rubric_levels[0].id
  }
}

data "opslevel_rubric_level" "first_level_by_name" {
  filter {
    field = "name"
    value = data.opslevel_rubric_levels.all.rubric_levels[0].name
  }
}

data "opslevel_rubric_level" "greatest_level_by_index" {
  filter {
    field = "index"
    value = max(data.opslevel_rubric_levels.all.rubric_levels[*].index...)
  }
}

resource "opslevel_rubric_level" "test" {
  description = var.description
  index       = var.index
  name        = var.name
}
