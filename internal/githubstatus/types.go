package githubstatus

type Page struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	UpdatedAt string `json:"updated_at"`
}

type StatusValue struct {
	Indicator   string `json:"indicator"`
	Description string `json:"description"`
}

type StatusResponse struct {
	Page   Page        `json:"page"`
	Status StatusValue `json:"status"`
}

type Component struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	Description *string `json:"description"`
	UpdatedAt   string  `json:"updated_at"`
}

type ComponentsResponse struct {
	Page       Page        `json:"page"`
	Components []Component `json:"components"`
}

type Incident struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Impact     string `json:"impact"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	StartedAt  string `json:"started_at"`
	ResolvedAt *string `json:"resolved_at"`
	Shortlink  string `json:"shortlink"`
}

type IncidentsResponse struct {
	Page      Page       `json:"page"`
	Incidents []Incident `json:"incidents"`
}

type ScheduledMaintenance struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	Impact         string  `json:"impact"`
	ScheduledFor   string  `json:"scheduled_for"`
	ScheduledUntil string  `json:"scheduled_until"`
	UpdatedAt      string  `json:"updated_at"`
	ResolvedAt     *string `json:"resolved_at"`
}

type ScheduledMaintenancesResponse struct {
	Page                  Page                   `json:"page"`
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
}
