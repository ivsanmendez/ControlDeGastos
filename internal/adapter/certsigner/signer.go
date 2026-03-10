package certsigner

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/youmark/pkcs8"
)

// Signer loads a .cer (X.509 certificate) and .key (DER-encoded encrypted PKCS#8
// private key, as used by Mexican SAT certificates) pair.
// The private key is decrypted on each Sign call using the provided password.
type Signer struct {
	encryptedKey []byte // raw encrypted key bytes (DER or PEM)
	certDER      []byte // DER-encoded certificate
}

// New creates a Signer from the given certificate and key file paths.
// If both paths are empty, it returns a no-op signer (Available() == false).
// The private key is NOT decrypted at startup — a password is required per Sign call.
func New(certPath, keyPath string) (*Signer, error) {
	if certPath == "" && keyPath == "" {
		return &Signer{}, nil
	}
	if certPath == "" || keyPath == "" {
		return nil, fmt.Errorf("certsigner: both SIGN_CERT_PATH and SIGN_KEY_PATH must be set")
	}

	certDER, err := loadCertDER(certPath)
	if err != nil {
		return nil, fmt.Errorf("certsigner: load certificate: %w", err)
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("certsigner: read key file: %w", err)
	}

	return &Signer{encryptedKey: keyBytes, certDER: certDER}, nil
}

// Sign decrypts the private key with password, then computes SHA-256 hash of data
// and signs with RSA PKCS#1 v1.5.
func (s *Signer) Sign(data []byte, password string) ([]byte, error) {
	if len(s.encryptedKey) == 0 {
		return nil, fmt.Errorf("certsigner: signer not configured")
	}

	key, err := decryptPrivateKey(s.encryptedKey, password)
	if err != nil {
		return nil, fmt.Errorf("certsigner: decrypt key: %w", err)
	}

	hash := sha256.Sum256(data)
	return rsa.SignPKCS1v15(nil, key, crypto.SHA256, hash[:])
}

// Certificate returns the DER-encoded X.509 certificate.
func (s *Signer) Certificate() []byte {
	return s.certDER
}

// Available reports whether signing is configured (cert + key loaded).
func (s *Signer) Available() bool {
	return len(s.encryptedKey) > 0 && len(s.certDER) > 0
}

// loadCertDER reads a .cer file. It handles both PEM-wrapped and raw DER.
func loadCertDER(path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Try PEM decode first
	block, _ := pem.Decode(raw)
	if block != nil && block.Type == "CERTIFICATE" {
		return block.Bytes, nil
	}

	// Assume raw DER — validate by parsing
	if _, err := x509.ParseCertificate(raw); err != nil {
		return nil, fmt.Errorf("file is neither valid PEM nor DER certificate: %w", err)
	}
	return raw, nil
}

// decryptPrivateKey tries to decrypt an encrypted private key using the password.
// It handles: DER-encoded encrypted PKCS#8 (SAT .key), PEM-wrapped encrypted PKCS#8,
// and falls back to unencrypted PKCS#8/PKCS#1 (PEM or DER).
func decryptPrivateKey(raw []byte, password string) (*rsa.PrivateKey, error) {
	derBytes := raw

	// If PEM-wrapped, unwrap first
	if block, _ := pem.Decode(raw); block != nil {
		derBytes = block.Bytes
	}

	// Try encrypted PKCS#8 (SAT format: DER-encoded encrypted PKCS#8)
	parsed, err := pkcs8.ParsePKCS8PrivateKey(derBytes, []byte(password))
	if err == nil {
		rsaKey, ok := parsed.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not RSA")
		}
		return rsaKey, nil
	}

	// Fallback: try stdlib unencrypted PKCS#8
	stdParsed, stdErr := x509.ParsePKCS8PrivateKey(derBytes)
	if stdErr == nil {
		rsaKey, ok := stdParsed.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("PKCS#8 key is not RSA")
		}
		return rsaKey, nil
	}

	// Fallback: try unencrypted PKCS#1
	rsaKey, pkcs1Err := x509.ParsePKCS1PrivateKey(derBytes)
	if pkcs1Err == nil {
		return rsaKey, nil
	}

	return nil, fmt.Errorf("failed to decrypt/parse key (encrypted PKCS#8: %v)", err)
}
