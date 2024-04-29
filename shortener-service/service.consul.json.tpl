{
  "address": "${LOCAL_IP}",
  "name": "shortener",
  "id": "shortener-${UUID}",
  "port": 8080,
  "checks": [
    {
      "name": "Shortener http check",
      "http": "http://${LOCAL_IP}:8080/health",
      "method": "GET",
      "header": { "Content-Type": ["application/json"] },
      "interval": "10s",
      "timeout": "5s"
    }
  ]
}