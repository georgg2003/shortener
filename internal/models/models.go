package models

type DeleteUserURLsTask struct {
	ShortIDs []string
	UserID   int64
}
