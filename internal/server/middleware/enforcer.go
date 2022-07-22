package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/auth"
)

func EnforcerHandler(app *application.App) func(next http.Handler) http.Handler {
	enforcer := newEnforcer(app)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result, err := enforcer.Enforce(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else if !result {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type apiEnforcer struct {
	enforcer *casbin.Enforcer
}

func newEnforcer(app *application.App) apiEnforcer {
	return apiEnforcer{app.Enforcer}
}

func (a *apiEnforcer) Enforce(r *http.Request) (bool, error) {
	caller, ok := r.Context().Value(ContextCaller).(auth.Caller)
	if !ok {
		return false, nil
	}

	result, err := a.enforcer.Enforce(caller.ID.String(), r.URL.Path, r.Method)
	if err != nil {
		return false, err
	}

	return result, nil
}
