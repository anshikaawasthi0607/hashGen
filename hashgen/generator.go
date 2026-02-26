package hashgen

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"
)

func Generate(input string) (string, error) {
	if err := validateInput(input); err != nil {
		return "", err
	}

	digest := sha256.Sum256([]byte(input))

	num := binary.BigEndian.Uint64(digest[:8]) //slicing

	return padOrTruncate(encodeBase62(num), hashLength), nil
}

func validateInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}

	for _, r := range input {
		isLower := r >= 'a' && r <= 'z'
		isUpper := r >= 'A' && r <= 'Z'
		isDigit := r >= '0' && r <= '9'
		if !isLower && !isUpper && !isDigit {
			return fmt.Errorf("input must contain only alphanumeric characters (a-z, A-Z, 0-9)")
		}
	}

	return nil
}