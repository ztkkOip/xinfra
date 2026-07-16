package wayne

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	HeaderService   = "X-Wayne-Service"
	HeaderTimestamp = "X-Wayne-Timestamp"
	HeaderNonce     = "X-Wayne-Nonce"
	HeaderSignature = "X-Wayne-Signature"

	SignaturePrefix = "sha256="
)

var (
	ErrMissingSecret = errors.New("wayne service api secret key is empty")
	ErrMissingHeader = errors.New("wayne service header value is empty")
)

type SignedHeaders struct {
	Service   string
	Timestamp string
	Nonce     string
	Signature string
}

func BodySHA256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func SigningPayload(method, uri, timestamp, nonce string, body []byte) string {
	return strings.Join([]string{
		strings.ToUpper(strings.TrimSpace(method)),
		uri,
		timestamp,
		nonce,
		BodySHA256Hex(body),
	}, "\n")
}

func Sign(secret, method, uri, timestamp, nonce string, body []byte) (string, error) {
	secret = secretValue(secret)
	if secret == "" {
		return "", ErrMissingSecret
	}
	payload := SigningPayload(method, uri, timestamp, nonce, body)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return SignaturePrefix + hex.EncodeToString(mac.Sum(nil)), nil
}

func Verify(secret, signature, method, uri, timestamp, nonce string, body []byte) bool {
	expected, err := Sign(secret, method, uri, timestamp, nonce, body)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(normalizeSignature(signature)), []byte(expected)) == 1
}

func BuildSignedHeaders(service, secret, method, uri string, body []byte, now time.Time) (SignedHeaders, error) {
	service = strings.TrimSpace(service)
	if service == "" {
		return SignedHeaders{}, fmt.Errorf("%w: service", ErrMissingHeader)
	}
	timestamp := strconv.FormatInt(now.Unix(), 10)
	nonce, err := NewNonce()
	if err != nil {
		return SignedHeaders{}, err
	}
	signature, err := Sign(secret, method, uri, timestamp, nonce, body)
	if err != nil {
		return SignedHeaders{}, err
	}
	return SignedHeaders{
		Service:   service,
		Timestamp: timestamp,
		Nonce:     nonce,
		Signature: signature,
	}, nil
}

func (h SignedHeaders) Apply(req *http.Request) {
	req.Header.Set(HeaderService, h.Service)
	req.Header.Set(HeaderTimestamp, h.Timestamp)
	req.Header.Set(HeaderNonce, h.Nonce)
	req.Header.Set(HeaderSignature, h.Signature)
}

func NewNonce() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func normalizeSignature(signature string) string {
	signature = strings.TrimSpace(signature)
	if strings.HasPrefix(signature, SignaturePrefix) {
		return signature
	}
	return SignaturePrefix + signature
}

func secretValue(secret string) string {
	secret = strings.TrimSpace(secret)
	if _, value, ok := strings.Cut(secret, ":"); ok {
		return strings.TrimSpace(value)
	}
	return secret
}
