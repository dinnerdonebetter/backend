# Dev Environment

## Setup

As much as was possible with my skills and the tools I had, I have made this thing entirely managed by Terraform. As with
all things infra/ops related, there are a few things that need to be done manually. Namely, creating the IAM principal
for Terraform to deploy with. Terraform will create a bunch of service accounts on initial deploy, but it won't be able 
to impersonate them unless you run some gcloud commands:

```bash
gcloud iam service-accounts add-iam-policy-binding outbound-emailer-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding data-changes-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding api-server@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding meal-plan-finalizer-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding mp-grocery-list-init-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding meal-plan-task-create-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding meal-plan-task-create-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding search-indexer-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding search-indexer-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding search-data-index-scheduler@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding email-prober@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
gcloud iam service-accounts add-iam-policy-binding webhook-executor-worker@dinner-done-better-dev.iam.gserviceaccount.com --member serviceAccount:terraform-cloud@dinner-done-better-dev.iam.gserviceaccount.com --role roles/iam.serviceAccountUser
```

You might get an error about not being the verified owner of a given domain. That's because you need to go to the Google Webmaster's admin interface thing and add the above terraform cloud service account as a verified owner.
