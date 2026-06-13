package request

type CreateReportRequest struct {
	Latitude    float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude   float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	Severity    string  `json:"severity" validate:"required,oneof=low medium high critical"`
}

type UpdateReportStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=verified cleaning resolved rejected"`
}

type ListReportQuery struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Status  string `query:"status"`
}

type NearbyQuery struct {
	Latitude  float64 `query:"lat" validate:"required,min=-90,max=90"`
	Longitude float64 `query:"lng" validate:"required,min=-180,max=180"`
	RadiusKm  float64 `query:"radius" validate:"required,min=0.1,max=50"`
}
