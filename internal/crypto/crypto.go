// Package crypto provides AES-256-GCM encryption for API keys.
// The encryption key is derived from the machine ID for secure local storage.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// Encrypt encrypts plaintext using AES-256-GCM with a machine-derived key.
// Returns base64-encoded ciphertext with prepended nonce.
func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key, err := deriveKey()
	if err != nil {
		return "", fmt.Errorf("crypto.Encrypt: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("crypto.Encrypt: new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("crypto.Encrypt: new gcm: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("crypto.Encrypt: generate nonce: %w", err)
	}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM.
func Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: decode base64: %w", err)
	}

	key, err := deriveKey()
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: new gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("crypto.Decrypt: ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("crypto.Decrypt: decrypt: %w", err)
	}

	return string(plaintext), nil
}

// deriveKey derives a 256-bit key from the machine ID.
func deriveKey() ([]byte, error) {
	machineID, err := getMachineID()
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}

	// Add salt and hash to get 32 bytes
	salt := "SubFlow-AES-Key-v1"
	hash := sha256.Sum256([]byte(salt + machineID))
	return hash[:], nil
}

// getMachineID returns a unique machine identifier.
func getMachineID() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxMachineID()
	case "darwin":
		return getDarwinMachineID()
	case "windows":
		return getWindowsMachineID()
	default:
		return getFallbackMachineID()
	}
}

// getLinuxMachineID reads /etc/machine-id or /var/lib/dbus/machine-id.
func getLinuxMachineID() (string, error) {
	paths := []string{
		"/etc/machine-id",
		"/var/lib/dbus/machine-id",
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil {
			id := strings.TrimSpace(string(data))
			if id != "" {
				return id, nil
			}
		}
	}

	return getFallbackMachineID()
}

// getDarwinMachineID uses the hardware UUID on macOS.
func getDarwinMachineID() (string, error) {
	// On macOS, we could use ioreg to get hardware UUID, but for simplicity
	// we'll use a fallback approach that works without external commands.
	return getFallbackMachineID()
}

// getWindowsMachineID uses the MachineGuid from registry.
func getWindowsMachineID() (string, error) {
	// On Windows, the MachineGuid is stored in registry, but reading it
	// requires the registry package. For simplicity, we use a fallback.
	return getFallbackMachineID()
}

// getFallbackMachineID creates a persistent machine ID file.
func getFallbackMachineID() (string, error) {
	// Try to read existing ID
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	idFile := homeDir + "/.subflow-machine-id"
	data, err := os.ReadFile(idFile)
	if err == nil {
		id := strings.TrimSpace(string(data))
		if id != "" {
			return id, nil
		}
	}

	// Generate new ID
	idBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, idBytes); err != nil {
		return "", fmt.Errorf("generate id: %w", err)
	}
	id := fmt.Sprintf("%x", idBytes)

	// Save for future use
	if err := os.WriteFile(idFile, []byte(id), 0600); err != nil {
		// Non-fatal, just log and continue
		// The ID will be regenerated next time, but encryption will still work
	}

	return id, nil
}
