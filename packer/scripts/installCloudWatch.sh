#!/bin/bash

# sudo apt-get install -y amazon-cloudwatch-agent
sudo apt-get update
sudo apt-get install -y wget
wget -qO - https://s3.amazonaws.com/amazoncloudwatch-agent/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb
sudo dpkg -i amazon-cloudwatch-agent.deb

