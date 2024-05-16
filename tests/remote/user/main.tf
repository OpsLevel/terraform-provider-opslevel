data "opslevel_users" "all" {}

data "opslevel_user" "first_user_by_email" {
  identifier = data.opslevel_users.all.users[0].email
}

data "opslevel_user" "first_user_by_id" {
  identifier = data.opslevel_users.all.users[0].id
}

resource "opslevel_user" "test" {
  email              = var.email
  name               = var.name
  role               = var.role
  skip_welcome_email = var.skip_welcome_email
}
