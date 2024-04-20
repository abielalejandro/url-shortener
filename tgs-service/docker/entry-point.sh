#!/bin/sh
set -e

export LOCAL_IP=$(ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}' | awk -F':' '{print $2}')
# generate random id
UUID=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
export UUID=$LOCAL_IP

envsubst < service.consul.hcl.tpl > /etc/consul.d/service.consul.hcl

if [[ -n "$API_TYPE" ]] || [[ "$API_TYPE"=="rpc" ]]; 
then  
  envsubst < service.consul.hcl.rpc.tpl > /etc/consul.d/service.consul.hcl
fi

consul agent -config-dir=/etc/consul.d/ > /tmp/consul-server.log 2>&1 &

iptables -t nat -A PREROUTING -p udp -m udp --dport 53 -j REDIRECT --to-ports 8600
iptables -t nat -A PREROUTING -p tcp -m tcp --dport 53 -j REDIRECT --to-ports 8600
iptables -t nat -A OUTPUT -d localhost -p udp -m udp --dport 53 -j REDIRECT --to-ports 8600
iptables -t nat -A OUTPUT -d localhost -p tcp -m tcp --dport 53 -j REDIRECT --to-ports 8600

./app
