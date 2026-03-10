package audit

import "context"

type AuditLog struct {
	Timestamp int64  `json:"ts"`
	Action    string `json:"action"`
	UserID    string `json:"user_id"`
	URL       string `json:"url"`
}

//go:generate go tool mockgen -destination ./mock/mock.go -package mock . AuditRepository
type AuditRepository interface {
	WriteLog(ctx context.Context, log AuditLog) error
}
