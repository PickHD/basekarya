package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

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
		suffix[i] = alphanum[rng.Intn(len(alphanum))]
	}

	return fmt.Sprintf("%s_%s", firstName, string(suffix))
}
