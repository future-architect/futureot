package occonfig

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"

	xray "contrib.go.opencensus.io/exporter/aws"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/Datadog/opencensus-go-exporter-datadog"
	honeycomb "github.com/honeycombio/opencensus-exporter/honeycomb"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/exporter/zipkin"
)

type Mode int

const (
	Trace Mode = 1
	Stats      = 2
)

type Config struct {
	ServiceName   string
	ServiceUrl    string
	HoneycombKey  string
	ConfigFile    string
	TraceExporter string
	TraceSampler  float64
	StatsExporter string
	ZPage         string
}

var getConfigFromCommandLine func() (*Config, error)

type OCConfig interface {
	Close()
	StartServer()
}

type occonfigImpl struct {
	finalizes   []func()
	startServer func()
}

func (f occonfigImpl) Close() {
	for _, finalizer := range f.finalizes {
		finalizer()
	}
}

func (f occonfigImpl) StartServer() {
	if f.startServer != nil {
		f.startServer()
	}
}

func getConfig() (*Config, error) {
	wd, _ := os.Getwd()
	config, err := initByEnvMap(os.Environ())
	if err != nil {
		return nil, err
	}
	config, err = readFiles(config, wd)
	if err != nil {
		return nil, err
	}
	if getConfigFromCommandLine != nil {
		commandConfig, err := getConfigFromCommandLine()
		if err != nil {
			return nil, err
		}
		commandConfig, err = readFiles(commandConfig, wd)
		if err != nil {
			return nil, err
		}
		config = mergeConfigs(config, commandConfig)
	}
	if config.TraceSampler < 0 {
		config.TraceSampler = 1.0
	}
	if config.ServiceName == "" {
		config.ServiceName = filepath.Base(os.Args[0])
	}
	return config, nil
}

func Init(mode Mode) (OCConfig, error) {
	finalizer := occonfigImpl{}
	config, err := getConfig()
	if err != nil {
		return finalizer, err
	}
	isStackDriverInitialized := false
	isDataDogInitialized := false
	isZPageInitialized := false
	if mode&Trace == Trace && config.TraceExporter != "" {
		exporter, err := SelectTraceExporter(config.TraceExporter)
		if err != nil {
			return finalizer, err
		}
		switch exporter.Type {
		case STACKDRIVER:
			{
				sd, err := stackdriver.NewExporter(stackdriver.Options{
					ProjectID: exporter.Host,
				})
				if err != nil {
					return finalizer, fmt.Errorf("Failed to create the GCP StackDriver exporter: %v", err)
				}
				finalizer.finalizes = append(finalizer.finalizes, func() {
					sd.Flush()
				})
				trace.RegisterExporter(sd)
				if mode&Stats == Stats && config.StatsExporter != "" {
					view.RegisterExporter(sd)
					isStackDriverInitialized = true
				}
			}
		case XRAY:
			{
				xe, err := xray.NewExporter(xray.WithVersion("latest"))
				if err != nil {
					return finalizer, fmt.Errorf("Failed to create the AWS X-Ray exporter: %v", err)
				}
				finalizer.finalizes = append(finalizer.finalizes, func() {
					xe.Flush()
				})
				trace.RegisterExporter(xe)
			}
		case DATADOG:
			{
				dd, err := datadog.NewExporter(datadog.Options{})
				if err != nil {
					return finalizer, fmt.Errorf("Failed to create the Datadog exporter: %v", err)
				}
				finalizer.finalizes = append(finalizer.finalizes, func() {
					dd.Stop()
				})
				trace.RegisterExporter(dd)
				if mode&Stats == Stats && config.StatsExporter != "" {
					view.RegisterExporter(dd)
					isDataDogInitialized = true
				}
			}
		case HONEYCOMB:
			{
				if config.HoneycombKey == "" {
					return finalizer, errors.New("Honeycomb Write Key is empty")
				}
				hc := honeycomb.NewExporter(config.HoneycombKey, "YOUR-DATASET-NAME")
				hc.SampleFraction = config.TraceSampler
				trace.RegisterExporter(hc)
			}
		case JAEGER:
			{
				je, err := jaeger.NewExporter(jaeger.Options{
					CollectorEndpoint: exporter.Host,
					Process: jaeger.Process{
						ServiceName: config.ServiceName,
					},
				})
				if err != nil {
					return finalizer, fmt.Errorf("Failed to create the Jaeger exporter: %v", err)
				}
				finalizer.finalizes = append(finalizer.finalizes, func() {
					je.Flush()
				})
				trace.RegisterExporter(je)
			}
		case ZIPKIN:
			{
				localEndpointURI := config.ServiceUrl
				reporterURI := exporter.Host
				serviceName := config.ServiceName

				localEndpoint, err := openzipkin.NewEndpoint(serviceName, localEndpointURI)
				if err != nil {
					return finalizer, fmt.Errorf("Failed to create Zipkin localEndpoint with URI %q error: %v", localEndpointURI, err)
				}

				reporter := zipkinHTTP.NewReporter(reporterURI)
				ze := zipkin.NewExporter(reporter, localEndpoint)

				trace.RegisterExporter(ze)
			}
		}

		switch config.TraceSampler {
		case 0.0:
			trace.ApplyConfig(trace.Config{
				DefaultSampler: trace.NeverSample(),
			})
		case 1.0:
			trace.ApplyConfig(trace.Config{
				DefaultSampler: trace.AlwaysSample(),
			})
		default:
			trace.ApplyConfig(trace.Config{
				DefaultSampler: trace.ProbabilitySampler(config.TraceSampler),
			})
		}
	}

	if mode&Stats == Stats && config.StatsExporter != "" {
		exporter, err := SelectStatsExporter(config.StatsExporter)
		if err != nil {
			return finalizer, err
		}
		switch exporter.Type {
		case STACKDRIVER:
			{
				if !isStackDriverInitialized {
					sd, err := stackdriver.NewExporter(stackdriver.Options{
						ProjectID: exporter.Host,
					})
					if err != nil {
						return finalizer, fmt.Errorf("Failed to create the GCP StackDriver exporter: %v", err)
					}
					finalizer.finalizes = append(finalizer.finalizes, func() {
						sd.Flush()
					})
					view.RegisterExporter(sd)
				}
			}
		case DATADOG:
			{
				if !isDataDogInitialized {
					dd, err := datadog.NewExporter(datadog.Options{})
					if err != nil {
						return finalizer, fmt.Errorf("Failed to create the Datadog exporter: %v", err)
					}
					finalizer.finalizes = append(finalizer.finalizes, func() {
						dd.Stop()
					})
					view.RegisterExporter(dd)
				}
			}
		case PROMETHEUS:
			{
				pe, err := prometheus.NewExporter(prometheus.Options{
					Namespace: config.ServiceName,
				})
				if err != nil {
					return finalizer, fmt.Errorf("Failed to create Prometheus exporter: %v", err)
				}

				view.RegisterExporter(pe)
				u, _ := url.Parse(exporter.Host)
				var z *url.URL
				if config.ZPage != "" {
					z, err = url.Parse(config.ZPage)
					if err != nil {
						return finalizer, fmt.Errorf("Failed to parse ZPage URL: %v", err)
					}
					if z.Port() == u.Port() {
						printZPageInformation(z)
						if z.Path == "/metrics" {
							return finalizer, fmt.Errorf("ZPage and Prometheus uses same endpoints: %v", err)
						}
						isZPageInitialized = true
					}
				}

				exit := make(chan bool)
				finalizer.finalizes = append(finalizer.finalizes, func() {
					exit <- true
				})
				go func() {
					mux := http.NewServeMux()
					mux.Handle("/metrics", pe)
					if z != nil {
						zpages.Handle(mux, z.Path)
					}
					if err := http.ListenAndServe(":"+u.Port(), mux); err != nil {
						log.Fatalf("Failed to run Prometheus /metrics endpoint: %v", err)
					}
					<-exit
				}()
			}
		}
	}
	if config.ZPage != "" && !isZPageInitialized {
		z, err := url.Parse(config.ZPage)
		if err != nil {
			return finalizer, fmt.Errorf("Failed to parse ZPage URL: %v", err)
		}
		exit := make(chan bool)
		finalizer.finalizes = append(finalizer.finalizes, func() {
			exit <- true
		})
		printZPageInformation(z)
		go func() {
			mux := http.NewServeMux()
			zpages.Handle(mux, z.Path)
			if err := http.ListenAndServe(":"+z.Port(), mux); err != nil {
				log.Fatalf("Failed to run ZPage %s endpoint: %v", z.Path, err)
			}
			<-exit
		}()
	}
	return finalizer, nil
}

func printZPageInformation(u *url.URL) {
	fmt.Fprintf(os.Stderr, "[OpenCensus] ZPage is initialized. The following URLs are available:\n")
	fmt.Fprintf(os.Stderr, "    http://localhost:%s%s/rpcz\n", u.Port(), u.Path)
	fmt.Fprintf(os.Stderr, "    http://localhost:%s%s/tracez\n", u.Port(), u.Path)
}
