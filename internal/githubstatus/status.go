package githubstatus

// IndicatorSeverity maps the page-level status indicator to a numeric severity.
// none=0, minor=1, major=2, critical=3, unknown=-1
func IndicatorSeverity(indicator string) float64 {
	switch indicator {
	case "none":
		return 0
	case "minor":
		return 1
	case "major":
		return 2
	case "critical":
		return 3
	default:
		return -1
	}
}

// ComponentSeverity maps component status to a numeric severity.
// operational=0, degraded_performance=1, partial_outage=2, major_outage=3, unknown=-1
func ComponentSeverity(status string) float64 {
	switch status {
	case "operational":
		return 0
	case "degraded_performance":
		return 1
	case "partial_outage":
		return 2
	case "major_outage":
		return 3
	default:
		return -1
	}
}

// ImpactSeverity maps incident impact to a numeric severity.
// none=0, minor=1, major=2, critical=3, unknown=-1
func ImpactSeverity(impact string) float64 {
	switch impact {
	case "none":
		return 0
	case "minor":
		return 1
	case "major":
		return 2
	case "critical":
		return 3
	default:
		return -1
	}
}
