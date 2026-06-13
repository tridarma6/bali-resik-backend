package response

import (
	"github.com/google/uuid"
)

type DashboardOverviewResponse struct {
	TotalPickups     int            `json:"total_pickups"`
	PickupsByStatus  map[string]int `json:"pickups_by_status"`
	TotalReports     int            `json:"total_reports"`
	ReportsByStatus  map[string]int `json:"reports_by_status"`
	TotalCitizens    int64          `json:"total_citizens"`
	TotalCollectors  int64          `json:"total_collectors"`
	PickupRate       float64        `json:"pickup_completion_rate"`
}

type MonthlyTrendResponse struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Count int `json:"count"`
}

type WasteTypeDistResponse struct {
	WasteType string `json:"waste_type"`
	Count     int    `json:"count"`
}

type SeverityDistResponse struct {
	Severity string `json:"severity"`
	Count    int    `json:"count"`
}

type RegionalStatResponse struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	TenantName   string    `json:"tenant_name"`
	TotalPickups int       `json:"total_pickups"`
	TotalReports int       `json:"total_reports"`
	TotalUsers   int64     `json:"total_users"`
}

type AnalyticsResponse struct {
	Overview  *DashboardOverviewResponse `json:"overview,omitempty"`
	PickupTrends  []MonthlyTrendResponse `json:"pickup_trends,omitempty"`
	ReportTrends  []MonthlyTrendResponse `json:"report_trends,omitempty"`
	WasteTypeDist []WasteTypeDistResponse `json:"waste_type_distribution,omitempty"`
	SeverityDist  []SeverityDistResponse  `json:"severity_distribution,omitempty"`
	NewUsers      []MonthlyTrendResponse `json:"new_users,omitempty"`
}
