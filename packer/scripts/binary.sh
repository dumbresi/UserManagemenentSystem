#!/bin/bash

sudo mkdir -p /usr/bin
sudo mv /tmp/webapp /usr/bin/webapp
sudo mv /tmp/.env /usr/bin/.env
sudo chown csye6225:csye6225 /usr/bin/webapp
sudo chmod +x /usr/bin/webapp
sudo mv /tmp/webapp.service /etc/systemd/system/webapp.service