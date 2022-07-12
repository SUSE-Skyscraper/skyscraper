package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func OktaAuthorizationHandler(app *application.App) func(next http.Handler) http.Handler {
	jwtVerifier := NewJwtVerifier(app.Config.Okta.Issuer, app.Config.Okta.ClientID)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(authHeader) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized: authorization header malformed")
				return
			}

			jwtToken := authHeader[1]
			claims, err := jwtVerifier.verifier.VerifyIdToken(jwtToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized")
				return
			}

			username := claims.Claims["sub"].(string)

			user, err := app.Repository.FindUserByUsername(r.Context(), username)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), User, user)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

type verifier interface {
	VerifyIdToken(jwt string) (*jwtverifier.Jwt, error)
}

type JwtVerifier struct {
	verifier verifier
}

func NewJwtVerifier(oktaIssuer string, oktaClientID string) JwtVerifier {
	toValidate := map[string]string{}
	toValidate["aud"] = "api://default"
	toValidate["cid"] = oktaClientID

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer:           oktaIssuer,
		ClaimsToValidate: toValidate,
	}

	v := jwtVerifierSetup.New()

	return JwtVerifier{
		verifier: v,
	}
}
