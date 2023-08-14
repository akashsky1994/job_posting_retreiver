
# Key setup and distribution
resource "tls_private_key" "jr_ssh_key" {
  algorithm = "RSA"
  rsa_bits = 4096
}

resource "local_file" "jr_pvt_key" {
  content = tls_private_key.jr_ssh_key.private_key_pem
  filename = "jr_ssh_key"
  file_permission = 0400
}

resource "digitalocean_ssh_key" "default" {
  name       = "jr_do_pvt_key"
  public_key = tls_private_key.jr_ssh_key.public_key_openssh
}