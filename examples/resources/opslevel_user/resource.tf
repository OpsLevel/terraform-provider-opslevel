resource "opslevel_user" "john" {
  force_send_invite  = true
  name               = "John Doe"
  email              = "john.doe@example.com"
  role               = "user" # or "admin"
  skip_welcome_email = true
}
