package wayne

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestBodySHA256Hex(t *testing.T) {
	got := BodySHA256Hex([]byte("hello"))
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got != want {
		t.Fatalf("BodySHA256Hex() = %q, want %q", got, want)
	}
}

func TestSigningPayload(t *testing.T) {
	body := []byte(`{"groupIds":[10]}`)
	got := SigningPayload("put", "/api/v1/internal/namespaces/1/users/2001/roles", "1721000000", "nonce-001", body)
	want := strings.Join([]string{
		"PUT",
		"/api/v1/internal/namespaces/1/users/2001/roles",
		"1721000000",
		"nonce-001",
		"cf296504b2a434969ee151c1a48aa07fabb64634c77e569af74260cf234080f2",
	}, "\n")
	if got != want {
		t.Fatalf("SigningPayload() = %q, want %q", got, want)
	}
}

func TestSignAndVerify(t *testing.T) {
	body := []byte(`{"groupIds":[10]}`)
	signature, err := Sign("test-secret", "PUT", "/api/v1/internal/namespaces/1/users/2001/roles", "1721000000", "nonce-001", body)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	want := "sha256=6da70ea095fde90b534d0677da62db867d0b8092f6b15ed86cb52d77571b5b1d"
	if signature != want {
		t.Fatalf("Sign() = %q, want %q", signature, want)
	}
	if !Verify("test-secret", signature, "PUT", "/api/v1/internal/namespaces/1/users/2001/roles", "1721000000", "nonce-001", body) {
		t.Fatal("Verify() rejected valid signature")
	}
	if !Verify("test-secret", strings.TrimPrefix(signature, SignaturePrefix), "PUT", "/api/v1/internal/namespaces/1/users/2001/roles", "1721000000", "nonce-001", body) {
		t.Fatal("Verify() rejected valid signature without prefix")
	}
	if Verify("test-secret", signature, "PUT", "/api/v1/internal/namespaces/1/users/2001/roles", "1721000000", "nonce-001", []byte(`{"groupIds":[11]}`)) {
		t.Fatal("Verify() accepted tampered body")
	}
}

func TestSignAcceptsServicePrefixedSecret(t *testing.T) {
	body := []byte(`{"groupIds":[10]}`)
	plainSecretSignature, err := Sign("test-secret", "PUT", "/api/v1/internal/namespaces/1/users/2001/roles", "1721000000", "nonce-001", body)
	if err != nil {
		t.Fatalf("Sign() with plain secret error = %v", err)
	}
	prefixedSecretSignature, err := Sign("xinfra:test-secret", "PUT", "/api/v1/internal/namespaces/1/users/2001/roles", "1721000000", "nonce-001", body)
	if err != nil {
		t.Fatalf("Sign() with prefixed secret error = %v", err)
	}
	if prefixedSecretSignature != plainSecretSignature {
		t.Fatalf("prefixed secret signature = %q, want %q", prefixedSecretSignature, plainSecretSignature)
	}
}

func TestBuildSignedHeadersAndApply(t *testing.T) {
	body := []byte(`{"groupIds":[10]}`)
	now := time.Unix(1721000000, 0)
	headers, err := BuildSignedHeaders("xinfra", "test-secret", "PUT", "/api/v1/internal/namespaces/1/users/2001/roles", body, now)
	if err != nil {
		t.Fatalf("BuildSignedHeaders() error = %v", err)
	}
	if headers.Service != "xinfra" {
		t.Fatalf("headers.Service = %q, want xinfra", headers.Service)
	}
	if headers.Timestamp != "1721000000" {
		t.Fatalf("headers.Timestamp = %q, want 1721000000", headers.Timestamp)
	}
	if headers.Nonce == "" {
		t.Fatal("headers.Nonce is empty")
	}
	if !Verify("test-secret", headers.Signature, "PUT", "/api/v1/internal/namespaces/1/users/2001/roles", headers.Timestamp, headers.Nonce, body) {
		t.Fatal("generated headers signature is invalid")
	}

	req, err := http.NewRequest(http.MethodPut, "http://wayne.example.com/api/v1/internal/namespaces/1/users/2001/roles", strings.NewReader(string(body)))
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	headers.Apply(req)
	if req.Header.Get(HeaderService) != headers.Service ||
		req.Header.Get(HeaderTimestamp) != headers.Timestamp ||
		req.Header.Get(HeaderNonce) != headers.Nonce ||
		req.Header.Get(HeaderSignature) != headers.Signature {
		t.Fatalf("Apply() did not write expected signed headers: %#v", req.Header)
	}
}
