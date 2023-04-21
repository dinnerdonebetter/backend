# API config

resource "google_secret_manager_secret" "api_service_config" {
  secret_id = "api_service_config"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "api_service_config" {
  secret = google_secret_manager_secret.api_service_config.id

  secret_data = file("${path.module}/service-config.json")
}

# Data changes topic

resource "google_secret_manager_secret" "data_changes_topic_name" {
  secret_id = "data_changes_topic_name"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "outbound_emails_topic_name" {
  secret_id = "outbound_emails_topic_name"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "data_changes_topic_name" {
  secret = google_secret_manager_secret.data_changes_topic_name.id

  secret_data = google_pubsub_topic.data_changes_topic.name
}

resource "google_secret_manager_secret_version" "outbound_emails_topic_name" {
  secret = google_secret_manager_secret.outbound_emails_topic_name.id

  secret_data = google_pubsub_topic.outbound_emails_topic.name
}

# API server cookie hash key

resource "random_string" "cookie_hash_key" {
  length  = 64
  special = false
}

resource "google_secret_manager_secret" "cookie_hash_key" {
  secret_id = "cookie_hash_key"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "cookie_hash_key" {
  secret = google_secret_manager_secret.cookie_hash_key.id

  secret_data = random_string.cookie_hash_key.result
}

# API server cookie block key

resource "random_string" "cookie_block_key" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "cookie_block_key" {
  secret_id = "cookie_block_key"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "cookie_block_key" {
  secret = google_secret_manager_secret.cookie_block_key.id

  secret_data = random_string.cookie_block_key.result
}

# API server PASETO local mode key

resource "random_string" "paseto_local_key" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "paseto_local_key" {
  secret_id = "paseto_local_key"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "paseto_local_key" {
  secret = google_secret_manager_secret.paseto_local_key.id

  secret_data = random_string.paseto_local_key.result
}

# External API services

variable "SEGMENT_API_TOKEN" {}
variable "SENDGRID_API_TOKEN" {}

# Sendgrid token

resource "google_secret_manager_secret" "sendgrid_api_token" {
  secret_id = "sendgrid_api_token"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "sendgrid_api_token" {
  secret = google_secret_manager_secret.sendgrid_api_token.id

  secret_data = var.SENDGRID_API_TOKEN
}

# Segment API token

resource "google_secret_manager_secret" "segment_api_token" {
  secret_id = "segment_api_token"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "segment_api_token" {
  secret = google_secret_manager_secret.segment_api_token.id

  secret_data = var.SEGMENT_API_TOKEN
}
