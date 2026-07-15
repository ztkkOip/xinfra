package sso

import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	samlHTTPRedirectBinding = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
	samlHTTPPostBinding     = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
)

type LoginConfig struct {
	EntityID       string
	ACSURL         string
	IDPMetadataURL string
	RelayState     string
}

func BuildLoginRedirect(cfg LoginConfig) (string, error) {
	if cfg.EntityID == "" {
		return "", errors.New("saml entity id is required")
	}
	if cfg.ACSURL == "" {
		return "", errors.New("saml acs url is required")
	}
	ssoURL, err := ReadIDPSSOURL(cfg.IDPMetadataURL)
	if err != nil {
		return "", err
	}
	request, err := buildAuthnRequest(cfg.EntityID, cfg.ACSURL, ssoURL)
	if err != nil {
		return "", err
	}
	encoded, err := deflateAndBase64(request)
	if err != nil {
		return "", err
	}
	redirectURL, err := url.Parse(ssoURL)
	if err != nil {
		return "", err
	}
	query := redirectURL.Query()
	query.Set("SAMLRequest", encoded)
	if cfg.RelayState != "" {
		query.Set("RelayState", cfg.RelayState)
	}
	redirectURL.RawQuery = query.Encode()
	return redirectURL.String(), nil
}

func ReadIDPSSOURL(metadataURL string) (string, error) {
	if metadataURL == "" {
		return "", errors.New("saml idp metadata url is required")
	}
	request, err := http.NewRequest(http.MethodGet, metadataURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Accept", "application/samlmetadata+xml, application/xml, text/xml, */*")
	request.Header.Set("User-Agent", "authserver-saml-metadata/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("saml idp metadata request failed: url=%s status=%s body=%q", request.URL.String(), response.Status, truncateForLog(string(raw), 300))
	}
	return readIDPSSOURL(raw)
}

func truncateForLog(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit] + "...(truncated)"
}

func readIDPSSOURL(raw []byte) (string, error) {
	var doc idpEntityDescriptor
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return "", err
	}
	for _, service := range doc.IDPSSODescriptor.SingleSignOnServices {
		if service.Binding == samlHTTPRedirectBinding && service.Location != "" {
			return service.Location, nil
		}
	}
	for _, service := range doc.IDPSSODescriptor.SingleSignOnServices {
		if service.Location != "" {
			return service.Location, nil
		}
	}
	return "", errors.New("saml idp metadata has no single sign on service")
}

func DecodeSAMLResponse(value, privateKeyFile string) (*SAMLDebugInfo, error) {
	if value == "" {
		return nil, errors.New("SAMLResponse is required")
	}
	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	info := &SAMLDebugInfo{
		RawXML: string(raw),
	}
	_ = xml.Unmarshal(raw, &info.Response)
	if info.Response.Assertion.Subject.NameID.Value == "" && info.Response.EncryptedAssertion.EncryptedData.CipherData.CipherValue != "" {
		assertionXML, assertion, err := decryptEncryptedAssertion(info.Response.EncryptedAssertion, privateKeyFile)
		if err != nil {
			return nil, err
		}
		info.DecryptedAssertionXML = assertionXML
		info.Response.EncryptedAssertion.Assertion = assertion
	}
	info.NameID = firstNameID(info.Response.Assertion.Subject.NameID.Value, info.Response.EncryptedAssertion.Assertion.Subject.NameID.Value)
	info.Attributes = make(map[string][]string)
	for _, attr := range info.Response.Assertion.AttributeStatement.Attributes {
		addSAMLAttributeValues(info.Attributes, attr)
	}
	for _, attr := range info.Response.EncryptedAssertion.Assertion.AttributeStatement.Attributes {
		addSAMLAttributeValues(info.Attributes, attr)
	}
	return info, nil
}

func addSAMLAttributeValues(attrs map[string][]string, attr attribute) {
	if attr.Name != "" {
		attrs[attr.Name] = append(attrs[attr.Name], attr.Values...)
	}
	if attr.FriendlyName != "" && attr.FriendlyName != attr.Name {
		attrs[attr.FriendlyName] = append(attrs[attr.FriendlyName], attr.Values...)
	}
}

func decryptEncryptedAssertion(encrypted encryptedAssertion, privateKeyFile string) (string, assertion, error) {
	key, err := readRSAPrivateKey(privateKeyFile)
	if err != nil {
		return "", assertion{}, err
	}
	encryptedKey, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encrypted.EncryptedData.EncryptedKey.CipherData.CipherValue))
	if err != nil {
		return "", assertion{}, fmt.Errorf("decode saml encrypted key: %w", err)
	}
	sessionKey, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, key, encryptedKey, nil)
	if err != nil {
		return "", assertion{}, fmt.Errorf("decrypt saml encrypted key: %w", err)
	}
	encryptedValue, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encrypted.EncryptedData.CipherData.CipherValue))
	if err != nil {
		return "", assertion{}, fmt.Errorf("decode saml encrypted assertion: %w", err)
	}
	plain, err := decryptAESCBC(encryptedValue, sessionKey)
	if err != nil {
		return "", assertion{}, err
	}
	assertionXML := string(plain)
	var out assertion
	if err := xml.Unmarshal(plain, &out); err != nil {
		return "", assertion{}, fmt.Errorf("parse decrypted saml assertion: %w", err)
	}
	return assertionXML, out, nil
}

func readRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("saml sp key file is required for encrypted assertion")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("saml sp key file has no PEM block")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("saml sp key is not an RSA private key")
	}
	return key, nil
}

func decryptAESCBC(value, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("init saml assertion cipher: %w", err)
	}
	if len(value) < block.BlockSize()*2 || len(value)%block.BlockSize() != 0 {
		return nil, errors.New("invalid saml encrypted assertion length")
	}
	iv := value[:block.BlockSize()]
	cipherText := value[block.BlockSize():]
	plain := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plain, cipherText)
	plain, err = pkcs7Unpad(plain, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return plain, nil
}

func pkcs7Unpad(value []byte, blockSize int) ([]byte, error) {
	if len(value) == 0 || len(value)%blockSize != 0 {
		return nil, errors.New("invalid saml assertion padding length")
	}
	padding := int(value[len(value)-1])
	if padding == 0 || padding > blockSize || padding > len(value) {
		return nil, errors.New("invalid saml assertion padding")
	}
	for _, b := range value[len(value)-padding:] {
		if int(b) != padding {
			return nil, errors.New("invalid saml assertion padding bytes")
		}
	}
	return value[:len(value)-padding], nil
}

func (i *SAMLDebugInfo) JSON() string {
	data, _ := json.MarshalIndent(map[string]any{
		"id":             i.Response.ID,
		"in_response_to": i.Response.InResponseTo,
		"issuer":         i.Response.Issuer,
		"destination":    i.Response.Destination,
		"name_id":        i.NameID,
		"attributes":     i.Attributes,
	}, "", "  ")
	return string(data)
}

func buildAuthnRequest(entityID, acsURL, destination string) ([]byte, error) {
	req := authnRequest{
		XMLName:                     xml.Name{Local: "samlp:AuthnRequest"},
		XmlnsSAMLp:                  samlProtocolNS,
		XmlnsSAML:                   "urn:oasis:names:tc:SAML:2.0:assertion",
		ID:                          "_" + randomHex(20),
		Version:                     "2.0",
		IssueInstant:                time.Now().UTC().Format(time.RFC3339),
		Destination:                 destination,
		AssertionConsumerServiceURL: acsURL,
		ProtocolBinding:             samlHTTPPostBinding,
		Issuer: issuer{
			Value: entityID,
		},
		NameIDPolicy: nameIDPolicy{
			Format:      "urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
			AllowCreate: true,
		},
	}
	return xml.Marshal(req)
}

func deflateAndBase64(raw []byte) (string, error) {
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return "", err
	}
	if _, err := writer.Write(raw); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func randomHex(size int) string {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func firstNameID(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

type idpEntityDescriptor struct {
	IDPSSODescriptor struct {
		SingleSignOnServices []struct {
			Binding  string `xml:"Binding,attr"`
			Location string `xml:"Location,attr"`
		} `xml:"SingleSignOnService"`
	} `xml:"IDPSSODescriptor"`
}

type authnRequest struct {
	XMLName                     xml.Name     `xml:"samlp:AuthnRequest"`
	XmlnsSAMLp                  string       `xml:"xmlns:samlp,attr"`
	XmlnsSAML                   string       `xml:"xmlns:saml,attr"`
	ID                          string       `xml:"ID,attr"`
	Version                     string       `xml:"Version,attr"`
	IssueInstant                string       `xml:"IssueInstant,attr"`
	Destination                 string       `xml:"Destination,attr"`
	AssertionConsumerServiceURL string       `xml:"AssertionConsumerServiceURL,attr"`
	ProtocolBinding             string       `xml:"ProtocolBinding,attr"`
	Issuer                      issuer       `xml:"saml:Issuer"`
	NameIDPolicy                nameIDPolicy `xml:"samlp:NameIDPolicy"`
}

type issuer struct {
	Value string `xml:",chardata"`
}

type nameIDPolicy struct {
	Format      string `xml:"Format,attr"`
	AllowCreate bool   `xml:"AllowCreate,attr"`
}

type SAMLDebugInfo struct {
	RawXML                string
	DecryptedAssertionXML string
	NameID                string
	Attributes            map[string][]string
	Response              samlResponse
}

type samlResponse struct {
	ID                 string             `xml:"ID,attr"`
	InResponseTo       string             `xml:"InResponseTo,attr"`
	Destination        string             `xml:"Destination,attr"`
	Issuer             string             `xml:"Issuer"`
	Assertion          assertion          `xml:"Assertion"`
	EncryptedAssertion encryptedAssertion `xml:"EncryptedAssertion"`
}

type encryptedAssertion struct {
	Assertion     assertion     `xml:"Assertion"`
	EncryptedData encryptedData `xml:"EncryptedData"`
}

type encryptedData struct {
	EncryptionMethod struct {
		Algorithm string `xml:"Algorithm,attr"`
	} `xml:"EncryptionMethod"`
	EncryptedKey encryptedKey `xml:"KeyInfo>EncryptedKey"`
	CipherData   cipherData   `xml:"CipherData"`
}

type encryptedKey struct {
	EncryptionMethod struct {
		Algorithm string `xml:"Algorithm,attr"`
	} `xml:"EncryptionMethod"`
	CipherData cipherData `xml:"CipherData"`
}

type cipherData struct {
	CipherValue string `xml:"CipherValue"`
}

type assertion struct {
	Subject            subject            `xml:"Subject"`
	AttributeStatement attributeStatement `xml:"AttributeStatement"`
}

type subject struct {
	NameID nameID `xml:"NameID"`
}

type nameID struct {
	Value string `xml:",chardata"`
}

type attributeStatement struct {
	Attributes []attribute `xml:"Attribute"`
}

type attribute struct {
	Name         string   `xml:"Name,attr"`
	FriendlyName string   `xml:"FriendlyName,attr"`
	Values       []string `xml:"AttributeValue"`
}
