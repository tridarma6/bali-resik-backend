package impl

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/dto/response"
	repo "github.com/indim/bali-resik-backend/internal/interfaces/repository"
)

type AnalyticsUseCaseImpl struct {
	analyticsRepo repo.AnalyticsRepository
	log           *logrus.Logger
}

func NewAnalyticsUseCase(
	analyticsRepo repo.AnalyticsRepository,
	log *logrus.Logger,
) *AnalyticsUseCaseImpl {
	return &AnalyticsUseCaseImpl{
		analyticsRepo: analyticsRepo,
		log:           log,
	}
}

const defaultTrendMonths = 12

func (uc *AnalyticsUseCaseImpl) GetDashboard(tenantID uuid.UUID) (*response.AnalyticsResponse, error) {
	overview, err := uc.analyticsRepo.GetOverview(tenantID)
	if err != nil {
		uc.log.WithError(err).Error("failed to get dashboard overview")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve dashboard data", err)
	}

	pickupTrends, err := uc.analyticsRepo.GetPickupTrends(tenantID, defaultTrendMonths)
	if err != nil {
		uc.log.WithError(err).Error("failed to get pickup trends")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve pickup trends", err)
	}

	reportTrends, err := uc.analyticsRepo.GetReportTrends(tenantID, defaultTrendMonths)
	if err != nil {
		uc.log.WithError(err).Error("failed to get report trends")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve report trends", err)
	}

	wasteDist, err := uc.analyticsRepo.GetPickupWasteTypeDistribution(tenantID)
	if err != nil {
		uc.log.WithError(err).Error("failed to get waste type distribution")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve waste type distribution", err)
	}

	severityDist, err := uc.analyticsRepo.GetReportSeverityDistribution(tenantID)
	if err != nil {
		uc.log.WithError(err).Error("failed to get severity distribution")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve severity distribution", err)
	}

	newUsers, err := uc.analyticsRepo.GetNewUsersOverTime(tenantID, defaultTrendMonths)
	if err != nil {
		uc.log.WithError(err).Error("failed to get new users trend")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve new user data", err)
	}

	return &response.AnalyticsResponse{
		Overview: &response.DashboardOverviewResponse{
			TotalPickups:    overview.TotalPickups,
			PickupsByStatus: overview.PickupsByStatus,
			TotalReports:    overview.TotalReports,
			ReportsByStatus: overview.ReportsByStatus,
			TotalCitizens:   overview.TotalCitizens,
			TotalCollectors: overview.TotalCollectors,
			PickupRate:      overview.PickupRate,
		},
		PickupTrends:  toMonthlyTrendResponse(pickupTrends),
		ReportTrends:  toMonthlyTrendResponse(reportTrends),
		WasteTypeDist: toWasteTypeDistResponse(wasteDist),
		SeverityDist:  toSeverityDistResponse(severityDist),
		NewUsers:      toMonthlyTrendResponse(newUsers),
	}, nil
}

func (uc *AnalyticsUseCaseImpl) GetRegionalStats() ([]response.RegionalStatResponse, error) {
	stats, err := uc.analyticsRepo.GetRegionalStats()
	if err != nil {
		uc.log.WithError(err).Error("failed to get regional stats")
		return nil, errors.NewDomainError("INTERNAL_ERROR", "Failed to retrieve regional statistics", err)
	}

	result := make([]response.RegionalStatResponse, 0, len(stats))
	for _, s := range stats {
		result = append(result, response.RegionalStatResponse{
			TenantID:     s.TenantID,
			TenantName:   s.TenantName,
			TotalPickups: s.TotalPickups,
			TotalReports: s.TotalReports,
			TotalUsers:   s.TotalUsers,
		})
	}

	return result, nil
}

func toMonthlyTrendResponse(trends []repo.MonthlyTrend) []response.MonthlyTrendResponse {
	result := make([]response.MonthlyTrendResponse, 0, len(trends))
	for _, t := range trends {
		result = append(result, response.MonthlyTrendResponse{
			Year:  t.Year,
			Month: t.Month,
			Count: t.Count,
		})
	}
	return result
}

func toWasteTypeDistResponse(dist []repo.WasteTypeDistribution) []response.WasteTypeDistResponse {
	result := make([]response.WasteTypeDistResponse, 0, len(dist))
	for _, d := range dist {
		result = append(result, response.WasteTypeDistResponse{
			WasteType: d.WasteType,
			Count:     d.Count,
		})
	}
	return result
}

func toSeverityDistResponse(dist []repo.ReportSeverityDistribution) []response.SeverityDistResponse {
	result := make([]response.SeverityDistResponse, 0, len(dist))
	for _, d := range dist {
		result = append(result, response.SeverityDistResponse{
			Severity: d.Severity,
			Count:    d.Count,
		})
	}
	return result
}
