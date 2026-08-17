package githubstatus

import "testing"

func TestIndicatorSeverity(t *testing.T) {
	tests := map[string]float64{
		"none":     0,
		"minor":    1,
		"major":    2,
		"critical": 3,
		"unknown":  -1,
	}

	for input, want := range tests {
		if got := IndicatorSeverity(input); got != want {
			t.Fatalf("IndicatorSeverity(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestComponentSeverity(t *testing.T) {
	tests := map[string]float64{
		"operational":          0,
		"degraded_performance":   1,
		"partial_outage":         2,
		"major_outage":           3,
		"unknown":                -1,
	}

	for input, want := range tests {
		if got := ComponentSeverity(input); got != want {
			t.Fatalf("ComponentSeverity(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestImpactSeverity(t *testing.T) {
	tests := map[string]float64{
		"none":     0,
		"minor":    1,
		"major":    2,
		"critical": 3,
		"unknown":  -1,
	}

	for input, want := range tests {
		if got := ImpactSeverity(input); got != want {
			t.Fatalf("ImpactSeverity(%q) = %v, want %v", input, got, want)
		}
	}
}
