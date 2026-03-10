package models

type DeleteUserURLsTask struct {
	ShortIDs []string
	UserID   int64
}

type Stats struct {
	UsersCount int64
	URLsCount  int64
}
