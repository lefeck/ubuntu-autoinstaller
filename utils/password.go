package utils

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

var sha512CryptPattern = regexp.MustCompile(`^\$6\$[A-Za-z0-9./]+\$[A-Za-z0-9./]+$`)

// IsSHA512Crypt checks if the given password string already looks like a SHA-512 crypt hash ($6$...)
func IsSHA512Crypt(s string) bool {
	return sha512CryptPattern.MatchString(strings.TrimSpace(s))
}

// HashSHA512Crypt uses openssl to generate a SHA-512 crypt hash for the given password.
// It requires that `openssl` is available in PATH.
func HashSHA512Crypt(plain string) (string, error) {
	if plain == "" {
		return "", errors.New("password is empty")
	}

	cmd := exec.Command("openssl", "passwd", "-6", plain)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(strings.TrimSpace(stderr.String()))
		}
		return "", err
	}
	out := strings.TrimSpace(stdout.String())
	if out == "" {
		return "", errors.New("empty hash output from openssl")
	}
	return out, nil
}
