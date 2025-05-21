package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
)

func main() {
	// Generate a 2048-bit RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Save private key in PEM format
	privateFile, _ := os.Create("private_key.pem")
	defer privateFile.Close()
	pem.Encode(privateFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Save public key in PEM format
	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	publicFile, _ := os.Create("public_key.pem")
	defer publicFile.Close()
	pem.Encode(publicFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	// Print public key parameters for JWKS
	pub := privateKey.PublicKey
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes())

	jwk := map[string]string{
		"kty": "RSA",
		"kid": "mykey",
		"use": "sig",
		"alg": "RS256",
		"n":   n,
		"e":   e,
	}

	jwkJSON, _ := json.MarshalIndent(map[string][]map[string]string{
		"keys": {jwk},
	}, "", "  ")
	fmt.Println("JWKS:\n", string(jwkJSON))
}
