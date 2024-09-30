resource "opslevel_user" "john" {
  name               = "John Doe"
  email              = "john.doe@example.com"
  role               = "user" # or "admin"
  send_invite        = true
  skip_welcome_email = true
}
