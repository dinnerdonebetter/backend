resource "google_project_service" "cloud_run_api" {
  service = "run.googleapis.com"
}
