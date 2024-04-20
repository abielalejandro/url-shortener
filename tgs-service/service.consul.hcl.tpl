service = {
  name = "tgs-service"
  id = "tgs-service-${UUID}"
  port = 8080
  check = {
    id = "tgs-check"
    http = "http://${LOCAL_IP}:8080/health",
    method = "GET",
    interval = "30s",
    timeout = "5s"
  }
}