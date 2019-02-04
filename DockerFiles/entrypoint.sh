#!/bin/bash
sed -i 's/{{ redis_host }}/127.0.0.1/g' /etc/redis/redis.conf
sed -i 's/{{ redis_memory }}/50M/g' /etc/redis/redis.conf
sed -i 's/{{ redis_host }}/127.0.0.1/g' /etc/gome/devices.json
sed -i 's/{{ slack_secret }}/testSecret/g' /etc/gome/secrets.json
sed -i 's/{{ aws_id  }}/testSecret/g' /etc/gome/secrets.json
sed -i 's/{{ aws_secret }}/testSecret/g' /etc/gome/secrets.json
sed -i 's/{{ aws_region }}/testSecret/g' /etc/gome/secrets.json
sed -i 's/{{ aws_token }}/testSecret/g' /etc/gome/secrets.json
sed -i 's/{{ aws_queue_url }}/testSecret/g' /etc/gome/secrets.json

service redis-server restart

cd /go/src/github.com/rebelit/gome
go build -o main .
cp main /etc/gome/

systemctl daemon-reload