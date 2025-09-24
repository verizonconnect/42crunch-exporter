package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/verizonconnect/42crunch-exporter/internal/exporter"

	"github.com/prometheus/client_golang/prometheus"
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

func main() {
	prometheus.MustRegister(prometheus.NewGoCollector())
	var (
		webConfig           = webflag.AddFlags(kingpin.CommandLine, ":9916")
		metricsPath         = kingpin.Flag("web.metrics-path", "Path under which to expose metrics").Default("/metrics").String()
		crunchAddress       = kingpin.Flag("42c-address", fmt.Sprintf("42Crunch server address (can also be set with $%s)", envAddress)).Default("https://platform.42crunch.com").Envar(envAddress).String()
		crunchAPIKey        = kingpin.Flag("42c-api-key", fmt.Sprintf("42Crunch API key (can also be set with $%s)", envAPIKey)).Envar(envAPIKey).Required().String()
		collectionInclRegex = kingpin.Flag("42c-collection-regex", fmt.Sprintf("Regex which will include only specific 42Crunch API collections. (can also be set with $%s)", env42cCollectionRegex)).Envar(env42cCollectionRegex).String()
	)

	kingpin.Version(version.Print(exporter.Namespace + "_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("Starting exporter", "namespace", exporter.Namespace, "version", version.Info())
	logger.Info("Build context", "context", version.BuildContext())
	client, err := crunch.NewClient(*crunchAddress, crunch.WithAPIKey(*crunchAPIKey))
	if err != nil {
		logger.Error("Error creating client", "err", err)
		os.Exit(1)
	}

	e := exporter.Exporter{
		Client: client,
		Logger: logger,
		Config: exporter.ExporterConfig{CollectionInclRegex: *collectionInclRegex},
	}

	if metricsPath != nil {
		http.HandleFunc(*metricsPath, e.HandlerFunc())
	} else {
		logger.Error("metricsPath is nil, cannot register metrics handler")
		os.Exit(1)
	}
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

	srv := &http.Server{}
	go func() {
		if err := web.ListenAndServe(srv, webConfig, logger); err != http.ErrServerClosed {
			logger.Error("Error starting HTTP server", "err", err)
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			logger.Info("Received SIGTERM, shutting down HTTP server gracefully...")
			if err := srv.Shutdown(context.Background()); err != nil {
				logger.Error("Error during server shutdown", "err", err)
				os.Exit(1)
			}
			os.Exit(0)
		case <-srvc:
			os.Exit(1)
		}
	}
}
