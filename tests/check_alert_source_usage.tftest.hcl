variables {
  check_alert_source_usage = "opslevel_check_alert_source_usage"

  # -- check_alert_source_usage fields --
  # required fields
  alert_type = "datadog"

  # optional fields
  alert_name_predicate = {
    type  = "exists"
    value = null
  }

  # -- check base fields --
  # required fields
  category = null
  level    = null
  name     = "TF Test Check Alert Source Usage"

  # optional fields
  enable_on = null
  enabled   = true
  filter    = null
  notes     = "Notes on Alert Source Usage Check"
  owner     = null
}

run "from_filter_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./opslevel_modules/modules/filter"
  }
}

run "from_rubric_category_module" {
  command = plan

  module {
    source = "./opslevel_modules/modules/rubric_category"
  }
}

run "from_rubric_level_module" {
  command = plan

  module {
    source = "./opslevel_modules/modules/rubric_level"
  }
}

run "from_team_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./opslevel_modules/modules/team"
  }
}

run "resource_check_alert_source_usage_create_with_all_fields" {

  variables {
    alert_type           = var.alert_type
    alert_name_predicate = var.alert_name_predicate
    category             = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on            = var.enable_on
    enabled              = var.enabled
    filter               = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/check/alert_source_usage"
  }

  assert {
    condition = alltrue([
      can(opslevel_check_alert_source_usage.this.category),
      can(opslevel_check_alert_source_usage.this.description),
      can(opslevel_check_alert_source_usage.this.enable_on),
      can(opslevel_check_alert_source_usage.this.enabled),
      can(opslevel_check_alert_source_usage.this.filter),
      can(opslevel_check_alert_source_usage.this.id),
      can(opslevel_check_alert_source_usage.this.level),
      can(opslevel_check_alert_source_usage.this.name),
      can(opslevel_check_alert_source_usage.this.notes),
      can(opslevel_check_alert_source_usage.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.check_alert_source_usage)
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.category == var.category
    error_message = "wrong category of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = startswith(opslevel_check_alert_source_usage.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.check_alert_source_usage)
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.level == var.level
    error_message = "wrong level ID of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_alert_source_usage)
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_alert_source_usage resource"
  }

}

run "resource_check_alert_source_usage_update_unset_optional_fields" {

  variables {
    alert_name_predicate = null
    category             = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on            = null
    enabled              = null
    filter               = null
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    notes = null
    owner = null
  }

  module {
    source = "./opslevel_modules/modules/check/alert_source_usage"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.alert_name_predicate == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enable_on == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enabled == false
    error_message = "expected 'false' default for 'enabled' in opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.filter == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.notes == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.owner == null
    error_message = var.error_expected_null_field
  }


}

run "resource_check_alert_source_usage_update_all_fields" {

  variables {
    alert_type           = var.alert_type
    alert_name_predicate = var.alert_name_predicate
    category             = run.from_rubric_category_module.all.rubric_categories[0].id
    enable_on            = var.enable_on
    enabled              = var.enabled
    filter               = run.from_filter_module.all.filters[0].id
    level = element([
      for lvl in run.from_rubric_level_module.all.rubric_levels :
      lvl.id if lvl.index == max(run.from_rubric_level_module.all.rubric_levels[*].index...)
    ], 0)
    name  = var.name
    notes = var.notes
    owner = run.from_team_module.all.teams[0].id
  }

  module {
    source = "./opslevel_modules/modules/check/alert_source_usage"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.category == var.category
    error_message = "wrong category of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enable_on == var.enable_on
    error_message = "wrong enable_on of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.enabled == var.enabled
    error_message = "wrong enabled of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.filter == var.filter
    error_message = "wrong filter ID of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.level == var.level
    error_message = "wrong level ID of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.check_alert_source_usage)
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.notes == var.notes
    error_message = "wrong notes of opslevel_check_alert_source_usage resource"
  }

  assert {
    condition     = opslevel_check_alert_source_usage.this.owner == var.owner
    error_message = "wrong owner ID of opslevel_check_alert_source_usage resource"
  }

}
