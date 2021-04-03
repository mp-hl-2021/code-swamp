package usecases

type Interface interface {
	IssueToken(userId uint) (string, error)
	UserIdByToken(token string) (uint, error)
}