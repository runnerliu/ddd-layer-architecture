package auth

// Auth 鉴权接口
type Auth interface {
	CheckPermission(action, resource, domain string) error
}
