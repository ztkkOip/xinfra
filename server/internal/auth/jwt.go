package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint64 `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"admin"`
	jwt.RegisteredClaims
}

type OAuthCodeClaims struct {
	UserID      uint64 `json:"uid"`
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	Scope       string `json:"scope"`
	Nonce       string `json:"nonce,omitempty"`
	jwt.RegisteredClaims
}

type IDTokenClaims struct {
	UserID        uint64 `json:"uid"`
	Username      string `json:"preferred_username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Nonce         string `json:"nonce,omitempty"`
	jwt.RegisteredClaims
}

func Sign(secret, issuer string, ttl time.Duration, userID uint64, username, email string, isAdmin bool) (string, string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)
	tokenID := randomID()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Issuer:    issuer,
			Subject:   username,
			Audience:  []string{"authserver"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return signed, tokenID, expiresAt, err
}

func SignOAuthCode(secret, issuer string, ttl time.Duration, userID uint64, clientID, redirectURI, scope, nonce string) (string, string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)
	codeID := randomID()
	claims := OAuthCodeClaims{
		UserID:      userID,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		Nonce:       nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        codeID,
			Issuer:    issuer,
			Subject:   clientID,
			Audience:  []string{"oauth-code"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return signed, codeID, expiresAt, err
}

func SignIDToken(privateKey *rsa.PrivateKey, keyID, issuer, clientID string, ttl time.Duration, userID uint64, username, email, name, nonce string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)
	subject := email
	if subject == "" {
		subject = username
	}
	claims := IDTokenClaims{
		UserID:        userID,
		Username:      username,
		Email:         email,
		EmailVerified: email != "",
		Name:          name,
		Nonce:         nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   subject,
			Audience:  []string{clientID},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if keyID != "" {
		token.Header["kid"] = keyID
	}
	signed, err := token.SignedString(privateKey)
	return signed, expiresAt, err
}

func Parse(secret string, tokenValue string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected jwt signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt claims")
	}
	return claims, nil
}

func ParseOAuthCode(secret string, tokenValue string) (*OAuthCodeClaims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &OAuthCodeClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected jwt signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*OAuthCodeClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid oauth code claims")
	}
	if !claimStringsContains(claims.Audience, "oauth-code") {
		return nil, errors.New("invalid oauth code audience")
	}
	return claims, nil
}

func claimStringsContains(values jwt.ClaimStrings, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
