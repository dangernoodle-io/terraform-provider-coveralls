terraform {
  required_providers {
    coveralls = {
      source = "dangernoodle-io/coveralls"
    }
  }
}

provider "coveralls" {
  token = "coveralls-api-token"
}
