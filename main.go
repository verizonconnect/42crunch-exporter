package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/mieliespoor/42crunch-exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	crunch "github.com/verizonconnect/42crunch-client-go"
	"go.uber.org/zap"
)

const (
	envAddress string = "42_CRUNCH_ADDR"
	envAPIKey  string = "42_CRUNCH_API_KEY"
)

func init() {
	prometheus.MustRegister(version.NewCollector(exporter.Namespace + "_exporter"))
}

func main() {
	var (
		webConfig     = webflag.AddFlags(kingpin.CommandLine, ":9916")
		metricsPath   = kingpin.Flag("web.metrics-path", "Path under which to expose metrics").Default("/metrics").String()
		crunchAddress = kingpin.Flag("42crunch.address", fmt.Sprintf("42Crunch server address (can also be set with $%s)", envAddress)).Default("https://platform.42crunch.com").Envar(envAddress).String()
		crunchAPIKey  = kingpin.Flag("42crunch.api-key", fmt.Sprintf("42Crunch API key (can also be set with $%s)", envAPIKey)).Envar(envAPIKey).Required().String()
	)

	kingpin.Version(version.Print(exporter.Namespace + "_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	logger.With()

	logger.Info("Starting exporter",
		zap.String("namespace", exporter.Namespace),
		zap.String("version", version.Info()),
	)
	logger.Info("Build context",
		zap.String("context", version.BuildContext()),
	)

	client, err := crunch.NewClient(*crunchAddress, crunch.WithAPIKey(*crunchAPIKey))
	if err != nil {
		logger.Error("Error creating client",
			zap.Error(err),
		)
		os.Exit(1)
	}

	e := exporter.Exporter{
		Client: client,
		Logger: logger,
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
			logger.Error("Error starting HTTP server", zap.Error(err))
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			logger.Info("Received SIGTERM, exiting gracefully...")
			os.Exit(0)
		case <-srvc:
			os.Exit(1)
		}
	}
}
