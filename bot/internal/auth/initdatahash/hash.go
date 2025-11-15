package initdatahash

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const secretPrefix = "WebAppData"

// Compute builds the data-check string from init data values and returns the expected hex hash.
func Compute(values url.Values, botToken string) (expectedHash string, dataCheckString string, err error) {
	if botToken == "" {
		return "", "", fmt.Errorf("bot token is required")
	}
	dataCheck := BuildDataCheckString(values)
	secretKey := computeHMAC([]byte(secretPrefix), botToken)
	expected := hex.EncodeToString(computeHMAC(secretKey, dataCheck))
	return expected, dataCheck, nil
}

// BuildDataCheckString reproduces Telegram/MAX data_check_string rules.
func BuildDataCheckString(values url.Values) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		if key == "hash" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	return strings.Join(pairs, "\n")
}

// CheckFreshness validates that auth_date is not older than maxAge.
func CheckFreshness(raw string, maxAge time.Duration, now time.Time) error {
	if raw == "" {
		return fmt.Errorf("auth_date missing")
	}
	ts, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("auth_date parse: %w", err)
	}
	issuedAt := time.Unix(ts, 0)
	if now.Sub(issuedAt) > maxAge {
		return fmt.Errorf("init data expired")
	}
	return nil
}

// HashesEqual compares normalized hex hashes using constant time.
func HashesEqual(provided, expected string) bool {
	a := NormalizeHash(provided)
	b := NormalizeHash(expected)
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// NormalizeHash trims and lowercases a hex string.
func NormalizeHash(hash string) string {
	return strings.ToLower(strings.TrimSpace(hash))
}

func computeHMAC(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}
