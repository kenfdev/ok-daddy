#!/bin/bash

curl https://releases.rancher.com/install-docker/17.03.sh | sh

usermod -aG docker ubuntu

mkdir -p /home/ubuntu/.ssh
chown ubuntu:ubuntu /home/ubuntu/.ssh

swapoff -a