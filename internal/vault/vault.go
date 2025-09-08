package vault

import (
	"database/sql"
	"errors"

	"crypto/rand"

	"golang.org/x/crypto/nacl/secretbox"
)

// VaultSecret represents a stored privileged credential
type VaultSecret struct {
	ID     int
	Name   string
	Data   []byte // encrypted
	UserID int    // who added it
}

// Encryption key (32 bytes) — replace with env var in production
var key [32]byte

func init() {
	// Generate random key for demo purposes (in prod, load from env or vault)
	if _, err := rand.Read(key[:]); err != nil {
		panic(err)
	}
}

// EncryptSecret encrypts plaintext using secretbox
func EncryptSecret(plaintext string) ([]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}
	encrypted := secretbox.Seal(nonce[:], []byte(plaintext), &nonce, &key)
	return encrypted, nil
}

// DecryptSecret decrypts data using secretbox
func DecryptSecret(data []byte) (string, error) {
	if len(data) < 24 {
		return "", errors.New("invalid data")
	}
	var nonce [24]byte
	copy(nonce[:], data[:24])
	decrypted, ok := secretbox.Open(nil, data[24:], &nonce, &key)
	if !ok {
		return "", errors.New("decryption failed")
	}
	return string(decrypted), nil
}

// StoreSecret stores a new credential in the vault
func StoreSecret(db *sql.DB, userID int, name, secret string) error {
	encrypted, err := EncryptSecret(secret)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO vault(name, data, user_id) VALUES($1,$2,$3)", name, encrypted, userID)
	return err
}

// RetrieveSecret retrieves and decrypts a credential by name
func RetrieveSecret(db *sql.DB, name string) (string, error) {
	var data []byte
	err := db.QueryRow("SELECT data FROM vault WHERE name=$1", name).Scan(&data)
	if err != nil {
		return "", err
	}
	return DecryptSecret(data)
}
