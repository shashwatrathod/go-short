package middleware

import (
	"context"
	"net/http"

	"github.com/shashwatrathod/url-shortner/core"
	"github.com/shashwatrathod/url-shortner/db"
	"github.com/shashwatrathod/url-shortner/db/dao"
)

type AppEnv struct {
	DBManager *db.ConnectionManager
	UrlAliasDao dao.UrlAliasDao
	AliasingStrategy core.AliasingStrategy
}

func NewAppEnv(dbManager *db.ConnectionManager) *AppEnv {
	return &AppEnv{
        DBManager:   dbManager,
        UrlAliasDao: dao.NewUrlAliasDao(dbManager),
		AliasingStrategy: core.NewSimpleAliasingStrategy(),
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