job "dummyMetricsApp2" {
  datacenters = ["dc1"]
  type        = "service"
  priority    = 70

  reschedule {
    delay          = "30s"
    delay_function = "constant"
    unlimited      = true
  }

  group "dummyMetricsApp2" {
    count = 1

    restart {
      interval = "60s"
      attempts = 10
      delay    = "5s"
      mode     = "delay"
    }

    network {
      mode = "bridge"
      port "http" {
        to = 9999
      }
    }

    service {
      name = "dummy-metrics-app2"
      port = "http"
      provider = "nomad"
  
      meta {
        metrics_port = "${NOMAD_HOST_PORT_http}"
      }
      check {
        type     = "http"
        name     = "metrics"
        path     = "/metrics"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "dummyMetricsApp" {
      driver = "docker"

      config {
        image       = "ghcr.io/oavner/dummy-metrics-app:7c44d57497ddd55fed07e3287ab7b898a29745c7"
      }

      resources {
        memory     = 512 
        cpu        = 1000 # in MHZ
      }
    }
  }
}
