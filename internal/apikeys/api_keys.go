package apikeys

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/crypto/argon2"
)

const argon2Memory = 64 * 1024
const argon2Time = 1
const argon2Threads = 2
const argon2KeyLength = 32

func GenerateAPIKey() (string, string, error) {
	apiKeyBytes, err := generateRandomBytes(64)
	if err != nil {
		return "", "", err
	}
	apiKey := base64.RawURLEncoding.EncodeToString(apiKeyBytes)

	saltBytes, err := generateRandomBytes(16)
	if err != nil {
		return "", "", err
	}

	hash := argon2.IDKey([]byte(apiKey), saltBytes, argon2Time, argon2Memory, argon2Threads, argon2KeyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(saltBytes)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argon2Memory, argon2Time,
		argon2Threads, b64Salt, b64Hash)

	return encodedHash, apiKey, nil
}

func VerifyAPIKey(apiKey string, encodedHash string) (bool, error) {
	mem, time, p, salt, hash, err := decodeEncodedHash(encodedHash)
	if err != nil {
		return false, err
	}

	argon2Hash := argon2.IDKey([]byte(apiKey), salt, time, mem, p, argon2KeyLength)

	match := subtle.ConstantTimeCompare(argon2Hash, hash)
	if match == 1 {
		return true, nil
	}

	return false, nil
}

func decodeEncodedHash(key string) (uint32, uint32, uint8, []byte, []byte, error) {
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

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
