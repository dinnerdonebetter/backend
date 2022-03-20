module github.com/prixfixeco/api_server/cmd/functions/meal_plan_finalizer

// these have to be at 1.16 per Cloud Functions requirements: https://cloud.google.com/functions/docs/concepts/exec#runtimes
go 1.16

replace github.com/prixfixeco/api_server => ../../../

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.5.2
	github.com/prixfixeco/api_server v0.0.0-00010101000000-000000000000
)
