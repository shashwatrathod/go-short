package middleware

import (
	"context"
	"net/http"

	"github.com/shashwatrathod/url-shortner/cache"
	"github.com/shashwatrathod/url-shortner/core"
	"github.com/shashwatrathod/url-shortner/db"
	"github.com/shashwatrathod/url-shortner/db/dao"
)

type AppEnv struct {
	DBManager        *db.ConnectionManager
	UrlAliasDao      dao.UrlAliasDao
	AliasingStrategy core.AliasingStrategy
	CacheManager     cache.CacheManager
}

func NewAppEnv(dbManager *db.ConnectionManager, cacheManager cache.CacheManager) *AppEnv {
	return &AppEnv{
		DBManager:        dbManager,
		UrlAliasDao:      dao.NewUrlAliasDao(dbManager),
		AliasingStrategy: core.NewSimpleAliasingStrategy(),
		CacheManager:     cacheManager,
	}
}

// define a custom context key type for context injection
type contextKey string

// ContextAppEnvKey is the key used to store AppEnv in the context.
const ContextAppEnvKey contextKey = "appEnv"

func ContextMiddleware(appEnv *AppEnv) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextAppEnvKey, appEnv)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
