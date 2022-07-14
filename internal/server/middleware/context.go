package middleware

type key int

const (
	CurrentUser key = iota
	CloudAccount
	Tag
	User
)
