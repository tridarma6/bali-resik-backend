package repository

import (
	"github.com/google/uuid"
)

type DashboardOverview struct {
	TotalPickups     int            `json:"total_pickups"`
	PickupsByStatus  map[string]int `json:"pickups_by_status"`
	TotalReports     int            `json:"total_reports"`
	ReportsByStatus  map[string]int `json:"reports_by_status"`
	TotalCitizens    int64          `json:"total_citizens"`
	TotalCollectors  int64          `json:"total_collectors"`
	PickupRate       float64        `json:"pickup_completion_rate"`
}

type MonthlyTrend struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Count int `json:"count"`
}

type WasteTypeDistribution struct {
	WasteType string `json:"waste_type"`
	Count     int    `json:"count"`
}

type ReportSeverityDistribution struct {
	Severity string `json:"severity"`
	Count    int    `json:"count"`
}

type RegionalStat struct {
	TenantID     uuid.UUID `json:"tenant_id"`
	TenantName   string    `json:"tenant_name"`
	TotalPickups int       `json:"total_pickups"`
	TotalReports int       `json:"total_reports"`
	TotalUsers   int64     `json:"total_users"`
}

type AnalyticsRepository interface {
	GetOverview(tenantID uuid.UUID) (*DashboardOverview, error)
	GetPickupTrends(tenantID uuid.UUID, months int) ([]MonthlyTrend, error)
	GetReportTrends(tenantID uuid.UUID, months int) ([]MonthlyTrend, error)
	GetPickupWasteTypeDistribution(tenantID uuid.UUID) ([]WasteTypeDistribution, error)
	GetReportSeverityDistribution(tenantID uuid.UUID) ([]ReportSeverityDistribution, error)
	GetNewUsersOverTime(tenantID uuid.UUID, months int) ([]MonthlyTrend, error)
	GetRegionalStats() ([]RegionalStat, error)
}
