package jwt

import (
	"context"
	"fmt"
	"strings"

	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

type verifier interface {
	VerifyIdToken(jwt string) (*jwtverifier.Jwt, error)
}

type Verifier struct {
	verifier verifier
	app      *application.App
}

func NewVerifier(app *application.App) Verifier {
	toValidate := map[string]string{}
	toValidate["aud"] = "api://default"
	toValidate["cid"] = app.Config.Okta.ClientID

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer:           app.Config.Okta.Issuer,
		ClaimsToValidate: toValidate,
	}

	v := jwtVerifierSetup.New()

	return Verifier{
		verifier: v,
		app:      app,
	}
}

func (v *Verifier) Verify(ctx context.Context, authorizationHeader string) (auth.Caller, bool, error) {
	authHeader := strings.Split(authorizationHeader, "Bearer ")
	if len(authHeader) != 2 {
		return auth.Caller{}, false, nil
	}

	jwtToken := authHeader[1]
	claims, err := v.verifier.VerifyIdToken(jwtToken)
	if err != nil {
		return auth.Caller{}, false, err
	}

	username := claims.Claims["sub"].(string)
	user, err := v.app.Repo.FindUserByUsername(ctx, username)
	if err != nil {
		return auth.Caller{}, false, err
	}

	if !user.Active {
		return auth.Caller{}, false, fmt.Errorf("user is not active")
	}

	caller := auth.Caller{
		ID:   user.ID,
		Type: auth.CallerUser,
	}

	return caller, true, nil
}
