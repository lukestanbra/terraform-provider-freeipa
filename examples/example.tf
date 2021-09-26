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

data "freeipa_user" "jdoe" {
  username = "jdoe4733520369981736867"
}

output "jdoe_user" {
  value = data.freeipa_user.jdoe
}

resource "freeipa_user" "jdoe" {
  username   = "jdoe3977187971214729826"
  first_name = "John"
  last_name  = "Doe"
  shell      = "/bin/bash"
}

resource "freeipa_user" "lukestanbra" {
  username   = "lukestanbra"
  first_name = "Luke"
  last_name  = "Stanbra"
  job_title  = "DevOps"
}