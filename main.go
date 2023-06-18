package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/kit/log/level"
	"github.com/mieliespoor/42crunch-exporter/internal/exporter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	crunch "github.com/verizonconnect/42crunch-client-go"
)

const (
	envAddress            string = "42C_ADDR"
	envAPIKey             string = "42C_API_KEY"
	env42cCollectionRegex string = "42C_COLLECTION_REGEX"
)

func init() {
	prometheus.MustRegister(version.NewCollector(exporter.Namespace + "_exporter"))
}

func main() {
	var (
		format              = promlog.AllowedFormat{}
		webConfig           = webflag.AddFlags(kingpin.CommandLine, ":9916")
		metricsPath         = kingpin.Flag("web.metrics-path", "Path under which to expose metrics").Default("/metrics").String()
		crunchAddress       = kingpin.Flag("42c-address", fmt.Sprintf("42Crunch server address (can also be set with $%s)", envAddress)).Default("https://platform.42crunch.com").Envar(envAddress).String()
		crunchAPIKey        = kingpin.Flag("42c-api-key", fmt.Sprintf("42Crunch API key (can also be set with $%s)", envAPIKey)).Envar(envAPIKey).Required().String()
		collectionInclRegex = kingpin.Flag("collection-excl-regex", fmt.Sprintf("")).Envar(env42cCollectionRegex).String()
	)

	format.Set("json")
	promlogConfig := promlog.Config{
		Format: &format,
	}

	flag.AddFlags(kingpin.CommandLine, &promlogConfig)
	kingpin.Version(version.Print(exporter.Namespace + "_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(&promlogConfig)

	level.Info(logger).Log("msg", fmt.Sprintf("Starting %s_exporter %s", exporter.Namespace, version.Info()))
	level.Info(logger).Log("msg", fmt.Sprintf("Build context %s", version.BuildContext()))

	client, err := crunch.NewClient(*crunchAddress, crunch.WithAPIKey(*crunchAPIKey))
	if err != nil {
		level.Error(logger).Log("msg", "Error creating client", "err", err)
		os.Exit(1)
	}

	e := exporter.Exporter{
		Client: client,
		Logger: logger,
		Config: exporter.ExporterConfig{CollectionInclRegex: collectionInclRegex},
	}

	http.HandleFunc(*metricsPath, e.HandlerFunc())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
						 <head><title>42Crunch Exporter</title></head>
						 <body>
						 <h1>42Crunch Exporter</h1>
						 <p><a href='` + *metricsPath + `'>Metrics</a></p>
						 </body>
						 </html>`))
	})

	srvc := make(chan struct{})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		srv := &http.Server{}
		if err := web.ListenAndServe(srv, webConfig, logger); err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			os.Exit(0)
		case <-srvc:
			os.Exit(1)
		}
	}
}
