package constants

const (
	WasteTypeOrganic     = "organic"
	WasteTypeAnorganic   = "anorganic"
	WasteTypeMixed       = "mixed"
	WasteTypeElectronic  = "electronic"
	WasteTypeHazardous   = "hazardous"

	PickupStatusPending   = "pending"
	PickupStatusAssigned  = "assigned"
	PickupStatusInProgress = "in_progress"
	PickupStatusCompleted = "completed"
	PickupStatusCancelled = "cancelled"

	ReportStatusReported  = "reported"
	ReportStatusVerified  = "verified"
	ReportStatusCleaning  = "cleaning"
	ReportStatusResolved  = "resolved"
	ReportStatusRejected  = "rejected"

	NotificationTypePickup    = "pickup"
	NotificationTypeReport    = "report"
	NotificationTypeReward    = "reward"
	NotificationTypeSystem    = "system"

	RewardTypeEarn   = "earn"
	RewardTypeRedeem = "redeem"

	RoleSuperAdmin    = "super_admin"
	RoleAdminKabupaten = "admin_kabupaten"
	RoleCitizen       = "citizen"
	RoleCollector     = "collector"

	RegionTypeKota      = "kota"
	RegionTypeKabupaten = "kabupaten"
)

var ValidWasteTypes = []string{
	WasteTypeOrganic, WasteTypeAnorganic, WasteTypeMixed,
	WasteTypeElectronic, WasteTypeHazardous,
}

var ValidPickupStatuses = []string{
	PickupStatusPending, PickupStatusAssigned, PickupStatusInProgress,
	PickupStatusCompleted, PickupStatusCancelled,
}

var ValidReportStatuses = []string{
	ReportStatusReported, ReportStatusVerified, ReportStatusCleaning,
	ReportStatusResolved, ReportStatusRejected,
}
