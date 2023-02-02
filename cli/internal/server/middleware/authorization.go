package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth/apikeys"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth/jwt"

	"github.com/go-chi/render"
)

func BearerAuthorizationHandler(app *application.App) func(next http.Handler) http.Handler {
	verifier := apikeys.NewVerifier(app)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			match, err := verifier.VerifyScim(r.Context(), authorizationHeader)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			} else if !match {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

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
