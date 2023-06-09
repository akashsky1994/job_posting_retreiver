variable "do_token" {
    description = "Digital Ocean Token"
    type = string
    sensitive = true
} #digital ocean personal token

variable "pvt_key" {
    description = "Private key"
    type = string
    sensitive = true
} # droplet private key location

variable "public_ip" {}