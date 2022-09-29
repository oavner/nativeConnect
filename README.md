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

### Bulding The Connections Logger
The connections logger is a simple golang program that can send multiple http requests to a slice of URLs concurrently. Using this tool we could set our desired environment variables in order to test a single or multiple connections to a single or multiple urls, and since the requests are sent concurrently there is no need to worry about runtime, it happens in a second.

Here is an example for testing multiple connections to multiple urls (4 sessions opened for each url):


The logger posts each connection data logs in `json` to `stdout` for now but it could be changed to `stderr` later on using the log golang module:

At the end of each run the looger also posts a simple human readable report that summs up all the connections state:

A more complexted logic could be handled in the report phase later on such as how many of the requests for a single url succeeded in ratio to how many were sent. Notice that the connections logger is desined as a croned task that should run in intervals, this design could be changed to run as an infinite loop that posts logs and reports all the time.

