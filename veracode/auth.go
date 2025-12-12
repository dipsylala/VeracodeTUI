package veracode

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const veracodeRequestVersionString = "vcode_request_version_1"

func GenerateAuthHeader(apiKeyID, apiKeySecret, httpMethod, requestURL string) (string, error) {
	// Parse the URL to get the path and query
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Get the current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(1000000)
	timestampStr := strconv.FormatInt(timestamp, 10)

	// Generate a random nonce (16 bytes)
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Build the data string for signing
	data := fmt.Sprintf("id=%s&host=%s&url=%s&method=%s",
		apiKeyID,
		parsedURL.Hostname(),
		parsedURL.RequestURI(),
		httpMethod,
	)

	// Decode the API key secret from hex
	keyBytes, err := hex.DecodeString(apiKeySecret)
	if err != nil {
		return "", fmt.Errorf("failed to decode API key secret: %w", err)
	}

	// Calculate signature using the correct HMAC chain
	signature := calculateSignature(keyBytes, nonce, []byte(timestampStr), []byte(data))

	// Build the authorization header
	authHeader := fmt.Sprintf("VERACODE-HMAC-SHA-256 id=%s,ts=%s,nonce=%X,sig=%X",
		apiKeyID,
		timestampStr,
		nonce,
		signature,
	)

	return authHeader, nil
}

func calculateSignature(key, nonce, timestamp, data []byte) []byte {
	// First HMAC: HMAC-SHA256(nonce, key)
	encryptedNonce := hmac256(nonce, key)

	// Second HMAC: HMAC-SHA256(timestamp, encryptedNonce)
	encryptedTimestamp := hmac256(timestamp, encryptedNonce)

	// Third HMAC: HMAC-SHA256(version_string, encryptedTimestamp)
	signingKey := hmac256([]byte(veracodeRequestVersionString), encryptedTimestamp)

	// Fourth HMAC: HMAC-SHA256(data, signingKey)
	return hmac256(data, signingKey)
}

func hmac256(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func NormalizeURL(rawURL string) string {
	// Remove trailing slashes
	rawURL = strings.TrimRight(rawURL, "/")
	return rawURL
}
