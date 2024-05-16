resource "opslevel_user" "john" {
  name               = "John Doe"
  email              = "john.doe@example.com"
  role               = "user" # or "admin"
  skip_welcome_email = true
}
