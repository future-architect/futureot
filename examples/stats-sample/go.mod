module stats_sample

go 1.12

replace github.com/shibukawa/occonfig => ../..

require (
	github.com/shibukawa/occonfig v0.0.0-00010101000000-000000000000
	go.opencensus.io v0.19.2
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
