package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/auth"
	"github.com/suse-skyscraper/skyscraper/internal/auth/apikeys"
	"github.com/suse-skyscraper/skyscraper/internal/auth/jwt"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func AuthorizationHandler(app *application.App) func(next http.Handler) http.Handler {
	authorizer := newAuthorizer(app)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := authHeaders{
				Authorization: r.Header.Get("Authorization"),
				APIKeyID:      r.Header.Get("X-API-Key"),
				APIKeySecret:  r.Header.Get("X-API-Secret"),
			}

			caller, match, err := authorizer.authorize(r.Context(), headers)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			} else if !match {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized")
				return
			}

			ctx := context.WithValue(r.Context(), ContextCaller, caller)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

type authHeaders struct {
	Authorization string `json:"authorization"`
	APIKeyID      string `json:"api_key_id"`
	APIKeySecret  string `json:"api_key_secret"`
}

type apiAuthorizer struct {
	jwtVerifier jwt.Verifier
	apiVerifier apikeys.Verifier
}

func newAuthorizer(app *application.App) apiAuthorizer {
	jwtVerifier := jwt.NewVerifier(app)
	apiVerifier := apikeys.NewVerifier(app)

	return apiAuthorizer{jwtVerifier, apiVerifier}
}

func (a *apiAuthorizer) authorize(ctx context.Context, headers authHeaders) (auth.Caller, bool, error) {
	if headers.Authorization != "" {
		return a.jwtVerifier.Verify(ctx, headers.Authorization)
	}

	if headers.APIKeyID != "" && headers.APIKeySecret != "" {
		return a.apiVerifier.Verify(ctx, headers.APIKeyID, headers.APIKeySecret)
	}

	return auth.Caller{}, false, nil
}
