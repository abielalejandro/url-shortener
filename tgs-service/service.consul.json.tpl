{
  "address": "${LOCAL_IP}",
  "name": "tgs",
  "id": "tgs-${UUID}",
  "port": 8080,
  "checks": [
    {
      "name": "TGS http check",
      "http": "http://${LOCAL_IP}:8080/health",
      "method": "GET",
      "header": { "Content-Type": ["application/json"] },
      "interval": "10s",
      "timeout": "5s"
    }
  ]
}