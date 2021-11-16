package vo

type LogVo struct {
	Message    string `json:"message"`
	UserID     int    `json:"user_id"`
	Username   string `json:"username"`
	StatusCode int    `json:"status_code"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
}
