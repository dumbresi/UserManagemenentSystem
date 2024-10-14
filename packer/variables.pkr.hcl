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

variable "aws_profile" {
  type    = string
  default = "dev"
}

variable "source_ami" {
  type    = string
  default = "ami-0866a3c8686eaeeba"
}

variable "ssh_username" {
  type    = string
  default = "ubuntu"
}
