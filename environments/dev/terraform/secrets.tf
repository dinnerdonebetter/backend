variable "SEGMENT_API_TOKEN" {}
variable "SENDGRID_API_TOKEN" {}

resource "random_string" "cookie_hash_key" {
  length  = 64
  special = false
}

resource "random_string" "cookie_block_key" {
  length  = 32
  special = false
}

resource "random_string" "paseto_local_key" {
  length  = 32
  special = false
}

resource "kubernetes_secret_v1" "config_auth" {
  metadata {
    namespace = local.kubernetes_namespace
    name      = "config.auth"
  }

  data = {
    "cookie_hash_key"  = random_string.cookie_hash_key.result
    "cookie_block_key" = random_string.cookie_block_key.result
    "paseto_local_key" = random_string.paseto_local_key.result
  }
}

resource "kubernetes_secret_v1" "config_sendgrid" {
  metadata {
    namespace = local.kubernetes_namespace
    name      = "config.third-party.sendgrid"
  }

  data = {
    "sendgrid_api_token" = var.SENDGRID_API_TOKEN
  }
}

resource "kubernetes_secret_v1" "config_segment" {
  metadata {
    namespace = local.kubernetes_namespace
    name      = "config.third-party.segment"
  }

  data = {
    "segment_api_token" = var.SEGMENT_API_TOKEN
  }
}