package apikeys

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/auth"
	"golang.org/x/crypto/argon2"
)

type Verifier struct {
	App *application.App
}

func NewVerifier(app *application.App) Verifier {
	return Verifier{
		App: app,
	}
}

func (v *Verifier) Verify(ctx context.Context, keyID, keySecret string) (auth.Caller, bool, error) {
	id, err := uuid.Parse(keyID)
	if err != nil {
		return auth.Caller{}, false, nil
	}

	apiKey, err := v.App.Repository.FindAPIKey(ctx, id)
	if err != nil {
		return auth.Caller{}, false, err
	}

	match, err := CompareArgon2Hash(keySecret, apiKey.Encodedhash)
	if err != nil {
		return auth.Caller{}, false, err
	}

	caller := auth.Caller{
		ID:   apiKey.ID,
		Type: auth.CallerAPIKey,
	}

	return caller, match, nil
}

func (v *Verifier) VerifyScim(ctx context.Context, authorizationHeader string) (bool, error) {
	bearer := strings.Split(authorizationHeader, "Bearer ")
	if len(bearer) != 2 {
		return false, nil
	}
	token := bearer[1]

	apiKey, err := v.App.Repository.FindScimAPIKey(ctx)
	if err != nil && err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return CompareArgon2Hash(token, apiKey.Encodedhash)
}

func CompareArgon2Hash(key string, encodedHash string) (bool, error) {
	mem, time, p, salt, hash, err := DecodeArgon2Hash(encodedHash)
	if err != nil {
		return false, err
	}

	argon2Hash := argon2.IDKey([]byte(key), salt, time, mem, p, argon2KeyLength)

	match := subtle.ConstantTimeCompare(argon2Hash, hash)
	if match == 1 {
		return true, nil
	}

	return false, nil
}

func DecodeArgon2Hash(key string) (uint32, uint32, uint8, []byte, []byte, error) {
	regex := regexp.MustCompile(`^\$argon2id\$v=\d+\$m=(\d+),t=(\d+),p=(\d+)\$(.*)\$(.*)$`)
	if !regex.MatchString(key) {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid encoded hash")
	}

	groups := regex.FindStringSubmatch(key)

	m, err := strconv.Atoi(groups[1])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	t, err := strconv.Atoi(groups[2])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	p, err := strconv.Atoi(groups[3])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(groups[4])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	hash, err := base64.RawStdEncoding.Strict().DecodeString(groups[5])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	return uint32(m), uint32(t), uint8(p), salt, hash, nil
}
