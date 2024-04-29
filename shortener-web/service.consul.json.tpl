{
  "address": "${LOCAL_IP}",
  "name": "shortener-web",
  "id": "shortener-web-${UUID}",
  "port": 8080,
  "checks": [
    {
      "name": "Shortener web http check",
      "http": "http://${LOCAL_IP}:8080/health",
      "method": "GET",
      "header": { "Content-Type": ["application/json"] },
      "interval": "10s",
      "timeout": "5s"
    }
  ]
}