package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

var NebiusMonitoring = dashboard.NewDashboardBuilder("Nebius Monitoring").
	Uid("nebius-monitoring").
	Tags([]string{"Nebius", "Monitoring"}).
	Refresh("1m").
	Time("now-15m", "now").
	Timezone("browser").
	Readonly().
	WithVariable(
		DatasourceVar,
	).
	Description("Dashboard provides an overview of the Nebius Monitoring. https://docs.nebius.com/observability/monitoring").
	Link(dashboard.NewDashboardLinkBuilder("Docs").
		Type(dashboard.DashboardLinkTypeLink).
		Url("https://docs.nebius.com/observability").
		TargetBlank(true).
		Icon("doc"),
	).
	Link(dashboard.NewDashboardLinkBuilder("GitHub").
		Type(dashboard.DashboardLinkTypeLink).
		Url("https://github.com/nebius/observability").
		TargetBlank(true).
		Icon("external link"),
	).
	WithRow(
		dashboard.NewRowBuilder("Ingest").
			GridPos(dashboard.GridPos{H: 1, W: 24, X: 0, Y: 0}).
			Id(1),
	).
	WithPanel(
		timeseries.NewPanelBuilder().
			Title("Ingest requests").
			Datasource(DatasourceRef).
			Description("Number of metrics ingestion requests per second").
			Unit("reqps").
			GridPos(dashboard.GridPos{H: 7, W: 12, X: 0, Y: 1}).
			WithTarget(
				prometheus.NewDataqueryBuilder().
					Expr(`sum (rate(requests_total{type="write"}[$__rate_interval]))`).
					LegendFormat("Requests").
					RefId("Requests rate").
					Range(),
			).
			WithTarget(
				prometheus.NewDataqueryBuilder().
					Expr(`requests_limits{type="monitoring.write.throughput.requests"}`).
					LegendFormat("Requests limit").
					RefId("Requests limit").
					Range(),
			).
			OverrideByQuery("Requests limit", []dashboard.DynamicConfigValue{
				func() dashboard.DynamicConfigValue {
					col, _ := dashboard.NewFieldColorBuilder().
						Mode(dashboard.FieldColorModeId("fixed")).
						FixedColor("dark-red").
						Build()
					return dashboard.DynamicConfigValue{Id: "color", Value: col}
				}(),
				{Id: "custom.fillOpacity", Value: 0},
				{Id: "custom.hideFrom", Value: map[string]bool{"legend": true, "tooltip": false, "viz": false}},
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).WithOverride(
			dashboard.MatcherConfig{Id: "byRegexp", Options: ".*"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			},
		),
	).
	WithPanel(
		timeseries.NewPanelBuilder().
			Title("Ingest requests errors").
			Datasource(DatasourceRef).
			Description("Number of failed metrics ingestion requests per second, by status code").
			Unit("reqps").
			GridPos(dashboard.GridPos{H: 7, W: 12, X: 12, Y: 1}).
			WithTarget(
				prometheus.NewDataqueryBuilder().
					Expr(`sum by (status_code) (rate(requests_total{status_code!~"2.*", type="write"}[$__rate_interval])) or on() vector(0)`).
					LegendFormat("{{status_code}}").
					RefId("Write errors by status").
					Range(),
			).
			OverrideByName("400", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "400: Bad Request"},
				fixedColor("#FFF176"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).
			OverrideByName("401", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "401: Unauthorized"},
				fixedColor("#FFB3B8"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).
			OverrideByName("403", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "403: Forbidden"},
				fixedColor("#FF9E80"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).
			OverrideByName("404", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "404: Not Found"},
				fixedColor("#FFB347"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).
			OverrideByName("408", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "408: Request Timeout"},
				fixedColor("#FFD966"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).
			OverrideByName("409", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "409: Conflict"},
				fixedColor("#FFE066"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).
			OverrideByName("422", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "422: Unprocessable Entity"},
				fixedColor("#FFC1E3"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).
			OverrideByName("429", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "429: Too Many Requests"},
				fixedColor("#FFEE58"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).WithOverride(
			dashboard.MatcherConfig{Id: "byRegexp", Options: ".*"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			},
		),
	).
	WithRow(
		dashboard.NewRowBuilder("Read").
			GridPos(dashboard.GridPos{H: 1, W: 24, X: 0, Y: 8}).
			Id(8),
	).
	WithPanel(
		timeseries.NewPanelBuilder().
			Title("Read requests").
			Datasource(DatasourceRef).
			Description("Number of metrics read requests per second").
			Unit("reqps").
			GridPos(dashboard.GridPos{H: 7, W: 12, X: 0, Y: 9}).
			WithTarget(
				prometheus.NewDataqueryBuilder().
					Expr(`sum (rate(requests_total{type="read"}[$__rate_interval]))`).
					LegendFormat("Requests").
					RefId("Requests").
					Range(),
			).
			WithTarget(
				prometheus.NewDataqueryBuilder().
					Expr(`requests_limits{type="monitoring.read.throughput.requests"}`).
					LegendFormat("Limit").
					RefId("Limit").
					Range(),
			).
			OverrideByQuery("Limit", []dashboard.DynamicConfigValue{
				fixedColor("dark-red"),
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			}).WithOverride(
			dashboard.MatcherConfig{Id: "byRegexp", Options: ".*"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			},
		),
	).
	WithPanel(
		timeseries.NewPanelBuilder().
			Title("Read requests errors").
			Datasource(DatasourceRef).
			Description("Number of failed metrics read requests per second, by status code").
			Unit("reqps").
			GridPos(dashboard.GridPos{H: 7, W: 12, X: 12, Y: 9}).
			WithTarget(
				prometheus.NewDataqueryBuilder().
					Expr(`sum by (status_code) (rate(requests_total{status_code!~"2.*", type="read"}[$__rate_interval])) or on() vector(0)`).
					LegendFormat("{{status_code}}").
					RefId("Read errors by status").
					Range(),
			).
			OverrideByName("400", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "400: Bad Request"},
				fixedColor("#FFF176"),
			}).
			OverrideByName("401", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "401: Unauthorized"},
				fixedColor("#FFB3B8"),
			}).
			OverrideByName("403", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "403: Forbidden"},
				fixedColor("#FF9E80"),
			}).
			OverrideByName("404", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "404: Not Found"},
				fixedColor("#FFB347"),
			}).
			OverrideByName("408", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "408: Request Timeout"},
				fixedColor("#FFD966"),
			}).
			OverrideByName("409", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "409: Conflict"},
				fixedColor("#FFE066"),
			}).
			OverrideByName("422", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "422: Unprocessable Entity"},
				fixedColor("#FFC1E3"),
			}).
			OverrideByName("429", []dashboard.DynamicConfigValue{
				{Id: "displayName", Value: "429: Too Many Requests"},
				fixedColor("#FFEE58"),
			}).WithOverride(
			dashboard.MatcherConfig{Id: "byRegexp", Options: ".*"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			},
		),
	).
	WithRow(
		dashboard.NewRowBuilder("Workload").
			GridPos(dashboard.GridPos{H: 1, W: 24, X: 0, Y: 16}).
			Id(5),
	).
	WithPanel(
		timeseries.NewPanelBuilder().
			Title("Samples ingestion rate").
			Description("Number of samples ingested per second, by type").
			Datasource(DatasourceRef).
			Unit("rowsps").
			GridPos(dashboard.GridPos{H: 8, W: 12, X: 0, Y: 17}).
			WithTarget(
				prometheus.NewDataqueryBuilder().
					Expr(`sum by(type) (rate(samples_total{}[$__rate_interval])) OR on() vector(0)`).
					LegendFormat("{{type}}").
					RefId("Samples rate by type").
					Range(),
			).WithOverride(
			dashboard.MatcherConfig{Id: "byRegexp", Options: ".*"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.drawStyle", Value: "line"},
				{Id: "custom.showPoints", Value: "never"},
				{Id: "custom.lineWidth", Value: 1},
			},
		),
	)

func fixedColor(col string) dashboard.DynamicConfigValue {
	c, _ := dashboard.NewFieldColorBuilder().
		Mode(dashboard.FieldColorModeId("fixed")).
		FixedColor(col).
		Build()
	return dashboard.DynamicConfigValue{Id: "color", Value: c}
}
