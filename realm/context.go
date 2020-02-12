package realm

import "context"

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// realmIDKey is the key for realmID values in Contexts
// It is unexported; clients use realm.NewContext() and realm.FromContext() instead of using this key directly
var realmIDKey key

// NewContext returns a new Context that carries the realmID
func NewContext(ctx context.Context, realmID string) context.Context {
	return context.WithValue(ctx, realmIDKey, realmID)
}

// FromContext returns the realmID stored in ctx, if any
func FromContext(ctx context.Context) (string, bool) {
	realmID, ok := ctx.Value(realmIDKey).(string)
	return realmID, ok
}
