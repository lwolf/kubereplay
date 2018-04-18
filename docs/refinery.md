# Refinery

RefineryCRD is one of the core components of Kubereplay.
It is the `server` part in terms of [goreplay distributed configuration](https://github.com/buger/goreplay/wiki/Distributed-configuration)  
During population it creates deployment and service used for receiving traffic from workloads and sending it to the configured output.

 
```yaml
apiVersion: kubereplay.lwolf.org/v1alpha1
kind: Refinery
metadata:
  name: refinery-example
spec:
  # Gor uses dynamic worker scaling by default.  Enter a number to run a set number of workers.
  # number > 0 indicates static value
  workers: 1  # -output-http-workers int
  # Specify HTTP request/response timeout. By default 5s. Example: --output-http-timeout 30s (default 5s)
  timeout: "5s" # -output-http-timeout duration
  output:
    file:
      enabled: false
      # Write incoming requests to file
      filename: "test.log"  # -output-file value

      # The flushed chunk is appended to existence file or not.
      append: true  # -output-file-append

      # Interval for forcing buffer flush to the file, default: 1s. (default 1s)
      flushInterval: "1s"  # -output-file-flush-interval duration

      # The length of the chunk queue. Default: 256 (default 256)
      queueSize: 256  # -output-file-queue-limit int

      # Size of each chunk. Default: 32mb (default 33554432)
      fileLimit: "32mb"  # -output-file-size-limit value

    tcp:
      enabled: false
      # Used for internal communication between Gor instances. Example:
      uri: "replay.kubernetes:28020" # -output-tcp value

    stdout: # --output-stdout
      enabled: true

    http:
      enabled: false
      # Forwards incoming requests to given http address.
      uri: "http://staging.com" # -output-http value

      # Enables http debug output.
      debug: true  # -output-http-debug

      # HTTP response buffer size, all data after this size will be discarded.
      responseBuffer: 1 # -output-http-response-buffer int

    elasticsearch:
      enabled: false
      # Send request and response stats to ElasticSearch:
      uri: es_host:api_port/index_name # -output-http-elasticsearch string

    kafka:
      enabled: false
      # Read request and response stats from Kafka:
      uri: "192.168.0.1:9092,192.168.0.2:9092" # -output-kafka-host string

      # If turned on, it will serialize messages from GoReplay text format to JSON.
      json: false  # -output-kafka-json-format

      # kafka topic
      topic: "kafka-log"  # -output-kafka-topic string


```