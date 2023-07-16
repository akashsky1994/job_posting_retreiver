
data "digitalocean_ssh_key" "ssh_key" {
  name = "jr_ssh_key"
}

data "digitalocean_reserved_ip" "myip" {
  ip_address = var.public_ip
}

resource "digitalocean_droplet" "jobretreiver" {
  image = "ubuntu-22-10-x64"
  name = "jobretreiver"
  region = "nyc3"
  size = "s-1vcpu-1gb"
  ssh_keys = [
    data.digitalocean_ssh_key.ssh_key.id
  ]
  connection {
    host = self.ipv4_address
    user = "root"
    type = "ssh"
    private_key = var.pvt_key
    timeout = "2m"
  }
  
  provisioner "file" {
    source      = "./external_scripts/setup_project.sh"
    destination = "/tmp/setup_project.sh"
  }
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      # install nginx
      "sudo apt-get update",
      "sudo apt install -y nginx",
      "chmod +x /tmp/setup_project.sh",
      "/tmp/setup_project.sh",
    ]
  }
}

# resource "digitalocean_reserved_ip" "jobretreiverip" {
#   droplet_id = digitalocean_droplet.jobretreiver.id
#   region = digitalocean_droplet.jobretreiver.region
# }

resource "digitalocean_reserved_ip_assignment" "jobretreiverip" {
  ip_address = data.digitalocean_reserved_ip.myip.ip_address
  droplet_id = digitalocean_droplet.jobretreiver.id
}

data "digitalocean_volume" "db_storage" {
  name   = "db-volume"
  region = "nyc3"
}
resource "digitalocean_volume_attachment" "db_volume" {
  droplet_id = digitalocean_droplet.jobretreiver.id
  volume_id  = data.digitalocean_volume.db_storage.id
}