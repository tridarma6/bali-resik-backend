package request

type CreatePickupRequest struct {
	WasteType     string  `json:"waste_type" validate:"required,oneof=organic anorganic mixed electronic hazardous"`
	Latitude      float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude     float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Address       string  `json:"address" validate:"required,min=5,max=500"`
	ScheduledDate string  `json:"scheduled_date,omitempty"`
	Notes         string  `json:"notes,omitempty" validate:"omitempty,max=500"`
}

type UpdatePickupStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=in_progress completed cancelled"`
}

type AssignCollectorRequest struct {
	CollectorID string `json:"collector_id" validate:"required,uuid"`
}

type ListPickupQuery struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Status  string `query:"status"`
}
