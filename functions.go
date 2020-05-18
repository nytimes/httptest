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

// This file contains the definitions for all the functions that can be used in dynamic headers
// To add a new function:
// 1) Define the regex used to identify it in the YAML
// 2) Define the function itself
// 3) Add logic to ProcessDynamicHeader to match the input string with the function, get the arguments, and call it

package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/youmark/pkcs8"
)

// These are the regexes used to identify which function needs to be called to get the header value
const nowPattern = `^now\(\)$`
const signStringPKCS8Pattern = `^signStringPKCS8\(([ -~]+)\)$`

// Function definitions (used to get dynamic header values)

// Returns the number of seconds since the Unix epoch
func now() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// Constructs a string from parts (delimited by newlines), signs it with the passphrase-encrypted PKCS #8 private key, and returns the signature in base64
func signStringPKCS8(keyEnvVar, passphraseEnvVar string, parts ...string) (string, error) {
	// Get the key and passphrase from environment variables
	key := os.Getenv(keyEnvVar)
	passphrase := os.Getenv(passphraseEnvVar)
	if key == "" || passphrase == "" {
		return "", fmt.Errorf("cannot call signStringPKCS8 function; key and passphrase must be defined in environment")
	}

	// Construct the string to sign
	stringToSign := buildString(parts...)

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

// Helper functions (not called to get dynamic header values)

func ProcessDynamicHeader(function string, existingHeaders map[string]string) (string, error) {
	// now
	if matches, _ := matches(nowPattern, function); matches {
		return now(), nil;
	}

	// signHeader
	if matches, args := matches(signStringPKCS8Pattern, function); matches {
		argSlice := parseArgs(args)

		key := parseArg(argSlice[0], existingHeaders)
		passphrase := parseArg(argSlice[1], existingHeaders)
		var parts []string
		for i := 2; i < len(argSlice); i++ {
			parts = append(parts, parseArg(argSlice[i], existingHeaders))
		}

		return signStringPKCS8(key, passphrase, parts...)
	}

	// default
	return "", fmt.Errorf("could not match %s; unknown function", function)
}

// Returns whether the subject matches the regex, as well as the arguments found in the match (if any)
func matches(regex, subject  string) (bool, string) {
	re := regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(subject)
	if len(matches) == 2 {
		return true, matches[1]
	} else if len(matches) == 1 {
		return true, ""
	} else {
		return false, ""
	}
}

func parseArgs(str string) []string {
	var args []string
	var inString bool
	var lastChar rune
	var buffer strings.Builder
	for _, currentChar := range str {
		if inString {
			buffer.WriteRune(currentChar)
			if currentChar == '"' && lastChar != '\\' {
				inString = false
			}
		} else {
			if currentChar != ',' {
				buffer.WriteRune(currentChar)
			} else {
				args = append(args, strings.TrimSpace(buffer.String()))
				buffer.Reset()
			}
		}
	}
	args = append(args, strings.TrimSpace(buffer.String()))
	return args
}

// Determine if the arg is a literal (string or number) or refers to a previously set header
func parseArg(arg string, existingHeaders map[string]string) string {
	if (arg[0] == '"' && arg[len(arg)-1] == '"') || (arg[0] == '\'' && arg[len(arg)-1] == '\'') {
		return removeQuotes(arg)
	} else {
		matches, err := regexp.MatchString(`^-?[0-9]*\.?[0-9]*$`, arg)
		if matches && err == nil {
			return arg
		} else {
			return existingHeaders[arg]
		}
	}
}

func removeQuotes(stringWithQuotes string) string {
	return stringWithQuotes[1:len(stringWithQuotes)-1]
}

func buildString(parts ...string) string {
	var buffer strings.Builder
	for _, part := range parts {
		buffer.WriteString(part)
		buffer.WriteRune('\n')
	}
	return buffer.String()
}
