packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

source "amazon-ebs" "ubuntu-ebs" {
  ami_name      = "${var.ami_name} - ${local.timestamp}"
  subnet_id     = "${var.subnet_id}"
  instance_type = "${var.instance_type}"
  region        = "${var.ami_region}"
  profile       = "${var.aws_profile}"
  source_ami    = "${var.source_ami}"
  ssh_username  = "${var.ssh_username}"
  ami_users     = "${var.ami_shared_users}"
}

build {
  name    = "sid-ubuntu-24.04-lts-ami"
  sources = ["source.amazon-ebs.ubuntu-ebs"]

  provisioner "shell" {
    script = "./scripts/vm_setup.sh"
  }

  provisioner "shell" {
    script = "./scripts/createUser.sh"
  }

  provisioner "file" {
    source      = "../webapp"
    destination = "/tmp/webapp"
  }

  provisioner "file" {
    source      = "webapp.service"
    destination = "/tmp/webapp.service"
  }

  provisioner "shell" {
    script = "./scripts/binary.sh"
  }

  provisioner "shell" {
    script = "./scripts/startAppService.sh"
  }

}


variable "source_ami" {
  type    = string
  default = "ami-0866a3c8686eaeeba"
}

variable "ami_name" {
  type        = string
  description = "this is the name of the AMI"
  default="my_ami"
}

variable "instance_type" {
  type = string
  default="t2.micro"
}

variable "ami_region" {
  type = string
  default="us-east-1"
}

variable "aws_profile" {
  type = string
  default="dev"
}

variable "ssh_username" {
  type = string
  default="ubuntu"
}

variable "ami_shared_users" {
  type = list(string)
  default=["920372991622"]
}

variable "subnet_id" {
  type    = string
  default = "subnet-06ddfbabda19fc6b2"
}
