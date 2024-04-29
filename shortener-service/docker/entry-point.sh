#!/bin/sh

export LOCAL_IP=$(ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}' | awk -F':' '{print $2}')
# generate random id
UUID=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
export UUID=$(echo "${LOCAL_IP//./-}")

function deregister_runner() {
  local ID="shortener-${UUID}"
  echo "finalizando contenedor ${ID}"
  curl --request PUT http://consul.service.consul:8500/v1/agent/service/deregister/$ID
}

trap "deregister_runner" EXIT
trap "deregister_runner" SIGINT
trap "deregister_runner" SIGTERM
trap "deregister_runner" SIGKILL
trap "deregister_runner" SIGQUIT

envsubst < service.consul.json.tpl > service.consul.json

if [[ -n "$API_TYPE" ]] && [[ "$API_TYPE" == "grpc" ]]; 
then  
  envsubst < service.consul.grpc.json.tpl > service.consul.json
fi

curl --request PUT --data @service.consul.json http://consul.service.consul:8500/v1/agent/service/register

./app & wait