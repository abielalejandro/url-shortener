service = {
  name = "shortener-service-rpc"
  id = "shortener-service-${UUID}"
  port = 9000
  check = {
    id = "shortener-check-rpc"
    tcp = "${LOCAL_IP}:9000",
    interval = "30s",
    timeout = "5s"
  }
}