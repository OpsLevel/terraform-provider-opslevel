run "datasource_rubric_categories_all" {

  variables {
    datasource_type = "opslevel_rubric_categories"
  }

  assert {
    condition     = can(data.opslevel_rubric_categories.all.rubric_categories)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_rubric_categories.all.rubric_categories) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_rubric_category_first" {

  variables {
    datasource_type = "opslevel_rubric_category"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_rubric_category.first_category_by_id.id),
      can(data.opslevel_rubric_category.first_category_by_id.name),
    ])
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_rubric_category.first_category_by_id.id == data.opslevel_rubric_categories.all.rubric_categories[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_rubric_category.first_category_by_id.name == data.opslevel_rubric_categories.all.rubric_categories[0].name
    error_message = replace(var.wrong_name_error, "TYPE", var.datasource_type)
  }

}
