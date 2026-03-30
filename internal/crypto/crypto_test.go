package crypto

import (
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
	}{
		{"empty string", ""},
		{"simple text", "hello world"},
		{"api key format", "sk-1234567890abcdef"},
		{"unicode", "こんにちは世界"},
		{"long text", "This is a longer piece of text that should also encrypt and decrypt correctly without any issues."},
		{"special chars", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := Encrypt(tc.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			decrypted, err := Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if decrypted != tc.plaintext {
				t.Errorf("Round trip failed: got %q, want %q", decrypted, tc.plaintext)
			}
		})
	}
}

func TestEncrypt_DifferentCiphertext(t *testing.T) {
	// Each encryption should produce different ciphertext due to random nonce
	plaintext := "test-api-key-12345"

	encrypted1, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("First Encrypt failed: %v", err)
	}

	encrypted2, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Second Encrypt failed: %v", err)
	}

	if encrypted1 == encrypted2 {
		t.Error("Two encryptions of same plaintext should produce different ciphertext")
	}

	// Both should decrypt to same plaintext
	decrypted1, _ := Decrypt(encrypted1)
	decrypted2, _ := Decrypt(encrypted2)

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Error("Both ciphertexts should decrypt to original plaintext")
	}
}

func TestDecrypt_InvalidInput(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"invalid base64", "not-valid-base64!!!"},
		{"too short", "YWJj"}, // "abc" in base64, too short for nonce
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Decrypt(tc.input)
			if err == nil {
				t.Error("Expected error for invalid input")
			}
		})
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	plaintext := "sensitive-data"
	encrypted, _ := Encrypt(plaintext)

	// Tamper with the ciphertext by changing a character
	tampered := []byte(encrypted)
	if len(tampered) > 10 {
		tampered[10] = 'X'
	}

	_, err := Decrypt(string(tampered))
	if err == nil {
		t.Error("Expected error for tampered ciphertext")
	}
}

func TestGetMachineID_Consistent(t *testing.T) {
	// Machine ID should be consistent across calls
	id1, err := getMachineID()
	if err != nil {
		t.Fatalf("First getMachineID failed: %v", err)
	}

	id2, err := getMachineID()
	if err != nil {
		t.Fatalf("Second getMachineID failed: %v", err)
	}

	if id1 != id2 {
		t.Error("Machine ID should be consistent across calls")
	}

	if id1 == "" {
		t.Error("Machine ID should not be empty")
	}
}

func TestDeriveKey_Consistent(t *testing.T) {
	key1, err := deriveKey()
	if err != nil {
		t.Fatalf("First deriveKey failed: %v", err)
	}

	key2, err := deriveKey()
	if err != nil {
		t.Fatalf("Second deriveKey failed: %v", err)
	}

	if len(key1) != 32 {
		t.Errorf("Key should be 32 bytes, got %d", len(key1))
	}

	if string(key1) != string(key2) {
		t.Error("Derived key should be consistent")
	}
}
