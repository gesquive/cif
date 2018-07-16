package certs

import (
	"bytes"
	"crypto"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"strings"
)

// Certificate is an in-memory representation of a certificate.
type Certificate struct {
	Issuer            string
	Subject           string
	Label             string
	Serial            string
	MD5Fingerprint    string
	SHA1Fingerprint   string
	SHA256Fingerprint string
	PEMBlock          *pem.Block
}

// nameToString converts name into a string representation containing the
// CommonName, Organization and OrganizationalUnit.
func nameToString(name pkix.Name) string {
	ret := ""
	if len(name.CommonName) > 0 {
		ret += "CN=" + name.CommonName
	}

	if org := strings.Join(name.Organization, "/"); len(org) > 0 {
		if len(ret) > 0 {
			ret += " "
		}
		ret += "O=" + org
	}

	if orgUnit := strings.Join(name.OrganizationalUnit, "/"); len(orgUnit) > 0 {
		if len(ret) > 0 {
			ret += " "
		}
		ret += "OU=" + orgUnit
	}

	return ret
}

func fingerprintString(hashFunc crypto.Hash, data []byte) string {
	hash := hashFunc.New()
	hash.Write(data)
	digest := hash.Sum(nil)

	hex := fmt.Sprintf("%x", digest)
	ret := ""
	for len(hex) > 0 {
		if len(ret) > 0 {
			ret += ":"
		}
		todo := 2
		if len(hex) < todo {
			todo = len(hex)
		}
		ret += hex[:todo]
		hex = hex[todo:]
	}

	return ret
}

func getCertName(subject pkix.Name) string {
	ret := ""
	if len(subject.CommonName) > 0 {
		ret = subject.CommonName
	}
	return ret
}

// WriteCert will format a certificate object
func WriteCert(cert *Certificate) string {
	var certInfo bytes.Buffer

	certInfo.WriteString("\n")

	// format we want:
	// # Issuer: CN=DigiCert Global Root CA O=DigiCert Inc OU=www.digicert.com
	// # Subject: CN=DigiCert Global Root CA O=DigiCert Inc OU=www.digicert.com
	// # Label: "DigiCert Global Root CA"
	// # Serial: 10944719598952040374951832963794454346
	// # MD5 Fingerprint: 79:e4:a9:84:0d:7d:3a:96:d7:c0:4f:e2:43:4c:89:2e
	// # SHA1 Fingerprint: a8:98:5d:3a:65:e5:e5:c4:b2:d7:d6:6d:40:c6:dd:2f:b1:9c:54:36
	// # SHA256 Fingerprint: 43:48:a0:e9:44:4c:78:cb:26:5e:05:8d:5e:89:44:b4:d8:4f:96:62:bd:26:db:25:7f:89:34:a4:43:c7:01:61
	// -----BEGIN CERTIFICATE-----
	// ...
	certInfo.WriteString("# Issuer: " + cert.Issuer + "\n")
	certInfo.WriteString("# Subject: " + cert.Subject + "\n")
	certInfo.WriteString("# Label: " + cert.Label + "\n")
	certInfo.WriteString("# Serial: " + cert.Serial + "\n")
	certInfo.WriteString("# MD5 Fingerprint: " + cert.MD5Fingerprint + "\n")
	certInfo.WriteString("# SHA1 Fingerprint: " + cert.SHA1Fingerprint + "\n")
	certInfo.WriteString("# SHA256 Fingerprint: " + cert.SHA256Fingerprint + "\n")
	certInfo.Write(pem.EncodeToMemory(cert.PEMBlock))

	return certInfo.String()
}
