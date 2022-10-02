job "nginx" {
  datacenters = ["dc1"]
  type        = "sysbatch"
  priority    = 70

  reschedule {
    delay          = "30s"
    delay_function = "constant"
    unlimited      = true
  }

  group "nginx" {
    count = 1

    restart {
      interval = "60s"
      attempts = 10
      delay    = "5s"
      mode     = "delay"
    }

    network {
      mode = "host"
    }

    task "connectionsLogger" {
      driver = "docker"

      config {
        image       = "ghcr.io/oavner/connections-logger:82fa42b7c7bc892e37af28802de4021408ebe640"
      }

      resources {
        memory     = 512 
        cpu        = 1000 # in MHZ
      }
    }
  }
}
