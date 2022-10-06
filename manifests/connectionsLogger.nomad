job "connectionsLogger" {
  datacenters = ["demo"]
  type        = "sysbatch"
  priority    = 70

  reschedule {
    delay          = "30s"
    delay_function = "constant"
    unlimited      = true
  }

   periodic {
    // Launch every 20 seconds
    cron = "*/20 * * * * * *"

    // Allow overlapping runs.
    prohibit_overlap = false
  }

  group "connectionsLoggers" {
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

      env {
        METHOD = "GET"
        URLS = '["http://10.155.0.113:22172/", "http://10.155.0.113:20667/"]'
        MAX_SESSIONS = "4"
      }

      config {
        image       = "ghcr.io/oavner/connections-logger:77be1b9233b2bc7a1cbdc29c91b8e88844c3bef8"
      }

      resources {
        memory     = 512 
        cpu        = 1000 # in MHZ
      }
    }
  }
}
