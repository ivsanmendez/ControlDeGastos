package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

// Claims are the JWT payload fields.
type Claims struct {
	UserID int64     `json:"uid"`
	Email  string    `json:"email"`
	Role   user.Role `json:"role"`
	jwtlib.RegisteredClaims
}

// Issuer implements user.TokenIssuer and provides JWT parsing.
type Issuer struct {
	secret []byte
	ttl    time.Duration
}

func NewIssuer(secret string) *Issuer {
	return &Issuer{
		secret: []byte(secret),
		ttl:    15 * time.Minute,
	}
}

func (i *Issuer) Issue(userID int64, email string, role user.Role) (user.TokenPair, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(i.ttl)),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(i.secret)
	if err != nil {
		return user.TokenPair{}, fmt.Errorf("sign jwt: %w", err)
	}

	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return user.TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return user.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: hex.EncodeToString(refreshBytes),
	}, nil
}

// Parse validates a JWT string and returns the claims.
func (i *Issuer) Parse(tokenString string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return i.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
