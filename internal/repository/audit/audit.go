package audit

type AuditLog struct {
	Timestamp int `json:"ts"`
	Action    string
	UserID    string `json:"user_id"`
	URL       string `json:"url"`
}

//go:generate mockgen -destination ./mock/mock.go -package mock . AuditRepository
type AuditRepository interface {
	WriteLog(AuditLog) error
}
