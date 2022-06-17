resource "opslevel_user" "john" {
  name  = "John Doe"
  email = "john.doe@example.com"
  role  = "user"
}

resource "opslevel_user" "ken" {
  name  = "Ken Doe"
  email = "ken.doe@example.com"
  role  = "admin"
}