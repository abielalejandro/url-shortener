#!/bin/sh
set -e

IP=$(ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}' | awk -F':' '{print $2}')
export LOCAL_IP=$IP
# generate random id
UUID=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)

envsubst < service.consul.hcl.tpl > service.consul.hcl
consul services register service.consul.hcl
# consul connect envoy -sidecar-for tgs-service

./app
