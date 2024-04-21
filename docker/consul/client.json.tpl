{
    "datacenter": "dc1",
    "node_name": "consul-client-${UUID}",
    "data_dir": "/opt/consul",
    "leave_on_terminate": true,
    "log_level": "INFO",
    "rejoin_after_leave": true,
    "disable_keyring_file": true,
    "start_join": ["10.50.0.2"],
    "retry_join": ["10.50.0.2"],
    "server": false,
    "connect": {
      "enabled": true
    },
    "ports": {
      "grpc": 8502
    } ,
    "service": {
      "id": "dns",
      "name": "dns",
      "tags": ["primary"],
      "address": "localhost",
      "port": 8600,
      "enable_tag_override": false,
      "check": {
        "id": "dns",
        "name": "Consul DNS TCP on port 8600",
        "tcp": "localhost:8600",
        "interval": "10s",
        "timeout": "1s"
      }
    }      
  }
  