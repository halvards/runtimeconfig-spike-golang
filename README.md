# runtimeconfig-spike-golang

Sample code for storing and retrieving values from the
[Google Cloud Runtime Configuration API](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/)
v1beta1 using Go.

See
[halvards/runtimeconfig-spike-nodejs](https://github.com/halvards/runtimeconfig-spike-nodejs)
for a version using Node.js and TypeScript.

## Authentication

This sample code uses
[Application Default Credentials](https://cloud.google.com/docs/authentication/production)
to find the credentials used to authenticate to the API endpoints.

If your application runs in an environment where a default service account
isn't available (e.g., on your local machine), you can supply credentials
by using the `gcloud auth` command from the
[Google Cloud SDK](https://cloud.google.com/sdk/docs/)
or by using a service account file.

### Option 1: gcloud auth

Run this command to generate a credentials file in
`${HOME}/.config/gcloud/application_default_credentials.json`:

    gcloud auth application-default login

### Option 2: Service account file

Alternatively, you can supply a service account file. The service account
must have the "Project Editor" role to run this code sample, as this
is required to create a config. To do this, set the
`GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to the full
path to your service account credentials file:

    export GOOGLE_APPLICATION_CREDENTIALS="$(pwd)/service-account.json"

## Project ID

This sample code uses libraries that can auto-detect your project ID if you
supply a service account file or run the code in an environment that supports
auto-detecting the project ID (such as App Engine, Compute Engine,
Kubernetes Engine, or Cloud Functions. If need to supply a project ID, or if
you want to override the detected project ID you can set the
`GOOGLE_CLOUD_PROJECT` environment variable:

    export GOOGLE_CLOUD_PROJECT="$(gcloud config list --format 'value(core.project)')"

## References

- [Runtime Configurator Fundamentals](https://cloud.google.com/deployment-manager/runtime-configurator/)
- [Setting and Getting Data](https://cloud.google.com/deployment-manager/runtime-configurator/set-and-get-variables)
- [Google Cloud Runtime Configuration API](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/)
- [Cloud Identity and Access Management (IAM) Overview](https://cloud.google.com/iam/docs/overview)
- [Understanding Service Accounts](https://cloud.google.com/iam/docs/understanding-service-accounts)
- [Understanding Roles](https://cloud.google.com/iam/docs/understanding-roles)
- [Runtime Configurator Access Control Options](https://cloud.google.com/deployment-manager/runtime-configurator/access-control)
- [IAM Policy](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/v1beta1/Policy)
- [Google APIs Client Library for Go](https://github.com/google/google-api-go-client/blob/master/README.md)
- [OAuth2 for Go](https://github.com/golang/oauth2/blob/master/README.md)

## Disclaimer

This is not an officially supported Google product.
