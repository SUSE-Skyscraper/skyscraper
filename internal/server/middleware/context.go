package middleware

type key int

const (
	ContextCaller key = iota
	ContextCloudAccount
	ContextTag
	ContextUser
)
