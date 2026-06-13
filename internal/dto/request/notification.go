package request

type ListNotificationQuery struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Read    string `query:"read"`
}
