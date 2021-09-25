terraform {
  required_providers {
    freeipa = {
      source = "hashicorp.com/lukestanbra/freeipa"
    }
  }
}

provider "freeipa" {
  username = "admin"
  password = "password"
  host     = "ipa.example.test"
}

data "freeipa_user" "admin" {
  username = "admin"
}

output "admin_user" {
  value = data.freeipa_user.admin
}

/* resource "freeipa_user" "lukestanbra" {
  username = "lukestanbra"
  firstname = "Puke"
  lastname = "Stanbra"
} */

resource "freeipa_user" "jdoe" {
  username  = "jdoe3977187971214729826"
  firstname = "John"
  lastname  = "Doe"
}