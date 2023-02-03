package middleware

type key int

const (
	_ key = iota
	ContextCaller
	ContextCloudAccount
	ContextTag
	ContextUser
	ContextAPIKey
	ContextOrganizationalUnit
	ContextTenant
)
