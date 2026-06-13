package request

type ListCollectorAppQuery struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Status  string `query:"status"`
}

type ReviewCollectorAppRequest struct {
	Notes string `json:"admin_notes" validate:"omitempty"`
}
