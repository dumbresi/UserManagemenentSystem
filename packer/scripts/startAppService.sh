#!/bin/bash

sudo systemctl daemon-reload

echo "daemon reload done"

sudo systemctl enable webapp.service

echo "enable webapp done"

sudo systemctl start webapp.service

echo "systemctl start webapp done"