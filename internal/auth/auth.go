package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("supersecretkey") // Replace with env var in production

type User struct {
	ID           int
	Username     string
	PasswordHash string
	TOTPSecret   string
	Role         string
}

// Register a new user with username, password, role
func Register(db *sql.DB, username, password, role string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate TOTP secret
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MiniIDAM",
		AccountName: username,
	})
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     username,
		PasswordHash: string(hash),
		TOTPSecret:   secret.Secret(),
		Role:         role,
	}

	// Insert into DB
	_, err = db.Exec("INSERT INTO users(username, password_hash, totp_secret, role) VALUES($1,$2,$3,$4)",
		user.Username, user.PasswordHash, user.TOTPSecret, user.Role)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login with username, password, and TOTP code
func Login(db *sql.DB, username, password, totpCode string) (string, error) {
	var user User
	err := db.QueryRow("SELECT id, password_hash, totp_secret, role FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.PasswordHash, &user.TOTPSecret, &user.Role)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Check TOTP
	if !totp.Validate(totpCode, user.TOTPSecret) {
		return "", errors.New("invalid TOTP code")
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT parses token string and returns user claims
func ValidateJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
