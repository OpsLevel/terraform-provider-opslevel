variables {
  rubric_category_one   = "opslevel_rubric_category"
  rubric_categories_all = "opslevel_rubric_categories"

  # required fields
  name = "TF Rubric Category"
}

run "resource_rubric_category_create_with_all_fields" {

  module {
    source = "./rubric_category"
  }

  assert {
    condition = alltrue([
      can(opslevel_rubric_category.test.id),
      can(opslevel_rubric_category.test.last_updated),
      can(opslevel_rubric_category.test.name),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.rubric_category_one)
  }

  assert {
    condition     = startswith(opslevel_rubric_category.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.rubric_category_one)
  }

  assert {
    condition     = opslevel_rubric_category.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.rubric_category_one)
  }

}

run "resource_rubric_category_update_all_fields" {

  variables {
    name = "${var.name} updated"
  }

  module {
    source = "./rubric_category"
  }

  assert {
    condition     = opslevel_rubric_category.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.rubric_category_one)
  }

}

run "datasource_rubric_categories_all" {

  module {
    source = "./rubric_category"
  }

  assert {
    condition     = can(data.opslevel_rubric_categories.all.rubric_categories)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.rubric_categories_all)
  }

  assert {
    condition     = length(data.opslevel_rubric_categories.all.rubric_categories) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.rubric_categories_all)
  }

}

run "datasource_rubric_category_first" {

  module {
    source = "./rubric_category"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_rubric_category.first_category_by_id.id),
      can(data.opslevel_rubric_category.first_category_by_id.name),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.rubric_category_one)
  }

  assert {
    condition     = data.opslevel_rubric_category.first_category_by_id.id == data.opslevel_rubric_categories.all.rubric_categories[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.rubric_category_one)
  }

  assert {
    condition     = data.opslevel_rubric_category.first_category_by_id.name == data.opslevel_rubric_categories.all.rubric_categories[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.rubric_category_one)
  }

}
