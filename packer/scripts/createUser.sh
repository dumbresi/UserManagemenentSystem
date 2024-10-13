#!/bin/bash

sudo groupadd csye6225
sudo useradd -r -g csye6225 csye6225 -s /usr/sbin/nologin 
# sudo usermod -a -G csye6225 csye6225