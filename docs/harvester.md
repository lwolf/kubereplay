# Harvester

HarvesterCRD  is one of the core components of Kubereplay.
It is used to configure which deployments should be controlled by Kubereplay.
Based on selector in Harvester spec Kubereplay will add GoReplay-sidecar to matching deployments.


```yaml
apiVersion: kubereplay.lwolf.org/v1alpha1
kind: Harvester
metadata:
  name: harvester-example
spec:
  # Set percentage of instances to capture traffic, from 0 to 100
  # 100 represents all instances
  segment: 70
  # Configure goreplay to listen to this port
  app_port: 8080
  # Select instances based on selector
  selector:
    app: kubereplay
    module: test
  # Name of the Refinery to send traffic to
  refinery: "refinery-example"

```