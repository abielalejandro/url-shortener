{
  "address": "${LOCAL_IP}",
  "name": "tgs",
  "id": "tgs-${UUID}",
  "port": 9000,
  "checks": [
    {
      "name": "TGS grpc check",
      "tcp": "${LOCAL_IP}:9000",
      "interval": "10s",
      "timeout": "5s"
    }
  ]
}