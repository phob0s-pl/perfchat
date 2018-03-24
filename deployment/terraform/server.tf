resource "google_compute_instance" "default" {
  count        = "1"
  name         = "aserver"
  machine_type = "g1-small"
  zone         = "europe-west1-b"

  project = "livechat-dev1234"

boot_disk {
   initialize_params {
     image = "debian-cloud/debian-9"
   }
 }

  metadata {
    ssh-keys = "jakubj00:${file("~/.ssh/gcpx.pub")}"
    hostname = "aserver"
  }

  network_interface {
    network = "default"

    access_config {}
  }
}