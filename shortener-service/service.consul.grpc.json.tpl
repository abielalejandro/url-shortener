{
  "address": "${LOCAL_IP}",
  "name": "shortener",
  "id": "shortener-${UUID}",
  "port": 9000,
  "checks": [
    {
      "name": "Shortener grpc check",
      "tcp": "${LOCAL_IP}:9000",
      "interval": "10s",
      "timeout": "5s"
    }
  ]
}