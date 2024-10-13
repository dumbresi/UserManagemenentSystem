variable "ami_name" {
  type        = string
  description = "this is the name of the AMI"
  default     = "my-default-ami"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

variable "ami_region" {
  type    = string
  default = "us-east-2"
}
