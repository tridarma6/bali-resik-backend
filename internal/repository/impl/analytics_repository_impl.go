package impl

import (
	"time"

	"github.com/google/uuid"
	"github.com/indim/bali-resik-backend/internal/interfaces/repository"
	"gorm.io/gorm"
)

type AnalyticsRepositoryImpl struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepositoryImpl {
	return &AnalyticsRepositoryImpl{db: db}
}

func (r *AnalyticsRepositoryImpl) GetOverview(tenantID uuid.UUID) (*repository.DashboardOverview, error) {
	result := &repository.DashboardOverview{
		PickupsByStatus: make(map[string]int),
		ReportsByStatus: make(map[string]int),
	}

	var totalPickups int64
	r.db.Model(&struct{}{}).Table("pickup_requests").
		Where("tenant_id = ?", tenantID).
		Count(&totalPickups)
	result.TotalPickups = int(totalPickups)

	type statusCount struct {
		Status string
		Count  int
	}
	var pickupStatuses []statusCount
	r.db.Model(&struct{}{}).Table("pickup_requests").
		Select("status, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("status").
		Scan(&pickupStatuses)
	for _, s := range pickupStatuses {
		result.PickupsByStatus[s.Status] = s.Count
	}

	var totalReports int64
	r.db.Model(&struct{}{}).Table("waste_reports").
		Where("tenant_id = ?", tenantID).
		Count(&totalReports)
	result.TotalReports = int(totalReports)

	var reportStatuses []statusCount
	r.db.Model(&struct{}{}).Table("waste_reports").
		Select("status, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("status").
		Scan(&reportStatuses)
	for _, s := range reportStatuses {
		result.ReportsByStatus[s.Status] = s.Count
	}

	r.db.Model(&struct{}{}).Table("users").
		Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
		Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
		Where("users.tenant_id = ? AND roles.name = ?", tenantID, "citizen").
		Count(&result.TotalCitizens)

	r.db.Model(&struct{}{}).Table("users").
		Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
		Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
		Where("users.tenant_id = ? AND roles.name = ?", tenantID, "collector").
		Count(&result.TotalCollectors)

	if totalPickups > 0 {
		var completed int64
		r.db.Model(&struct{}{}).Table("pickup_requests").
			Where("tenant_id = ? AND status = ?", tenantID, "completed").
			Count(&completed)
		result.PickupRate = float64(completed) / float64(totalPickups) * 100
	}

	return result, nil
}

func (r *AnalyticsRepositoryImpl) GetPickupTrends(tenantID uuid.UUID, months int) ([]repository.MonthlyTrend, error) {
	dateCutoff := time.Now().AddDate(0, -months, 0)

	var trends []repository.MonthlyTrend
	err := r.db.Model(&struct{}{}).Table("pickup_requests").
		Select("EXTRACT(YEAR FROM created_at)::int as year, EXTRACT(MONTH FROM created_at)::int as month, count(*)::int as count").
		Where("tenant_id = ? AND created_at >= ?", tenantID, dateCutoff).
		Group("year, month").
		Order("year ASC, month ASC").
		Scan(&trends).Error

	return trends, err
}

func (r *AnalyticsRepositoryImpl) GetReportTrends(tenantID uuid.UUID, months int) ([]repository.MonthlyTrend, error) {
	dateCutoff := time.Now().AddDate(0, -months, 0)

	var trends []repository.MonthlyTrend
	err := r.db.Model(&struct{}{}).Table("waste_reports").
		Select("EXTRACT(YEAR FROM created_at)::int as year, EXTRACT(MONTH FROM created_at)::int as month, count(*)::int as count").
		Where("tenant_id = ? AND created_at >= ?", tenantID, dateCutoff).
		Group("year, month").
		Order("year ASC, month ASC").
		Scan(&trends).Error

	return trends, err
}

func (r *AnalyticsRepositoryImpl) GetPickupWasteTypeDistribution(tenantID uuid.UUID) ([]repository.WasteTypeDistribution, error) {
	var dist []repository.WasteTypeDistribution
	err := r.db.Model(&struct{}{}).Table("pickup_requests").
		Select("waste_type, count(*)::int as count").
		Where("tenant_id = ?", tenantID).
		Group("waste_type").
		Order("count DESC").
		Scan(&dist).Error

	return dist, err
}

func (r *AnalyticsRepositoryImpl) GetReportSeverityDistribution(tenantID uuid.UUID) ([]repository.ReportSeverityDistribution, error) {
	var dist []repository.ReportSeverityDistribution
	err := r.db.Model(&struct{}{}).Table("waste_reports").
		Select("severity, count(*)::int as count").
		Where("tenant_id = ?", tenantID).
		Group("severity").
		Order("count DESC").
		Scan(&dist).Error

	return dist, err
}

func (r *AnalyticsRepositoryImpl) GetNewUsersOverTime(tenantID uuid.UUID, months int) ([]repository.MonthlyTrend, error) {
	dateCutoff := time.Now().AddDate(0, -months, 0)

	var trends []repository.MonthlyTrend
	err := r.db.Model(&struct{}{}).Table("users").
		Select("EXTRACT(YEAR FROM created_at)::int as year, EXTRACT(MONTH FROM created_at)::int as month, count(*)::int as count").
		Where("tenant_id = ? AND created_at >= ?", tenantID, dateCutoff).
		Group("year, month").
		Order("year ASC, month ASC").
		Scan(&trends).Error

	return trends, err
}

func (r *AnalyticsRepositoryImpl) GetRegionalStats() ([]repository.RegionalStat, error) {
	var stats []repository.RegionalStat
	err := r.db.Model(&struct{}{}).Table("tenants").
		Select(`
			tenants.id as tenant_id,
			tenants.name as tenant_name,
			COALESCE(pickup_counts.total, 0) as total_pickups,
			COALESCE(report_counts.total, 0) as total_reports,
			COALESCE(user_counts.total, 0) as total_users
		`).
		Joins(`LEFT JOIN (
			SELECT tenant_id, count(*)::int as total FROM pickup_requests GROUP BY tenant_id
		) pickup_counts ON pickup_counts.tenant_id = tenants.id`).
		Joins(`LEFT JOIN (
			SELECT tenant_id, count(*)::int as total FROM waste_reports GROUP BY tenant_id
		) report_counts ON report_counts.tenant_id = tenants.id`).
		Joins(`LEFT JOIN (
			SELECT tenant_id, count(*)::int as total FROM users GROUP BY tenant_id
		) user_counts ON user_counts.tenant_id = tenants.id`).
		Order("tenants.name ASC").
		Scan(&stats).Error

	return stats, err
}
