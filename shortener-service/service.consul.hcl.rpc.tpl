service = {
  name = "shortener"
  id = "shortener-${UUID}"
  port = 9000
  check = {
    id = "shortener-check-rpc"
    tcp = "${LOCAL_IP}:9000",
    interval = "30s",
    timeout = "5s"
  }
}