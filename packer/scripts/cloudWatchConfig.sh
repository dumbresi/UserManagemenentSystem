#!/bin/bash

sudo mv /tmp/cloudwatch-config.json /opt/cloudwatch-config.json
sudo mkdir -p /var/log/webapp/
sudo chown csye6225:csye6225 /var/log/webapp/