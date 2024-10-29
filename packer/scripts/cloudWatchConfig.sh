#!/bin/bash

cat <<EOF | sudo tee /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json
{
  "agent": {
      "metrics_collection_interval": 10,
      "logfile": "/var/logs/amazon-cloudwatch-agent.log"
  },
  "logs": {
      "logs_collected": {
          "files": {
              "collect_list": [
                  {
                      "file_path": "/var/log/tomcat9/csye6225.log",
                      "log_group_name": "csye6225",
                      "log_stream_name": "webapp"
                  }
              ]
          }
      }
  },
  "metrics": {
    "metrics_collected": {
       "statsd": {
          "service_address": ":8125",
          "metrics_collection_interval": 15,
          "metrics_aggregation_interval": 300
       }
    }
  }
}
EOF
