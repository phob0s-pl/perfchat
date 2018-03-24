provider "google" {
  credentials = "${file("~/.config/gcloud/terraform-admin.json")}"
  project     = "livechat-dev1234"
  region      = "europe-west1-b"
}