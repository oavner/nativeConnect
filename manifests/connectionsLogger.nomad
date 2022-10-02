job "connectionsLogger" {
  datacenters = ["dc1"]
  type        = "sysbatch"
  // priority    = 70

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

      config {
        image       = "ghcr.io/oavner/connections-logger:3023d3a1258e3c5b6b5f858ef87906363cb329fb"
      }

      resources {
        memory     = 512 
        cpu        = 1000 # in MHZ
      }
    }
  }
}
