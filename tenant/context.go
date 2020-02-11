package tenant

import "context"

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// tenantIDKey is the key for tenantID values in Contexts
// It is unexported; clients use tenant.NewContext() and tenant.FromContext() instead of using this key directly
var tenantIDKey key

// NewContext returns a new Context that carries the tenantID
func NewContext(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// FromContext returns the tenantID stored in ctx, if any
func FromContext(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(tenantIDKey).(string)
	return tenantID, ok
}
