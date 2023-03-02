terraform {
  required_providers {
    msgraph = {
      source = "hsheppard/msgraph"
    }
  }
}

provider "msgraph" {
  # example configuration here
}

data "msgraph_users" "users" {}

output "my_users" {
	value = data.msgraph_users.users
}

