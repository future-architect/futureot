# Prometheus

1. Rewrite IP address in prometheus.yml to your local IP (localhost is not accessible from docker container).
2. Launch Prometheus
   ```sh
   $ docker run --name prometheus --rm -p 9090:9090 -v /abs-path/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
   ```
3. Launch your program
   ```sh
   $ OC_STATS_EXPORTER=p8s://:8888 go run main.go
   ```
4. open :9090 with your browser.

# Graphite

1. Launch graphite
   ```sh
   docker run -d --name graphite --rm -p 8888:80 -p 2003-2004:2003-2004 graphiteapp/graphite-statsd
   ```
2. Launch your program
   ```sh
   $ OC_STATS_EXPORTER=graphite go run main.go
   ```
3. open :8888 with your browser.

