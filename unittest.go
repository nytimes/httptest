package main

import (
  "testing"
)


func TestFormatKey(t *testing.T) {
  key1 := `-----BEGIN PRIVATE KEY----- abcdef ghijkl mnopqr -----END PRIVATE KEY-----`
  key1ExpectedResult := `-----BEGIN PRIVATE KEY-----
abcdef
ghijkl
mnopqr
-----END PRIVATE KEY-----`

  key2 := `-----BEGIN ENCRYPTED PRIVATE KEY----- abcdef ghijkl mnopqr -----END ENCRYPTED PRIVATE KEY-----`
  key2ExpectedResult := `-----BEGIN ENCRYPTED PRIVATE KEY-----
abcdef
ghijkl
mnopqr
-----END ENCRYPTED PRIVATE KEY-----`

  result1, _ := FormatKey(key1, false)
  if result1 != key1ExpectedResult {
    t.Errorf("Expected:\n%s\nFor input:\n%s\nGot:\n%s\n\n", key1ExpectedResult, key1, result1)
  }

  result2, _ := FormatKey(key2, true)
  if result2 != key2ExpectedResult {
    t.Errorf("Expected:\n%s\nFor input:\n%s\nGot:\n%s\n\n", key2ExpectedResult, key2, result2)
  }

  _, err1 := FormatKey(key1, true)
  if err1 == nil {
    t.Errorf("Expected failure when passing decrypted key and asserting it's encrypted")
  }

  _, err2 := FormatKey(key2, false)
  if err2 == nil {
    t.Errorf("Expected failure when passing encrypted key and asserting it's decrypted")
  }
}
