package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	_ "os"

	"golang.org/x/crypto/openpgp"
)

// verifyKeyPair verifies if the given public and private keys match.
func verifyKeyPair(publicKeyPath, privateKeyPath string) error {
	// Read the public key
	publicKeyData, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %v", err)
	}

	// Read the private key
	privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %v", err)
	}

	// Parse the public key
	publicKeyRing, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(publicKeyData))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}

	// Parse the private key
	privateKeyRing, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(privateKeyData))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// Create a test message
	message := []byte("This is a test message to verify the key pair.")

	// Sign the message with the private key
	var signedMessage bytes.Buffer
	err = openpgp.ArmoredDetachSign(&signedMessage, privateKeyRing[0], bytes.NewReader(message), nil)
	if err != nil {
		return fmt.Errorf("failed to sign the message with the private key: %v", err)
	}

	// Verify the signature with the public key
	_, err = openpgp.CheckArmoredDetachedSignature(publicKeyRing, bytes.NewReader(message), &signedMessage)
	if err != nil {
		return fmt.Errorf("failed to verify the signature with the public key: %v", err)
	}

	return nil
}

func main() {
	// Paths to the public and private key files
	publicKeyPath := "pgp_public.key"
	privateKeyPath := "pgp_private.key"

	// Verify the key pair
	err := verifyKeyPair(publicKeyPath, privateKeyPath)
	if err != nil {
		log.Fatalf("Key pair verification failed: %v", err)
	}

	fmt.Println("The public and private keys match!")
}
