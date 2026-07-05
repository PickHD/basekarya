package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const alphanum = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateUsername(fullName string) string {
	parts := strings.Fields(strings.ToLower(fullName))
	firstName := ""
	if len(parts) > 0 {
		firstName = parts[0]
	}

	if len(firstName) > 8 {
		firstName = firstName[:8]
	}

	suffix := make([]byte, 6)
	for i := range suffix {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphanum))))
		suffix[i] = alphanum[n.Int64()]
	}

	return fmt.Sprintf("%s_%s", firstName, string(suffix))
}
