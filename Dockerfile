FROM ubuntu:20.04

RUN apt-get update && apt-get install -y \
    curl

RUN curl -sL \
    https://github.com/shashankm/kube-drain/releases/latest/download/kube-drain_0.2.0_Linux_x86_64.tar.gz \
    | tar -xzC /usr/local/bin kube-drain

COPY nodesfile /tmp/nodesfile
