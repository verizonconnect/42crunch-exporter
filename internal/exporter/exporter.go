package exporter

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

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
	Config ExporterConfig
}

type ExporterConfig struct {
	CollectionInclRegex *string
}

func (e *Exporter) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()

		if err := e.collectApiCollectionMetrics(r.Context(), registry); err != nil {
			level.Error(e.Logger).Log("err", err)
			http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
			return
		}

		if err := e.collectApiAuditMetrics(r.Context(), registry); err != nil {
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
		obj, err := regexp.Match(*e.Config.CollectionInclRegex, []byte(c.Description.Name))
		if err != nil {
			level.Error(e.Logger).Log("msg", "regex failed", "err", err)
		} else if obj {
			collectionInformation.With(prometheus.Labels{
				"id":   c.Description.Id,
				"name": c.Description.Name,
			}).Set(float64(c.Summary.ApiCount))
		} else {
			level.Debug(e.Logger).Log("msg", fmt.Sprintf("regex did not match for %s", c.Description.Name))
		}
	}

	return nil
}

func (e *Exporter) collectApiAuditMetrics(ctx context.Context, registry *prometheus.Registry) error {
	collections, err := e.Client.Collections.GetAll(ctx)
	if err != nil {
		level.Error(e.Logger).Log("msg", "api collections could not be retrieved", "err", err)
		return err
	}

	var (
		apiInformation = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "information"),
				Help: "Basic information about an API",
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
				Help: "The number of critical vulnerabilities per api based on the API Audit",
			},
			[]string{"id"},
		)
		apiAssessmentHighs = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_highs"),
				Help: "The number of high vulnerabilities per api based on the API Audit",
			},
			[]string{"id"},
		)
		apiAssessmentMediums = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_mediums"),
				Help: "The number of medium vulnerabilities per api based on the API Audit",
			},
			[]string{"id"},
		)
		apiAssessmentLows = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_lows"),
				Help: "The number of low vulnerabilities per api based on the API Audit",
			},
			[]string{"id"},
		)
		apiAssessmentInfos = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_infos"),
				Help: "The number of information messages per api based on the API Audit",
			},
			[]string{"id"},
		)
		apiAssessmentGrade = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_grade"),
				Help: "API Audit Assessment Grade",
			},
			[]string{"id"},
		)
		apiAssessmentErrors = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_errors"),
				Help: "The number of API errors",
			},
			[]string{"id"},
		)
		apiAssessmentValid = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_valid"),
				Help: "Indicating whether the api schema is valid",
			},
			[]string{"id"},
		)
		apiAssessmentLastAudit = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: prometheus.BuildFQName(Namespace, "api", "assessment_last_audit"),
				Help: "Last API Audit Assessment date, represented as a Unix timestamp",
			},
			[]string{"id"},
		)
	)

	registry.MustRegister(
		apiInformation,
		apiAssessmentCriticals,
		apiAssessmentHighs,
		apiAssessmentMediums,
		apiAssessmentLows,
		apiAssessmentInfos,
		apiAssessmentGrade,
		apiAssessmentErrors,
		apiAssessmentValid,
		apiAssessmentLastAudit,
	)

	for _, c := range collections.Items {
		obj, err := regexp.Match(*e.Config.CollectionInclRegex, []byte(c.Description.Name))
		if obj == false {
			continue
		}

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

			apiAssessmentHighs.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(api.Assessment.NumHighs))

			apiAssessmentMediums.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(api.Assessment.NumMediums))

			apiAssessmentLows.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(api.Assessment.NumLows))

			apiAssessmentInfos.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(api.Assessment.NumInfos))

			apiAssessmentErrors.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(api.Assessment.NumErrors))

			apiAssessmentGrade.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(api.Assessment.Grade))

			valid := 0
			if api.Assessment.IsValid {
				valid = 1
			}
			apiAssessmentValid.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(valid))

			unix := int64(0)
			if len(api.Assessment.Last) > 0 {
				last, _ := time.Parse(time.RFC3339, api.Assessment.Last)
				unix = last.Unix()
			}

			apiAssessmentLastAudit.With(
				prometheus.Labels{
					"id": api.Description.Id,
				}).Set(float64(unix))
		}
	}

	return nil
}
