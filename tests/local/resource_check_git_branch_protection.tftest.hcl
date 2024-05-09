mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_git_branch_protection" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_git_branch_protection.example.name == "foo"
    error_message = "wrong value name for opslevel_check_git_branch_protection.example"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.example.enabled == true
    error_message = "wrong value enabled on opslevel_check_git_branch_protection.example"
  }

  assert {
    condition     = can(opslevel_check_git_branch_protection.example.id)
    error_message = "id attribute missing from in opslevel_check_git_branch_protection.example"
  }

  assert {
    condition     = can(opslevel_check_git_branch_protection.example.owner)
    error_message = "owner attribute missing from in opslevel_check_git_branch_protection.example"
  }

  assert {
    condition     = can(opslevel_check_git_branch_protection.example.filter)
    error_message = "filter attribute missing from in opslevel_check_git_branch_protection.example"
  }

  assert {
    condition     = can(opslevel_check_git_branch_protection.example.category)
    error_message = "category attribute missing from in opslevel_check_git_branch_protection.example"
  }

  assert {
    condition     = can(opslevel_check_git_branch_protection.example.level)
    error_message = "level attribute missing from in opslevel_check_git_branch_protection.example"
  }

  assert {
    condition     = opslevel_check_git_branch_protection.example.notes == null
    error_message = "wrong value for notes in opslevel_check_git_branch_protection.example"
  }
}