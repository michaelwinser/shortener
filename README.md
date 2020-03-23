# shortener
A URL shortener designed for Cloud Run, storing data in Cloud Firestore and logging usage to BigQuery.

## Use

Assuming you set up this service on `my.domain.com`:

* Post `{"Key": "shortlink", "URL": "www.url.com/xxx"}` to `my.domain.com/` to store a shortlink
* Redirects are stored in Cloud Firestore
* Visit `my.domain.com/shortlink` to be redirected to `www.url.com/xxx`
* Usage (along with source IP and User-Agent) is stored in BigQuery table `Redirect.Usage`

## TODO

* Serve a nice HTML UI on `GET /` to create and post JSON requests
* Some kind of authentication

## Setup

Sign up for Google Cloud, install the Cloud SDK, and create a project.

Enable the APIs required:

```sh
$ gcloud services enable cloudbuild.googleapis.com run.googleapis.com firestore.googleapis.com
Operation "operations/..." finished successfully.
```

Grant your Cloud Build service account permission to create and deploy a Cloud Run service, by visiting the [Cloud Build Service Account Permissions Page](https://console.cloud.google.com/cloud-build/settings/service-account) and enabling Cloud Run.  Accept the popup which asks if you want to enable Service Account User permissions too.

Grant your Cloud Run service account permission to read/write to Firestore and BigQuery, by visiting the [IAM Settings Page](https://console.cloud.google.com/iam-admin/iam) and adding roles `BigQuery Data Owner` and `Cloud Datastore User` to your default Compute service account (`...-compute@developer.gserviceaccount.com`).

Configure a [Cloud Build trigger](https://console.cloud.google.com/cloud-build/triggers) to run on push to master using `cloudbuild.yaml`, then run the trigger.  (A source SHA is required to tag the container.)  This will build the container and deploy to Cloud Run.

To run from a custom domain like `my.domain.com`, visit the [Cloud Run Domain Mappings page](https://console.cloud.google.com/run/domains) and add a custom domain mapping.  Follow instructions for your domain registrar.