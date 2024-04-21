service = {
  name = "shortener-service"
  id = "shortener-service-${UUID}"
  port = 8080
  check = {
    id = "shortener-check"
    http = "http://${LOCAL_IP}:8080/health",
    method = "GET",
    interval = "30s",
    timeout = "5s"
  }
}