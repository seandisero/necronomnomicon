package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, authParams.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey(
		[]byte(password),
		salt,
		authParams.time,
		authParams.memory,
		authParams.threads,
		authParams.keyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(key)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		authParams.memory,
		authParams.time,
		authParams.threads,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func ValidatePassword(password, hash string) (bool, error) {
	hashArgs := strings.Split(hash, "$")
	if len(hashArgs) != 6 {
		return false, fmt.Errorf("malformed password hash")
	}

	var version int
	var memory int
	var time int
	var threads int

	_, err := fmt.Sscanf(hashArgs[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	_, err = fmt.Sscanf(
		hashArgs[3],
		"m=%d,t=%d,p=%d",
		&memory, &time, &threads,
	)
	if err != nil {
		return false, err
	}

	encodedSalt := hashArgs[4]
	encodedHash := hashArgs[5]

	decodedSalt, err := base64.RawStdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, err
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}

	if version != argon2.Version {
		return false, fmt.Errorf("mismatch argon2 version")
	}

	pHash := argon2.IDKey([]byte(password), []byte(decodedSalt),
		uint32(time),
		uint32(memory),
		uint8(threads),
		authParams.keyLength,
	)
	if subtle.ConstantTimeCompare(decodedHash, []byte(pHash)) == 1 {
		return true, nil
	}
	return false, nil
}
