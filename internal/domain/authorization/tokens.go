package authorization

import (
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var (
	atLifetime = 15 * time.Minute
	rtLifetime = 7 * 24 * time.Hour
)

type CustomClaims struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

func IsTokenExpired(claims *CustomClaims) bool {
	if claims.ExpiresAt == nil {
		return true
	}

	return claims.ExpiresAt.Before(time.Now().UTC())
}

func ParseToken(tokenString string) (*jwt.Token, *CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, fmt.Errorf("токен не валиден")
	}

	return token, claims, nil
}

func deleteByRefreshToken(db *sql.DB, rt string) error {
	_, err := db.Exec("DELETE FROM auth_token WHERE refresh = ?", rt)

	if err != nil {
		return err
	}

	return nil
}

func createTokens(db *sql.DB, userId, firstName, lastName string) (string, string, error) {
	issuedAt := time.Now()
	atExp := issuedAt.Add(atLifetime)
	rtExp := issuedAt.Add(rtLifetime)

	claims := CustomClaims{
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			Audience:  jwt.ClaimStrings{"web-front"},
			ExpiresAt: jwt.NewNumericDate(time.Unix(atExp.Unix(), 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(issuedAt.Unix(), 0)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenStr, err := accessToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Unix(rtExp.Unix(), 0)),
	})
	refreshTokenStr, err := refreshToken.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		return "", "", err
	}

	_, err = db.Exec(`
INSERT INTO auth_token (expired_at, access, refresh) VALUES (?, ?, ?)`,
		atExp.Unix(),
		accessTokenStr,
		refreshTokenStr,
	)

	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, err
}
