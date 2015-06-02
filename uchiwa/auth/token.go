package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/dgrijalva/jwt-go"
	"github.com/palourde/logger"
)

var (
	keyPEM    []byte
	pubKeyPEM []byte
)

// GetToken returns a string that contain the token
func GetToken(role *Role) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims["Role"] = role
	tokenString, err := t.SignedString(keyPEM)
	return tokenString, err
}

func initToken() {
	keyPair, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		logger.Fatalf("Could not generate the private key: %s", err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
	})
	pubKeyANS1, err := x509.MarshalPKIXPublicKey(&keyPair.PublicKey)
	if err != nil {
		logger.Fatalf("Could not generate the public key: %s", err)
	}
	pubKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyANS1,
	})
}
