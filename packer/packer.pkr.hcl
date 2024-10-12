packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "ubuntu-ebs" {
  ami_name      = "my-default-ami"
  instance_type = "t2.micro"
  region        = "us-east-1"
  profile       = "dev"
  source_ami    = "ami-0866a3c8686eaeeba"
  ssh_username  = "ubuntu"
}

build {
  name    = "sid-ubuntu-24.04-lts-ami"
  sources = ["source.amazon-ebs.ubuntu-ebs"]

  provisioner "shell" {
    script = "./scripts/vm_setup.sh"
  }

  provisioner "shell" {
    script = "./scripts/postgres_setup.sh"
  }

}

