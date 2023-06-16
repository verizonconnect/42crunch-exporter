package exporter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	crunch "github.com/verizonconnect/42crunch-client-go"
)

const (
	// Namespace is the metrics namespace of the exporter
	Namespace string = "fortytwo_crunch"
)

type Exporter struct {
	Client *crunch.Client
	Logger log.Logger
}

func (e *Exporter) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()

		if err := e.collectApiCollectionMetrics(r.Context(), registry); err != nil {
			level.Error(e.Logger).Log("err", err)
			http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
			return
		}

		if err := e.collectApiMetrics(r.Context(), registry); err != nil {
			level.Error(e.Logger).Log("err", err)
			http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
			return
		}

		// Serve
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func (e *Exporter) collectApiCollectionMetrics(ctx context.Context, registry *prometheus.Registry) error {
	var (
		collectionInformation = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "collection", "information"),
				Help: "Basic information about the api collection",
			},
			[]string{
				"id",
				"name",
			},
		)
	)

	level.Info(e.Logger).Log("msg", "collecting collection metrics...")
	registry.MustRegister(
		collectionInformation,
	)

	collections, err := e.Client.Collections.GetAll(ctx)
	if err != nil {
		level.Error(e.Logger).Log("msg", "api collection metrics collection failed...", "err", err)
		return err
	}

	for _, c := range collections.Items {
		collectionInformation.With(prometheus.Labels{
			"id":   c.Description.Id,
			"name": c.Description.Name,
		}).Set(float64(c.Summary.ApiCount))
	}

	return nil
}

func (e *Exporter) collectApiMetrics(ctx context.Context, registry *prometheus.Registry) error {
	collections, err := e.Client.Collections.GetAll(ctx)
	if err != nil {
		level.Error(e.Logger).Log("msg", "api collections could not be retrieved", "err", err)
		return err
	}

	var (
		apiInformation = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "information"),
				Help: "Basic information about the api collection",
			},
			[]string{
				"id",
				"collection_id",
				"name",
				"tags",
			},
		)
		apiAssessmentCriticals = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_criticals"),
			},
			[]string{"id"},
		)
		apiAssessmentGrade = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_grade"),
				Help: "API Audit Assessment Grade for the api",
			},
			[]string{"id"},
		)
	)

	registry.MustRegister(
		apiInformation,
		apiAssessmentCriticals,
		apiAssessmentGrade,
	)

	for _, c := range collections.Items {
		apiResult, err := e.Client.API.ListApis(ctx, c.Description.Id)
		if err != nil {
			level.Error(e.Logger).Log("msg", "collection apis could not be retrieved", "err", err)
			return err
		}

		for _, api := range apiResult.Items {
			apiTags := ","
			for _, t := range api.Tags {
				apiTags = apiTags + t.TagName + ","
			}

			apiInformation.With(prometheus.Labels{
				"id":            api.Description.Id,
				"name":          api.Description.Name,
				"collection_id": c.Description.Id,
				"tags":          apiTags,
			}).Set(1)

			apiAssessmentCriticals.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(api.Assessment.NumCriticals))

			apiAssessmentGrade.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(api.Assessment.Grade))
		}
	}

	return nil
}
