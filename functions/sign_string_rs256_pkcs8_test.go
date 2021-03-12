package functions

import (
	"testing"
)

func TestFormatKey(t *testing.T) {
	unencrypted := `-----BEGIN PRIVATE KEY----- abcdef ghijkl mnopqr -----END PRIVATE KEY-----`
	unencryptedExpected := `-----BEGIN PRIVATE KEY-----
abcdef
ghijkl
mnopqr
-----END PRIVATE KEY-----`

	encrypted := `-----BEGIN ENCRYPTED PRIVATE KEY----- abcdef ghijkl mnopqr -----END ENCRYPTED PRIVATE KEY-----`
	encryptedExpected := `-----BEGIN ENCRYPTED PRIVATE KEY-----
abcdef
ghijkl
mnopqr
-----END ENCRYPTED PRIVATE KEY-----`

	var FormatKeyTests = []struct {
		input     string
		encrypted bool
		expected  string
		err       error
	}{
		{unencrypted, false, unencryptedExpected, nil},
		{encrypted, true, encryptedExpected, nil},
		{unencrypted, true, "", errInvalidKeyFormat},
		{encrypted, false, "", errInvalidKeyFormat},
	}

	for _, tc := range FormatKeyTests {
		actual, err := formatKey(tc.input, tc.encrypted)
		if actual != tc.expected {
			t.Errorf("formatKey(%v, %v): expected %v, actual %v", tc.input, tc.encrypted, tc.expected, actual)
		}
		if err != tc.err {
			t.Errorf("formatKey(%v, %v): expected %v, got: %v", tc.input, tc.encrypted, tc.err, err)
		}
	}
}

func TestSignStringRS256PKCS8(t *testing.T) {
	testKey := `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDrmBH72WcZZiFr
llm9tZRaCXxKnZKGpXeLHTTyW+REAZ6vSO1siOxuQrv2sNqR3ZWEZ9HpKO1X6GNj
Ue+T+cBrnu5p3nQTTp0MnjKyu5qcJRf9zfZ3VeJkFvBHdfMH2t4RXoi3bYMSS37y
NaOu1uB7CBx7WEnWYHfIMoP36ksorkT4tCn/1meDmcXNl8UQrsqj7I4dkmGCt0xY
0l7wb4cWMjB6pq65+33Jc6eatladlFhXpTwyLBwhhHcAs8Uw5OA1TGTAHBHixr7A
qPNOwufDz2+Z2xSuLNgH1ySF2KsnJRS20NTCOO6G6JNFqpOiVQjpAukF7GTG4Z2S
jvAEyhIDAgMBAAECggEBAL2mCLQH6eqUQErvGQaR6P4hrKAUACPLh1PBCyIdvr7P
3wGTXyyDfG+14MFQ1GGfUgDn4h4jCAw/0eHdz1H7Nl5r7dfjbuUr31iM8JrYUjln
0sxIxCKETF3t6TZdSGoGUcUBqGSgD2bmxyYK79yKtOHVQbg49hdQSJwrrfgf7qir
NPtAESSeZ9oWECItub5D8wtkrTVHPY4l2zDfeED7jToIPmk1pvReox3VpWD3GbNC
ne/V/d0q0bTOSFH5T5iOM5tO4vdTjPKaZPqkjp+Z0VwKdrSGSB2aPFzxQ/ykCsf/
X0W2EdSs6pTI6pO11/2z/Jp0BvuAie6dX7TH6XdL2kECgYEA+ndVA6D38dSPR6bg
52iCI2NvGkmiT+jIuASDWESX3TOhfaJFWkdqjpiHr2hoVN6pc6II3hoaTe2n31vJ
j6+lDRwH2d6R6coK+2XUtESXwRx8Cm8FxWSV3a7AvXYa+kDGtH27a6tRAYYPAvMF
Gm+QwcIOMYBgA9FNdYlmbwss2gkCgYEA8MyennofezaXiQoEY+W3Nz1E8lpZwzJ1
oOTS0LV9fi9nGcAn8+lFS62R3yAYf+ZqZINbvjeMqFXUVmzK1asH01z1R6wmT0Zy
jNieW41q6JcocloIqnAbgCpSFYyWMbmAXdSNWQscVjKonVMXPNDN0Qef/BmSfGjH
+Fqky5dkfqsCgYEA05xCnWBAW3bk3vqlBZ4MZW27DpCrq6vW+XIGrmq1i9P1Wrng
sleoNXW1HYOushW1QNbjexK+qpxhuppH/ze80QifsXkT+lwTTzdHsE5LkIJKYl5O
l+lVnQfqG6hPPqO/vfqEgIErXYgv4qQD6cPcn3cemsAFXvRU5zsA6kycxlECgYBn
tZkzZCGe7ZpCWWAerldEoUzKnINAgMEMtMDfRutvp3beLlaGxJclyvGiia5Dl7eG
5tRijoY0EhNLzbtmXy0VqVmyrsApMIwxgTJi9/tthXzUE1bcIUCW6KNFyLD/ZYeV
4e+mxBRGQ7c/WwQNG1kpiAEtkM34ayCFJHUJgoCz1QKBgDzdje3424q5rm8AxJWU
q3ppmYZlnJasav2kpz0QMGNtQABxO5E7OBMNQ44ZFq6E3ACNCqH8pgU0tZ5PELh0
rkQyQoaL2NFUSot1zrKhmwRMuLs64OmLQ/rm7pc8FcslG9/5aBOJyeQhtCPqgMb6
rTNMAzgQN0C7ETNGsmREOgA+
-----END PRIVATE KEY-----`

	testKeyEncrypted := `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIFHzBJBgkqhkiG9w0BBQ0wPDAbBgkqhkiG9w0BBQwwDgQIH/QBn6fJJTcCAggA
MB0GCWCGSAFlAwQBKgQQKxmw9YZRaQH5oHCdsLyK8ASCBNB5EZBKYTT6WImc71iW
wi/cpUEXe8HehEKHMMxFR1JaO+7hMwfZ1h8UPoTwFq+PDaUewrWihBFV20D+3otq
nRhNb2mOIN7gs/mSfPxoOcuBqRWG8+LE2Zxz127EmHp0cJkNR9GITJU4ZR0Q4qxe
SohxFSrZ56GTOmQM/vi8icx7R9GwxKPS3JR98KaZNigKG0m6y4RlcXJoDTbB6f5y
xm7xmkSlFeUlPx5HPC8gdHOlEGptxUcPnwge4SrMpATiR3D1LJWm8DgPDNraKZek
TWp2TbbIKW70dI0XkK0J0t536zUqCAlVxRfpmnYmE4gs2EWLWVUkxO4n0HAtWYHU
xjEaMerHViD6oCilPr7af9cA4HVivUjt1GDf8K9mcYfGH9HjftHCaAElGKr4tmxe
qEjxALt2qwOmXKnmjvVd1VFXgt7l9k+IpuDLHp49ikftj6kK1twuRxp/B09Yjpq/
VrYd2Q6K9ulhbaX3kPms/BsIW947jxg5yqxtXFSPSWavFFNj9t1k4ImdirLVCQIY
JSRkUyPNMrMTyMBuOwDjW/tk7NRCK6bVpi/OTzPzYhZFIVF98qP4kT492/9pzghU
ICKjw4tj3vaaiO+akXk9AiUobNDYo+2fAkgKbgPHlMph/Rbgi5fAks4SFPNm25Iw
T/imkFApiirMeiPUKAwKf0n1mPbe2BeCXieoqYIno6r3efvE7V60omMovDI2oaGv
SbztTEaQdHzbGps6q9xQKEdSwBon9Au4+JKvDltr1wOgJaTcDKc0I5o45MvtNqCs
Zyow3marJQTEKZazbv6X+Pmh02X4q0CtlTBx2pik4Ca0oc93GpULUEhRwfCCVFGO
r+vcCcDb/IO4duyOxNsv9dEzBYHbPEZfHmH4swG+7n3wN+W+6slTzpyy17Do+/Be
B4fes7iD+kDCVU6utrqJA7Ea2YjM0xi+2EgYgRuoxgwCqFi0TNcWIdpDPImRDDLf
8hXSvnyuz09Qh7rz5NPcBBZIyFwft8tovYfrfnFjVZ47MrmCGEkTJWzzSp2ayVYD
hScpn89UqyIa6RKALTc8XPzK0f9y1lswaI8LEikIk2gNhU6TA+tA2dhQ8YhTZTCF
N8V970v5/dQ9Rqb0g7cp+a37ZBbueRfGimnimI9G1jf7Fu7HHugtrF7Xv6tl9yn+
wsh4j8xsqd8YmDRUUJZuiRD3yi1vJ94MQXtwiQQajCAnMi4/NIAkjLS7Gpe30pfS
LByPVxZK1J665ubUhT9Bzzdj72LW6NozHzJDb7/uEruxgHFUOa5Cq1qt5C/in0FE
GHFbANk/vULxoKSTIFaDrH7tOsTjS4HFReAQst0AXqt2BzNXEaUlduRoY3LRtFwf
o8PoZryioXb7BRC8DPSv4hMKXurhQSFon5s3G+bTIikrlNWSM3Z54iYV5bu1uDbd
ostmm+e2pw7WBiTd5iKDNp804ALXXPx2Shnw269qzTBTa/MAtAxChAvx6antbh1c
OtWlAYewEHg5HDFBERn7t497Lex54o/bW538g0Q0r9RBnXuUXP6FQLal8ll43ruo
UOxg1XGJ+PIc/XFrsFte8oF+Tp31NJjGtgYWuS8G1LQ59MlGsMRV3h4gW4Dp6LZa
N7zz5i2xfPWkNQ3pbD+vZmYsqg==
-----END ENCRYPTED PRIVATE KEY-----
`

	expectedSignature := "iimf7QvkUodn5vb9DTJ/vPjVFabVVKAvcZKRQHQXEGdfuQ9/kNEyecdKS6gC55yALND+20Om+nF/lVfi0EAP1MVcU/1cRecqmGDaZRp/Gz1D0ZH0pgNqEiCWqwd26UsAsxeszgMIXplAl0khb/tHSw5Ls4V3OuRLoZsTiTF7nMibA4qkjTTcm4Hcn/JZouRQUSjumoWuxEbS9G8MCIgi2arChovsYvRpiyRCXgW99hDRDiahvt6MRAhZZ1NhX6GCSgEptH/+36Te9jgD+s13FhRuryzZYSHtRmd2JwwxrLG5JRwm2vhDxhlRZEKRhGDVfxD61T4zHAe60ezM7D1rpA=="
	expectedSignatureEncrypted := "IBNAcgymjFO9ub9j/npc26UoIXfPMphIeFiv1LJVkAajdc6+bKDNsnGtA0Lgx5fXEZKjhzHhqQoY/hRJls2xuGjuQyxFv2vara64wRcsB4TV5cPuSlCe1I4Bgiz1XQ7a5pmtSAjoOe9jXuI0FKSFpAfaQVuQbN2+AINWrC4V/LZpKrLS0NsQ8bfJ3PJB+h6zWuqVEHgl6ZiVvf7hQL1ZmJIBMtOkPpsAaQDYpxzFhK+S01qW6oLqNcl6OodEiZOZhtaDhYtwr41irVxmLbXJZtTt5m6mQso4nGAPf6BZe26s6DauyR++ZaZT+enXq+7dCxI4EV+KyQ6bZffpu+57Uw=="

	var tests = []struct {
		existingHeaders map[string]string
		args            []string
		valid           string
	}{
		{
			map[string]string{},
			[]string{testKey, "", "hello world"},
			expectedSignature,
		},
		{
			map[string]string{},
			[]string{testKeyEncrypted, "password", "hello world"},
			expectedSignatureEncrypted,
		},
	}

	for _, tc := range tests {
		actual, err := SignStringRS256PKCS8(tc.existingHeaders, tc.args)
		if actual != tc.valid {
			t.Errorf("SignStringRS256PKCS8(%s): expected %v, actual %v", tc.args, tc.valid, actual)
		}
		if err != nil {
			t.Errorf("SignStringRS256PKCS8(%s): expected no error, got %v", tc.args, err)
		}
	}
}

func TestValidateSignStringRS256PKCS8(t *testing.T) {
	var tests = []struct {
		args  []string
		valid bool
	}{
		{[]string{}, false},
		{[]string{"key"}, false},
		{[]string{"key", "passphrase"}, false},
		{[]string{"key", ""}, false},
		{[]string{"key", "passphrase", "string"}, true},
		{[]string{"key", "", "string"}, true},
		{[]string{"key", "passphrase", "string", "string"}, true},
		{[]string{"key", "", "string", "string", "string", "string"}, true},
	}

	for _, tc := range tests {
		actual := validateSignStringRS256PKCS8(tc.args)
		if actual != tc.valid {
			t.Errorf("validateSignStringRS256PKCS8(%s): expected %v, actual %v", tc.args, tc.valid, actual)
		}
	}
}

func TestArgsToStringToSign(t *testing.T) {
	var tests = []struct {
		existingHeaders map[string]string
		args            []string
		expected        string
	}{
		{
			map[string]string{},
			[]string{""},
			"\n",
		},
		{
			map[string]string{},
			[]string{"my string to sign"},
			"my string to sign\n",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{"my string to sign"},
			"my string to sign\n",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{"x-previous-header"},
			"123\n",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{"x-previous-header", "my string"},
			"123\nmy string\n",
		},
	}

	for _, tc := range tests {
		actual := argsToStringToSign(tc.existingHeaders, tc.args)
		if actual != tc.expected {
			t.Errorf("argsToStringToSign(%s, %s): expected %s, actual %s", tc.existingHeaders, tc.args, tc.expected, actual)
		}
	}
}
