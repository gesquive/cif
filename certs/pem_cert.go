package certs

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
)

// DecodePEMBlock takes a single PEM file as provided by a user and decodes it into our intermediate certificate
// representation.
//
// The special function of this code is to handle the case that the user has accidentally provided us with a
// concatenated set of certificates. In this case, all certificates will be added to the trust store, with the
// label manipulated slightly to distinguish between them.
func DecodePEMBlock(data []byte) ([]*Certificate, error) {
	// Step one, decode the PEM file into all its constituent parts.
	blocks := make([]*pem.Block, 0)
	var p *pem.Block

	for data != nil && len(data) > 0 {
		p, data = pem.Decode(data)
		if p == nil {
			return nil, errors.New("Invalid PEM file.")
		}
		blocks = append(blocks, p)
	}

	// Now, for each block, parse the PEM certificate.
	parsedCerts := make([]*Certificate, 0, len(blocks))
	for i, block := range blocks {
		p, err := parsePEMCertificate(block)
		if err != nil {
			log.Printf("Failed to parse cert %v.\n", i)
			return nil, err
		}

		parsedCerts = append(parsedCerts, p)
	}

	return parsedCerts, nil
}

func parsePEMCertificate(p *pem.Block) (*Certificate, error) {
	// The decoded PEM file should be x509 data. We should therefore be able to pull that data out
	// using the x509 module.
	c, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		log.Printf("Invalid certificate: %v.\n", err)
		return nil, errors.New("Invalid certificate.")
	}

	// Transform this into our internal representation.
	parsed := &Certificate{
		nameToString(c.Issuer),
		nameToString(c.Subject),
		getCertName(c.Subject),
		c.SerialNumber.String(),
		fingerprintString(crypto.MD5, c.Raw),
		fingerprintString(crypto.SHA1, c.Raw),
		fingerprintString(crypto.SHA256, c.Raw),
		p,
	}
	return parsed, nil
}
