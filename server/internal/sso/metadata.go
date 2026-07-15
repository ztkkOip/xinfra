package sso

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"os"
	"strings"
)

const (
	samlMetadataNS = "urn:oasis:names:tc:SAML:2.0:metadata"
	samlProtocolNS = "urn:oasis:names:tc:SAML:2.0:protocol"
	xmlDsigNS      = "http://www.w3.org/2000/09/xmldsig#"
)

type MetadataConfig struct {
	EntityID string
	ACSURL   string
	CertFile string
}

func BuildSPMetadata(cfg MetadataConfig) ([]byte, error) {
	if cfg.EntityID == "" {
		return nil, errors.New("saml entity id is required")
	}
	if cfg.ACSURL == "" {
		return nil, errors.New("saml acs url is required")
	}
	cert, err := readCertificate(cfg.CertFile)
	if err != nil {
		return nil, err
	}

	doc := entityDescriptor{
		XMLName:  xml.Name{Local: "EntityDescriptor"},
		Xmlns:    samlMetadataNS,
		EntityID: cfg.EntityID,
		SPSSODescriptor: spSSODescriptor{
			ProtocolSupportEnumeration: samlProtocolNS,
			AuthnRequestsSigned:        true,
			WantAssertionsSigned:       true,
			KeyDescriptors: []keyDescriptor{
				newKeyDescriptor("signing", cert),
				newKeyDescriptor("encryption", cert),
			},
			NameIDFormats: []string{
				"urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
				"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
			},
			AssertionConsumerServices: []assertionConsumerService{
				{
					Binding:   "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
					Location:  cfg.ACSURL,
					Index:     0,
					IsDefault: true,
				},
			},
		},
	}

	body, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	out := append([]byte(xml.Header), body...)
	out = append(out, '\n')
	return out, nil
}

func readCertificate(path string) (string, error) {
	if path == "" {
		return "", errors.New("saml sp cert file is required")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return stripPEMCertificate(string(raw)), nil
	}
	if _, err := x509.ParseCertificate(block.Bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(block.Bytes), nil
}

func stripPEMCertificate(value string) string {
	value = strings.ReplaceAll(value, "-----BEGIN CERTIFICATE-----", "")
	value = strings.ReplaceAll(value, "-----END CERTIFICATE-----", "")
	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, " ", "")
	return value
}

func newKeyDescriptor(use string, cert string) keyDescriptor {
	return keyDescriptor{
		Use: use,
		KeyInfo: keyInfo{
			Xmlns: xmlDsigNS,
			X509Data: x509Data{
				X509Certificate: cert,
			},
		},
	}
}

type entityDescriptor struct {
	XMLName         xml.Name        `xml:"EntityDescriptor"`
	Xmlns           string          `xml:"xmlns,attr"`
	EntityID        string          `xml:"entityID,attr"`
	SPSSODescriptor spSSODescriptor `xml:"SPSSODescriptor"`
}

type spSSODescriptor struct {
	ProtocolSupportEnumeration string                     `xml:"protocolSupportEnumeration,attr"`
	AuthnRequestsSigned        bool                       `xml:"AuthnRequestsSigned,attr"`
	WantAssertionsSigned       bool                       `xml:"WantAssertionsSigned,attr"`
	KeyDescriptors             []keyDescriptor            `xml:"KeyDescriptor"`
	NameIDFormats              []string                   `xml:"NameIDFormat"`
	AssertionConsumerServices  []assertionConsumerService `xml:"AssertionConsumerService"`
}

type keyDescriptor struct {
	Use     string  `xml:"use,attr"`
	KeyInfo keyInfo `xml:"KeyInfo"`
}

type keyInfo struct {
	Xmlns    string   `xml:"xmlns,attr"`
	X509Data x509Data `xml:"X509Data"`
}

type x509Data struct {
	X509Certificate string `xml:"X509Certificate"`
}

type assertionConsumerService struct {
	Binding   string `xml:"Binding,attr"`
	Location  string `xml:"Location,attr"`
	Index     int    `xml:"index,attr"`
	IsDefault bool   `xml:"isDefault,attr"`
}
