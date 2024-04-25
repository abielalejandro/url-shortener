service = {
  name = "shortener-web"
  id = "shortener-web-${UUID}"
  port = 8080
  check = {
    id = "shortener-web-check"
    http = "http://${LOCAL_IP}:8080/health",
    method = "GET",
    interval = "30s",
    timeout = "5s"
  }
}