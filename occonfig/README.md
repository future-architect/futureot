# OpenCensus Configurator

It provides the following initialization options:

* By Environment Variables
* By JSON file
* By Commandline Options

## How to Use for Your Programs' Users

### By Environment Variables

* ``OC_SERVICE_NAME``: Service name that appears resulting report. Default value is the command name.

* ``OC_SERVICE_URL``:  Service URL some tracer filters result by the URL

* ``OC_CONFIG_JSON``: JSON file path for settings (see below)

* ``OC_TRACE_EXPORTER`` (required for tracing)

   * ``stackdriver://demo-project-id``: StackDriver
   * ``sd://demo-project-id`` : short form of stackdriver
   * ``datadog://localhost:8126`` or ``dd://localhost:8126`` : DataDog
   * ``datadog`` or ``dd`` : DataDog (default host:port is localhost:8126)
   * ``xray``: AWS X-Ray
   * ``jaeger://localhost:6831`` : Jaeger
   * ``jaeger://localhost`` : Jaeger (default port is 6831)
   * ``jaeger`` : Jaeger (default host is localhost:6831)
   * ``zipkin://localhost:9411/api/v2/spans`` : Zipkin (default port is 9411)
   * ``zipkin://localhost/api/v2/spans`` : Zipkin (default port is 9411)
   * ``zipkin://localhost`` : Zipkin (default port is 9411, default path is /api/v2/spans)
   * ``honeycomb`` : HoneyComb

* ``OC_TRACE_SAMPLER``

   * ``always``: Default value
   * ``never``: Never send trace
   * floating number (0-1): Probabilistic sampler

* ``OC_HONEYCOMB_WRITE_KEY``: honeycomb.io API key.If 
    the value starts ``file://``,
    this library searches local file.

### Typical Usage for commandline

#### Via Einvironment Variables

```bash
$ export OC_TRACE_EXPORTER=stackdriver://demo-project-id
$ export OC_SERVICE_NAME=my-service
$ ./your-program
```

#### Via flag

* Common Settings

   * ``-oc-service-name``: Service name
   * ``-oc-service-url``: Service name
   * ``-oc-config-json``: JSON file path for settings (see below)

* For tracer

   * ``-oc-honeycomb-write-key``: honeycomb.io write key file path
   * ``-oc-trace-exporter``: Exporter setting

```bash
# flag package support
$ ./your-program -oc-trace-exporter stackdriver://demo-project-id -oc-service-name my-service
```

#### Via kingpin.v2

* Common Settings

   * ``--oc-service-name``: Service name
   * ``--oc-service-url``: Service name
   * ``--oc-config-json``: JSON file path for settings (see below)

* For tracer

   * ``--oc-trace-exporter``: Exporter setting
   * ``--oc-honeycomb-write-key``: honeycomb.io write key file path

```
# kingpin.v2 package support
$ ./your-program --oc-trace-exporter=xray --oc-service-name=my-service
```

### JSON file format

You can pass setting file path via ``-oc-config-json`` (flag support),  ``--oc-config-json`` (kingpin.v2 support) options.

Extends specified base JSON.

```json
{
  "service-name": "my-awesome-service",
  "service-url":  "http://localhost:8080",
  "extends": "../config.json",
  "trace": {
    "exporter": "stackdriver://demo-project-id",
    "honeycomb-write-key": "honeycomb.key",
    "sampling": "always"
  }
}
```

## Priority of configs

1. Commnadline options
2. JSON file that is specified at ``extends`` in the file of commandline option ``--oc-config-json``
3. JSON file that is specified at commandline option ``--oc-config-json``
4. Environment variables
5. JSON file that is specified at ``extends`` in the file of environment variable ``OC_CONFIG_JSON``
6. JSON file that is specified at environment variable ``OC_CONFIG_JSON``

## Tool settings for Local development users

### Typical Usage for local Jaeger

#### Docker

```bash
$ docker run -d --name jaeger -p 14268:14268 -p 16686:16686 jaegertracing/all-in-one:1.12
$ OC_TRACE_EXPORTER=jaeger ./your-program
```

#### docker-compose

```yaml
version: '3'
services:
  jaeger:
    image: jaegertracing/all-in-one:1.12
    ports:
      - 16686:16686  # for web console
      - 14268:14268  # it is needed if your service is run out of docker-compose
  your-service:
    image: your-service
    ports:
      - 8080:8080
    environment:
      - OC_TRACE_EXPORTER=jaeger://jaeger
      - OC_SERVICE_URL=http://localhost:8080
    depends_on:
      - jaeger
```

### Typical Usage for local Zipkin

#### Docker

```bash
$ docker run -d -p 9411:9411 openzipkin/zipkin
$ OC_TRACE_EXPORTER=zipkin ./your-program
```

#### docker-compose

```yaml
version: '3'
services:
  zipkin:
    image: openzipkin/zipkin
    ports:
      - 9411:9411
  your-service:
    image: your-service
    ports:
      - 8080:8080
    environment:
      - OC_TRACE_EXPORTER=zipkin://zipkin
      - OC_SERVICE_URL=http://localhost:8080
    depends_on:
      - zipkin
```

## How to Use for Programmers

```go
package main

import (
	"flag"
	"github.com/shibukawa/occonfig"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Usage 1: Only support EnvVar
func init() {
	// configuration via Environment Variables
	finalizer, err := occonfig.Init(occonfig.Trace | occonfig.Stats)
	if err != nil {
		panic(err)
	}
	defer finalizer.Close()
	
	// start your logic from here
}

// Usage 2: Support EnvVar and flag as an option parser
func main() {
	// Call it befor flag.Parse()
	occonfig.UseFlag(occonfig.Trace)
	
	flag.Parse()
	
	// Then call Init()
	finalizer, err := occonfig.Init(occonfig.Trace)
	if err != nil {
		panic(err)
	}
	defer finalizer.Close()

	// start your logic from here
}

// Usage 3: Support EnvVar and kingpin.v2 as an option parser
func main() {
	// Call it befor flag.Parse()
	occonfig.UseKingpin(occonfig.Stats)
	
	kingpin.Parse()
	
	// Then call Init()
	finalizer, err := occonfig.Init(occonfig.Stats)
	if err != nil {
		panic(err)
	}
	defer finalizer.Close()

	// start your logic from here
}

```