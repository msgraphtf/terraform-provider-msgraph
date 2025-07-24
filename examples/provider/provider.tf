terraform {
  required_providers {
    msgraph = {
      source = "msgraphtf/msgraph"
    }
  }
}

provider "msgraph" {
  # Leave empty to use AzureCLI authentication
}

variable "test_user_password" {
  type = string
}

resource "msgraph_user" "test_user" {
  account_enabled = false
  display_name = "Test User"
  user_principal_name = "test_user@msgraphtf.onmicrosoft.com"
  mail_nickname = "test_user"
  password_profile = {
    password = var.test_user_password
  }
}

resource "msgraph_group" "test_group" {
  display_name = "test_group"
  mail_enabled = false
  mail_nickname = "test_group"
  security_enabled = true
}
