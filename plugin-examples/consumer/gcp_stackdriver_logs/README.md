
### GCP Stackdriver Logs

> it is a consumer plug-in for Ogi, it consumes logs from GCP Stackdriver and passes them on to Producer

#### Environment Variable Configuration

* `GOOGLE_PROJECT_ID` is GCP project-id

* `GOOGLE_LOG_TOPIC` is Stackdriver Log Topic from where events need to be read

* `GOOGLE_CLOUD_CREDENTIAL_FILE` is path to JSON credential file for Service Account from GCP project with access to Stackdriver

---
