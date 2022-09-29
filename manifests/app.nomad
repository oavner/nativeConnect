job "nginx" {
  datacenters = ["dc1"]
  type        = "service"
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
      mode = "bridge"
      port "http" {
        to = 80
      }
    }

    service {
      name = "nginx"
      port = 8888
      provider = "nomad"
  
      meta {
        metrics_port = "${NOMAD_HOST_PORT_http}"
      }
    }

    task "nginx" {
      driver = "docker"

      config {
        image       = "nginx"
      }

      resources {
        memory     = 512 
        cpu        = 1000 # in MHZ
      }
    }
  }
}
