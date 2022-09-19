### Testing Nomad Native Connect 

In order to test the new Nomad native connect feature I started by creating a job that doesnt use the Consul Service Mesh. In order to drop the usage of the Consul Service Mesh **I ommited  the entire `connect` stanza from the job template!** Instead, I added a `provider = "nomad"` entry at the `service` stanza like so:

```hcl
    service {
      name = "nginx"
      port = 8888
      provider = "nomad"
    }
```
Take a look at the [OTP job template](https://github.com/trigovision/store_deployment/blob/master/templates/services/otp/app.nomad.tpl) for an example of using the Consul connect, and compare it to the basic [nginx job template](https://github.com/oavner/nativeConnect/blob/main/app.nomad) in the current repo. As u can see, integrating into the new Native Connect only has a few requirments:
1. Nomad cluster version 1.3 and higher.
2. A `provider = "nomad"` entry at the service stanza.
3. Ommiting older consul `connect` stanzas.<br/><br/><br/>

Using the native connect of Nomad instead of Consul also comes with minor networking changes in the Noamd cluster:
1. Moving to a vxlan architecture. Instead of sidecar proxy containers connected to each task, all the containers running with the Native Connect share the same overlay network and theres no need for another sidecar proxy containers at all. That means less resource exhaustion **validate**
2. No more cherry picked ACL's for controlling each service's ingress or egress traffic. No sidecar proxies means less security- as mentioned before all running containers share the same overlay network and theoretically can communicate freely with each other.


