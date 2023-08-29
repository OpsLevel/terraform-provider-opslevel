data "opslevel_team" "devs" {
  alias = "devs"
}

resource "opslevel_secret" "my_secret" {
  alias = "secret-alias"
  owner = data.opslevel_team.devs.id
  value = "too_many_passwords"
}

resource "opslevel_secret" "my_secret_2" {
  alias = "secret-alias-2"
  owner = "devs"
  value = "0sd09wer0sdlkjwer90wer098sdfsewr"
}
