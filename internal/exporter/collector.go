package exporter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/dmdutra/github-outage-exporter/internal/githubstatus"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "github"

type Collector struct {
	client *githubstatus.Client
	logger *slog.Logger

	statusDesc        *prometheus.Desc
	componentDesc     *prometheus.Desc
	componentInfoDesc *prometheus.Desc
	incidentDesc      *prometheus.Desc
	maintenanceDesc   *prometheus.Desc
	scrapeSuccessDesc *prometheus.Desc
	scrapeDurationDesc *prometheus.Desc
}

func NewCollector(client *githubstatus.Client, logger *slog.Logger) *Collector {
	if logger == nil {
		logger = slog.Default()
	}

	return &Collector{
		client: client,
		logger: logger,
		statusDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "status_indicator_severity"),
			"Overall GitHub status indicator severity (0=none, 1=minor, 2=major, 3=critical).",
			[]string{"indicator", "description"},
			nil,
		),
		componentDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "component_status_severity"),
			"GitHub component status severity (0=operational, 1=degraded, 2=partial_outage, 3=major_outage).",
			[]string{"component_id", "component", "status"},
			nil,
		),
		componentInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "component_operational"),
			"Whether the GitHub component is fully operational (1) or not (0).",
			[]string{"component_id", "component", "status"},
			nil,
		),
		incidentDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "incident_active"),
			"Whether an unresolved GitHub incident is active (1) or not (0).",
			[]string{"incident_id", "name", "status", "impact"},
			nil,
		),
		maintenanceDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "maintenance_active"),
			"Whether a scheduled maintenance is currently active (1) or not (0).",
			[]string{"maintenance_id", "name", "status", "impact"},
			nil,
		),
		scrapeSuccessDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "success"),
			"Whether the last scrape of GitHub Status API succeeded.",
			nil,
			nil,
		),
		scrapeDurationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "duration_seconds"),
			"Duration of the last GitHub Status API scrape in seconds.",
			nil,
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.statusDesc
	ch <- c.componentDesc
	ch <- c.componentInfoDesc
	ch <- c.incidentDesc
	ch <- c.maintenanceDesc
	ch <- c.scrapeSuccessDesc
	ch <- c.scrapeDurationDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), c.client.Timeout())
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex

	var statusResp *githubstatus.StatusResponse
	var componentsResp *githubstatus.ComponentsResponse
	var incidentsResp *githubstatus.IncidentsResponse
	var maintenancesResp *githubstatus.ScheduledMaintenancesResponse
	var fetchErr error

	collectEndpoint := func(name string, fn func(context.Context) error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				mu.Lock()
				if fetchErr == nil {
					fetchErr = err
				}
				c.logger.Warn("failed to fetch endpoint", "endpoint", name, "error", err)
				mu.Unlock()
			}
		}()
	}

	collectEndpoint("status", func(ctx context.Context) error {
		resp, err := c.client.GetStatus(ctx)
		if err != nil {
			return err
		}
		statusResp = resp
		return nil
	})

	collectEndpoint("components", func(ctx context.Context) error {
		resp, err := c.client.GetComponents(ctx)
		if err != nil {
			return err
		}
		componentsResp = resp
		return nil
	})

	collectEndpoint("incidents", func(ctx context.Context) error {
		resp, err := c.client.GetUnresolvedIncidents(ctx)
		if err != nil {
			return err
		}
		incidentsResp = resp
		return nil
	})

	collectEndpoint("maintenances", func(ctx context.Context) error {
		resp, err := c.client.GetActiveMaintenances(ctx)
		if err != nil {
			return err
		}
		maintenancesResp = resp
		return nil
	})

	wg.Wait()

	duration := time.Since(start).Seconds()
	success := fetchErr == nil

	ch <- prometheus.MustNewConstMetric(c.scrapeDurationDesc, prometheus.GaugeValue, duration)
	if success {
		ch <- prometheus.MustNewConstMetric(c.scrapeSuccessDesc, prometheus.GaugeValue, 1)
	} else {
		ch <- prometheus.MustNewConstMetric(c.scrapeSuccessDesc, prometheus.GaugeValue, 0)
		return
	}

	if statusResp != nil {
		ch <- prometheus.MustNewConstMetric(
			c.statusDesc,
			prometheus.GaugeValue,
			githubstatus.IndicatorSeverity(statusResp.Status.Indicator),
			statusResp.Status.Indicator,
			statusResp.Status.Description,
		)
	}

	if componentsResp != nil {
		for _, component := range componentsResp.Components {
			labels := []string{component.ID, component.Name, component.Status}
			ch <- prometheus.MustNewConstMetric(
				c.componentDesc,
				prometheus.GaugeValue,
				githubstatus.ComponentSeverity(component.Status),
				labels...,
			)

			operational := 0.0
			if component.Status == "operational" {
				operational = 1
			}
			ch <- prometheus.MustNewConstMetric(
				c.componentInfoDesc,
				prometheus.GaugeValue,
				operational,
				labels...,
			)
		}
	}

	if incidentsResp != nil {
		for _, incident := range incidentsResp.Incidents {
			ch <- prometheus.MustNewConstMetric(
				c.incidentDesc,
				prometheus.GaugeValue,
				1,
				incident.ID,
				incident.Name,
				incident.Status,
				incident.Impact,
			)
		}
	}

	if maintenancesResp != nil {
		for _, maintenance := range maintenancesResp.ScheduledMaintenances {
			ch <- prometheus.MustNewConstMetric(
				c.maintenanceDesc,
				prometheus.GaugeValue,
				1,
				maintenance.ID,
				maintenance.Name,
				maintenance.Status,
				maintenance.Impact,
			)
		}
	}
}
