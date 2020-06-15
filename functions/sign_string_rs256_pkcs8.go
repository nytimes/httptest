package functions

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/youmark/pkcs8"
)

var errInvalidKeyFormat = errors.New("key in invalid format")

// SignStringRS256PKCS8 constructs a string from given args (delimited by newlines),
// signing it with the (possibly passphrase-encrypted) PKCS #8 private key, and returns the signature in base64.
func SignStringRS256PKCS8(existingHeaders map[string]string, args []string) (string, error) {
	if !validateSignStringRS256PKCS8(args) {
		return "", fmt.Errorf("error calling signStringRS256PKCS8; at least 3 arguments are needed (key, passphrase, and a string to sign)")
	}

	key, passphrase, err := argsToKeyPassphrase(args)
	if err != nil {
		return "", fmt.Errorf("error calling SignStringRS256PKCS8; %w", err)
	}

	// Construct the string to sign
	stringToSign := argsToStringToSign(existingHeaders, args)

	// Get the key in PEM format
	pemBlock, _ := pem.Decode([]byte(key))

	// Parse the key, decrypting it if necessary
	decryptedKey, err := pkcs8.ParsePKCS8PrivateKey(pemBlock.Bytes, []byte(passphrase))
	if err != nil {
		return "", fmt.Errorf("error calling signStringRS256PKCS8; unable to parse private key: %w", err)
	}

	// Convert decrypted key to RSA key
	var rsaKey *rsa.PrivateKey
	var ok bool
	rsaKey, ok = decryptedKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("error calling signStringRS256PKCS8; key is not an RSA key")
	}

	// Hash the string using SHA256
	hash := sha256.Sum256([]byte(stringToSign))

	// Sign the hashed header with the RSA key
	signature, err := rsa.SignPKCS1v15(nil, rsaKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("error calling signStringRS256PKCS8; could not sign header: %w", err)
	}

	// Return signature in base64 encoding
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Builds the string we want to sign from the args as literals and from existing headers.
func argsToStringToSign(existingHeaders map[string]string, args []string) string {
	var buffer strings.Builder

	for _, arg := range args {
		if value, ok := existingHeaders[arg]; ok {
			buffer.WriteString(value)
		} else {
			buffer.WriteString(arg)
		}
		buffer.WriteRune('\n')
	}

	return buffer.String()
}

// Gets the key and passphrase from the args.
func argsToKeyPassphrase(args []string) (string, string, error) {
	if args[0] == "" {
		return "", "", fmt.Errorf("error calling signStringRS256PKCS8; key is empty")
	}

	passphrase := args[1]

	key, err := formatKey(args[0], passphrase != "")
	if err != nil {
		return "", "", err
	}

	return key, passphrase, nil
}

// Validates the required number of arguments are provided (key, passphrase, and a string to sign).
func validateSignStringRS256PKCS8(args []string) bool {
	return len(args) > 2
}

// Fixes the issue of newlines being converted to spaces in multiline environment variables upon unmarshalling.
func formatKey(key string, encrypted bool) (string, error) {
	prefix := "-----BEGIN PRIVATE KEY-----"
	postfix := "-----END PRIVATE KEY-----"
	if encrypted {
		prefix = "-----BEGIN ENCRYPTED PRIVATE KEY-----"
		postfix = "-----END ENCRYPTED PRIVATE KEY-----"
	}

	if !strings.HasPrefix(key, prefix) {
		return "", errInvalidKeyFormat
	}

	dataStart := strings.Index(key, prefix) + len(prefix)
	dataEnd := strings.Index(key, postfix)

	data := key[dataStart:dataEnd]
	formattedData := strings.ReplaceAll(data, " ", "\n")

	return prefix + formattedData + postfix, nil
}
