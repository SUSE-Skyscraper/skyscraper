package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/suse-skyscraper/skyscraper-web/internal/application"
)

func AuthorizationHandler(conf application.Config) func(next http.Handler) http.Handler {
	jwtVerifier := NewJwtVerifier(conf.Okta.Issuer, conf.Okta.ClientID)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/*
				if r.Method == "OPTION" {
					next.ServeHTTP(w, r)
					return
				}*/

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

			userEmail := claims.Claims["sub"]
			ctx := context.WithValue(r.Context(), UserEmail, userEmail)
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
