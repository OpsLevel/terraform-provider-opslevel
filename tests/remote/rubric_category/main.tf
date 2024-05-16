data "opslevel_rubric_categories" "all" {}

data "opslevel_rubric_category" "first_category_by_id" {
  filter {
    field = "id"
    value = data.opslevel_rubric_categories.all.rubric_categories[0].id
  }
}

data "opslevel_rubric_category" "first_category_by_name" {
  filter {
    field = "name"
    value = data.opslevel_rubric_categories.all.rubric_categories[0].name
  }
}

resource "opslevel_rubric_category" "test" {
  name = var.name
}
