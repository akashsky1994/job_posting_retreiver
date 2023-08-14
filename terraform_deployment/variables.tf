variable "do_token" {
    description = "Digital Ocean Token"
    type = string
    sensitive = true
} #digital ocean personal token

variable "public_ip" {
    description = "Public IP"
    type = string
}