run "datasource_rubric_categories_all" {

  assert {
    condition     = length(data.opslevel_rubric_categories.all.rubric_categories) > 0
    error_message = "zero rubric_categories found in data.opslevel_rubric_categories"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_rubric_categories.all.rubric_categories[0].id),
      can(data.opslevel_rubric_categories.all.rubric_categories[0].name),
    ])
    error_message = "cannot set all expected rubric_category datasource fields"
  }

}

run "datasource_rubric_category_first" {

  assert {
    condition     = data.opslevel_rubric_category.first_category_by_id.id == data.opslevel_rubric_categories.all.rubric_categories[0].id
    error_message = "wrong ID on opslevel_rubric_category"
  }

  assert {
    condition     = data.opslevel_rubric_category.first_category_by_id.name == data.opslevel_rubric_categories.all.rubric_categories[0].name
    error_message = "wrong name on opslevel_rubric_category"
  }

}
