FROM golang:1.12beta2-stretch

RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y redis-server git

##redis
RUN rm -rf /etc/redis/redis.conf
COPY roles /etc/redis/redis.conf

##gome
RUN mkdir -p /etc/gome
RUN mkdir -p /go/src/github.com/rebelit/gome
COPY roles /etc/gome/devices.json
COPY roles /etc/gome/secrets.json
#ADD DockerFiles/entrypoint.sh /entrypoint.sh

#ENTRYPOINT /entrypoint.sh
#CMD "service gome-server restart"