resource "google_project_iam_custom_role" "meal_plan_finalizer_role" {
  role_id     = "meal_plan_finalizer_role"
  title       = "Meal Plan finalizer role"
  description = "An IAM role for the Meal Plan finalizer"
  permissions = [
    "secretmanager.versions.access",
    "cloudsql.instances.connect",
    "cloudsql.instances.get",
    "pubsub.topics.list",
    "pubsub.topics.publish",
    "pubsub.subscriptions.consume",
    "pubsub.subscriptions.create",
    "pubsub.subscriptions.delete",
  ]
}

locals {
  meal_plan_finalizer_database_username = "meal_plan_finalizer_db_user"
}

resource "google_pubsub_topic" "meal_plan_topic" {
  name = "meal_plan_finalization_work"
}

resource "google_cloud_scheduler_job" "meal_plan_finalization" {
  project = local.project_id
  region  = local.gcp_region
  name    = "meal-plan-finalizer-scheduler"

  schedule  = "* * * * *" # every minute
  time_zone = "America/Chicago"

  pubsub_target {
    topic_name = google_pubsub_topic.meal_plan_topic.id
    data       = base64encode("{}")
  }
}

resource "google_storage_bucket" "meal_plan_finalizer_bucket" {
  name     = "meal-plan-finalizer-cloud-function"
  location = "US"
}

data "archive_file" "meal_plan_finalizer_function" {
  type        = "zip"
  source_dir  = "${path.module}/meal_plan_finalizer_cloud_function"
  output_path = "${path.module}/meal_plan_finalizer_cloud_function.zip"
}

resource "google_storage_bucket_object" "meal_plan_finalizer_archive" {
  name   = format("meal_plan_finalizer_function-%s.zip", data.archive_file.meal_plan_finalizer_function.output_md5)
  bucket = google_storage_bucket.meal_plan_finalizer_bucket.name
  source = "${path.module}/meal_plan_finalizer_cloud_function.zip"
}

resource "google_service_account" "meal_plan_finalizer_user_service_account" {
  account_id   = "meal-plan-finalizer-worker"
  display_name = "Meal Plans Finalizer"
}

resource "google_project_iam_member" "meal_plan_finalizer_user" {
  project = local.project_id
  role    = google_project_iam_custom_role.meal_plan_finalizer_role.id
  member  = format("serviceAccount:%s", google_service_account.meal_plan_finalizer_user_service_account.email)
}

resource "random_password" "meal_plan_finalizer_user_database_password" {
  length           = 64
  special          = true
  override_special = "#$*-_=+[]"
}

resource "google_secret_manager_secret" "meal_plan_finalizer_user_database_password" {
  secret_id = "meal_plan_finalizer_user_database_password"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "meal_plan_finalizer_user_database_password" {
  secret = google_secret_manager_secret.meal_plan_finalizer_user_database_password.id

  secret_data = random_password.meal_plan_finalizer_user_database_password.result
}

resource "google_sql_user" "meal_plan_fimeal_plan_finalizer_user" {
  name     = local.meal_plan_finalizer_database_username
  instance = google_sql_database_instance.dev.name
  password = random_password.meal_plan_finalizer_user_database_password.result
}

resource "google_cloudfunctions_function" "meal_plan_finalizer" {
  name                = "meal-plan-finalizer"
  description         = "Meal Plan Finalizer"
  runtime             = local.go_runtime
  available_memory_mb = 128

  source_archive_bucket = google_storage_bucket.meal_plan_finalizer_bucket.name
  source_archive_object = google_storage_bucket_object.meal_plan_finalizer_archive.name
  service_account_email = google_service_account.meal_plan_finalizer_user_service_account.email

  entry_point = "FinalizeMealPlans"

  event_trigger {
    event_type = local.pubsub_topic_publish_event
    resource   = google_pubsub_topic.meal_plan_topic.name
  }

  environment_variables = {
    # TODO: use the meal_plan_finalizer_user for this, currently it has permission denied for accessing tables
    # https://dba.stackexchange.com/questions/53914/permission-denied-for-relation-table
    # https://www.postgresql.org/docs/13/sql-alterdefaultprivileges.html
    PRIXFIXE_DATABASE_USER                     = google_sql_user.api_user.name,
    PRIXFIXE_DATABASE_PASSWORD                 = random_password.api_user_database_password.result,
    PRIXFIXE_DATABASE_NAME                     = local.database_name,
    PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME = google_sql_database_instance.dev.connection_name,
    GOOGLE_CLOUD_SECRET_STORE_PREFIX           = format("projects/%d/secrets", data.google_project.project.number)
    GOOGLE_CLOUD_PROJECT_ID                    = data.google_project.project.project_id
  }

  #  secret_environment_variables = {
  #    key    = "PRIXFIXE_DATA_CHANGES_TOPIC_NAME"
  #    secret = google_secret_manager_secret.data_changes_topic_name.id
  #  }
  #
  #  secret_volumes = {
  #    mount_path = "/config/"
  #    secret     = google_secret_manager_secret.api_service_config.id
  #    versions = {
  #      path = "config.json"
  #    }
  #  }
}