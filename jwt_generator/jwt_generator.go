package main

import (
	"os"
	"time"
	"fmt"
	"encoding/pem"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Load RSA private key
	privateKeyData, err := os.ReadFile("/Users/abhigyanankur/Desktop/coding/pocketfm/key_generator/private_key.pem")
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic("invalid PEM format")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		panic(err)
	}

	// Create JWT with only required claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "my-issuer",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Optional header (must match Envoy config `kid`)
	token.Header["kid"] = "mykey"

	// Sign the JWT
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bearer " + signedToken)
}
