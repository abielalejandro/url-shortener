service = {
  name = "tgs-service-rpc"
  id = "tgs-service-${UUID}"
  port = 9000
  check = {
    id = "tgs-check-rpc"
    tcp = "${LOCAL_IP}:9000",
    method = "GET",
    interval = "30s",
    timeout = "5s"
  }
}