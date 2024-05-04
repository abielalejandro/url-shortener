service = {
  name = "tgs"
  id = "tgs-service-${UUID}"
  port = 9000
  check = {
    id = "tgs-check-rpc"
    tcp = "${LOCAL_IP}:9000",
    interval = "30s",
    timeout = "5s"
  }
}