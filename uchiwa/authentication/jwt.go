package authentication

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

const jwtToken = "jwtToken"

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// GetJWTFromContext retrieves the JWT Token from the request
func GetJWTFromContext(r *http.Request) *jwt.Token {
	if value := context.Get(r, jwtToken); value != nil {
		return value.(*jwt.Token)
	}
	return nil
}

// GetToken returns a string that contain the token
func GetToken(role *Role, username string) (string, error) {
	if username == "" {
		return "", errors.New("Could not generate a token for the user. Invalid username")
	}

	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims["Role"] = role
	t.Claims["Username"] = username

	if privateKey == nil {
		return "", errors.New("Could not generate a token for the user. Invalid private key")
	}

	tokenString, err := t.SignedString(privateKey)
	return tokenString, err
}

// generateKeyPair generates an RSA keypair of 2048 bits using a random rand.Reader
func generateKeyPair() *rsa.PrivateKey {
	keypair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Fatalf("Could not generate an RSA keypair: %s", err)
	}

	return keypair
}

// generateToken generates a private and public RSA keys
// in order to be used for the JWT signature
func generateToken() (*rsa.PrivateKey, *rsa.PublicKey) {
	logger.Debug("Generating new temporary RSA keys")
	privateKey := generateKeyPair()
	// Precompute some calculations
	privateKey.Precompute()
	publicKey := &privateKey.PublicKey

	return privateKey, publicKey
}

// initToken initializes the token by weither loading the keys from the
// filesystem with the loadToken() function or by generating temporarily
// ones with the generateToken() function
func initToken(a structs.Auth) {
	var err error
	privateKey, publicKey, err = loadToken(a)
	if err != nil {
		// At this point we need to generate temporary RSA keys
		logger.Debug(err)
		privateKey, publicKey = generateToken()
	}
}

// loadToken loads a private and public RSA keys from the filesystem
// in order to be used for the JWT signature
func loadToken(a structs.Auth) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	logger.Debug("Attempting to load the RSA keys from the filesystem")

	if a.PrivateKey == "" || a.PublicKey == "" {
		return nil, nil, errors.New("The paths to the private and public RSA keys were not provided")
	}

	// Read the files from the filesystem
	prv, err := ioutil.ReadFile(a.PrivateKey)
	if err != nil {
		logger.Fatalf("Unable to open the private key file: %v", err)
	}
	pub, err := ioutil.ReadFile(a.PublicKey)
	if err != nil {
		logger.Fatalf("Unable to open the public key file: %v", err)
	}

	// Parse the RSA keys
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(prv)
	if err != nil {
		logger.Fatalf("Unable to parse the private key: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		logger.Fatalf("Unable to parse the public key: %v", err)
	}

	logger.Info("Provided RSA keys successfully loaded")
	return privateKey, publicKey, nil
}

// setJWTIntoContext injects the JWT Token into the request for later use
func setJWTInContext(r *http.Request, token *jwt.Token) {
	context.Set(r, jwtToken, token)
}

// verifyJWT extracts and verifies the validity of the JWT
func verifyJWT(r *http.Request) (*jwt.Token, error) {
	token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Debugf("Unexpected signing method: %v", t.Header["alg"])
			return nil, errors.New("")
		}
		return publicKey, nil
	})

	if token == nil || err != nil {
		logger.Debug(err)
		return nil, errors.New("")
	}

	if !token.Valid {
		logger.Debug("Invalid JWT")
		return nil, errors.New("")
	}

	return token, nil
}
