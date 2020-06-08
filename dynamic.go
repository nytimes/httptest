// Copyright 2019 The New York Times Company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
  "crypto"
  "crypto/rsa"
  "crypto/sha256"
  "encoding/base64"
  "encoding/pem"
  "fmt"
  "strconv"
  "strings"
  "time"

  "github.com/youmark/pkcs8"
)

// Process the dynamic headers and add them to the map of all headers
func ProcessDynamicHeaders(dynamicHeaders map[string]Function, allHeaders map[string]string) error {
	for name, function := range dynamicHeaders {
		if _, present := allHeaders[name]; present {
			return fmt.Errorf("cannot process dynamic header %s; a header with that name is already defined", name)
		}
		fn := funcMap[function.Name]
		args := function.Args
    var err error
		allHeaders[name], err = fn(allHeaders, args)
		if err != nil {
			return err
		}
	}
	return nil
}

// Generic signature for any function that can resolve a dynamic header value
type resolveHeader func(existingHeaders map[string]string, args []string) (string, error)

// Map of strings to dynamic header functions
var funcMap = map[string]resolveHeader{
  "now": now,
  "signStringRS256": signStringRS256,
}

// Returns the number of seconds since the Unix epoch
func now(existingHeaders map[string]string, args []string) (string, error) {
	return strconv.FormatInt(time.Now().Unix(), 10), nil
}

// Constructs a string from args (delimited by newlines), signs it with the passphrase-encrypted PKCS #8 private key, and returns the signature in base64
func signStringRS256(existingHeaders map[string]string, args []string) (string, error) {
	// Get the key and passphrase from environment variables
	key := formatKey(args[0])
	passphrase := args[1]
	if key == "" || passphrase == "" {
		return "", fmt.Errorf("cannot call signStringRS256 function; key and passphrase must be defined in environment")
	}

	// Construct the string to sign
  var buffer strings.Builder
  buffer.WriteString(existingHeaders[args[2]])
  buffer.WriteRune('\n')
  buffer.WriteString(args[3])
  buffer.WriteRune('\n')
  buffer.WriteString(existingHeaders[args[4]])
  buffer.WriteRune('\n')
  buffer.WriteString(existingHeaders[args[5]])
  buffer.WriteRune('\n')
  stringToSign := buffer.String()

	// Get the key in PEM format
	pemBlock, _ := pem.Decode([]byte(key))

	// Decrypt the key with the passphrase
	decryptedKey, err := pkcs8.ParsePKCS8PrivateKey(pemBlock.Bytes, []byte(passphrase))
	if err != nil {
		return "", fmt.Errorf("unable to parse private key")
	}

	// Convert decrypted key to RSA key
	var rsaKey *rsa.PrivateKey
	var ok bool
	rsaKey, ok = decryptedKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("key is not an RSA key")
	}

	// Hash the string using SHA256
	hash := sha256.Sum256([]byte(stringToSign))

	// Sign the hashed header with the RSA key
	signature, err := rsa.SignPKCS1v15(nil, rsaKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("could not sign header")
	}

	// Return signature in base64 encoding
	return base64.StdEncoding.EncodeToString(signature), nil
}

// This function fixes the issue of newlines being converted to spaces in multiline environment variables upon unmarshalling
func formatKey(key string) string {
  prefix := "-----BEGIN ENCRYPTED PRIVATE KEY-----"
  postfix := "-----END ENCRYPTED PRIVATE KEY-----"
  dataStart := strings.Index(key, prefix) + len(prefix)
  dataEnd := strings.Index(key, postfix)

  data := key[dataStart:dataEnd]
  formattedData := strings.ReplaceAll(data, " ", "\n")

  return prefix + formattedData + postfix
}
