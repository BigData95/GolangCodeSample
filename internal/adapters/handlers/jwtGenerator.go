package handlers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/oauth2/google"
)

type credentials struct {
	SaEmail string `json:"client_email"`
}

func GenerateJWT() (string, error) {
	now := time.Now().Unix()
	var expiryLength int64 = 7200
	jsonFile, err := os.Open(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		fmt.Printf("Error opening SA FILE, err: ", err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var email credentials
	err = json.Unmarshal(byteValue, &email)
	if err != nil {
		fmt.Printf("Error unmarsh crendential from SA")
		return "", fmt.Errorf("Error : %w", err)
	}

	// Build the JWT payload.
	jwt := &ClaimSet{
		Iat: now,
		// expires after 'expiryLength' seconds.
		Exp: now + expiryLength,
		// Iss must match 'issuer' in the security configuration in your
		// swagger spec (e.g. service account email). It can be any string.
		Iss: email.SaEmail,
		// Aud must be either your Endpoints service name, or match the value
		// specified as the 'x-google-audience' in the OpenAPI document.
		Aud: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		// Sub and Email should match the service account's email address.
		Sub:           email.SaEmail,
		PrivateClaims: map[string]interface{}{"email": email.SaEmail},
	}
	jwsHeader := &Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}

	// Extract the RSA private key from the service account keyfile.
	sa, err := os.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		return "", fmt.Errorf("Could not read service account file: %w", err)
	}
	conf, err := google.JWTConfigFromJSON(sa)
	if err != nil {
		return "", fmt.Errorf("Could not parse service account JSON: %w", err)
	}
	block, _ := pem.Decode(conf.PrivateKey)
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("private key parse error: %w", err)
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	// Sign the JWT with the service account's private key.
	if !ok {
		return "", errors.New("private key failed rsa.PrivateKey type assertion")
	}
	response, err := Encode(jwsHeader, jwt, rsaKey)
	return response, err
}
