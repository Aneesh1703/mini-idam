package auth

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateTOTPSecret generates a new TOTP secret for a user
func GenerateTOTPSecret(username string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MiniIDAM",
		AccountName: username,
	})
	if err != nil {
		return "", err
	}
	return key.Secret(), nil
}

// ValidateTOTPCode validates a given TOTP code against the secret
func ValidateTOTPCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateCurrentTOTPCode generates the current valid TOTP code (for testing)
func GenerateCurrentTOTPCode(secret string) (string, error) {
	code, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", err
	}
	return code, nil
}
