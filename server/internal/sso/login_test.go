package sso

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestReadIDPSSOURLFetchesMetadataURL(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		_, _ = w.Write([]byte(testIDPMetadataXML(
			"https://idp.example.com/saml2/meta",
			"https://idp.example.com/sso/redirect",
			"https://idp.example.com/sso/post",
		)))
	}))
	defer server.Close()

	got, err := ReadIDPSSOURL(server.URL)
	if err != nil {
		t.Fatalf("ReadIDPSSOURL returned error: %v", err)
	}
	if got != "https://idp.example.com/sso/redirect" {
		t.Fatalf("unexpected sso url: %s", got)
	}
}

func TestReadIDPSSOURLReturnsHTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusBadGateway)
	}))
	defer server.Close()

	if _, err := ReadIDPSSOURL(server.URL); err == nil {
		t.Fatal("expected error for non-200 metadata response")
	}
}

func TestBuildLoginRedirectUsesRemoteMetadata(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testIDPMetadataXML(
			"https://idp.example.com/saml2/meta",
			"https://idp.example.com/sso/redirect",
			"https://idp.example.com/sso/post",
		)))
	}))
	defer server.Close()

	redirectURL, err := BuildLoginRedirect(LoginConfig{
		EntityID:       "https://sp.example.com/auth/api/v1/saml/metadata",
		ACSURL:         "https://sp.example.com/auth/api/v1/saml/acs",
		IDPMetadataURL: server.URL,
		RelayState:     "/auth/",
	})
	if err != nil {
		t.Fatalf("BuildLoginRedirect returned error: %v", err)
	}

	parsed, err := url.Parse(redirectURL)
	if err != nil {
		t.Fatalf("parse redirect url: %v", err)
	}
	if parsed.String() == "" {
		t.Fatal("redirect url is empty")
	}
	if parsed.Scheme != "https" || parsed.Host != "idp.example.com" || parsed.Path != "/sso/redirect" {
		t.Fatalf("unexpected redirect target: %s", redirectURL)
	}
	if parsed.Query().Get("SAMLRequest") == "" {
		t.Fatal("missing SAMLRequest query")
	}
	if parsed.Query().Get("RelayState") != "/auth/" {
		t.Fatalf("unexpected relay state: %s", parsed.Query().Get("RelayState"))
	}
}

func TestDecodeSAMLResponseDecryptsEncryptedAssertion(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	keyFile := writeTestPrivateKey(t, key)

	assertion := `<saml:Assertion xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
  <saml:Subject><saml:NameID>alice@example.com</saml:NameID></saml:Subject>
  <saml:AttributeStatement>
    <saml:Attribute Name="email"><saml:AttributeValue>alice@example.com</saml:AttributeValue></saml:Attribute>
  </saml:AttributeStatement>
</saml:Assertion>`
	sessionKey := []byte("0123456789abcdef")
	encryptedAssertion := encryptTestAssertion(t, []byte(assertion), sessionKey)
	encryptedKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, &key.PublicKey, sessionKey, nil)
	if err != nil {
		t.Fatalf("encrypt key: %v", err)
	}
	response := `<samlp:Response xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol">
  <saml:EncryptedAssertion xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
    <xenc:EncryptedData xmlns:xenc="http://www.w3.org/2001/04/xmlenc#">
      <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
        <xenc:EncryptedKey>
          <xenc:CipherData><xenc:CipherValue>` + base64.StdEncoding.EncodeToString(encryptedKey) + `</xenc:CipherValue></xenc:CipherData>
        </xenc:EncryptedKey>
      </ds:KeyInfo>
      <xenc:CipherData><xenc:CipherValue>` + base64.StdEncoding.EncodeToString(encryptedAssertion) + `</xenc:CipherValue></xenc:CipherData>
    </xenc:EncryptedData>
  </saml:EncryptedAssertion>
</samlp:Response>`

	info, err := DecodeSAMLResponse(base64.StdEncoding.EncodeToString([]byte(response)), keyFile)
	if err != nil {
		t.Fatalf("DecodeSAMLResponse returned error: %v", err)
	}
	if info.NameID != "alice@example.com" {
		t.Fatalf("unexpected name id: %q", info.NameID)
	}
	if got := info.Attributes["email"]; len(got) != 1 || got[0] != "alice@example.com" {
		t.Fatalf("unexpected email attribute: %#v", got)
	}
	if info.DecryptedAssertionXML == "" {
		t.Fatal("missing decrypted assertion xml")
	}
}

func TestDecodeSAMLResponseDecryptsGCMEncryptedAssertion(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	keyFile := writeTestPrivateKey(t, key)

	assertion := `<saml:Assertion xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
  <saml:Subject><saml:NameID>bob@example.com</saml:NameID></saml:Subject>
  <saml:AttributeStatement>
    <saml:Attribute Name="email"><saml:AttributeValue>bob@example.com</saml:AttributeValue></saml:Attribute>
  </saml:AttributeStatement>
</saml:Assertion>`
	sessionKey := []byte("0123456789abcdef0123456789abcdef")
	encryptedAssertion := encryptTestAssertionGCM(t, []byte(assertion), sessionKey)
	encryptedKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, &key.PublicKey, sessionKey, nil)
	if err != nil {
		t.Fatalf("encrypt key: %v", err)
	}
	response := `<samlp:Response xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol">
  <saml:EncryptedAssertion xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
    <xenc:EncryptedData xmlns:xenc="http://www.w3.org/2001/04/xmlenc#">
      <xenc:EncryptionMethod Algorithm="` + xmlEncAES256GCM + `"></xenc:EncryptionMethod>
      <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
        <xenc:EncryptedKey>
          <xenc:CipherData><xenc:CipherValue>` + base64.StdEncoding.EncodeToString(encryptedKey) + `</xenc:CipherValue></xenc:CipherData>
        </xenc:EncryptedKey>
      </ds:KeyInfo>
      <xenc:CipherData><xenc:CipherValue>` + base64.StdEncoding.EncodeToString(encryptedAssertion) + `</xenc:CipherValue></xenc:CipherData>
    </xenc:EncryptedData>
  </saml:EncryptedAssertion>
</samlp:Response>`

	info, err := DecodeSAMLResponse(base64.StdEncoding.EncodeToString([]byte(response)), keyFile)
	if err != nil {
		t.Fatalf("DecodeSAMLResponse returned error: %v", err)
	}
	if info.NameID != "bob@example.com" {
		t.Fatalf("unexpected name id: %q", info.NameID)
	}
	if got := info.Attributes["email"]; len(got) != 1 || got[0] != "bob@example.com" {
		t.Fatalf("unexpected email attribute: %#v", got)
	}
}

func writeTestPrivateKey(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "sp-*.key")
	if err != nil {
		t.Fatalf("create key file: %v", err)
	}
	err = pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	closeErr := file.Close()
	if err != nil {
		t.Fatalf("write key file: %v", err)
	}
	if closeErr != nil {
		t.Fatalf("close key file: %v", closeErr)
	}
	return file.Name()
}

func encryptTestAssertion(t *testing.T, plain, key []byte) []byte {
	t.Helper()
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("init cipher: %v", err)
	}
	plain = pkcs7Pad(plain, block.BlockSize())
	iv := bytes.Repeat([]byte{1}, block.BlockSize())
	out := make([]byte, len(iv)+len(plain))
	copy(out, iv)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out[len(iv):], plain)
	return out
}

func encryptTestAssertionGCM(t *testing.T, plain, key []byte) []byte {
	t.Helper()
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("init cipher: %v", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("init gcm: %v", err)
	}
	nonce := bytes.Repeat([]byte{2}, aead.NonceSize())
	cipherText := aead.Seal(nil, nonce, plain, nil)
	out := make([]byte, len(nonce)+len(cipherText))
	copy(out, nonce)
	copy(out[len(nonce):], cipherText)
	return out
}

func pkcs7Pad(value []byte, blockSize int) []byte {
	padding := blockSize - len(value)%blockSize
	return append(value, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func testIDPMetadataXML(entityID, redirectURL, postURL string) string {
	return `<?xml version="1.0"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="` + entityID + `">
  <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="` + redirectURL + `"></SingleSignOnService>
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="` + postURL + `"></SingleSignOnService>
  </IDPSSODescriptor>
</EntityDescriptor>`
}
